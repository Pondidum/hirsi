package config

import (
	"hirsi/enhancement"
	"hirsi/message"
	"os"
	"testing"
	"time"

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

func TestLiveConfig(t *testing.T) {
	content, err := os.Open("live.toml")
	if err == os.ErrNotExist {
		t.Skip("no live.toml found")
	}
	defer content.Close()
	require.NotEmpty(t, content)

	cfg, err := parse(".", content)
	require.NoError(t, err)

	require.NotEmpty(t, cfg.Enhancements)

	message := &message.Message{
		WrittenAt: time.Now(),
		Message:   "whaaat",
	}

	for _, e := range cfg.Enhancements {
		require.NoError(t, e.Enhance(message))
	}
}
