package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWorkflowBinaryNaming(t *testing.T) {
	path := filepath.Clean(filepath.Join("..", "..", ".github", "workflows", "build.yml"))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	text := string(content)
	required := []string{
		"binary_name=\"gcp-db-proxy-${GOOS}-${GOARCH}${ext}\"",
		"if [[ \"${GOOS}\" == \"windows\" ]]; then",
		"ext=\".exe\"",
		"path: dist/gcp-db-proxy-${{ matrix.target_os }}-${{ matrix.target_arch }}${{ matrix.target_os == 'windows' && '.exe' || '' }}",
	}
	for _, needle := range required {
		if !strings.Contains(text, needle) {
			t.Fatalf("missing binary naming behavior %q", needle)
		}
	}
}
