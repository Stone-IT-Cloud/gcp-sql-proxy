package specify

import (
	"fmt"
	"io"
)

type ProgressRenderer interface {
	Render(event ProgressEvent) error
}

type PlainTextRenderer struct {
	w io.Writer
}

func NewPlainTextRenderer(w io.Writer) *PlainTextRenderer {
	return &PlainTextRenderer{w: w}
}

func (r *PlainTextRenderer) Render(event ProgressEvent) error {
	if !event.Valid() {
		return fmt.Errorf("invalid progress event")
	}
	if event.Reason != "" {
		_, err := fmt.Fprintf(r.w, "[%s] %s: %s (%s)\n", event.Type, event.PhaseName, event.Message, event.Reason)
		return err
	}
	_, err := fmt.Fprintf(r.w, "[%s] %s: %s\n", event.Type, event.PhaseName, event.Message)
	return err
}

type RichRenderer struct {
	w io.Writer
}

func NewRichRenderer(w io.Writer) *RichRenderer {
	return &RichRenderer{w: w}
}

func (r *RichRenderer) Render(event ProgressEvent) error {
	if !event.Valid() {
		return fmt.Errorf("invalid progress event")
	}
	_, err := fmt.Fprintf(r.w, "• %s: %s\n", event.PhaseName, event.Message)
	return err
}
