package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyProgressHappyPath(t *testing.T) {
	var out bytes.Buffer
	cmd := specify.NewCommand(specify.NewPlainTextRenderer(&out), []string{"spec.md", "plan.md"})
	phases := []specify.ExecutionPhase{
		{ID: "hooks", Name: "hooks", Type: specify.PhaseHook},
		{ID: "generation", Name: "generation", Type: specify.PhaseGeneration},
	}
	if err := cmd.Run(context.Background(), phases); err != nil {
		t.Fatalf("run: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "[completed] completion: completed successfully") {
		t.Fatalf("missing completion summary: %s", got)
	}
}
