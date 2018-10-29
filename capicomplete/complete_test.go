package capicomplete_test

import (
	"testing"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/capicomplete"
	"github.com/stretchr/testify/assert"
)

func Test_when_nothing_is_entered_the_full_list_of_apis_is_displayed_with_alias(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
		{Name: "one", Alias: "1"},
		{Name: "two", Alias: "xx"},
	}}
	ac := autocomplete.Params{}

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one (1)", "two (xx)"}, act)
}
func Test_when_api_has_no_alias_it_is_ommited_from_that_response_item(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
		{Name: "one", Alias: "1"},
		{Name: "two"},
	}}
	ac := autocomplete.Params{}

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one (1)", "two"}, act)
}
func Test_when_only_one_item_returned_alias_is_not_included(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
		{Name: "one", Alias: "1"},
	}}
	ac := autocomplete.Params{}

	act := capicomplete.GenerateResponse(&ac, &p)

	assert.Equal(t, []string{"one"}, act)
}

func Test_when_partial_api_name_is_given_leading_submatches_are_returned(t *testing.T) {
	p := capi.Profile{APIs: []capi.API{
		{Name: "aaa", Alias: "xxx-alias"},
		{Name: "aaa-other"},
		{Name: "xxx-other"},
	}}

	t.Run("match with name", func(t *testing.T) {
		ac := autocomplete.Params{Word: "aa"}
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa (xxx-alias)", "aaa-other"}, act)
	})
	t.Run("match with alias", func(t *testing.T) {
		ac := autocomplete.Params{Word: "xx"}
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa (xxx-alias)", "xxx-other"}, act)
	})
	t.Run("when_no_matches_return_all", func(t *testing.T) {
		ac := autocomplete.Params{Word: "no-match"}
		act := capicomplete.GenerateResponse(&ac, &p)

		assert.Equal(t, []string{"aaa (xxx-alias)", "aaa-other", "xxx-other"}, act)
	})
}
