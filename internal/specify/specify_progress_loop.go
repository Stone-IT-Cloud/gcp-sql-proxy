package specify

import (
	"context"
	"time"
)

func StartHeartbeat(ctx context.Context, phase ExecutionPhase, interval time.Duration, emit func(ProgressEvent)) {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			emit(NewEvent(EventHeartbeat, phase, "still running"))
		}
	}
}
