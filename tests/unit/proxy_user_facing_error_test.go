package unit

import (
	"errors"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestUserFacingErrorMissingRolesIsActionableAndNonSecret(t *testing.T) {
	err := proxy.MissingCloudSQLOrIAPRolesError{
		MissingRoles: []string{"roles/cloudsql.client", "roles/iap.tunnelResourceAccessor"},
	}

	msg := proxy.UserFacingError(err)
	if !strings.Contains(msg, "contact DevOps") {
		t.Fatalf("message should instruct to contact DevOps, got %q", msg)
	}
	for _, role := range err.MissingRoles {
		if !strings.Contains(msg, role) {
			t.Fatalf("message missing role %q; got %q", role, msg)
		}
	}
	// Ensure we didn't accidentally echo token/credential material.
	for _, secretHint := range []string{"token", "JWT", "refresh"} {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(secretHint)) {
			t.Fatalf("message appears to contain secret material hint %q; got %q", secretHint, msg)
		}
	}
}

func TestUserFacingErrorPermissionCheckUnavailableIsUserFriendly(t *testing.T) {
	err := &proxy.PermissionCheckUnavailableError{Cause: errors.New("network")}
	msg := proxy.UserFacingError(err)
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	if !strings.Contains(strings.ToLower(msg), "permission check unavailable") {
		t.Fatalf("unexpected message: %q", msg)
	}
}
