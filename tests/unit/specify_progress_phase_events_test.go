package unit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyPhaseEvents(t *testing.T) {
	var out bytes.Buffer
	cmd := specify.NewCommand(specify.NewPlainTextRenderer(&out), []string{"spec.md"})
	phase := specify.ExecutionPhase{ID: "hooks", Name: "hooks", Type: specify.PhaseHook}
	if err := cmd.EmitPhase(phase); err != nil {
		t.Fatalf("emit phase: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "[started] hooks") || !strings.Contains(got, "[completed] hooks") {
		t.Fatalf("missing phase transitions: %s", got)
	}
}
