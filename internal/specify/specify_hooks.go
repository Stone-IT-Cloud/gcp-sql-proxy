package specify

func HookTransitionEvents(phase ExecutionPhase) (ProgressEvent, ProgressEvent) {
	start := NewEvent(EventStarted, phase, "hook started")
	done := NewEvent(EventCompleted, phase, "hook completed")
	return start, done
}

func SkippedEvent(phase ExecutionPhase, reason string) ProgressEvent {
	ev := NewEvent(EventSkipped, phase, "step skipped")
	ev.Reason = reason
	return ev
}
