package integration

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelSuccessMessageLifecycle(t *testing.T) {
	initial := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", false, false)
	reconnect := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", true, false)
	if !strings.Contains(initial, "Tunnel connection established.") {
		t.Fatalf("missing initial readiness message")
	}
	if !strings.Contains(reconnect, "Tunnel connection re-established.") {
		t.Fatalf("missing reconnect readiness message")
	}
}
