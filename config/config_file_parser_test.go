package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	cfg, err := Parse(context.Background(), "test.toml")
	assert.NoError(t, err)

	assert.Equal(t, cfg.DbPath, "relative/file.db")

	assert.Len(t, cfg.Renderers, 2)
}
