package unit

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelSuccessMessageFormat(t *testing.T) {
	got := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", false, false)
	for _, want := range []string{
		"Tunnel connection established.",
		"Target instance: project:region:instance",
		"IP connectivity mode: public",
		"Host: 127.0.0.1",
		"Port: 5432",
		"Password: <db_password>",
		"PostgreSQL client:",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in message: %s", want, got)
		}
	}
}
