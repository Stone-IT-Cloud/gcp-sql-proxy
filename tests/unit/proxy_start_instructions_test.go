package unit

import (
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestConnectionInstructionsFormat(t *testing.T) {
	got := proxy.ConnectionInstructions(5432, "you@example.com")
	want := "Host: 127.0.0.1\nPort: 5432\nUser: you@example.com\nPassword: [LEAVE EMPTY]\n"
	if got != want {
		t.Fatalf("instructions mismatch.\nGot:\n%q\nWant:\n%q", got, want)
	}
}
