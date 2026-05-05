package unit

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

type connChanListener struct {
	addr   *net.TCPAddr
	conns  chan net.Conn
	closed chan struct{}
}

func newConnChanListener(port int) *connChanListener {
	return &connChanListener{
		addr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: port,
		},
		conns:  make(chan net.Conn, 10),
		closed: make(chan struct{}),
	}
}

func (l *connChanListener) Accept() (net.Conn, error) {
	select {
	case <-l.closed:
		return nil, net.ErrClosed
	default:
	}

	select {
	case c := <-l.conns:
		return c, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}

func (l *connChanListener) Close() error {
	select {
	case <-l.closed:
	default:
		close(l.closed)
	}
	return nil
}

func (l *connChanListener) Addr() net.Addr {
	return l.addr
}

type dialOutcome struct {
	conn net.Conn
	err  error
}

type fakeDialer struct {
	mu       sync.Mutex
	outcomes []dialOutcome
	calls    int
}

func (d *fakeDialer) Dial(ctx context.Context, instance string, opts ...cloudsqlconn.DialOption) (net.Conn, error) {
	_ = ctx
	_ = instance
	_ = opts

	d.mu.Lock()
	defer d.mu.Unlock()
	if d.calls >= len(d.outcomes) {
		return nil, errors.New("unexpected dial call")
	}
	out := d.outcomes[d.calls]
	d.calls++
	return out.conn, out.err
}

func (d *fakeDialer) Close() error { return nil }

func waitUntil(cond func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

var _ proxy.CloudSQLDialer = (*fakeDialer)(nil)
