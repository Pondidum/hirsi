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
	obs, err := renderer.NewObsidianRenderer("/home/andy/dev/projects/hirsi/obsidian-dev/hirsi-dev")
	if err != nil {
		return nil, err
	}

	return &Config{
		DbPath: "/home/andy/.local/hirsi/hirsi.db",
		Enhancements: []enhancement.Enhancement{
			&enhancement.PwdEnhancement{},
		},
		Renderers: []renderer.Renderer{
			// &renderer.JsonRenderer{Writer: os.Stderr},
			obs,
		},
	}, nil

}
