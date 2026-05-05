package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/config"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var oauthAuthGlobalsMu sync.Mutex

func oauthRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestStartupWithoutOAuthCredentialsStillRuns(t *testing.T) {
	home := t.TempDir()
	cmd := exec.Command("go", "run", "./cmd/sql-proxy", "--instance", "project:region:instance", "--port", "55433")
	cmd.Dir = oauthRepoRoot(t)
	cmd.Env = append(os.Environ(), "HOME="+home, "USERPROFILE="+home)

	if err := cmd.Start(); err != nil {
		t.Fatalf("start command: %v", err)
	}

	time.Sleep(700 * time.Millisecond)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("interrupt process: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("process did not exit after interrupt")
	}
}

func TestOAuthPathDoesNotOverrideSecurityDefaults(t *testing.T) {
	oauthAuthGlobalsMu.Lock()
	defer oauthAuthGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if _, err := config.InitWithArgs([]string{"--instance", "project:region:instance"}); err != nil {
		t.Fatalf("config.InitWithArgs() error = %v", err)
	}
	if _, err := auth.GetClient(context.Background()); !errors.Is(err, auth.ErrMissingCredentials) {
		t.Fatalf("auth.GetClient() error = %v, want %v", err, auth.ErrMissingCredentials)
	}

	// Ensure auth bootstrap doesn't mutate unrelated security defaults in config state.
	if viper.IsSet("iam_authn") || viper.IsSet("private_tunnel") {
		t.Fatal("unexpected IAM/private tunnel override detected")
	}
}

func TestFirstRunOAuthFlowPersistsTokenAndShowsSuccess(t *testing.T) {
	oauthAuthGlobalsMu.Lock()
	defer oauthAuthGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	var callbackBody string
	restore := auth.SetTestHooks(
		func(authURL string) error {
			u, err := url.Parse(authURL)
			if err != nil {
				return err
			}
			redirect, _ := url.QueryUnescape(u.Query().Get("redirect_uri"))
			state := u.Query().Get("state")
			cb := fmt.Sprintf("%s/?code=test-code&state=%s", redirect, state)
			client := &http.Client{Timeout: 1 * time.Second}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, cb, nil)
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			callbackBody = string(body)
			return nil
		},
		func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
			return &oauth2.Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				Expiry:       time.Now().Add(time.Hour),
			}, nil
		},
	)
	defer restore()

	client, err := auth.GetClient(context.Background())
	if err != nil {
		t.Fatalf("GetClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("expected authenticated client, got nil")
	}
	if !strings.Contains(callbackBody, "Success! You can close this window.") {
		t.Fatalf("unexpected callback success body: %q", callbackBody)
	}
}

func TestValidTokenSkipsBrowserFlow(t *testing.T) {
	oauthAuthGlobalsMu.Lock()
	defer oauthAuthGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	tokenPath := filepath.Join(home, ".sql-proxy", "token.json")
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	token := oauth2.Token{
		AccessToken:  "cached-access",
		RefreshToken: "cached-refresh",
		Expiry:       time.Now().Add(time.Hour),
	}
	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("open token: %v", err)
	}
	if err := json.NewEncoder(f).Encode(&token); err != nil {
		_ = f.Close()
		t.Fatalf("encode token: %v", err)
	}
	_ = f.Close()

	var browserCalls int32
	restore := auth.SetTestHooks(
		func(authURL string) error {
			atomic.AddInt32(&browserCalls, 1)
			return nil
		},
		nil,
	)
	defer restore()

	client, err := auth.GetClient(context.Background())
	if err != nil {
		t.Fatalf("GetClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("expected authenticated client, got nil")
	}
	if atomic.LoadInt32(&browserCalls) != 0 {
		t.Fatalf("expected no browser flow, got %d calls", browserCalls)
	}
}

func TestCallbackStateMismatchFailsAuth(t *testing.T) {
	oauthAuthGlobalsMu.Lock()
	defer oauthAuthGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	restore := auth.SetTestHooks(
		func(authURL string) error {
			u, err := url.Parse(authURL)
			if err != nil {
				return err
			}
			redirect, _ := url.QueryUnescape(u.Query().Get("redirect_uri"))
			client := &http.Client{Timeout: 1 * time.Second}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/?code=test-code&state=wrong-state", redirect), nil)
			if err != nil {
				return err
			}
			_, err = client.Do(req)
			return err
		},
		func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
			return &oauth2.Token{AccessToken: "unused", Expiry: time.Now().Add(time.Hour)}, nil
		},
	)
	defer restore()

	_, err := auth.GetClient(context.Background())
	if !errors.Is(err, auth.ErrStateMismatch) {
		t.Fatalf("expected ErrStateMismatch, got %v", err)
	}
}

func TestCallbackPortFallbackFrom8080(t *testing.T) {
	oauthAuthGlobalsMu.Lock()
	defer oauthAuthGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	block, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Skipf("cannot reserve 8080 in current environment: %v", err)
	}
	defer block.Close()

	var fallbackPort int
	restore := auth.SetTestHooks(
		func(authURL string) error {
			u, err := url.Parse(authURL)
			if err != nil {
				return err
			}
			redirect, _ := url.QueryUnescape(u.Query().Get("redirect_uri"))
			ru, err := url.Parse(redirect)
			if err != nil {
				return err
			}
			if ru.Port() == "8080" {
				return errors.New("expected fallback callback port, got 8080")
			}
			fmt.Sscanf(ru.Port(), "%d", &fallbackPort)
			state := u.Query().Get("state")
			client := &http.Client{Timeout: 1 * time.Second}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/?code=test-code&state=%s", redirect, state), nil)
			if err != nil {
				return err
			}
			_, err = client.Do(req)
			return err
		},
		func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
			return &oauth2.Token{AccessToken: "fallback", Expiry: time.Now().Add(time.Hour)}, nil
		},
	)
	defer restore()

	client, err := auth.GetClient(context.Background())
	if err != nil {
		t.Fatalf("GetClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("expected authenticated client, got nil")
	}
	if fallbackPort == 0 || fallbackPort == 8080 {
		t.Fatalf("expected non-8080 fallback port, got %d", fallbackPort)
	}
}
