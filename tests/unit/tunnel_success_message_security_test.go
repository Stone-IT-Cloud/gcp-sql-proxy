package unit

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelSuccessMessageUsesPlaceholders(t *testing.T) {
	got := proxy.TunnelSuccessMessage(5432, "project:region:instance", "user@example.com", false, false)
	if strings.Contains(strings.ToLower(got), "password: [leave empty]") {
		t.Fatalf("legacy password placeholder leaked: %s", got)
	}
	for _, token := range []string{"<db_password>", "<db_name>", "<db_user>"} {
		if !strings.Contains(got, token) {
			t.Fatalf("missing placeholder %q in %s", token, got)
		}
	}
}
