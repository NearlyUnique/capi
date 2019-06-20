package postman

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_APISet_basic_property_values(t *testing.T) {
	pm := ItemConfig{
		Info: Info{Name: "api-name"},
		Item: []Item{
			{
				Name: "cmd-name",
				Request: Request{
					URL: URL{
						Raw: "some-url?a-query1=a-query-value1&a-query2=a-query-value2",
						Query: []Value{
							{Key: "a-query1", Value: "a-query-value1"},
							{Key: "a-query2", Value: "a-query-value2"},
						},
					},
					Method: "a-method",
					Header: []Value{
						{Key: "a-header1", Value: "a-header-value1"},
						{Key: "a-header2", Value: "a-header-value2"},
					},
					Body: Body{
						Raw: `{
	"attrib1": "attrib-value1"
}`,
					},
				},
			},
		},
	}
	set, err := pm.ToAPISet()
	require.NoError(t, err)

	t.Run(".Name from .Info.Name", func(t *testing.T) {
		require.Equal(t, 1, len(set.APIs))
		assert.Equal(t, "api-name", set.APIs[0].Name)
	})
	t.Run("For this test there is exactly one command", func(t *testing.T) {
		require.Equal(t, 1, len(set.APIs))
		require.Equal(t, 1, len(set.APIs[0].Commands))

		command := set.APIs[0].Commands[0]

		t.Run(".Cmd[].Name from .Items[].Name", func(t *testing.T) {
			assert.Equal(t, "cmd-name", command.Name)
		})
		t.Run(".Cmd[].Path from .Items[].Request.URL.Raw without the query", func(t *testing.T) {
			assert.Equal(t, "some-url", command.Path)
		})
		t.Run(".Cmd[].Method from .Items[].Request.Method", func(t *testing.T) {
			assert.Equal(t, "a-method", command.Method)
		})
		t.Run(".Cmd[].Header from .Items[].Request.Header", func(t *testing.T) {
			require.Equal(t, 2, len(command.Header))
			assert.Contains(t, command.Header, "a-header1")
			assert.Equal(t, "a-header-value1", command.Header["a-header1"][0])
			assert.Contains(t, command.Header, "a-header2")
			assert.Equal(t, "a-header-value2", command.Header["a-header2"][0])
		})
		t.Run(".Cmd[].Query from .Items[].Request.URL.Query", func(t *testing.T) {
			require.Equal(t, 2, len(command.Query))
			assert.Contains(t, command.Query, "a-query1")
			assert.Equal(t, "a-query-value1", command.Query["a-query1"][0])
			assert.Contains(t, command.Query, "a-query2")
			assert.Equal(t, "a-query-value2", command.Query["a-query2"][0])
		})
		t.Run(".Cmd[].Body from .Items[].Request.Body.Raw and is json parsable", func(t *testing.T) {
			require.NotNil(t, command.Body)
			var actual struct{ Attrib1 string }
			err = json.Unmarshal(command.Body.Data, &actual)
			require.NoError(t, err)
			assert.Equal(t, actual.Attrib1, "attrib-value1")
		})
	})
}

func Test_when_a_map_key_is_duplicated_in_postman_it_will_also_duplicated_in_APISet(t *testing.T) {
	pm := ItemConfig{
		Item: []Item{{
			Request: Request{
				Header: []Value{
					{Key: "a-duplicate-h", Value: "a-header-value1"},
					{Key: "a-duplicate-h", Value: "a-header-value2"},
				},
				URL: URL{
					Query: []Value{
						{Key: "a-duplicate-q", Value: "a-query-value1"},
						{Key: "a-duplicate-q", Value: "a-query-value2"},
					},
				},
			},
		}},
	}
	set, err := pm.ToAPISet()
	require.NoError(t, err)
	require.Equal(t, 1, len(set.APIs))
	require.Equal(t, 1, len(set.APIs[0].Commands))
	command := set.APIs[0].Commands[0]

	t.Run("for headers", func(t *testing.T) {
		assert.Equal(t, 1, len(command.Header))
		assert.Contains(t, command.Header, "a-duplicate-h")
		assert.Equal(t, 2, len(command.Header["a-duplicate-h"]))
		assert.Contains(t, command.Header["a-duplicate-h"], "a-header-value1")
		assert.Contains(t, command.Header["a-duplicate-h"], "a-header-value2")
	})

	t.Run("for queries", func(t *testing.T) {
		assert.Equal(t, 1, len(command.Query))
		assert.Contains(t, command.Query, "a-duplicate-q")
		assert.Equal(t, 2, len(command.Query["a-duplicate-q"]))
		assert.Contains(t, command.Query["a-duplicate-q"], "a-query-value1")
		assert.Contains(t, command.Query["a-duplicate-q"], "a-query-value2")
	})
}
func Test_postman_variables_are_converted_to_capi_variables(t *testing.T) {
	pm := ItemConfig{
		Item: []Item{
			{
				Request: Request{
					URL: URL{
						Raw:   "{{a-url}}",
						Query: []Value{{Key: "a-query", Value: "{{a-query-value}}"}},
					},
					Header: []Value{{Key: "a-header", Value: "{{a-header-value}}"}},
					Body: Body{
						Raw: `["{{attrib-value}}"]`,
					},
				},
			},
		},
	}
	set, err := pm.ToAPISet()
	require.NoError(t, err)

	command := set.APIs[0].Commands[0]
	assert.Equal(t, "{a-url}", command.Path)
	assert.Equal(t, "{a-query-value}", command.Query["a-query"][0])
	assert.Equal(t, "{a-header-value}", command.Header["a-header"][0])
	assert.Equal(t, `["{attrib-value}"]`, string(command.Body.Data))
}

func Test_spaces_in_command_names_are_replaced_with_underscore(t *testing.T) {
	pm := ItemConfig{
		Item: []Item{{
			Name: "name with spaces",
		}},
	}
	set, err := pm.ToAPISet()
	require.NoError(t, err)

	assert.Equal(t, "name_with_spaces", set.APIs[0].Commands[0].Name)

}
