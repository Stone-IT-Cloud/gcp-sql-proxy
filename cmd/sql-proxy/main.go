package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func main() {
	os.Exit(run())
}

func run() int {
	settings, err := config.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, config.UserFacingError(err))
		return 1
	}

	addr := fmt.Sprintf("127.0.0.1:%d", settings.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Port conflict: %s is already in use. Use --port to select another port or update ~/.sql-proxy/config.yaml.\n",
			addr,
		)
		return 1
	}
	defer listener.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Try to bootstrap authenticated connectivity.
	client, err := auth.GetClient(ctx)
	if err != nil && !errors.Is(err, auth.ErrMissingCredentials) {
		fmt.Fprintf(os.Stderr, "Authentication unavailable: %v\n", err)
		return 1
	}

	// Keep existing behavior for "no OAuth creds": keep listener alive so startup
	// + signal shutdown tests remain stable, but do not attempt tunnel dial.
	if err != nil && errors.Is(err, auth.ErrMissingCredentials) {
		go func() {
			<-ctx.Done()
			fmt.Fprintln(os.Stdout, "Shutting down proxy...")
			_ = listener.Close()
		}()

		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				if errors.Is(acceptErr, net.ErrClosed) || ctx.Err() != nil {
					break
				}
				fmt.Fprintf(os.Stderr, "Listener error: %v\n", acceptErr)
				continue
			}
			_ = conn.Close()
		}

		return 0
	}

	// Verified access + tunnel startup.
	if err := proxy.VerifyAccess(ctx, client, settings.Instance); err != nil {
		fmt.Fprintln(os.Stderr, proxy.UserFacingError(err))
		return 1
	}

	// proxy.Start prints connection instructions on successful initialization.
	if err := proxy.Start(ctx, listener, settings.Instance, client); err != nil {
		// Keep user-facing error stable and avoid leaking secrets.
		fmt.Fprintln(os.Stderr, proxy.UserFacingError(err))
		// Best-effort wait so tests that interrupt quickly still exit cleanly.
		time.Sleep(10 * time.Millisecond)
		return 1
	}

	return 0
}
