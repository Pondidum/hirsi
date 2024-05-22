package config

import (
	"hirsi/enhancement"
	"hirsi/renderer"
	"os"
)

type Config struct {
	DbPath       string
	Enhancements []enhancement.Enhancement
	Renderers    []renderer.Renderer
}

var AppConfig = &Config{
	DbPath: "/home/andy/.local/hirsi/hirsi.db",
	Enhancements: []enhancement.Enhancement{
		&enhancement.PwdEnhancement{},
	},
	Renderers: []renderer.Renderer{
		&renderer.JsonRenderer{Writer: os.Stderr},
	},
}
