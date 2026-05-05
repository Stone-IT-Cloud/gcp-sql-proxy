package unit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestStartRelayBidirectionalAndDialFailureIsolation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := newConnChanListener(55433)

	var instructions bytes.Buffer
	restoreWriter := proxy.SetTestInstructionsWriter(&instructions)
	defer restoreWriter()

	restoreEmail := proxy.SetTestPrincipalEmail(func(ctx context.Context, _ *http.Client) (string, error) {
		_ = ctx
		return "user@example.com", nil
	})
	defer restoreEmail()

	localClient1, localStart1 := net.Pipe()
	defer localClient1.Close()

	localClient2, localStart2 := net.Pipe()
	defer localClient2.Close()

	remoteServer2, remoteDialer2 := net.Pipe()
	defer remoteServer2.Close()

	dialer := &fakeDialer{
		outcomes: []dialOutcome{
			{conn: nil, err: errors.New("dial failed")},
			{conn: remoteDialer2, err: nil},
		},
	}

	restoreDialer := proxy.SetTestDialerFactory(func(ctx context.Context, httpClient *http.Client) (proxy.CloudSQLDialer, error) {
		_ = ctx
		_ = httpClient
		return dialer, nil
	})
	defer restoreDialer()

	// Start proxy.
	startErrCh := make(chan error, 1)
	go func() { startErrCh <- proxy.Start(ctx, listener, "project:region:instance", &http.Client{}) }()

	// Enqueue two local connections concurrently.
	listener.conns <- localStart1
	listener.conns <- localStart2

	// Wait for instruction block to be emitted.
	ok := waitUntil(func() bool {
		return strings.Contains(instructions.String(), "Password: [LEAVE EMPTY]") &&
			strings.Contains(instructions.String(), "User: user@example.com")
	}, 2*time.Second)
	if !ok {
		t.Fatal("timed out waiting for connection instructions")
	}

	// Verify bidirectional relay on the second connection continues even if the first dial fails.
	msgLocalToRemote := []byte("hello")
	msgRemoteToLocal := []byte("world")

	tryRelay := func(local net.Conn) bool {
		_ = remoteServer2.SetReadDeadline(time.Now().Add(750 * time.Millisecond))
		_ = local.SetReadDeadline(time.Now().Add(750 * time.Millisecond))

		remoteReadCh := make(chan []byte, 1)
		go func() {
			buf := make([]byte, len(msgLocalToRemote))
			_, err := io.ReadFull(remoteServer2, buf)
			if err != nil {
				remoteReadCh <- nil
				return
			}
			remoteReadCh <- buf
		}()

		if _, err := local.Write(msgLocalToRemote); err != nil {
			return false
		}

		var gotRemote []byte
		select {
		case gotRemote = <-remoteReadCh:
		case <-time.After(2 * time.Second):
			return false
		}
		if !bytes.Equal(gotRemote, msgLocalToRemote) {
			return false
		}

		// Reverse direction: remote -> local
		localReadCh := make(chan []byte, 1)
		go func() {
			buf := make([]byte, len(msgRemoteToLocal))
			_, err := io.ReadFull(local, buf)
			if err != nil {
				localReadCh <- nil
				return
			}
			localReadCh <- buf
		}()

		if _, err := remoteServer2.Write(msgRemoteToLocal); err != nil {
			return false
		}

		var gotLocal []byte
		select {
		case gotLocal = <-localReadCh:
		case <-time.After(2 * time.Second):
			return false
		}
		return bytes.Equal(gotLocal, msgRemoteToLocal)
	}

	ok = tryRelay(localClient1) || tryRelay(localClient2)
	if !ok {
		t.Fatal("expected at least one relay direction to succeed despite dial failure")
	}

	// Trigger shutdown and ensure Start exits.
	cancel()
	select {
	case err := <-startErrCh:
		if err != nil {
			t.Fatalf("Start returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Start shutdown")
	}

	out := instructions.String()
	for _, mustContain := range []string{
		"Host: 127.0.0.1",
		"Port: 55433",
		"User: user@example.com",
		"Password: [LEAVE EMPTY]",
	} {
		if !strings.Contains(out, mustContain) {
			t.Fatalf("instructions missing %q; got:\n%s", mustContain, out)
		}
	}

}
