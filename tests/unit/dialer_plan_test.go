package unit

import (
	"testing"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/proxy"
)

func TestDefaultDialerPlanEnforcesIAMAuthAndPrivateIP(t *testing.T) {
	plan := proxy.DefaultDialerPlan()
	if !plan.UseIAMAuthN {
		t.Fatal("UseIAMAuthN = false, want true")
	}
	if !plan.UsePrivateIP {
		t.Fatal("UsePrivateIP = false, want true")
	}
}
