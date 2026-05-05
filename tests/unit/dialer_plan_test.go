package unit

import (
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestDefaultDialerPlanDefaultsToPublicUnlessPrivateRequested(t *testing.T) {
	plan := proxy.DefaultDialerPlan(false)
	if !plan.UseIAMAuthN {
		t.Fatal("UseIAMAuthN = false, want true")
	}
	if plan.UsePrivateIP {
		t.Fatal("UsePrivateIP = true, want false")
	}

	privatePlan := proxy.DefaultDialerPlan(true)
	if !privatePlan.UsePrivateIP {
		t.Fatal("UsePrivateIP = false, want true")
	}
}
