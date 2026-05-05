package integration

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelSuccessInstructionPayloadCompleteness(t *testing.T) {
	got := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", false, false)
	for _, token := range []string{
		"Host: 127.0.0.1",
		"Port: 5432",
		"Target instance: project:region:instance",
		"PostgreSQL client:",
	} {
		if !strings.Contains(got, token) {
			t.Fatalf("missing %q in payload: %s", token, got)
		}
	}
}
