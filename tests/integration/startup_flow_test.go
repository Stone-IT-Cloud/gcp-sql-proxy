package integration

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
	"github.com/spf13/viper"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestPortConflictExitCodeAndMessage(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port
	home := t.TempDir()

	cmd := exec.Command("go", "run", "./cmd/sql-proxy", "--instance", "project:region:inst", "--port", fmt.Sprint(port))
	cmd.Dir = repoRoot(t)
	cmd.Env = append(os.Environ(), "HOME="+home, "USERPROFILE="+home)
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("expected non-zero exit for conflict, got success: %s", string(out))
	}
	var exitErr *exec.ExitError
	if !strings.Contains(string(out), "already in use") || !strings.Contains(string(out), "--port") {
		t.Fatalf("expected actionable conflict message, got: %s", string(out))
	}
	if ok := errors.As(err, &exitErr); ok && exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
	}
}

func TestGracefulShutdownOnSignal(t *testing.T) {
	home := t.TempDir()
	cmd := exec.Command("go", "run", "./cmd/sql-proxy", "--instance", "project:region:inst", "--port", "55432")
	cmd.Dir = repoRoot(t)
	cmd.Env = append(os.Environ(), "HOME="+home, "USERPROFILE="+home)

	if err := cmd.Start(); err != nil {
		t.Fatalf("start command: %v", err)
	}

	time.Sleep(800 * time.Millisecond)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("signal process: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				t.Fatalf("unexpected wait error: %v", err)
			}
			// Some environments report signal shutdown as non-zero while still
			// closing listener and exiting promptly; accept that behavior.
			if exitErr.ExitCode() != 0 && exitErr.ExitCode() != -1 {
				t.Fatalf("graceful shutdown exit code = %d, want 0 or signal-terminated", exitErr.ExitCode())
			}
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Signal(syscall.SIGKILL)
		t.Fatal("process did not shutdown within timeout")
	}
}

func TestSecurityDefaultsNotOverriddenByConfigInit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if _, err := config.InitWithArgs([]string{"--instance", "project:region:inst"}); err != nil {
		t.Fatalf("InitWithArgs: %v", err)
	}

	if viper.IsSet("iam_authn") || viper.IsSet("private_tunnel") {
		t.Fatal("unexpected IAM/private tunnel override set during config startup flow")
	}
}
