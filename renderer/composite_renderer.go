package renderer

import (
	"hirsi/message"
	"sync"

	"golang.org/x/sync/errgroup"
)

type CompositeRenderer struct {
	renderers []Renderer
}

func NewCompositeRendererM(m map[string]Renderer) *CompositeRenderer {
	renderers := make([]Renderer, 0, len(m))
	for _, r := range m {
		renderers = append(renderers, r)
	}

	return &CompositeRenderer{
		renderers: renderers,
	}
}

func NewCompositeRenderer(renderers []Renderer) *CompositeRenderer {
	return &CompositeRenderer{
		renderers: renderers,
	}
}

func (r *CompositeRenderer) Render(message *message.Message) error {
	wg := sync.WaitGroup{}
	wg.Add(len(r.renderers))

	g := errgroup.Group{}

	for _, r := range r.renderers {
		renderer := r
		g.Go(func() error {
			return renderer.Render(message)
		})
	}

	return g.Wait()
}
