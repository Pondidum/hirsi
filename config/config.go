package config

import (
	"bytes"
	"context"
	"fmt"
	"hirsi/enhancement"
	"hirsi/renderer"
	"io/fs"
	"os"
	"path"
)

type Config struct {
	DbPath       string
	Enhancements []enhancement.Enhancement
	Renderers    map[string]renderer.Renderer
}

type MachineFS struct{}

func (fs MachineFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs MachineFS) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func CreateConfig(ctx context.Context) (*Config, error) {

	// note MachineFS is needed as doing `os.DirFS("/")` will return an FS which doesn't work with
	// absolute paths, which is quite irritating.
	content, err := findConfigFile(&realEnvironment{}, &MachineFS{})
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(ctx, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	if cfg.DbPath, err = dbPath(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type environment interface {
	GetEnv(key string) string
	GetHome() (string, error)
	GetPwd() (string, error)
}

type realEnvironment struct{}

func (r *realEnvironment) GetEnv(key string) string { return os.Getenv(key) }
func (r *realEnvironment) GetHome() (string, error) { return os.UserHomeDir() }
func (r *realEnvironment) GetPwd() (string, error)  { return os.Getwd() }

func findConfigFile(env environment, f fs.ReadFileFS) ([]byte, error) {
	filepath := env.GetEnv("HIRSI_CONFIG")
	if filepath != "" {
		content, err := f.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		return content, nil
	}

	pwd, err := env.GetPwd()
	if err != nil {
		return nil, err
	}

	if content, err := f.ReadFile(path.Join(pwd, "hirsi.toml")); err == nil {
		return content, nil
	}

	configDir, err := configPath(env)
	if err != nil {
		return nil, err
	}

	filepath = path.Join(configDir, "hirsi/hirsi.toml")
	content, err := f.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return content, nil
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
