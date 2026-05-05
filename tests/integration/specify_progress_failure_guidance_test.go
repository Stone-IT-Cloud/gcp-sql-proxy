package integration

import (
	"errors"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyProgressFailureGuidance(t *testing.T) {
	phase := specify.ExecutionPhase{ID: "validation", Name: "validation", Type: specify.PhaseValidation}
	ev := specify.FailureEvent(phase, errors.New("bad checklist"))
	if !strings.Contains(ev.Reason, "Fix validation issues") {
		t.Fatalf("unexpected guidance: %s", ev.Reason)
	}
}
