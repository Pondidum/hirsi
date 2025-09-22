package enhancement

import (
	"hirsi/message"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
)

func TestTagAddConfig(t *testing.T) {
	cfg := `
[[enhancement]]
type = "tag-add"
check = "pwd"
condition = "equals"
[enhancement.values]
"/home/andy/dev/projects/homelab" = ["homelab"]
"/home/andy/dev/projects/trace" = ["trace", "otel"]
	`

	config := &struct {
		Enhancement []TagAddConfig
	}{}

	_, err := toml.Decode(cfg, config)
	require.NoError(t, err)

	tagAddCfg := config.Enhancement[0]
	require.Equal(t, "pwd", tagAddCfg.Check)
	require.Equal(t, "equals", tagAddCfg.Condition)
	require.Equal(t, map[string][]string{
		"/home/andy/dev/projects/homelab": []string{"homelab"},
		"/home/andy/dev/projects/trace":   []string{"trace", "otel"},
	}, tagAddCfg.Values)
}

func TestTagAdd(t *testing.T) {

	cases := []struct {
		condition string
		tags      map[string]string
		expected  map[string]string
		err       error
	}{
		{
			condition: "equals",
			tags:      map[string]string{},
			expected:  map[string]string{},
		},
		{
			condition: "equals",
			tags:      map[string]string{"test": "first"},
			expected:  map[string]string{"test": "first", "set": "1"},
		},
		{
			condition: "equals",
			tags:      map[string]string{"test": "second"},
			expected:  map[string]string{"test": "second", "set": "2"},
		},
		{
			condition: "equals",
			tags:      map[string]string{"test": "first", "set": "2"},
			expected:  map[string]string{"test": "first", "set": "1"},
		},
		{
			condition: "prefix",
			tags:      map[string]string{},
			expected:  map[string]string{},
		},
		{
			condition: "prefix",
			tags:      map[string]string{"test": "first/prefix"},
			expected:  map[string]string{"test": "first/prefix", "set": "1"},
		},
		{
			condition: "prefix",
			tags:      map[string]string{"test": "second/prefix"},
			expected:  map[string]string{"test": "second/prefix", "set": "2"},
		},
		{
			condition: "prefix",
			tags:      map[string]string{"test": "first/prefix", "set": "2"},
			expected:  map[string]string{"test": "first/prefix", "set": "1"},
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {

			enh := &TagAddEnhancement{
				field:     "test",
				condition: tc.condition,
				values: map[string][]string{
					"first":  {"set", "1"},
					"second": {"set", "2"},
				},
			}

			m := &message.Message{
				Tags: tc.tags,
			}

			err := enh.Enhance(m)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, m.Tags)
		})
	}
}
