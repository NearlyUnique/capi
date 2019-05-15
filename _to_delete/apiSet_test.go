package _to_delete_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/NearlyUnique/capi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_when_loaded_the_api_set_contains_apis(t *testing.T) {
	var r io.Reader = strings.NewReader(`{
	"envPrefix": "a_prefix",
    "apis": [{
      "name": "api-name",
      "baseUrl": "https://example.com/root_path",
      "commands": [ {
          "name": "cmd-name", "method": "GET", "path": "/some/path",
          "header": { "a-header":"a-header-value" }
        } ]
	} ]
}`)
	p, err := capi.ParseAPI(ioutil.NopCloser(r))

	require.NoError(t, err)

	assert.Equal(t, "a_prefix", p.EnvPrefix)
	api := p.APIs[0]
	assert.Equal(t, "api-name", api.Name)
	assert.Equal(t, "https://example.com/root_path", api.BaseURL)
	cmd := api.Commands[0]
	assert.Equal(t, "cmd-name", cmd.Name)
	assert.Equal(t, "GET", cmd.Method)
	assert.Equal(t, "/some/path", cmd.Path)
	assert.Equal(t, "a-header-value", cmd.Header["a-header"])
}

func Test_an_api_can_be_selected_by_name(t *testing.T) {
	set := capi.APISet{
		APIs: []capi.API{
			{Name: "first"},
			{Name: "second"},
		},
	}
	t.Run("when name matches exactly", func(t *testing.T) {
		actual, err := set.SelectAPI("second")
		assert.NoError(t, err)
		assert.Equal(t, "second", actual.Name)
	})
	t.Run("returns an error if not found", func(t *testing.T) {
		actual, err := set.SelectAPI("unknown")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
	t.Run("returns an error if empty", func(t *testing.T) {
		actual, err := set.SelectAPI("")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func Test_a_command_can_be_selected_by_name(t *testing.T) {
	p := capi.APISet{
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
	api, err := p.SelectAPI("first")
	require.NoError(t, err)

	t.Run("when name matches exactly", func(t *testing.T) {
		actual, err := api.SelectCommand("cmd1")
		assert.NoError(t, err)
		assert.Equal(t, "cmd1", actual.Name)
		assert.Equal(t, "/one", actual.Path)
	})
	t.Run("returns an error if not found", func(t *testing.T) {
		actual, err := api.SelectCommand("unknown")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
	t.Run("returns an error no name", func(t *testing.T) {
		_, err := api.SelectCommand("")
		assert.Error(t, err)
	})
}

func Test_list_params_for_a_command(t *testing.T) {
	t.Run("params for a command include path variables", func(t *testing.T) {
		cmd := capi.Command{Path: "/{arg1}/literal/{arg2}"}

		actual := cmd.ListParams()

		assert.Contains(t, actual, "arg1")
		assert.Contains(t, actual, "arg2")
		assert.Equal(t, 2, len(actual))
	})
	t.Run("params for a command include header keys", func(t *testing.T) {
		cmd := capi.Command{
			Header: map[string]string{
				"header1":           "",
				"underscore_header": "",
				"dash-header":       "",
			},
		}

		actual := cmd.ListParams()

		assert.Contains(t, actual, "header1")
		assert.Contains(t, actual, "underscore_header")
		assert.Contains(t, actual, "dash-header")
		assert.Equal(t, 3, len(actual))
	})
	t.Run("params for a command include header values", func(t *testing.T) {
		cmd := capi.Command{
			Header: map[string]string{
				"h0": "",
				"h1": "{arg1}",
				"h2": "prefix {arg2}",
				"h3": "surround {arg3} with",
			},
		}

		actual := cmd.ListParams()

		// unavoidable header keys
		assert.Contains(t, actual, "h0")
		//values
		assert.Contains(t, actual, "arg1")
		assert.Contains(t, actual, "arg2")
		assert.Contains(t, actual, "arg2")
		assert.Equal(t, 4, len(actual))
	})
	t.Run("params will only appear once even it duplicated in definitions", func(t *testing.T) {
		cmd := capi.Command{
			Path: "/{duplicate}/any/{path1}",
			Header: map[string]string{
				"header1":   "{duplicate}",
				"duplicate": "",
			},
		}

		actual := cmd.ListParams()

		assert.Contains(t, actual, "path1")
		assert.Contains(t, actual, "duplicate")
		assert.Equal(t, 2, len(actual))
	})
	t.Run("params are read from default headers", func(t *testing.T) {
		api := capi.API{
			DefaultHeader: map[string]string{
				"from-default":       "",
				"default-with-param": "{default-param}",
			},
		}
		cmd := capi.Command{
			Path: "/any",
			Header: map[string]string{
				"cmd-header": "",
			},
		}
		api.Add(&cmd)

		actual := cmd.ListParams()

		assert.Contains(t, actual, "cmd-header")
		assert.Contains(t, actual, "from-default")
		assert.Contains(t, actual, "default-param")
		assert.Equal(t, 3, len(actual))
	})
	t.Run("params can come from data", func(t *testing.T) {
		cmd := capi.Command{
			Path: "/any",
			Data: []byte(`{"any":"{v1}","list":["{v2}"]}`),
		}

		actual := cmd.ListParams()

		assert.Contains(t, actual, "v1")
		assert.Contains(t, actual, "v2")
	})
}

func Test_lookup_value_by_name(t *testing.T) {
	local := capi.LocalEnv{
		EnvPrefix: "XXX_",
		Environ: func() []string {
			return []string{
				"SOME_ENV_VAR_KEY=a value",
				"XXX_VALUE_DEFINED_WITH_PREFIX=another value",
				"xxx_prefix_has_precedence=user prefix has precedence",
				"prefix_has_precedence=this is never found",
			}
		},
	}

	testData := []struct {
		name, key, expected string
	}{
		{"key as defined", "SOME_ENV_VAR_KEY", "a value"},
		{"lower case", "some_env_var_key", "a value"},
		{"mixed case", "SOME_env_var_key", "a value"},

		{"with prefix", "VALUE_DEFINED_WITH_PREFIX", "another value"},

		{"unknown env variable", "unknown_env_variable", ""},

		{"user prefix has precedence", "prefix_has_precedence", "user prefix has precedence"},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {

			value := local.Lookup(td.key)

			assert.Equal(t, td.expected, value)
		})
	}
}

func Test_data_is_loaded_as_raw_json(t *testing.T) {
	buf := []byte(`{"name":"cmd_name", "data": { "sk":"{arg}", "ik": 12, "bk": true, "ak":[1] }}`)

	var actual capi.Command
	err := json.Unmarshal(buf, &actual)

	var rx = regexp.MustCompile(`{(?P<Name>[a-zA-Z0-9-_]+)}`)
	rx.ReplaceAllFunc(actual.Data, func(buf []byte) []byte {
		v := string(buf)
		_ = v
		return buf
	})

	require.NoError(t, err)
	require.Equal(t, "cmd_name", actual.Name)
	assert.Contains(t, string(actual.Data), `"sk"`)
}
