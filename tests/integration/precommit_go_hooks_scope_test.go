package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrecommitGoHookScopingAndCoverage(t *testing.T) {
	content, err := os.ReadFile(filepath.Clean(filepath.Join("..", "..", ".pre-commit-config.yaml")))
	if err != nil {
		t.Fatalf("read .pre-commit-config.yaml: %v", err)
	}

	text := string(content)
	required := []string{
		"id: golangci-lint",
		"id: gofmt",
		"id: goimports",
		"id: go-mod-tidy",
		"id: go-test",
		"files: '^(go\\.mod|go\\.sum|.*\\.go)$'",
	}
	for _, needle := range required {
		if !strings.Contains(text, needle) {
			t.Fatalf("missing required Go hook content %q", needle)
		}
	}
}
