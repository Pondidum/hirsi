package pipeline

import (
	"hirsi/filter"
	"hirsi/message"
	"hirsi/renderer"
)

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		Name: name,
	}
}

type Pipeline struct {
	Name     string
	filters  []filter.Filter
	renderer renderer.Renderer // either a single renderer, or a composite renderer
}

func (p *Pipeline) Handle(message *message.Message) error {
	for _, filter := range p.filters {
		if !filter(message) {
			return nil
		}
	}

	return p.renderer.Render(message)
}

func (p *Pipeline) AddFilters(filters ...filter.Filter) {
	p.filters = append(p.filters, filters...)
}

func (p *Pipeline) SetRenderers(renderers []renderer.Renderer) {
	if len(renderers) == 1 {
		p.renderer = renderers[0]
	} else {
		p.renderer = renderer.NewCompositeRenderer(renderers)
	}
}
