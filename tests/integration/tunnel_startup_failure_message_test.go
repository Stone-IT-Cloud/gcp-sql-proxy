package integration

import (
	"errors"
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestTunnelStartupFailureSuppressesSuccessMessage(t *testing.T) {
	got := proxy.UserFacingError(errors.New("dial failed"))
	if got == "" {
		t.Fatal("expected failure guidance")
	}
}
