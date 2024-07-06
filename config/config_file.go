package config

import (
	"fmt"
	"hirsi/enhancement"
	"hirsi/renderer"
	"io"

	"github.com/BurntSushi/toml"
)

type ConfigFile struct {
	Renderer map[string]toml.Primitive
	Filter   map[string]toml.Primitive
	Pipeline []struct {
		Filters   []string
		Renderers []string
	}
}

func Parse(reader io.Reader) (*Config, error) {
	cfg := &ConfigFile{}
	meta, err := toml.NewDecoder(reader).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	sinks := map[string]renderer.Renderer{}

	for name, rendererCfg := range cfg.Renderer {

		resource := &ResourceType{}
		if err := meta.PrimitiveDecode(rendererCfg, resource); err != nil {
			return nil, err
		}

		factory, found := rendererFactories[resource.Type]
		if !found {
			return nil, fmt.Errorf("no renderer factory found for %s", resource.Type)
		}

		renderer, err := factory(meta, rendererCfg)
		if err != nil {
			return nil, err
		}

		sinks[name] = renderer
	}

	config := &Config{
		Enhancements: []enhancement.Enhancement{
			&enhancement.PwdEnhancement{},
		},
	}
	for _, pipeline := range cfg.Pipeline {

		sink := buildSink(sinks, pipeline.Renderers)

		for _, filterName := range pipeline.Filters {

			filter := cfg.Filter[filterName]
			resource := &ResourceType{}
			if err := meta.PrimitiveDecode(filter, resource); err != nil {
				return nil, err
			}

			factory, found := filterFactories[resource.Type]
			if !found {
				return nil, fmt.Errorf("no filter type called %s found", resource.Type)
			}

			if sink, err = factory(meta, filter, sink); err != nil {
				return nil, err
			}
		}

		config.Renderers = append(config.Renderers, sink)
	}

	return config, nil

}

func buildSink(allSinks map[string]renderer.Renderer, names []string) renderer.Renderer {

	if len(names) == 1 {
		return allSinks[names[0]]
	}

	sinks := make([]renderer.Renderer, len(names))
	for i, name := range names {
		sinks[i] = allSinks[name]
	}

	return renderer.NewCompositeRenderer(sinks)
}

var rendererFactories = map[string]func(meta toml.MetaData, p toml.Primitive) (renderer.Renderer, error){
	"obsidian": func(meta toml.MetaData, p toml.Primitive) (renderer.Renderer, error) {

		cfg := &obsidianConfig{}
		if err := meta.PrimitiveDecode(p, cfg); err != nil {
			return nil, err
		}

		return renderer.NewObsidianRenderer(cfg.Path, cfg.Titles...)
	},
	"log": func(meta toml.MetaData, p toml.Primitive) (renderer.Renderer, error) {

		cfg := &logConfig{}
		if err := meta.PrimitiveDecode(p, cfg); err != nil {
			return nil, err
		}

		return renderer.NewLogRenderer(cfg.Path)
	},
}

var filterFactories = map[string]func(meta toml.MetaData, p toml.Primitive, sink renderer.Renderer) (renderer.Renderer, error){
	"tag-with-prefix": func(meta toml.MetaData, p toml.Primitive, sink renderer.Renderer) (renderer.Renderer, error) {

		cfg := &tagFilterConfig{}
		if err := meta.PrimitiveDecode(p, cfg); err != nil {
			return nil, err
		}

		filter := renderer.HasTagWithPrefix(cfg.Tag, cfg.Prefix)

		return renderer.NewFilteredRenderer(filter, sink), nil
	},
	"tag-without-prefix": func(meta toml.MetaData, p toml.Primitive, sink renderer.Renderer) (renderer.Renderer, error) {

		cfg := &tagFilterConfig{}
		if err := meta.PrimitiveDecode(p, cfg); err != nil {
			return nil, err
		}

		filter := renderer.HasTagWithoutPrefix(cfg.Tag, cfg.Prefix)

		return renderer.NewFilteredRenderer(filter, sink), nil
	},
}

type ResourceType struct {
	Type string
}

type obsidianConfig struct {
	ResourceType

	Path   string
	Titles []string
}

type logConfig struct {
	ResourceType

	Path string
}

type tagFilterConfig struct {
	Tag    string
	Prefix string
}

// type RendererConfig struct {
// 	Name string
// 	Type string

// }
