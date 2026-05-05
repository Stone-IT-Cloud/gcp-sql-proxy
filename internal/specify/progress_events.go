package specify

import "time"

type PhaseType string

const (
	PhaseHook       PhaseType = "hook"
	PhaseGeneration PhaseType = "generation"
	PhaseValidation PhaseType = "validation"
	PhaseCompletion PhaseType = "completion"
	PhaseFailure    PhaseType = "failure"
)

type EventType string

const (
	EventStarted   EventType = "started"
	EventHeartbeat EventType = "heartbeat"
	EventCompleted EventType = "completed"
	EventSkipped   EventType = "skipped"
	EventFailed    EventType = "failed"
)

type ExecutionPhase struct {
	ID        string
	Name      string
	Type      PhaseType
	StartedAt time.Time
	EndedAt   time.Time
}

type ProgressEvent struct {
	Type      EventType
	PhaseID   string
	PhaseName string
	Message   string
	Timestamp time.Time
	Reason    string
}

func NewEvent(eventType EventType, phase ExecutionPhase, message string) ProgressEvent {
	return ProgressEvent{
		Type:      eventType,
		PhaseID:   phase.ID,
		PhaseName: phase.Name,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func (e ProgressEvent) Valid() bool {
	if e.PhaseID == "" || e.Message == "" {
		return false
	}
	if (e.Type == EventSkipped || e.Type == EventFailed) && e.Reason == "" {
		return false
	}
	return true
}
