package unit

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestStartInitializesDialerUsingInjectedFactory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := newConnChanListener(55435)

	restoreWriter := proxy.SetTestInstructionsWriter(io.Discard)
	defer restoreWriter()

	restoreEmail := proxy.SetTestPrincipalEmail(func(_ context.Context, _ *http.Client) (string, error) {
		return "user@example.com", nil
	})
	defer restoreEmail()

	httpClient := &http.Client{}
	calledCh := make(chan struct{}, 1)

	restoreDialer := proxy.SetTestDialerFactory(func(givenCtx context.Context, givenClient *http.Client) (proxy.CloudSQLDialer, error) {
		_ = givenCtx
		if givenClient != httpClient {
			t.Fatalf("httpClient pointer mismatch: got %p, want %p", givenClient, httpClient)
		}
		calledCh <- struct{}{}
		return &dummyDialer{}, nil
	})
	defer restoreDialer()

	done := make(chan error, 1)
	go func() { done <- proxy.Start(ctx, listener, "project:region:instance", httpClient) }()

	select {
	case <-calledCh:
		// OK
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for dialer factory to be called")
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected Start to return nil after cancellation, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Start shutdown")
	}
}
