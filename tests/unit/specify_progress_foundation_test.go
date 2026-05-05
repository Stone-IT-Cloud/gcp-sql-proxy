package unit

import (
	"bytes"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

type failingRenderer struct{}

func (failingRenderer) Render(event specify.ProgressEvent) error {
	_ = event
	return assertErr("rich failed")
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

func TestProgressEventValidationAndFallback(t *testing.T) {
	phase := specify.ExecutionPhase{ID: "generation", Name: "generation", Type: specify.PhaseGeneration}
	ev := specify.NewEvent(specify.EventStarted, phase, "phase started")
	if !ev.Valid() {
		t.Fatal("expected event to be valid")
	}

	var out bytes.Buffer
	fallback := specify.NewFallbackRenderer(failingRenderer{}, specify.NewPlainTextRenderer(&out))
	if err := fallback.Render(ev); err != nil {
		t.Fatalf("fallback render: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("expected fallback output")
	}
}
