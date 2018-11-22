package capi_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/NearlyUnique/capi"
)

func Test_http_request_can_be_made(t *testing.T) {
	p := capi.Profile{
		APIs: []capi.API{
			{
				Name:    "an_api",
				BaseURL: "https://postman-echo.com",
				DefaultHeader: map[string]string{
					"any-default-header":  "default value",
					"default-header-key":  "before {any-arg} after",
					"overwritten-default": "-overwrite-this-",
				},
				Commands: []capi.Command{
					{
						Name:   "api_cmd",
						Path:   "/post",
						Method: "POST",
						Header: map[string]string{
							"content-type": "application/json",
						},
					}, {
						Name: "api_cmd_args",
						Path: "/sub/{v1}/path/{v2}",
						Header: map[string]string{
							"command-header-key":  "{any-arg}",
							"overwritten-default": "overwritten: {any-arg}",
						},
					},
				},
			},
		},
	}

	args := []string{"cli_cmd", "an_api", "api_cmd"}
	cmd, err := capi.Prepare(p, args)

	t.Run("for valid arguments the cmd is prepared", func(t *testing.T) {
		assert.NoError(t, err)
		require.NotNil(t, cmd)
		assert.Equal(t, "/post", cmd.Path)
	})
	t.Run("a request can be sent", func(t *testing.T) {
		req, err := capi.CreateRequest(cmd)
		assert.NoError(t, err)
		assert.Equal(t, "/post", req.URL.Path)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})
	t.Run("literal url profile generates the same url", func(t *testing.T) {
		req, err := capi.CreateRequest(cmd)
		assert.NoError(t, err)

		assert.Equal(t, "/post", req.URL.Path)

	})
	t.Run("url parameters are replaced", func(t *testing.T) {
		args := []string{"cli_cmd", "an_api", "api_cmd_args", "--v1=first", "--v2", "second"}
		cmd, err := capi.Prepare(p, args)
		require.NoError(t, err)
		req, err := capi.CreateRequest(cmd)
		assert.NoError(t, err)

		assert.Equal(t, "/sub/first/path/second", req.URL.Path)

	})
	t.Run("headers with value expansion are added", func(t *testing.T) {
		args := []string{"cli_cmd", "an_api", "api_cmd_args", "--any-arg", "value1"}
		cmd, err := capi.Prepare(p, args)
		require.NoError(t, err)

		req, err := capi.CreateRequest(cmd)
		require.NoError(t, err)

		assert.Equal(t, "value1", req.Header.Get("command-header-key"))
		assert.Equal(t, "before value1 after", req.Header.Get("default-header-key"))
		assert.Equal(t, "overwritten: value1", req.Header.Get("overwritten-default"))
	})
	t.Run("for a POST, when data has a value it is sent", func(t *testing.T) {
		base := p.APIs[0].BaseURL
		defer func() {
			p.APIs[0].Commands[0].Data = nil
			p.APIs[0].BaseURL = base
		}()
		// ^^^ fix test, probably better to have separate test data ^^^

		const example = `{"a_key":"a_value"}`
		var actual []byte
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			actual, err = ioutil.ReadAll(r.Body)
			require.NoError(t, err)
		}))

		p.APIs[0].Commands[0].Data = []byte(example)
		p.APIs[0].BaseURL = ts.URL

		args := []string{"cli_cmd", "an_api", "api_cmd"}
		cmd, err := capi.Prepare(p, args)
		require.NoError(t, err)

		req, err := capi.CreateRequest(cmd)
		require.NoError(t, err)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, http.StatusOK)

		assert.Equal(t, len(actual), len(example))
		assert.Contains(t, string(actual), `"a_key"`)
	})
}
