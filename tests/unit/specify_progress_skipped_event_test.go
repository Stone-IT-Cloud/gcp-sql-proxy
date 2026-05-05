package unit

import (
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifySkippedEvent(t *testing.T) {
	phase := specify.ExecutionPhase{ID: "hooks", Name: "hooks", Type: specify.PhaseHook}
	ev := specify.SkippedEvent(phase, "optional hook disabled")
	if ev.Type != specify.EventSkipped {
		t.Fatalf("type = %s, want skipped", ev.Type)
	}
	if ev.Reason == "" || !ev.Valid() {
		t.Fatalf("expected valid skipped event with reason")
	}
}
