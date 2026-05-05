package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWorkflowHasRequiredTriggersAndMatrix(t *testing.T) {
	path := filepath.Clean(filepath.Join("..", "..", ".github", "workflows", "build.yml"))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	text := string(content)
	required := []string{
		"push:",
		"pull_request:",
		"workflow_dispatch:",
		"branches:",
		"- main",
		"tags:",
		"- \"v*\"",
		"strategy:",
		"fail-fast: false",
		"runs-on: ${{ matrix.os }}",
		"target_os: linux",
		"target_os: windows",
		"target_os: darwin",
		"target_arch: amd64",
		"target_arch: arm64",
	}
	for _, needle := range required {
		if !strings.Contains(text, needle) {
			t.Fatalf("workflow missing %q", needle)
		}
	}
}
