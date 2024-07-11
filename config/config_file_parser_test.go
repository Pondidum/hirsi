package config

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	content, err := os.ReadFile("test.toml")
	assert.NoError(t, err)

	cfg, err := Parse(context.Background(), bytes.NewReader(content))
	assert.NoError(t, err)

	assert.Len(t, cfg.Renderers, 2)
}
