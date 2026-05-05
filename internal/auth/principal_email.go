package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrMissingPrincipalEmail = errors.New("missing principal email")

type userInfoEmailResponse struct {
	Email string `json:"email"`
}

// PrincipalEmail resolves a human-readable mailbox for the authenticated principal
// using the OAuth userinfo email endpoint.
func PrincipalEmail(ctx context.Context, httpClient *http.Client) (string, error) {
	if httpClient == nil {
		return "", fmt.Errorf("principal email: missing http client: %w", ErrMissingPrincipalEmail)
	}

	// Keep a bounded latency so startup doesn't hang on identity enrichment.
	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return "", fmt.Errorf("principal email: build request: %w: %w", err, ErrMissingPrincipalEmail)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("principal email: call userinfo endpoint: %w: %w", err, ErrMissingPrincipalEmail)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("principal email: unexpected userinfo status: %s: %w", resp.Status, ErrMissingPrincipalEmail)
	}

	var parsed userInfoEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("principal email: decode response: %w: %w", err, ErrMissingPrincipalEmail)
	}
	if parsed.Email == "" {
		return "", ErrMissingPrincipalEmail
	}

	return parsed.Email, nil
}
