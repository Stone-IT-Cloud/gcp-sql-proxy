package unit

import (
	"errors"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyFailureEventIncludesGuidance(t *testing.T) {
	phase := specify.ExecutionPhase{ID: "generation", Name: "generation", Type: specify.PhaseGeneration}
	ev := specify.FailureEvent(phase, errors.New("boom"))
	if ev.Type != specify.EventFailed {
		t.Fatalf("type = %s, want failed", ev.Type)
	}
	if !strings.Contains(ev.Reason, "rerun /speckit-specify") {
		t.Fatalf("missing next-step guidance: %q", ev.Reason)
	}
}
