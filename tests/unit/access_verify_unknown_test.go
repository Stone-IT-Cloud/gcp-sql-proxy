package unit

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestVerifyAccessTransportErrorsReturnPermissionCheckUnavailable(t *testing.T) {
	restore := proxy.SetTestInstancesGet(func(_ context.Context, _ *http.Client, _, _ string) error {
		return errors.New("dial timeout")
	})
	defer restore()

	err := proxy.VerifyAccess(context.Background(), &http.Client{}, "project:region:instance")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var u *proxy.PermissionCheckUnavailableError
	if !errors.As(err, &u) {
		t.Fatalf("expected PermissionCheckUnavailableError, got %T: %v", err, err)
	}
	if u == nil {
		t.Fatal("expected typed error instance")
	}
}
