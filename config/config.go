package config

import (
	"hirsi/enhancement"
	"hirsi/renderer"
)

type Config struct {
	DbPath       string
	Enhancements []enhancement.Enhancement
	Renderers    []renderer.Renderer
}

func CreateConfig() (*Config, error) {
	epicObsidian, err := renderer.NewObsidianRenderer("/home/andy/dev/epic/obsidian")
	if err != nil {
		return nil, err
	}

	epicPlain, err := renderer.NewLogRenderer("/home/andy/dev/epic/log-plain")
	if err != nil {
		return nil, err
	}

	plain, err := renderer.NewLogRenderer("/home/andy/log")
	if err != nil {
		return nil, err
	}

	return &Config{
		DbPath: "/home/andy/.local/hirsi/hirsi.db",
		Enhancements: []enhancement.Enhancement{
			&enhancement.PwdEnhancement{},
		},
		Renderers: []renderer.Renderer{
			renderer.NewFilteredRenderer(
				renderer.HasTagWithPrefix("pwd", "/home/andy/dev/epic"),
				renderer.NewCompositeRenderer([]renderer.Renderer{epicObsidian, epicPlain}),
			),
			renderer.NewFilteredRenderer(
				renderer.HasTagWithoutPrefix("pwd", "/home/andy/dev/epic"),
				plain,
			),
		},
	}, nil

}
