package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrecommitIncludesBaselineHooks(t *testing.T) {
	content, err := os.ReadFile(filepath.Clean(filepath.Join("..", "..", ".pre-commit-config.yaml")))
	if err != nil {
		t.Fatalf("read .pre-commit-config.yaml: %v", err)
	}

	text := string(content)
	baselineHooks := []string{
		"id: trailing-whitespace",
		"id: end-of-file-fixer",
		"id: check-yaml",
	}
	for _, hook := range baselineHooks {
		if !strings.Contains(text, hook) {
			t.Fatalf("missing baseline hook %q", hook)
		}
	}
}
