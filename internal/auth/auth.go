package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthScope         = "https://www.googleapis.com/auth/sqlservice.admin"
	defaultRedirectURL = "http://localhost:8080"
	defaultTokenName   = "token.json"
	authFlowTimeout    = 2 * time.Minute
)

var (
	// Inject these via ldflags or environment in deployment.
	OAuthClientID     = ""
	OAuthClientSecret = ""
)

var (
	ErrMissingCode        = errors.New("missing authorization code")
	ErrStateMismatch      = errors.New("oauth state mismatch")
	ErrMissingCredentials = errors.New("missing oauth client credentials")
	ErrAuthTimeout        = errors.New("oauth authentication timed out")
)

var hooksMu sync.RWMutex

var openBrowserFn = openBrowser
var exchangeTokenFn = func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
	return cfg.Exchange(ctx, code)
}

// SetTestHooks overrides browser and token-exchange behaviors for tests.
// It returns a restore function that resets defaults.
func SetTestHooks(
	openFn func(string) error,
	exchangeFn func(context.Context, *oauth2.Config, string) (*oauth2.Token, error),
) func() {
	hooksMu.Lock()
	defer hooksMu.Unlock()

	originalOpen := openBrowserFn
	originalExchange := exchangeTokenFn
	if openFn != nil {
		openBrowserFn = openFn
	}
	if exchangeFn != nil {
		exchangeTokenFn = exchangeFn
	}
	return func() {
		hooksMu.Lock()
		defer hooksMu.Unlock()
		openBrowserFn = originalOpen
		exchangeTokenFn = originalExchange
	}
}

func oauthConfig(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     OAuthClientID,
		ClientSecret: OAuthClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{oauthScope},
		RedirectURL:  redirectURL,
	}
}

// BuildOAuthConfig returns the runtime OAuth configuration for authentication flow setup.
func BuildOAuthConfig(redirectURL string) (*oauth2.Config, error) {
	if !hasCredentials() {
		return nil, ErrMissingCredentials
	}
	return oauthConfig(redirectURL), nil
}

func hasCredentials() bool {
	return OAuthClientID != "" && OAuthClientSecret != ""
}

func tokenPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".sql-proxy", defaultTokenName), nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	if err := json.NewDecoder(f).Decode(&tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func ensureTokenPermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Mode().Perm() != 0o600 {
		if err := os.Chmod(path, 0o600); err != nil {
			return fmt.Errorf("set token file permissions: %w", err)
		}
	}
	return nil
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create token directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open token file: %w", err)
	}
	if err := json.NewEncoder(f).Encode(token); err != nil {
		_ = f.Close()
		return fmt.Errorf("encode token: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close token file: %w", err)
	}
	if err := ensureTokenPermissions(path); err != nil {
		return err
	}
	return nil
}

func markInvalidToken(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	backup := fmt.Sprintf("%s.invalid-%d", path, time.Now().Unix())
	if err := os.Rename(path, backup); err == nil {
		return
	}
	_ = os.Remove(path)
}

func generateState() (string, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate oauth state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func getTokenFromWeb(ctx context.Context, cfg *oauth2.Config) (*oauth2.Token, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("start callback listener: %w", err)
		}
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	flowCfg := *cfg
	flowCfg.RedirectURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	state, err := generateState()
	if err != nil {
		return nil, err
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
			http.Error(w, "Authentication failed. You can close this window.", http.StatusBadRequest)
			select {
			case errCh <- fmt.Errorf("oauth callback error: %s", oauthErr):
			default:
			}
			return
		}
		receivedState := r.URL.Query().Get("state")
		if receivedState == "" || receivedState != state {
			http.Error(w, "Authentication failed due to invalid session state.", http.StatusBadRequest)
			select {
			case errCh <- ErrStateMismatch:
			default:
			}
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code.", http.StatusBadRequest)
			select {
			case errCh <- ErrMissingCode:
			default:
			}
			return
		}
		_, _ = io.WriteString(w, "Success! You can close this window.")
		select {
		case codeCh <- code:
		default:
		}
	})

	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(ln)
	}()

	authURL := flowCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	hooksMu.RLock()
	openFn := openBrowserFn
	exchangeFn := exchangeTokenFn
	hooksMu.RUnlock()

	if err := openFn(authURL); err != nil {
		fmt.Fprintf(os.Stderr, "Could not open browser automatically. Open this URL manually:\n%s\n", authURL)
	}

	var code string
	timer := time.NewTimer(authFlowTimeout)
	defer timer.Stop()

	select {
	case code = <-codeCh:
	case err = <-errCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = server.Shutdown(shutdownCtx)
		cancel()
		return nil, err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = server.Shutdown(shutdownCtx)
		cancel()
		return nil, ctx.Err()
	case <-timer.C:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = server.Shutdown(shutdownCtx)
		cancel()
		return nil, ErrAuthTimeout
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = server.Shutdown(shutdownCtx)
	cancel()

	tok, err := exchangeFn(ctx, &flowCfg, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code for token: %w", err)
	}
	return tok, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform for auto-open: %s", runtime.GOOS)
	}
	return cmd.Start()
}

// GetClient returns an authenticated HTTP client when OAuth credentials are configured.
func GetClient(ctx context.Context) (*http.Client, error) {
	if !hasCredentials() {
		return nil, ErrMissingCredentials
	}
	cfg := oauthConfig(defaultRedirectURL)
	path, err := tokenPath()
	if err != nil {
		return nil, err
	}

	var tok *oauth2.Token
	tok, err = tokenFromFile(path)
	if err == nil {
		if permErr := ensureTokenPermissions(path); permErr != nil {
			return nil, permErr
		}
	}
	if err != nil || !tok.Valid() {
		if err == nil && !tok.Valid() {
			markInvalidToken(path)
		}
		if err != nil && !os.IsNotExist(err) {
			markInvalidToken(path)
		}
		tok, err = getTokenFromWeb(ctx, cfg)
		if err != nil {
			return nil, err
		}
		if err := saveToken(path, tok); err != nil {
			return nil, err
		}
	}

	return cfg.Client(ctx, tok), nil
}
