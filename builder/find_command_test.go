package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
)

func Test_find_api_in_set(t *testing.T) {
	t.Run("an_empty_api_set_will_return_nil", func(t *testing.T) {
		set := builder.APISet{}

		actual, err := set.FindAPI("any")

		assert.Error(t, err)
		assert.Empty(t, actual)
	})
	t.Run("find_supports_case_insensitive_substring_match", func(t *testing.T) {
		set := builder.APISet{APIs: []builder.API{
			{Name: "one_two_three", BaseURL: "http://1.org"},
			{Name: "DOG_CAT_BIRD", BaseURL: "http://2.org"},
			{Name: "AppleBananaCherry", BaseURL: "http://4.org"},
			{Name: "cat", BaseURL: "http://4.org"},
		}}

		actual, err := set.FindAPI("cat")

		require.NoError(t, err)
		require.Equal(t, 2, len(actual))

		assertResult := func(expected string) {
			if actual[0].BaseURL != expected &&
				actual[1].BaseURL != expected {
				assert.Failf(t, "missing '%s'", expected)
			}
		}
		assertResult("http://2.org")
		assertResult("http://4.org")
	})
}

func Test_find_command_in_api(t *testing.T) {
	t.Run("an_empty_api_will_return_nil", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{
				{Name: "any"},
			},
		}

		actual, err := set.APIs[0].FindCommand("any")

		assert.Error(t, err)
		assert.Empty(t, actual)
	})
	t.Run("find_supports_case_insensitive_substring_match", func(t *testing.T) {
		set := builder.APISet{APIs: []builder.API{{Name: "any",
			Commands: []builder.Command{
				{Name: "DOG_CAT_BIRD", Path: "/1"},
				{Name: "one_two_three", Path: "/2"},
				{Name: "cat", Path: "/3"},
				{Name: "AppleBananaCherry", Path: "/4"},
			},
		}}}

		actual, err := set.APIs[0].FindCommand("cat")

		require.NoError(t, err)
		require.Equal(t, 2, len(actual))

		assertResult := func(expected string) {
			if actual[0].Path != expected &&
				actual[1].Path != expected {
				assert.Failf(t, "missing '%s'", expected)
			}
		}
		assertResult("/1")
		assertResult("/3")
	})
}
