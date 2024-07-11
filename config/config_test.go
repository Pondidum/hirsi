package config

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConfigFile(t *testing.T) {

	testCases := []struct {
		name              string
		hirsiConfigEnvVar string
		xdgEnvVar         string
		localConfig       bool
		expectedErr       error
		expectedContent   string
	}{
		{
			name:              "HIRSI_CONFIG specified, and exists",
			hirsiConfigEnvVar: "/env/path/hirsi.toml",
			expectedContent:   "i am from the env var",
		},
		{
			name:              "HIRSI_CONFIG specified, and doesn't exist",
			hirsiConfigEnvVar: "/env/but/not/exists/hirsi.toml",
			expectedErr:       os.ErrNotExist,
		},
		{
			name:            "only XDG specified, and exists",
			xdgEnvVar:       "/user/home/custom/xdg",
			expectedContent: "I am from xdg subpath",
		},
		{
			name:        "only XDG specified, and doesn't exist",
			xdgEnvVar:   "/user/home/custom/no/xdg",
			expectedErr: os.ErrNotExist,
		},
		{
			name:              "HIRSI_CONFIG specified, and exists, with local config",
			hirsiConfigEnvVar: "/env/path/hirsi.toml",
			localConfig:       true,
			expectedContent:   "i am from the env var",
		},
		{
			name:              "HIRSI_CONFIG specified, and doesn't exist, with local config",
			hirsiConfigEnvVar: "/env/but/not/exists/hirsi.toml",
			localConfig:       true,
			expectedErr:       os.ErrNotExist,
		},
		{
			name:            "only XDG specified, and exists, with local config",
			xdgEnvVar:       "/user/home/custom/xdg",
			localConfig:     true,
			expectedContent: "I am from pwd",
		},
		{
			name:            "only XDG specified, and doesn't exist, with local config",
			xdgEnvVar:       "/user/home/custom/no/xdg",
			localConfig:     true,
			expectedContent: "I am from pwd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			env := &testEnvironment{
				env: map[string]string{
					"HIRSI_CONFIG":    tc.hirsiConfigEnvVar,
					"XDG_CONFIG_HOME": tc.xdgEnvVar,
				},
				home: "/user/home",
				pwd:  "/src",
			}

			fs := testFs{
				"/env/path/hirsi.toml":                   "i am from the env var",
				"/user/home/custom/xdg/hirsi/hirsi.toml": "I am from xdg subpath",
			}

			if tc.localConfig {
				fs["/src/hirsi.toml"] = "I am from pwd"
			}

			content, err := findConfigFile(env, fs)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}

			assert.Equal(t, tc.expectedContent, string(content))

		})
	}
}

type testFs map[string]string

func (fs testFs) ReadFile(name string) ([]byte, error) {
	if content, found := fs[name]; found {
		return []byte(content), nil
	}

	return nil, os.ErrNotExist

}

func (fs testFs) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented")
}

type testEnvironment struct {
	env map[string]string

	home      string
	homeError error

	pwd      string
	pwdError error
}

func (e *testEnvironment) GetEnv(key string) string {
	if val, found := e.env[key]; found {
		return val
	}
	return ""
}

func (e *testEnvironment) GetHome() (string, error) {
	if e.homeError != nil {
		return "", e.homeError
	}

	return e.home, nil
}

func (e *testEnvironment) GetPwd() (string, error) {
	if e.pwdError != nil {
		return "", e.pwdError
	}

	return e.pwd, nil
}
