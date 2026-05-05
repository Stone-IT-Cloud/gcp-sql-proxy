package auth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// CloudSQLConnectorTokenSources returns OAuth token sources suitable for
// cloudsqlconn.WithIAMAuthNTokenSources(adminTokenSource, iamLoginTokenSource).
//
// The returned token sources derive from the persisted OAuth refresh token
// and must support:
// - SQL Admin API calls for connector admin leg
// - SQL IAM login token minting leg
func CloudSQLConnectorTokenSources(ctx context.Context) (adminTokenSource oauth2.TokenSource, iamLoginTokenSource oauth2.TokenSource, err error) {
	if !hasCredentials() {
		return nil, nil, ErrMissingCredentials
	}

	path, err := tokenPath()
	if err != nil {
		return nil, nil, fmt.Errorf("cloudsql connector token sources: resolve token path: %w", err)
	}

	tok, err := tokenFromFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cloudsql connector token sources: read token file: %w", err)
	}
	if tok == nil || tok.RefreshToken == "" {
		return nil, nil, fmt.Errorf("cloudsql connector token sources: missing refresh token: %w", ErrMissingCredentials)
	}

	// OAuth refresh tokens are consent-scoped. Ensure admin token sources request
	// the exact scopes required by WithIAMAuthNTokenSources contract.
	adminCfg := &oauth2.Config{
		ClientID:     OAuthClientID,
		ClientSecret: OAuthClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{sqlAdminScope, cloudPlatformScope},
	}
	iamLoginCfg := &oauth2.Config{
		ClientID:     OAuthClientID,
		ClientSecret: OAuthClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{sqlLoginScope},
	}

	adminTokenSource = adminCfg.TokenSource(ctx, tok)
	iamLoginTokenSource = iamLoginCfg.TokenSource(ctx, tok)
	return adminTokenSource, iamLoginTokenSource, nil
}
