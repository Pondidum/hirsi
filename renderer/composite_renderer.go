package renderer

import (
	"hirsi/message"
	"sync"

	"golang.org/x/sync/errgroup"
)

type CompositeRenderer struct {
	renderers []Renderer
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
