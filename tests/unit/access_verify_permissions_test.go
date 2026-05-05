package unit

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
	"google.golang.org/api/googleapi"
)

func TestVerifyAccessPermissionDeniedReturnsTypedError(t *testing.T) {
	restore := proxy.SetTestInstancesGet(func(_ context.Context, _ *http.Client, _ string, _ string) error {
		return &googleapi.Error{
			Code:    http.StatusForbidden,
			Message: "forbidden",
		}
	})
	defer restore()

	err := proxy.VerifyAccess(context.Background(), &http.Client{}, "project:region:instance")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var denied proxy.MissingCloudSQLOrIAPRolesError
	if !errors.As(err, &denied) {
		t.Fatalf("expected MissingCloudSQLOrIAPRolesError, got %T: %v", err, err)
	}

	msg := proxy.UserFacingError(err)
	if msg == "" {
		t.Fatal("expected non-empty user-facing error")
	}
	for _, mustContain := range []string{"contact DevOps", "roles/cloudsql.client", "roles/iap.tunnelResourceAccessor"} {
		if !strings.Contains(msg, mustContain) {
			t.Fatalf("user-facing message missing %q; got %q", mustContain, msg)
		}
	}
}
