package config

import (
	"context"
	"hirsi/enhancement"
	"hirsi/message"
	"hirsi/renderer"
	"hirsi/tracing"
	"io/fs"
	"path"
)

type Config struct {
	DbPath       string
	Tracing      *tracing.TraceConfiguration
	Enhancements []enhancement.Enhancement
	Renderers    map[string]renderer.Renderer
}


func CreateConfig(ctx context.Context) (*Config, error) {

	// note MachineFS is needed as doing `os.DirFS("/")` will return an FS which doesn't work with
	// absolute paths, which is quite irritating.
	filepath, err := findConfigFile(&realEnvironment{}, &MachineFS{})
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(ctx, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func findConfigFile(env environment, f fs.StatFS) (string, error) {
	filepath := env.GetEnv("HIRSI_CONFIG")
	if filepath != "" {
		if _, err := f.Stat(filepath); err != nil {
			return "", err
		}

		return filepath, nil
	}

	pwd, err := env.GetPwd()
	if err != nil {
		return "", err
	}

	filepath = path.Join(pwd, "hirsi.toml")
	if _, err := f.Stat(filepath); err == nil {
		return filepath, nil
	}

	configDir, err := configPath(env)
	if err != nil {
		return "", err
	}

	filepath = path.Join(configDir, "hirsi/hirsi.toml")
	if _, err := f.Stat(filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

func configPath(env environment) (string, error) {
	if xdgData := env.GetEnv("XDG_CONFIG_HOME"); xdgData != "" {
		return xdgData, nil
	}

	home, err := env.GetHome()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config"), nil

}
