package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWorkflowArtifactUploadConfiguration(t *testing.T) {
	path := filepath.Clean(filepath.Join("..", "..", ".github", "workflows", "build.yml"))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	text := string(content)
	required := []string{
		"uses: actions/upload-artifact@v6",
		"name: gcp-db-proxy-${{ matrix.target_os }}-${{ matrix.target_arch }}${{ matrix.target_os == 'windows' && '.exe' || '' }}",
		"retention-days: 14",
	}
	for _, needle := range required {
		if !strings.Contains(text, needle) {
			t.Fatalf("missing artifact behavior %q", needle)
		}
	}
}
