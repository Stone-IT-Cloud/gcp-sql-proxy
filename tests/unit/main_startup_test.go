package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadmeDocumentsMissingToolRemediation(t *testing.T) {
	path := filepath.Clean(filepath.Join("..", "..", "README.md"))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	text := string(content)
	required := []string{
		"Missing required tools are treated as failures",
		"goimports",
		"golangci-lint",
	}
	for _, needle := range required {
		if !strings.Contains(text, needle) {
			t.Fatalf("README.md missing required guidance %q", needle)
		}
	}
}
