package capicomplete_test

import (
	"testing"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/capicomplete"
	"github.com/stretchr/testify/assert"
)

func Test_when_nothing_is_entered_the_full_list_of_apis_is_displayed(t *testing.T) {
	p := capi.APISet{APIs: []capi.API{
		{Name: "one"},
		{Name: "two"},
	}}
	ac := autocomplete.Mock("any ", "")

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one", "two"}, act)
}

func Test_when_partial_api_name_is_given_leading_submatches_are_returned(t *testing.T) {
	p := capi.APISet{APIs: []capi.API{
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
	p := capi.APISet{APIs: []capi.API{
		{Name: "aaa"},
		{Name: "bbb",
			Commands: []capi.Command{
				{Name: "cmd1", Path: "/one"},
				{Name: "cmd2", Path: "/two"},
			},
		},
	}}
	ac := autocomplete.Mock("app-name bbb c", "")
	options := capicomplete.GenerateResponse(&ac, &p)

	assert.Contains(t, options, "cmd1")
}
func Test_all_params_are_listed_with_double_dash_prefix(t *testing.T) {
	p := capi.APISet{APIs: []capi.API{
		{Name: "an_api",
			Commands: []capi.Command{
				{Name: "a_cmd", Path: "/one/{first_arg}", Header: map[string]string{"header1": "any", "header2": "any"}},
			},
		},
	}}

	t.Run("with no user filter", func(t *testing.T) {
		ac := autocomplete.Mock("any an_api a_cmd ", "")

		actual := capicomplete.GenerateResponse(&ac, &p)

		assert.Contains(t, actual, "--first_arg")
		assert.Contains(t, actual, "--header1")
		assert.Contains(t, actual, "--header2")
	})
	t.Run("with ambiguous user filter", func(t *testing.T) {
		ac := autocomplete.Mock("any an_api a_cmd --head", "")

		actual := capicomplete.GenerateResponse(&ac, &p)

		assert.Contains(t, actual, "--header1")
		assert.Contains(t, actual, "--header2")
	})
	t.Run("user filter that can only fit one result returns single flag", func(t *testing.T) {
		ac := autocomplete.Mock("any an_api a_cmd --first", "")

		actual := capicomplete.GenerateResponse(&ac, &p)

		assert.Contains(t, actual, "--first_arg")
	})
}

func Test_header_special_case(t *testing.T) {
	p := capi.APISet{APIs: []capi.API{{
		Name: "an_api",
		Commands: []capi.Command{{
			Name: "a_cmd",
			Path: "/",
			Header: map[string]string{
				"header1": "any",
				"header2": "{second_arg}",
			},
		}},
	}},
	}

	t.Run("when value in header contains arg, the header is not treated as a filter", func(t *testing.T) {
		ac := autocomplete.Mock("any an_api a_cmd ", "")

		actual := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, 2, len(actual))

		assert.NotContains(t, actual, "--header2")
		assert.Contains(t, actual, "--second_arg")
		assert.Contains(t, actual, "--header1")
	})
}
