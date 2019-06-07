package capicomplete_test

import (
	"testing"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/capicomplete"
	"github.com/stretchr/testify/assert"
)

func Test_when_nothing_is_entered_the_full_list_of_apis_is_displayed(t *testing.T) {
	p := builder.APISet{APIs: []builder.API{
		{Name: "one"},
		{Name: "two"},
	}}
	ac := autocomplete.Mock("any ", "")

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one", "two"}, act)
}

func Test_when_partial_api_name_is_given_leading_submatches_are_returned(t *testing.T) {
	p := builder.APISet{APIs: []builder.API{
		{Name: "aaa"},
		{Name: "aaa-other"},
		{Name: "xxx-other"},
	}}

	t.Run("match with name", func(t *testing.T) {
		ac := autocomplete.Mock("any aa", "")
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa", "aaa-other"}, act)
	})
	t.Run("when no matches return none", func(t *testing.T) {
		ac := autocomplete.Mock("any no-match", "")
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Empty(t, act)
	})
}

func Test_searching_for_command_within_api(t *testing.T) {
	p := builder.APISet{APIs: []builder.API{
		{Name: "any-api",
			Commands: []builder.Command{
				{Name: "cmd1", Path: "/one"},
				{Name: "cmd2", Path: "/two"},
			},
		},
		{Name: "any-other-api"},
	}}
	t.Run("if the api exists the matching commands are returned", func(t *testing.T) {
		ac := autocomplete.Mock("app-name any-api c", "")
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, 2, len(act))
		assert.Contains(t, act, "cmd1")
		assert.Contains(t, act, "cmd2")
	})
	t.Run("if the api does not exist an error is returned", func(t *testing.T) {
		ac := autocomplete.Mock("app-name no-such-api c", "")
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, 2, len(act))
		assert.Contains(t, act, "error")
	})
	t.Run("if the api search is ambiguous an error is returned", func(t *testing.T) {
		ac := autocomplete.Mock("app-name any cmd1", "")
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, 2, len(act))
		assert.Contains(t, act, "error")
	})
}

func Test_all_params_are_listed_with_double_dash_prefix(t *testing.T) {
	t.Run("with no user filter", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "an_api",
				BaseURL: "{base_url_value}",
				Commands: []builder.Command{{
					Name:   "a_cmd",
					Path:   "/one/{path_value}",
					Header: map[string][]string{"header1": {"{header_value}"}},
					Query:  map[string][]string{"query1": {"{query_value}"}},
					Body:   &builder.CommandBody{Data: []byte(`"{body_value}"`)},
				}},
			}},
		}
		ac := autocomplete.Mock("any an_api a_cmd ", "")

		actual := capicomplete.GenerateResponse(&ac, &set)

		assert.Equal(t, 5, len(actual))
		assert.Contains(t, actual, "--base_url_value")
		assert.Contains(t, actual, "--path_value")
		assert.Contains(t, actual, "--header_value")
		assert.Contains(t, actual, "--query_value")
		assert.Contains(t, actual, "--body_value")
	})
	t.Run("with ambiguous user filter", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "an_api",
				BaseURL: "{value_one}",
				Commands: []builder.Command{{
					Name: "a_cmd",
					Path: "/one/{value_two}",
				}},
			}},
		}

		ac := autocomplete.Mock("any an_api a_cmd --value", "")

		actual := capicomplete.GenerateResponse(&ac, &set)

		assert.Equal(t, 2, len(actual))
		assert.Contains(t, actual, "--value_one")
		assert.Contains(t, actual, "--value_two")
	})
	t.Run("user filter that can only fit one result returns single flag", func(t *testing.T) {
		set := builder.APISet{
			APIs: []builder.API{{
				Name:    "an_api",
				BaseURL: "{value_one}",
				Commands: []builder.Command{{
					Name: "a_cmd",
					Path: "/one/{value_two}",
				}},
			}},
		}

		ac := autocomplete.Mock("any an_api a_cmd --value_o", "")

		actual := capicomplete.GenerateResponse(&ac, &set)

		assert.Equal(t, 1, len(actual))
		assert.Contains(t, actual, "--value_one")
	})
}
