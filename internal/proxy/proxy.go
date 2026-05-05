package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
)

// CloudSQLDialer is the minimal interface Start needs.
type CloudSQLDialer interface {
	Dial(ctx context.Context, instance string, opts ...cloudsqlconn.DialOption) (net.Conn, error)
	Close() error
}

var (
	principalEmailFn = auth.PrincipalEmail

	// dialerFactoryFn creates the cloudsqlconn dialer.
	// It is overridden in unit tests to avoid network/cloud dependencies.
	dialerFactoryFn = func(ctx context.Context, httpClient *http.Client, usePrivateIP bool) (CloudSQLDialer, error) {
		adminTokenSource, iamLoginTokenSource, err := auth.CloudSQLConnectorTokenSources(ctx)
		if err != nil {
			return nil, err
		}

		dialOpts := []cloudsqlconn.DialOption{}
		if usePrivateIP {
			dialOpts = append(dialOpts, cloudsqlconn.WithPrivateIP())
		}

		// IMPORTANT: even though httpClient is used for SQL Admin pre-flight,
		// IAM DB auth requires explicit token sources for the connector.
		return cloudsqlconn.NewDialer(
			ctx,
			cloudsqlconn.WithIAMAuthN(),
			cloudsqlconn.WithDefaultDialOptions(dialOpts...),
			cloudsqlconn.WithIAMAuthNTokenSources(adminTokenSource, iamLoginTokenSource),
			cloudsqlconn.WithHTTPClient(httpClient),
		)
	}

	instructionsWriter io.Writer = os.Stdout
)

// SetTestDialerFactory overrides the dialer factory used by Start.
func SetTestDialerFactory(fn func(ctx context.Context, httpClient *http.Client, usePrivateIP bool) (CloudSQLDialer, error)) (restore func()) {
	prev := dialerFactoryFn
	if fn != nil {
		dialerFactoryFn = fn
	}
	return func() {
		dialerFactoryFn = prev
	}
}

// SetTestPrincipalEmail overrides PrincipalEmail resolution for unit tests.
func SetTestPrincipalEmail(fn func(ctx context.Context, httpClient *http.Client) (string, error)) (restore func()) {
	prev := principalEmailFn
	if fn != nil {
		principalEmailFn = fn
	}
	return func() {
		principalEmailFn = prev
	}
}

// SetTestInstructionsWriter overrides the stdout writer for connection instructions.
func SetTestInstructionsWriter(w io.Writer) (restore func()) {
	prev := instructionsWriter
	if w != nil {
		instructionsWriter = w
	}
	return func() {
		instructionsWriter = prev
	}
}

// Start starts a local TCP listener accept loop and forwards each incoming
// connection to the target Cloud SQL instance using the Cloud SQL Go connector.
func Start(ctx context.Context, listener net.Listener, instance string, httpClient *http.Client, usePrivateIP bool) error {
	if ctx == nil {
		return errors.New("Start: missing ctx")
	}
	if listener == nil {
		return errors.New("Start: missing listener")
	}
	if instance == "" {
		return errors.New("Start: missing instance")
	}

	plan := DefaultDialerPlan(usePrivateIP)
	if !plan.UseIAMAuthN {
		return fmt.Errorf("Start: invalid dialer plan (IAM=%v, PrivateIP=%v)", plan.UseIAMAuthN, plan.UsePrivateIP)
	}

	// Resolve and emit operator-ready connection instructions before accepting clients.
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return fmt.Errorf("Start: listener must be TCP, got %T", listener.Addr())
	}

	email, err := principalEmailFn(ctx, httpClient)
	if err != nil {
		return err
	}

	// Ensure the dialer is closed on shutdown.
	dialer, err := dialerFactoryFn(ctx, httpClient, plan.UsePrivateIP)
	if err != nil {
		return err
	}
	defer func() { _ = dialer.Close() }()

	if err := emitTunnelReadyMessage(instructionsWriter, tcpAddr.Port, instance, email, false, plan.UsePrivateIP); err != nil {
		return fmt.Errorf("Start: write connection instructions: %w", err)
	}

	// Closing the listener on ctx cancellation unblocks Accept().
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = listener.Close()
		case <-done:
		}
	}()
	defer close(done)

	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			if ctx.Err() != nil || errors.Is(acceptErr, net.ErrClosed) {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Proxy listener accept error: %v\n", acceptErr)
			// Prevent tight CPU spin if the listener enters a persistent error state.
			time.Sleep(50 * time.Millisecond)
			continue
		}

		go func(local net.Conn) {
			defer local.Close()

			remote, dialErr := dialer.Dial(ctx, instance)
			if dialErr != nil {
				// Close local immediately; accept loop continues unaffected.
				return
			}
			defer remote.Close()

			// Emit full guidance on each successful reconnect/dial event.
			_ = emitTunnelReadyMessage(instructionsWriter, tcpAddr.Port, instance, email, true, plan.UsePrivateIP)

			// Ensure ctx cancellation terminates active relays promptly.
			connDone := make(chan struct{})
			defer close(connDone)
			go func() {
				select {
				case <-ctx.Done():
					_ = local.Close()
					_ = remote.Close()
				case <-connDone:
				}
			}()

			relayBidirectional(local, remote)
		}(conn)
	}
}

func relayBidirectional(local net.Conn, remote net.Conn) {
	// Bidirectional relay. This is intentionally symmetric and bounded by conn closure.
	var wg sync.WaitGroup
	var closeOnce sync.Once
	wg.Add(2)

	copyFn := func(dst, src net.Conn) {
		defer wg.Done()
		_, _ = io.Copy(dst, src)
		// Unblock the opposite copy direction if one side exits first.
		closeOnce.Do(func() {
			_ = local.Close()
			_ = remote.Close()
		})
	}

	go copyFn(local, remote)
	copyFn(remote, local)

	wg.Wait()
}
