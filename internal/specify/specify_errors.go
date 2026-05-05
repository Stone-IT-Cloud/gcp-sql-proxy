package specify

import "fmt"

func FailureNextStep(phaseID string) string {
	switch phaseID {
	case "hooks":
		return "Review hook output and rerun /speckit-specify."
	case "generation":
		return "Check spec inputs and rerun /speckit-specify."
	case "validation":
		return "Fix validation issues and rerun /speckit-specify."
	default:
		return "Review logs and retry."
	}
}

func FailureEvent(phase ExecutionPhase, err error) ProgressEvent {
	ev := NewEvent(EventFailed, phase, fmt.Sprintf("phase failed: %v", err))
	ev.Reason = FailureNextStep(phase.ID)
	return ev
}
