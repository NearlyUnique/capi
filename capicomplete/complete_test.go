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
	ac := autocomplete.Params{}

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
		ac := autocomplete.Params{Word: "aa"}
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa", "aaa-other"}, act)
	})
	t.Run("when_no_matches_return_all", func(t *testing.T) {
		ac := autocomplete.Params{Word: "no-match"}
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa", "aaa-other", "xxx-other"}, act)
	})
}
