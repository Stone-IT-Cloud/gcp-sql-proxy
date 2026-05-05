package specify

import (
	"context"
	"fmt"
)

type Command struct {
	renderer ProgressRenderer
	artifacts []string
}

func NewCommand(renderer ProgressRenderer, artifacts []string) *Command {
	return &Command{renderer: renderer, artifacts: artifacts}
}

func (c *Command) EmitPhase(phase ExecutionPhase) error {
	if err := c.renderer.Render(NewEvent(EventStarted, phase, "phase started")); err != nil {
		return err
	}
	return c.renderer.Render(NewEvent(EventCompleted, phase, "phase completed"))
}

func (c *Command) EmitCompletion() error {
	phase := ExecutionPhase{ID: "completion", Name: "completion", Type: PhaseCompletion}
	msg := "completed successfully"
	if len(c.artifacts) > 0 {
		msg = fmt.Sprintf("completed successfully (artifacts: %d)", len(c.artifacts))
	}
	return c.renderer.Render(NewEvent(EventCompleted, phase, msg))
}

func (c *Command) Run(ctx context.Context, phases []ExecutionPhase) error {
	for _, p := range phases {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := c.EmitPhase(p); err != nil {
			return err
		}
	}
	return c.EmitCompletion()
}
