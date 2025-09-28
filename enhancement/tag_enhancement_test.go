package enhancement

import (
	"hirsi/message"
	"slices"
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
		tags      []message.Tag
		expected  []message.Tag
		err       error
	}{
		{
			condition: "equals",
		},
		{
			condition: "equals",
			tags:      []message.Tag{{"test", "first"}},
			expected:  []message.Tag{{"test", "first"}, {"set", "1"}},
		},
		{
			condition: "equals",
			tags:      []message.Tag{{"test", "second"}},
			expected:  []message.Tag{{"test", "second"}, {"set", "2"}},
		},
		{
			condition: "equals",
			tags:      []message.Tag{{"test", "first"}, {"set", "2"}},
			expected:  []message.Tag{{"test", "first"}, {"set", "2"}, {"set", "1"}},
		},
		{
			condition: "prefix",
			tags:      []message.Tag{},
			expected:  []message.Tag{},
		},
		{
			condition: "prefix",
			tags:      []message.Tag{{"test", "first/prefix"}},
			expected:  []message.Tag{{"test", "first/prefix"}, {"set", "1"}},
		},
		{
			condition: "prefix",
			tags:      []message.Tag{{"test", "second/prefix"}},
			expected:  []message.Tag{{"test", "second/prefix"}, {"set", "2"}},
		},
		{
			condition: "prefix",
			tags:      []message.Tag{{"test", "first/prefix"}, {"set", "2"}},
			expected:  []message.Tag{{"test", "first/prefix"}, {"set", "2"}, {"set", "1"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.condition, func(t *testing.T) {

			enh := &TagAddEnhancement{
				field:     "test",
				condition: tc.condition,
				replacements: map[string][]string{
					"first":  {"set", "1"},
					"second": {"set", "2"},
				},
			}

			m := &message.Message{
				Tags: message.NewTagsFrom(tc.tags),
			}

			err := enh.Enhance(m)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, slices.Collect(slices.Values(tc.expected)), slices.Collect(m.Tags.All()))
		})
	}
}
