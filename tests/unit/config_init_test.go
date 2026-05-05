package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
)

func TestInitWithArgsPrecedence(t *testing.T) {
	tests := []struct {
		name         string
		configBody   string
		args         []string
		wantPort     int
		wantInstance string
	}{
		{
			name:         "defaults when config missing",
			args:         []string{"--instance", "project:region:inst"},
			wantPort:     5432,
			wantInstance: "project:region:inst",
		},
		{
			name:         "config used when flags missing",
			configBody:   "port: 6000\ninstance: cfg:region:inst\n",
			args:         []string{},
			wantPort:     6000,
			wantInstance: "cfg:region:inst",
		},
		{
			name:         "flags override config",
			configBody:   "port: 6000\ninstance: cfg:region:inst\n",
			args:         []string{"--port", "7000", "--instance", "flag:region:inst"},
			wantPort:     7000,
			wantInstance: "flag:region:inst",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)

			if tt.configBody != "" {
				cfgDir := filepath.Join(home, ".sql-proxy")
				if err := os.MkdirAll(cfgDir, 0o755); err != nil {
					t.Fatalf("create cfg dir: %v", err)
				}
				cfgFile := filepath.Join(cfgDir, "config.yaml")
				if err := os.WriteFile(cfgFile, []byte(tt.configBody), 0o644); err != nil {
					t.Fatalf("write config: %v", err)
				}
			}

			got, err := config.InitWithArgs(tt.args)
			if err != nil {
				t.Fatalf("InitWithArgs() error = %v", err)
			}
			if got.Port != tt.wantPort {
				t.Fatalf("port = %d, want %d", got.Port, tt.wantPort)
			}
			if got.Instance != tt.wantInstance {
				t.Fatalf("instance = %q, want %q", got.Instance, tt.wantInstance)
			}
		})
	}
}
