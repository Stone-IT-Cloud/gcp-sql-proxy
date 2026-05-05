package integration

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

type chanListener struct {
	addr   *net.TCPAddr
	conns  chan net.Conn
	closed chan struct{}
}

func newChanListener(port int) *chanListener {
	return &chanListener{
		addr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: port,
		},
		conns:  make(chan net.Conn, 10),
		closed: make(chan struct{}),
	}
}

func (l *chanListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.conns:
		return c, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}

func (l *chanListener) Close() error {
	select {
	case <-l.closed:
		return nil
	default:
		close(l.closed)
		return nil
	}
}

func (l *chanListener) Addr() net.Addr { return l.addr }

type singleDialDialer struct {
	remote net.Conn
}

func (d *singleDialDialer) Dial(ctx context.Context, instance string, opts ...any) (net.Conn, error) {
	_ = ctx
	_ = instance
	_ = opts
	return d.remote, nil
}

func (d *singleDialDialer) Close() error { return nil }

func TestProxyTunnelShutdownClosesRelayConnections(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	listener := newChanListener(55436)

	remoteServerConn, remoteDialerConn := net.Pipe()
	localClientConn, localStartConn := net.Pipe()

	restoreEmail := proxy.SetTestPrincipalEmail(func(_ context.Context, _ *http.Client) (string, error) {
		return "user@example.com", nil
	})
	defer restoreEmail()

	restoreWriter := proxy.SetTestInstructionsWriter(ioDiscard{})
	defer restoreWriter()

	restoreDialer := proxy.SetTestDialerFactory(func(_ context.Context, _ *http.Client) (proxy.CloudSQLDialer, error) {
		// Provide the dialed remote connection for the accepted local socket.
		return &shutdownDialer{remote: remoteDialerConn}, nil
	})
	defer restoreDialer()

	startErrCh := make(chan error, 1)
	go func() { startErrCh <- proxy.Start(ctx, listener, "project:region:instance", &http.Client{}) }()

	// Feed one accepted connection.
	listener.conns <- localStartConn

	// Give the relay goroutines a moment to start.
	time.Sleep(50 * time.Millisecond)

	cancel()

	select {
	case err := <-startErrCh:
		if err != nil {
			t.Fatalf("expected Start to return nil, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Start shutdown")
	}

	// Connections should be closed by ctx cancellation.
	_ = localClientConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1)
	_, readErr := localClientConn.Read(buf)
	if readErr == nil {
		t.Fatal("expected local client read to fail after shutdown")
	}

	_ = remoteServerConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, readErr = remoteServerConn.Read(buf)
	if readErr == nil {
		t.Fatal("expected remote server read to fail after shutdown")
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }

type shutdownDialer struct {
	remote net.Conn
}

func (d *shutdownDialer) Dial(ctx context.Context, instance string, opts ...cloudsqlconn.DialOption) (net.Conn, error) {
	_ = ctx
	_ = instance
	_ = opts
	return d.remote, nil
}

func (d *shutdownDialer) Close() error { return nil }
