package unit

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestStartAcceptLoopTerminatesOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := newConnChanListener(55434)

	restoreWriter := proxy.SetTestInstructionsWriter(io.Discard)
	defer restoreWriter()

	restoreEmail := proxy.SetTestPrincipalEmail(func(ctx context.Context, _ *http.Client) (string, error) {
		_ = ctx
		return "user@example.com", nil
	})
	defer restoreEmail()

	// Avoid cloud dialer initialization.
	restoreDialer := proxy.SetTestDialerFactory(func(ctx context.Context, _ *http.Client, usePrivateIP bool) (proxy.CloudSQLDialer, error) {
		_ = ctx
		_ = usePrivateIP
		return &dummyDialer{}, nil
	})
	defer restoreDialer()

	startErrCh := make(chan error, 1)
	go func() { startErrCh <- proxy.Start(ctx, listener, "project:region:instance", &http.Client{}, false) }()

	// Cancel shortly after startup to force listener close and Accept() unblocking.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-startErrCh:
		if err != nil {
			t.Fatalf("expected Start to return nil on shutdown, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Start shutdown")
	}
}

type dummyDialer struct{}

func (d *dummyDialer) Dial(ctx context.Context, instance string, opts ...cloudsqlconn.DialOption) (net.Conn, error) {
	_ = ctx
	_ = instance
	_ = opts
	return nil, nil
}

func (d *dummyDialer) Close() error { return nil }
