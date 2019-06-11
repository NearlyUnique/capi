package builder_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_http_request_can_be_made(t *testing.T) {
	set := builder.APISet{
		APIs: []builder.API{
			{
				Name:    "an_api",
				BaseURL: "https://postman-echo.com",
				Header: map[string]builder.StringOrList{
					"any-default-header":  {"default value"},
					"default-header-key":  {"before {any-arg} after"},
					"overwritten-default": {"-overwrite-this-"},
				},
				Commands: []builder.Command{
					{
						Name:   "api_cmd",
						Path:   "/post",
						Method: "POST",
						Header: map[string]builder.StringOrList{
							"content-type": {"application/json"},
						},
					}, {
						Name: "api_cmd_args",
						Path: "/sub/{v1}/path/{v2}",
						Header: map[string]builder.StringOrList{
							"command-header-key":  {"{any-arg}"},
							"overwritten-default": {"overwritten: {any-arg}"},
						},
					},
				},
			},
		},
	}
	api, err := set.FindAPI("an_api")
	require.NoError(t, err)
	cmd, err := api[0].FindCommand("api_cmd")
	require.NoError(t, err)
	aCommand := cmd[0]
	cmd, err = api[0].FindCommand("api_cmd_args")
	require.NoError(t, err)
	cmdWithArgs := cmd[0]

	t.Run("a request can be sent", func(t *testing.T) {
		req, err := aCommand.CreateRequest()
		assert.NoError(t, err)
		assert.Equal(t, "/post", req.URL.Path)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})
	t.Run("literal url apiSet generates the same url", func(t *testing.T) {
		req, err := aCommand.CreateRequest()
		assert.NoError(t, err)

		assert.Equal(t, "/post", req.URL.Path)

	})
	t.Run("url parameters are replaced", func(t *testing.T) {
		req, err := cmdWithArgs.CreateRequest(
			fakeSource("v1", "first"),
			fakeSource("v2", "second"),
		)
		assert.NoError(t, err)

		assert.Equal(t, "/sub/first/path/second", req.URL.Path)
	})
	t.Run("headers with value expansion are added", func(t *testing.T) {

		req, err := cmdWithArgs.CreateRequest(
			fakeSource("any-arg", "value1"),
		)
		require.NoError(t, err)

		assert.Equal(t, "value1", req.Header.Get("command-header-key"))
		assert.Equal(t, "before value1 after", req.Header.Get("default-header-key"))
		assert.Equal(t, "overwritten: value1", req.Header.Get("overwritten-default"))
	})
	t.Run("for a POST, when data has a value it is sent", func(t *testing.T) {
		base := set.APIs[0].BaseURL
		defer func() {
			aCommand.Body = nil
			set.APIs[0].BaseURL = base
		}()
		// ^^^ fix test, probably better to have separate test data ^^^

		const example = `{"a_key":"a_value","substitute":"{arg1}"}`
		var actual []byte
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { _ = r.Body.Close() }()
			var err error
			actual, err = ioutil.ReadAll(r.Body)
			require.NoError(t, err)
		}))

		aCommand.Body = &builder.CommandBody{Data: []byte(example)}
		set.APIs[0].BaseURL = ts.URL

		req, err := aCommand.CreateRequest(
			fakeSource("arg1", "value1"),
		)
		require.NoError(t, err)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, http.StatusOK)

		assert.Equal(t, len(example), len(actual))
		assert.Contains(t, string(actual), `"a_key"`)
		assert.Contains(t, string(actual), `"value1"`)
	})
}
