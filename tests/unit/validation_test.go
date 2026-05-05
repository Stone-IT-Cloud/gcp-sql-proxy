package unit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
)

func TestInitWithArgsValidationFailures(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		cfgBody string
		wantErr error
	}{
		{
			name:    "missing instance",
			args:    []string{},
			wantErr: config.ErrMissingInstance,
		},
		{
			name:    "invalid port",
			args:    []string{"--instance", "x:y:z", "--port", "70000"},
			wantErr: config.ErrInvalidPort,
		},
		{
			name:    "malformed config",
			args:    []string{"--instance", "x:y:z"},
			cfgBody: "port: [",
			wantErr: config.ErrMalformedConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)

			if tt.cfgBody != "" {
				cfgDir := filepath.Join(home, ".sql-proxy")
				if err := os.MkdirAll(cfgDir, 0o755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(tt.cfgBody), 0o644); err != nil {
					t.Fatalf("write cfg: %v", err)
				}
			}

			_, err := config.InitWithArgs(tt.args)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserFacingErrorMessages(t *testing.T) {
	cases := []struct {
		err      error
		contains string
	}{
		{err: config.ErrMissingInstance, contains: "Missing instance"},
		{err: config.ErrInvalidPort, contains: "Invalid port"},
		{err: config.ErrMalformedConfig, contains: "Invalid configuration"},
	}

	for _, tc := range cases {
		got := config.UserFacingError(tc.err)
		if !strings.Contains(got, tc.contains) {
			t.Fatalf("message %q does not contain %q", got, tc.contains)
		}
	}
}
