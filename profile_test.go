package capi_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NearlyUnique/capi"
	"github.com/stretchr/testify/require"
)

func Test_when_loaded_the_profile_contains_apis(t *testing.T) {
	var r io.Reader = strings.NewReader(`{
	"envPrefix": "a_prefix",
    "apis": [{
      "name": "a-name",
      "baseUrl": {"env1": "http://localhost:8080", "env2": "https://example.com/root_path"},
      "commands": [ {
          "name": "a-name", "method": "GET", "path": "/some/path",
          "header": { "any-header":"" }
        } ]
	} ]
}`)
	p, err := capi.LoadAPI(ioutil.NopCloser(r))

	require.NoError(t, err)

	assert.Equal(t, "a_prefix", p.EnvPrefix)
	assert.Equal(t, "https://example.com/root_path", p.APIs[0].BaseURL["env2"])
}

func Test_an_api_can_be_selected_by_name(t *testing.T) {
	p := capi.Profile{
		APIs: []capi.API{
			{Name: "first"},
			{Name: "second"},
		},
	}
	t.Run("when name matches exactly", func(t *testing.T) {
		actual, err := p.SelectAPI([]string{"cmd-name", "second"})
		assert.NoError(t, err)
		assert.Equal(t, "second", actual.Name)
	})
	t.Run("returns an error if not found", func(t *testing.T) {
		actual, err := p.SelectAPI([]string{"cmd-name", "unknown"})
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
	t.Run("returns an error if insufficient args", func(t *testing.T) {
		_, err := p.SelectAPI([]string{})
		assert.Error(t, err)
		_, err = p.SelectAPI([]string{"cmd-name"})
		assert.Error(t, err)
	})
}

func Test_a_command_can_be_selected_by_name(t *testing.T) {
	p := capi.Profile{
		APIs: []capi.API{
			{Name: "first",
				Commands: []capi.Command{
					{Name: "cmd1", Path: "/one"},
					{Name: "cmd1", Path: "/two"},
				},
			},
			{Name: "second"},
		},
	}
	api, err := p.SelectAPI([]string{"cmd", "first"})
	require.NoError(t, err)

	t.Run("when name matches exactly", func(t *testing.T) {
		actual, err := p.SelectCommand(api, []string{"cmd-name", "any", "cmd1"})
		assert.NoError(t, err)
		assert.Equal(t, "cmd1", actual.Name)
		assert.Equal(t, "/one", actual.Path)
	})
	t.Run("returns an error if not found", func(t *testing.T) {
		actual, err := p.SelectCommand(api, []string{"cmd-name", "any", "unknown"})
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
	t.Run("returns an error if insufficient args", func(t *testing.T) {
		_, err := p.SelectCommand(api, []string{})
		assert.Error(t, err)
		_, err = p.SelectCommand(api, []string{"cmd-name"})
		assert.Error(t, err)
		_, err = p.SelectCommand(api, []string{"cmd-name", "api-name"})
		assert.Error(t, err)
	})
	t.Run("returns an error if api is nil", func(t *testing.T) {
		_, err := p.SelectCommand(nil, []string{"any", "arg", "list"})
		assert.Error(t, err)
	})
}
