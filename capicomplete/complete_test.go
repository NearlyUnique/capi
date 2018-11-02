package capicomplete_test

import (
	"testing"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/capicomplete"
	"github.com/stretchr/testify/assert"
)

func Test_when_nothing_is_entered_the_full_list_of_apis_is_displayed(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
		{Name: "one"},
		{Name: "two"},
	}}
	ac := autocomplete.Mock("any ", "")

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one", "two"}, act)
}

func Test_when_partial_api_name_is_given_leading_submatches_are_returned(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
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
	p := capi.Profile{APIs: []capi.API{
		{Name: "aaa"},
		{Name: "bbb",
			Commands: []capi.Command{
				{Name: "cmd1", Path: "/one"},
				{Name: "cmd2", Path: "/two"},
			},
		},
	}}
	ac := autocomplete.Params{
		Line: "app-name bbb c",
		//     0123456789 123
		Point:    13,
		Word:     "c",
		PrevWord: "bbb",
		Type:     9,
		Key:      "9",
	}
	options := capicomplete.GenerateResponse(&ac, &p)

	assert.Contains(t, options, "cmd1")
}
