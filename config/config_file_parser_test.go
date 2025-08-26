package config

import (
	"hirsi/enhancement"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
)

func TestConfigParsing(t *testing.T) {
	cfg, err := Parse("test.toml")
	require.NoError(t, err)

	require.Equal(t, cfg.DbPath, "relative/file.db")

	require.Len(t, cfg.Pipelines, 2)

	first := cfg.Pipelines[0]
	require.Equal(t, "only-secret", first.Name)

	second := cfg.Pipelines[1]
	require.Equal(t, "not-secret", second.Name)
}

func TestParsingEnhancements(t *testing.T) {
	cfg := `
[[enhancement]]
type = "tag-add"
check = "pwd"
condition = "equals"
[[enhancement.values]]
"/home/andy/dev/projects/homelab" = ["homelab"]
"/home/andy/dev/projects/trace" = ["trace", "otel"]
	`

	cf := &ConfigFile{}
	meta, err := toml.Decode(cfg, cf)
	require.NoError(t, err)

	all, err := parseEnhancements(cf, meta)
	require.NoError(t, err)
	require.IsType(t, &enhancement.TagAddEnhancement{}, all[0])
}
