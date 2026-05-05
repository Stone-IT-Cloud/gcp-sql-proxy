package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrecommitLocalHooksBlockOnFailure(t *testing.T) {
	content, err := os.ReadFile(filepath.Clean(filepath.Join("..", "..", ".pre-commit-config.yaml")))
	if err != nil {
		t.Fatalf("read .pre-commit-config.yaml: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "id: go-test") || !strings.Contains(text, "entry: go test -short ./...") {
		t.Fatal("go-test hook not configured as required")
	}
	if !strings.Contains(text, "id: go-mod-tidy") || !strings.Contains(text, "entry: go mod tidy") {
		t.Fatal("go-mod-tidy hook not configured as required")
	}
	if !strings.Contains(text, "pass_filenames: false") {
		t.Fatal("expected pass_filenames: false for local system hooks")
	}
}
