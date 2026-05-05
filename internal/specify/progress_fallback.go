package specify

import "fmt"

type FallbackRenderer struct {
	primary  ProgressRenderer
	fallback ProgressRenderer
}

func NewFallbackRenderer(primary ProgressRenderer, fallback ProgressRenderer) *FallbackRenderer {
	return &FallbackRenderer{
		primary:  primary,
		fallback: fallback,
	}
}

func (r *FallbackRenderer) Render(event ProgressEvent) error {
	if err := r.primary.Render(event); err != nil {
		fallbackEvent := event
		fallbackEvent.Message = fmt.Sprintf("%s (fallback plain text)", event.Message)
		return r.fallback.Render(fallbackEvent)
	}
	return nil
}
