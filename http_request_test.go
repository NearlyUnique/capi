package capi_test

import (
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
					},
				},
			},
		},
	}
	var actualRequest http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequest = *r
		w.WriteHeader(http.StatusTeapot)
	}))

	p.APIs[0].BaseURL = ts.URL

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

		resp, err := ts.Client().Do(req)
		assert.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, "/post", actualRequest.URL.String())
		assert.Equal(t, "POST", actualRequest.Method)
		assert.Contains(t, "application/json", actualRequest.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusTeapot, resp.StatusCode)
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

}
