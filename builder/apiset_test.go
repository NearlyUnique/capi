package builder_test

import (
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_an_apiset_is_empty_by_default(t *testing.T) {
	b := builder.APISet{}
	assert.Equal(t, len(b.APIs), 0)
}

func Test_adding_apis(t *testing.T) {
	t.Run("returns error when name is not set", func(t *testing.T) {
		b := builder.APISet{}
		err := b.Add(builder.API{BaseURL: "any"})
		assert.Error(t, err)
	})
	t.Run("returns error when baseURL is not set", func(t *testing.T) {
		b := builder.APISet{}
		err := b.Add(builder.API{Name: "any"})
		assert.Error(t, err)
	})
	t.Run("after adding an API it can be listed", func(t *testing.T) {
		b := builder.APISet{}
		err := b.Add(builder.API{Name: "the name", BaseURL: "the url"})

		require.NoError(t, err)
		assert.Equal(t, len(b.APIs), 1)
		assert.Equal(t, b.APIs[0].Name, "the name")
		assert.Equal(t, b.APIs[0].BaseURL, "the url")
	})
	t.Run("after adding, the listed API links to the APISet", func(t *testing.T) {
		apiSet := builder.APISet{}
		err := apiSet.Add(builder.API{Name: "the name", BaseURL: "the url"})

		require.NoError(t, err)
		assert.Equal(t, *apiSet.APIs[0].Set, apiSet)
	})
	t.Run("an api initially has no default headers", func(t *testing.T) {
		api := builder.API{Name: "the name", BaseURL: "the url"}

		assert.Empty(t, api.Header)
	})
	t.Run("an api initially has no default query string", func(t *testing.T) {
		api := builder.API{Name: "the name", BaseURL: "the url"}

		assert.Empty(t, api.Query)
	})
	t.Run("an api initially has no commands", func(t *testing.T) {
		api := builder.API{Name: "the name", BaseURL: "the url"}

		assert.Empty(t, api.Commands)
	})
}

func Test_an_api_has_commands(t *testing.T) {
	t.Run("Adding a command to an api without a name generates an error", func(t *testing.T) {
		api := builder.API{}

		err := api.Add(builder.Command{})

		assert.Error(t, err)
	})
	t.Run("A command added to an api can be listed", func(t *testing.T) {
		api := builder.API{}
		err := api.Add(builder.Command{Name: "any name"})
		require.NoError(t, err)

		assert.Equal(t, len(api.Commands), 1)
		assert.Equal(t, api.Commands[0].Name, "any name")
	})
	t.Run("A command added to an api links to that api", func(t *testing.T) {
		api := builder.API{}
		err := api.Add(builder.Command{Name: "any name"})
		require.NoError(t, err)

		assert.Equal(t, api, *api.Commands[0].API)
	})
	t.Run("When a command has no Method set it is defaulted to GET", func(t *testing.T) {
		api := builder.API{}
		err := api.Add(builder.Command{Name: "any name"})
		require.NoError(t, err)

		assert.Equal(t, "GET", api.Commands[0].Method)
	})
	t.Run("a command initially has no default headers", func(t *testing.T) {
		api := builder.Command{Name: "the name"}

		assert.Empty(t, api.Header)
	})
	t.Run("a command initially has no default query string", func(t *testing.T) {
		api := builder.Command{Name: "the name"}

		assert.Empty(t, api.Query)
	})
}
func Test_an_http_request_is_generated_from_a_command_with_default_GET_method(t *testing.T) {
	set := builder.APISet{}
	err := set.Add(builder.API{Name: "any", BaseURL: "http://example.com"})
	require.NoError(t, err)

	err = set.APIs[0].Add(builder.Command{Name: "any cmd"})
	require.NoError(t, err)

	req, err := set.APIs[0].Commands[0].CreateRequest()

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", req.URL.String())
	assert.Equal(t, "GET", req.Method)
	assert.Empty(t, req.Header)
}
func Test_an_http_request_is_generated_from_a_command_with_command_path_added(t *testing.T) {
	set := builder.APISet{}
	err := set.Add(builder.API{Name: "any", BaseURL: "http://example.com"})
	require.NoError(t, err)

	err = set.APIs[0].Add(builder.Command{
		Name:   "any cmd",
		Path:   "/some-segment",
		Method: "PUT",
	})
	require.NoError(t, err)

	req, err := set.APIs[0].Commands[0].CreateRequest()

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/some-segment", req.URL.String())
	assert.Equal(t, "PUT", req.Method)
	assert.Empty(t, req.Header)
}
