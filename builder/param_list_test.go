package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
)

func Test_all_components_are_listed(t *testing.T) {
	set := builder.APISet{
		APIs: []builder.API{{
			Name:    "an_api",
			BaseURL: "{base_url_value}",
			Header:  map[string][]string{"api_header": {"{api_header_value}"}},
			Query:   map[string][]string{"api_query": {"{api_query_value}"}},
			Commands: []builder.Command{{
				Name:   "a_cmd",
				Path:   "/one/{path_value}",
				Header: map[string][]string{"header1": {"{header_value}"}},
				Query:  map[string][]string{"query1": {"{query_value}"}},
				Body:   &builder.CommandBody{Data: []byte(`"{body_value}"`)},
			}},
		}},
	}
	api, err := set.FindAPI("an_api")
	require.NoError(t, err)
	cmd, err := api[0].FindCommand("a_cmd")
	require.NoError(t, err)

	actual := builder.ListParams(cmd[0])

	assert.Equal(t, 7, len(actual))
	assert.Contains(t, actual, "base_url_value")
	assert.Contains(t, actual, "path_value")
	assert.Contains(t, actual, "api_header_value")
	assert.Contains(t, actual, "api_query_value")
	assert.Contains(t, actual, "header_value")
	assert.Contains(t, actual, "query_value")
	assert.Contains(t, actual, "body_value")
}

func Test_a_parameter_should_only_appear_once(t *testing.T) {
	set := builder.APISet{
		APIs: []builder.API{{
			Name:    "an_api",
			BaseURL: "{some_param}",
			Commands: []builder.Command{{
				Name: "a_cmd",
				Path: "/one/{some_param}",
			}},
		}},
	}
	api, err := set.FindAPI("an_api")
	require.NoError(t, err)
	cmd, err := api[0].FindCommand("a_cmd")
	require.NoError(t, err)

	actual := builder.ListParams(cmd[0])

	assert.Equal(t, 1, len(actual))
	assert.Contains(t, actual, "some_param")
}
