package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigParsing(t *testing.T) {
	cfg, err := Parse("test.toml")
	require.NoError(t, err)

	require.Equal(t, cfg.DbPath, "relative/file.db")

	require.Len(t, cfg.Pipelines, 2)

	first := cfg.Pipelines[0]
	require.Equal(t, "only-secret", first.Name)
	// require.Len(t, first.enhancements, 1)
	// require.Len(t, first.filters, 1)
	// require.NotNil(t, first.renderers)
	// require.IsType(t, &log.LogRenderer{}, first.renderers)

	second := cfg.Pipelines[1]
	require.Equal(t, "not-secret", second.Name)
	// require.Len(t, second.enhancements, 1)
	// require.Len(t, second.filters, 1)
	// require.NotNil(t, second.renderers)
	// require.IsType(t, &renderer.CompositeRenderer{}, second.renderers)
}
