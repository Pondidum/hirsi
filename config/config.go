package config

import (
	"bytes"
	"hirsi/enhancement"
	"hirsi/renderer"
	"os"
	"path"
)

type Config struct {
	DbPath       string
	Enhancements []enhancement.Enhancement
	Renderers    []renderer.Renderer
}

func CreateConfig() (*Config, error) {
	return createConfigFromEnvironment()
}

func createConfigFromEnvironment() (*Config, error) {
	filepath := os.Getenv("HIRSI_CONFIG")
	if filepath == "" {

		if xdgData := os.Getenv("XDG_CONFIG_HOME"); xdgData == "" {
			filepath = path.Join(xdgData, "hirsi/hirsi.toml")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}

			filepath = path.Join(home, ".config/hirsi/hirsi.toml")
		}
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	if cfg.DbPath, err = dbPath(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func dbPath() (string, error) {
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData == "" {
		return path.Join(xdgData, "hirsi/hirsi.db"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".local/share/hirsi/hirsi.db"), nil
}
