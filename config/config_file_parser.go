package config

import (
	"context"
	"fmt"
	"hirsi/enhancement"
	"hirsi/renderer"
	"path"

	"github.com/BurntSushi/toml"
)

type ConfigFile struct {
	directoryPath string

	Storage struct {
		Path string
	}

	Renderer map[string]toml.Primitive
	Filter   map[string]toml.Primitive
	Pipeline map[string]struct {
		Filters   []string
		Renderers []string
	}
}

func (cf *ConfigFile) resolveFile(filepath string) string {
	if path.IsAbs(filepath) {
		return filepath
	}

	return path.Join(cf.directoryPath, filepath)
}

func Parse(ctx context.Context, filepath string) (*Config, error) {

	cfg := &ConfigFile{
		directoryPath: path.Dir(filepath),
	}

	meta, err := toml.DecodeFile(filepath, &cfg)
	if err != nil {
		return nil, err
	}

	renderers, err := parsePipelines(ctx, meta, cfg)
	if err != nil {
		return nil, err
	}

	dbPath := cfg.Storage.Path
	if !path.IsAbs(dbPath) {
		dbPath = path.Join(cfg.directoryPath, dbPath)
	}

	config := &Config{
		DbPath:    dbPath,
		Renderers: renderers,
		Enhancements: []enhancement.Enhancement{
			&enhancement.PwdEnhancement{},
		},
	}

	return config, nil
}

func parsePipelines(ctx context.Context, meta toml.MetaData, cfg *ConfigFile) (map[string]renderer.Renderer, error) {

	rendererDefinitions, err := parseRenderers(ctx, meta, cfg)
	if err != nil {
		return nil, err
	}

	renderers := make(map[string]renderer.Renderer, len(cfg.Pipeline))

	for name, pipeline := range cfg.Pipeline {

		sink := buildRenderer(rendererDefinitions, pipeline.Renderers)

		for _, filterName := range pipeline.Filters {

			filter, found := cfg.Filter[filterName]
			if !found {
				return nil, fmt.Errorf("no filter called %s found", filterName)
			}
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

		renderers[name] = sink
	}

	return renderers, nil
}

func parseRenderers(ctx context.Context, meta toml.MetaData, cfg *ConfigFile) (map[string]renderer.Renderer, error) {
	renderers := map[string]renderer.Renderer{}

	for name, rendererCfg := range cfg.Renderer {

		resource := &ResourceType{}
		if err := meta.PrimitiveDecode(rendererCfg, resource); err != nil {
			return nil, err
		}

		factory, found := rendererFactories[resource.Type]
		if !found {
			return nil, fmt.Errorf("no renderer factory found for %s", resource.Type)
		}

		renderer, err := factory(ctx, cfg, func(target any) error {
			return meta.PrimitiveDecode(rendererCfg, target)
		})

		if err != nil {
			return nil, err
		}

		renderers[name] = renderer
	}

	return renderers, nil
}

func buildRenderer(allSinks map[string]renderer.Renderer, names []string) renderer.Renderer {

	if len(names) == 1 {
		return allSinks[names[0]]
	}

	sinks := make([]renderer.Renderer, len(names))
	for i, name := range names {
		sinks[i] = allSinks[name]
	}

	return renderer.NewCompositeRenderer(sinks)
}

func obsidianFactory(ctx context.Context, cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error) {

	c := &obsidianConfig{}
	if err := decode(c); err != nil {
		return nil, err
	}

	r, err := renderer.NewObsidianRenderer(cfg.resolveFile(c.Path))
	if err != nil {
		return nil, err
	}

	r.AddTitles(c.Titles)
	r.PopulateAutoLinker(ctx)

	return r, nil
}

func logFactory(ctx context.Context, cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error) {

	c := &logConfig{}
	if err := decode(c); err != nil {
		return nil, err
	}

	return renderer.NewLogRenderer(cfg.resolveFile(c.Path))
}

var rendererFactories = map[string]func(ctx context.Context, cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error){
	"obsidian": obsidianFactory,
	"log":      logFactory,
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
