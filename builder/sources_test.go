package builder_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sources(t *testing.T) {
	t.Run("base url parameters are replaced", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name: "the_name", BaseURL: "{some_url}",
				Commands: []builder.Command{{Name: "anyName"}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("some_url", "http://example.com"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)
		assertStartsWith(t, "http://example.com", req.URL.String())
	})
	t.Run("path parameters are replaced", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name: "the_name", BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Path: "/any/{arg1}",
				}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("arg1", "value_one"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)
		assertEndsWith(t, "/any/value_one", req.URL.String())
	})
	t.Run("if a source is a nil func it is ignored", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name: "the_name", BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Path: "/any/{arg1}",
				}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest(
			nil,
			fakeSource("arg1", "value_one"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)
		assertEndsWith(t, "/any/value_one", req.URL.String())
	})
	t.Run("all url parameters are replaced", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "{some_url}",

				Commands: []builder.Command{{
					Name: "anyName",
					Path: "/any/{arg1}",
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("some_url", "https://example.com"),
			fakeSource("arg1", "value_one"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)
		assert.Equal(t, "https://example.com/any/value_one", req.URL.String())
	})
	t.Run("default headers are inserted with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Header: map[string]builder.StringOrList{
					"h1": {"hv1a", "hv1b"},
					"h2": {"{some_header}"},
				},
				Commands: []builder.Command{{Name: "anyName"}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("some_header", "some_value"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)
		// note that Header.Get("key") only returns the FIRST value
		// when keys are added they are converted to UPPERCASE
		assert.Equal(t, []string{"hv1a", "hv1b"}, req.Header["H1"])
		assert.Equal(t, []string{"some_value"}, req.Header["H2"])
	})
	t.Run("command headers are inserted with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Header: map[string]builder.StringOrList{
						"h3": {"{command_header}"},
					},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("command_header", "cmd_header_value"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, "cmd_header_value", req.Header.Get("h3"))
	})
	t.Run("command headers completely override are inserted with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Header: map[string]builder.StringOrList{
					"h4": {"default_value"},
				},
				Commands: []builder.Command{{
					Name: "anyName",
					Header: map[string]builder.StringOrList{
						"h4": {"overridden {command_override}"},
					},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("command_override", "from_source"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, "overridden from_source", req.Header.Get("h4"))
	})
	t.Run("default api query is inserted with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:     "the_name",
				BaseURL:  "{some_url}",
				Query:    map[string]builder.StringOrList{"anyKey": {"{some_arg}"}},
				Commands: []builder.Command{{Name: "anyName"}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("some_arg", "arg_value"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Contains(t, req.URL.RawQuery, "anyKey=arg_value")
	})
	t.Run("command query is inserted", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Query: map[string]builder.StringOrList{
						"simple": {"a_value"},
					},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest()

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Contains(t, req.URL.RawQuery, "simple=a_value")
	})
	t.Run("command query is inserted with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Query: map[string]builder.StringOrList{
						"q3": {"{query_value}"},
					},
				}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest(
			fakeSource("query_value", "replacement_value"),
		)

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Contains(t, req.URL.RawQuery, "q3=replacement_value")
	})
	t.Run("command query overrides default api query", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Query: map[string]builder.StringOrList{
					"q": {"api_value"},
				},
				Commands: []builder.Command{{
					Name: "anyName",
					Query: map[string]builder.StringOrList{
						"q": {"command_value"},
					},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest()

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Contains(t, req.URL.RawQuery, "q=command_value")
	})
	t.Run("if no method is set, GET is used", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name: "the_name", BaseURL: "https://example.com",
				Commands: []builder.Command{{
					Name: "anyName",
				}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest()

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, http.MethodGet, req.Method)
	})
	t.Run("the method defined in the command is used in the request", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name: "the_name", BaseURL: "https://example.com",
				Commands: []builder.Command{{
					Name:   "anyName",
					Method: http.MethodPatch,
				}},
			}},
		}
		req, err := firstCmd(t, set).CreateRequest()

		assert.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, http.MethodPatch, req.Method)
	})
	t.Run("if a body is set, it is included in the request", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Body: &builder.CommandBody{Data: []byte(`some data`)},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest()

		assert.NoError(t, err)
		require.NotNil(t, req)
		require.NotNil(t, req.Body)

		buf, err := ioutil.ReadAll(req.Body)
		defer func() { _ = req.Body.Close() }()

		assert.Equal(t, "some data", string(buf))
	})
	t.Run("if a body is set, it is included in the request with replacement text", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "the_name",
				BaseURL: "http://any.org",
				Commands: []builder.Command{{
					Name: "anyName",
					Body: &builder.CommandBody{Data: []byte(`{"any":"{the_arg}"}`)},
				}},
			}},
		}

		req, err := firstCmd(t, set).CreateRequest(fakeSource("the_arg", "replacement_value"))

		assert.NoError(t, err)
		require.NotNil(t, req)
		require.NotNil(t, req.Body)

		buf, err := ioutil.ReadAll(req.Body)
		defer func() { _ = req.Body.Close() }()

		assert.Equal(t, `{"any":"replacement_value"}`, string(buf))
	})
}

func Test_api_command_search(t *testing.T) {
	set := builder.APISet{
		APIs: []builder.API{
			{
				Name:    "the_api",
				BaseURL: "https://one.example.com",
				Commands: []builder.Command{
					{Name: "the_command", Path: "/path"},
					{Name: "other_command"},
				},
			},
			{
				Name:    "other_api",
				BaseURL: "https://two.example.com",
				Commands: []builder.Command{
					{Name: "the_command"},
				},
			},
		},
	}
	t.Run("zero api matches returns error", func(t *testing.T) {
		actual, err := set.CreateRequest("unknown", "")

		assert.Error(t, err)
		_, ok := err.(builder.NotFound)
		assert.True(t, ok)
		assert.Nil(t, actual)
	})
	t.Run("more than one api matches returns error", func(t *testing.T) {
		actual, err := set.CreateRequest("api", "")

		assert.Error(t, err)
		_, ok := err.(builder.NotFound)
		assert.True(t, ok)
		assert.Nil(t, actual)
	})
	t.Run("zero command matches returns error", func(t *testing.T) {
		actual, err := set.CreateRequest("the_api", "unknown")

		assert.Error(t, err)
		_, ok := err.(builder.NotFound)
		assert.True(t, ok)
		assert.Nil(t, actual)
	})
	t.Run("more than one command matches returns error", func(t *testing.T) {
		actual, err := set.CreateRequest("the_api", "command")

		assert.Error(t, err)
		_, ok := err.(builder.NotFound)
		assert.True(t, ok)
		assert.Nil(t, actual)
	})
	t.Run("one api and one command returns http request", func(t *testing.T) {
		actual, err := set.CreateRequest("the_api", "the_command")

		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.Equal(t, "https://one.example.com/path", actual.URL.String())
	})
}

func Test_env_source_is_case_insensitive(t *testing.T) {
	keys := []string{"ONE", "one", "One", "onE", "oNe"}
	testData := []struct {
		envVar, value string
	}{
		{"one=singleValue", "singleValue"},
		{"one=has=equals", "has=equals"},
		{"one=endsWithEquals=", "endsWithEquals="},
		{"one==startsWithEquals", "=startsWithEquals"},
		{"one=has space", "has space"},
	}
	for _, item := range testData {
		envSrc := builder.NewEnvVarsSource([]string{item.envVar})
		for _, k := range keys {
			assert.Equal(t, item.value, envSrc(k))
		}
	}
}
