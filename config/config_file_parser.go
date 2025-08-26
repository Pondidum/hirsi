package config

import (
	"context"
	"fmt"
	"hirsi/enhancement"
	"hirsi/filter"
	"hirsi/pipeline"
	"hirsi/renderer"
	"hirsi/renderer/log"
	"hirsi/renderer/obsidian"
	"hirsi/tracing"
	"path"

	"github.com/BurntSushi/toml"
)

type ConfigFile struct {
	directoryPath string

	Tracing *tracing.TraceConfiguration
	Storage struct {
		Path string
	}

	Enhancements []toml.Primitive `toml:"enhancement"`
	Pipelines    []PipelineConfig `toml:"pipe"`
}

type PipelineConfig struct {
	Name      string
	Filters   []toml.Primitive `toml:"filter"`
	Renderers []toml.Primitive `toml:"renderer"`
}

type ResourceType struct {
	Type string
}

func Parse(filepath string) (*Config, error) {

	cf := &ConfigFile{
		directoryPath: path.Dir(filepath),
	}

	meta, err := toml.DecodeFile(filepath, &cf)
	if err != nil {
		return nil, err
	}

	enhancements, err := parseEnhancements(cf, meta)
	if err != nil {
		return nil, err
	}

	pipelines := make([]*pipeline.Pipeline, len(cf.Pipelines))
	for i, pipeline := range cf.Pipelines {
		if pipelines[i], err = parsePipeline(cf, meta, pipeline); err != nil {
			return nil, err
		}
	}

	dbPath := cf.Storage.Path
	if !path.IsAbs(dbPath) {
		dbPath = path.Join(cf.directoryPath, dbPath)
	}

	config := &Config{
		DbPath:       dbPath,
		Tracing:      cf.Tracing,
		Enhancements: enhancements,
		Pipelines:    pipelines,
	}

	return config, nil
}

var filterFactories = map[string]func(decode func(target any) error) (filter.Filter, error){
	"tag-with-prefix": func(decode func(target any) error) (filter.Filter, error) {
		conf := struct {
			Tag    string
			Prefix string
		}{}
		if err := decode(&conf); err != nil {
			return nil, err
		}
		return filter.HasTagWithPrefix(conf.Tag, conf.Prefix), nil
	},
	"tag-without-prefix": func(decode func(target any) error) (filter.Filter, error) {
		conf := struct {
			Tag    string
			Prefix string
		}{}
		if err := decode(&conf); err != nil {
			return nil, err
		}
		return filter.HasTagWithoutPrefix(conf.Tag, conf.Prefix), nil
	},
}

func parseEnhancements(cf *ConfigFile, meta toml.MetaData) ([]enhancement.Enhancement, error) {
	enhancements := make([]enhancement.Enhancement, len(cf.Enhancements))

	for i, ec := range cf.Enhancements {
		r := &ResourceType{}
		if err := meta.PrimitiveDecode(ec, r); err != nil {
			return nil, err
		}

		builder, found := enhancement.GetBuilder(r.Type)
		if !found {
			return nil, fmt.Errorf("no enhancement type called '%s' found", r.Type)
		}

		if err := meta.PrimitiveDecode(ec, builder); err != nil {
			return nil, err
		}

		enhancements[i] = builder.Build()
	}

	return enhancements, nil
}

func parsePipeline(cf *ConfigFile, meta toml.MetaData, pc PipelineConfig) (*pipeline.Pipeline, error) {
	p := pipeline.NewPipeline(pc.Name)

	for _, fc := range pc.Filters {
		resource := &ResourceType{}
		if err := meta.PrimitiveDecode(fc, resource); err != nil {
			return nil, err
		}
		factory, found := filterFactories[resource.Type]
		if !found {
			return nil, fmt.Errorf("no filter of type '%s' found", resource.Type)
		}

		filter, err := factory(func(target any) error { return meta.PrimitiveDecode(fc, target) })
		if err != nil {
			return nil, err
		}

		p.AddFilters(filter)
	}

	renderers := make([]renderer.Renderer, len(pc.Renderers))
	for i, rc := range pc.Renderers {
		resource := &ResourceType{}
		if err := meta.PrimitiveDecode(rc, resource); err != nil {
			return nil, err
		}
		factory, found := rendererFactories[resource.Type]
		if !found {
			return nil, fmt.Errorf("no renderer of type '%s' found", resource.Type)
		}

		renderer, err := factory(cf, func(target any) error { return meta.PrimitiveDecode(rc, target) })
		if err != nil {
			return nil, err
		}
		renderers[i] = renderer
	}
	p.SetRenderers(renderers)

	return p, nil
}

var rendererFactories = map[string]func(cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error){
	"obsidian": obsidianFactory,
	"log":      logFactory,
}

func obsidianFactory(cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error) {

	c := &struct {
		ResourceType

		Path   string
		Titles []string
	}{}
	if err := decode(c); err != nil {
		return nil, err
	}

	r, err := obsidian.NewObsidianRenderer(resolveFile(cfg, c.Path))
	if err != nil {
		return nil, err
	}

	r.AddTitles(c.Titles)
	r.PopulateAutoLinker(context.TODO())

	return r, nil
}

func logFactory(cfg *ConfigFile, decode func(target any) error) (renderer.Renderer, error) {

	c := &struct {
		ResourceType
		Path string
	}{}
	if err := decode(c); err != nil {
		return nil, err
	}

	return log.NewLogRenderer(resolveFile(cfg, c.Path))
}

func resolveFile(cf *ConfigFile, filepath string) string {
	if path.IsAbs(filepath) {
		return filepath
	}

	return path.Join(cf.directoryPath, filepath)
}
