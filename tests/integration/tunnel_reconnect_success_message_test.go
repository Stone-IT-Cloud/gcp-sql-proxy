package integration

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelReconnectSuccessMessage(t *testing.T) {
	got := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", true, false)
	if !strings.Contains(got, "Tunnel connection re-established.") {
		t.Fatalf("missing reconnect marker: %s", got)
	}
}
