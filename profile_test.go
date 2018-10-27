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
      "alias": "name",
      "baseUrl": {"env1": "http://localhost:8080", "env2": "https://example.com/root_path"},
      "commands": [ {
          "name": "a-name", "method": "GET", "path": "/some/path",
          "header": [ "any-header" ]
        } ]
	} ]
}`)
	p, err := capi.LoadAPI(ioutil.NopCloser(r))

	require.NoError(t, err)

	assert.Equal(t, "a_prefix", p.EnvPrefix)
	assert.Equal(t, "https://example.com/root_path", p.APIs[0].BaseURL["env2"])
}
