package unit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
	"golang.org/x/oauth2"
)

var authGlobalsMu sync.Mutex

func TestGetClientWithoutCredentials(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	_, err := auth.GetClient(context.Background())
	if !errors.Is(err, auth.ErrMissingCredentials) {
		t.Fatalf("expected ErrMissingCredentials, got %v", err)
	}
}

func TestBuildOAuthConfigHasExpectedRedirectAndScope(t *testing.T) {
	authGlobalsMu.Lock()
	defer authGlobalsMu.Unlock()

	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	cfg, err := auth.BuildOAuthConfig("http://localhost:8080")
	if err != nil {
		t.Fatalf("BuildOAuthConfig() error = %v", err)
	}
	if cfg.RedirectURL != "http://localhost:8080" {
		t.Fatalf("redirect url = %q, want http://localhost:8080", cfg.RedirectURL)
	}
	want := []string{
		"https://www.googleapis.com/auth/sqlservice.admin",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/sqlservice.login",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	if len(cfg.Scopes) != len(want) {
		t.Fatalf("scopes length = %d, want %d (%v)", len(cfg.Scopes), len(want), cfg.Scopes)
	}
	for i := range want {
		if cfg.Scopes[i] != want[i] {
			t.Fatalf("scopes[%d] = %q, want %q (all scopes: %v)", i, cfg.Scopes[i], want[i], cfg.Scopes)
		}
	}
}

func TestTokenFilePermissions(t *testing.T) {
	authGlobalsMu.Lock()
	defer authGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	path := filepath.Join(home, ".sql-proxy", "token.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	seedToken, err := json.Marshal(oauth2.Token{
		AccessToken:  "seed-access",
		RefreshToken: "seed-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("marshal token: %v", err)
	}
	if err := os.WriteFile(path, seedToken, 0o644); err != nil {
		t.Fatalf("write token file: %v", err)
	}

	auth.OAuthClientID = "test-client"
	auth.OAuthClientSecret = "test-secret"
	defer func() {
		auth.OAuthClientID = ""
		auth.OAuthClientSecret = ""
	}()

	restore := auth.SetTestHooks(
		func(url string) error { return nil },
		func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
			return &oauth2.Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				TokenType:    "Bearer",
				Expiry:       time.Now().Add(time.Hour),
			}, nil
		},
	)
	defer restore()

	client, err := auth.GetClient(context.Background())
	if err != nil {
		t.Fatalf("GetClient returned unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected GetClient to return a client")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat token: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestInvalidTokenRenamedBeforeReauth(t *testing.T) {
	authGlobalsMu.Lock()
	defer authGlobalsMu.Unlock()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	tokenPath := filepath.Join(home, ".sql-proxy", "token.json")
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("write invalid token: %v", err)
	}

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
			state := u.Query().Get("state")
			client := &http.Client{Timeout: 3 * time.Second}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/?code=test-code&state=%s", redirect, state), nil)
			if err != nil {
				return err
			}
			_, err = client.Do(req)
			return err
		},
		func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
			return &oauth2.Token{
				AccessToken:  "new-access",
				RefreshToken: "new-refresh",
				Expiry:       time.Now().Add(time.Hour),
			}, nil
		},
	)
	defer restore()

	client, err := auth.GetClient(context.Background())
	if err != nil {
		t.Fatalf("GetClient returned unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected GetClient to return a client")
	}

	entries, err := os.ReadDir(filepath.Dir(tokenPath))
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	foundInvalidBackup := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "token.json.invalid-") {
			foundInvalidBackup = true
			break
		}
	}
	if !foundInvalidBackup {
		t.Fatal("expected invalid token backup file to be created")
	}
}
