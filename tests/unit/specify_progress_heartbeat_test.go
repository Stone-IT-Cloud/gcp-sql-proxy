package unit

import (
	"context"
	"testing"
	"time"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/specify"
)

func TestSpecifyHeartbeatInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	phase := specify.ExecutionPhase{ID: "generation", Name: "generation", Type: specify.PhaseGeneration}
	ch := make(chan specify.ProgressEvent, 2)
	go specify.StartHeartbeat(ctx, phase, 10*time.Millisecond, func(event specify.ProgressEvent) {
		ch <- event
	})

	select {
	case ev := <-ch:
		if ev.Type != specify.EventHeartbeat {
			t.Fatalf("type = %s, want heartbeat", ev.Type)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting heartbeat")
	}
}
