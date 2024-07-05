package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	content, err := os.ReadFile("test.toml")
	assert.NoError(t, err)

	cfg, err := Parse(string(content))
	assert.NoError(t, err)

	assert.Len(t, cfg.Renderers, 2)
}
