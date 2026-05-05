package integration

import (
	"bytes"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

type alwaysFailRenderer struct{}

func (alwaysFailRenderer) Render(event specify.ProgressEvent) error {
	_ = event
	return errFallback("render fail")
}

type errFallback string

func (e errFallback) Error() string { return string(e) }

func TestSpecifyProgressFallbackToPlainText(t *testing.T) {
	var out bytes.Buffer
	r := specify.NewFallbackRenderer(alwaysFailRenderer{}, specify.NewPlainTextRenderer(&out))
	phase := specify.ExecutionPhase{ID: "generation", Name: "generation", Type: specify.PhaseGeneration}
	ev := specify.NewEvent(specify.EventHeartbeat, phase, "still running")
	if err := r.Render(ev); err != nil {
		t.Fatalf("render: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("expected fallback plain text output")
	}
}
