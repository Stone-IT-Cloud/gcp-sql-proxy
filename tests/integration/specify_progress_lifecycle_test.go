package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyProgressLifecycle(t *testing.T) {
	var out bytes.Buffer
	renderer := specify.NewPlainTextRenderer(&out)
	cmd := specify.NewCommand(renderer, []string{"spec.md"})
	phase := specify.ExecutionPhase{ID: "generation", Name: "generation", Type: specify.PhaseGeneration}

	ctx, cancel := context.WithCancel(context.Background())
	hbDone := make(chan struct{})
	go func() {
		specify.StartHeartbeat(ctx, phase, 10*time.Millisecond, func(event specify.ProgressEvent) {
			_ = renderer.Render(event)
		})
		close(hbDone)
	}()

	if err := cmd.EmitPhase(phase); err != nil {
		t.Fatalf("emit phase: %v", err)
	}
	skipped := specify.SkippedEvent(specify.ExecutionPhase{ID: "hooks", Name: "hooks", Type: specify.PhaseHook}, "optional")
	if err := renderer.Render(skipped); err != nil {
		t.Fatalf("render skipped: %v", err)
	}
	time.Sleep(25 * time.Millisecond)
	cancel()
	<-hbDone

	if err := cmd.EmitCompletion(); err != nil {
		t.Fatalf("completion: %v", err)
	}

	got := out.String()
	for _, needle := range []string{"[started]", "[heartbeat]", "[skipped]", "[completed] completion"} {
		if !strings.Contains(got, needle) {
			t.Fatalf("missing %q in output: %s", needle, got)
		}
	}
}
