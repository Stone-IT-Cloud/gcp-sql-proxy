package unit

import (
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestPostgresClientTemplateOrdering(t *testing.T) {
	got := proxy.BuildPostgresClientTemplate("127.0.0.1", 5432, "appdb", "appuser")
	wantOrder := []string{"host=127.0.0.1", "port=5432", "dbname=appdb", "user=appuser"}
	last := -1
	for _, token := range wantOrder {
		idx := strings.Index(got, token)
		if idx == -1 {
			t.Fatalf("missing token %q in %q", token, got)
		}
		if idx < last {
			t.Fatalf("token %q appears out of order in %q", token, got)
		}
		last = idx
	}
}
