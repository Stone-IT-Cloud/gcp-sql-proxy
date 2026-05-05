package unit

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
	"google.golang.org/api/googleapi"
)

func TestVerifyAccessNotFoundReturnsInstanceNotFoundError(t *testing.T) {
	restore := proxy.SetTestInstancesGet(func(_ context.Context, _ *http.Client, _, _ string) error {
		return &googleapi.Error{
			Code:    http.StatusNotFound,
			Message: "not found",
		}
	})
	defer restore()

	err := proxy.VerifyAccess(context.Background(), &http.Client{}, "project:region:instance")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var notFound proxy.InstanceNotFoundOrNotAccessibleError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected InstanceNotFoundOrNotAccessibleError, got %T: %v", err, err)
	}
}
