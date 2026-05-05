package unit

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelReadySuccessEventMessage(t *testing.T) {
	got := proxy.TunnelSuccessMessage(6543, "p:r:i", "user@example.com", false, false)
	if !strings.Contains(got, "Tunnel connection established.") {
		t.Fatalf("expected initial success message, got: %s", got)
	}
}
