package autocomplete_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_auto_complete_data_is_generated(t *testing.T) {
	args := []string{"complete-app", "any-app", "", "any-app"}
	env := []string{
		"OTHER_ENV=some value",
		"COMP_LINE=any-app",
		"COMP_TYPE=9",
		"COMP_POINT=8",
		"COMP_KEY=9",
	}
	ExpAc := autocomplete.Params{
		PrevWord: "any-app",
		Word:     "",
		Command:  "any-app",
		Key:      "9",
		Point:    8,
		Line:     "any-app",
		Type:     autocomplete.TypeNormal,
	}
	t.Run("must have 3 string arguments to cli", func(t *testing.T) {
		// zeroth argument is ALWAYS the name of the running program, e.g. the complete app
		assert.Nil(t, autocomplete.Prepare([]string{}, env))
		assert.Nil(t, autocomplete.Prepare([]string{"one"}, env))
		assert.Nil(t, autocomplete.Prepare([]string{"one", "two"}, env))
		assert.Nil(t, autocomplete.Prepare([]string{"one", "two", "three"}, env))
		assert.NotNil(t, autocomplete.Prepare([]string{"one", "two", "three", "four"}, env))
		assert.Nil(t, autocomplete.Prepare([]string{"one", "two", "three", "four", "five"}, env))
	})
	t.Run("must have required env var", func(t *testing.T) {

		autocomplete.LogHook = func(format string, args ...interface{}) {
			//t.Log(fmt.Sprintf(format, args...))
			//t.FailNow()
		}
		assert.NotNil(t, autocomplete.Prepare(args, env))

		remove := []string{"COMP_LINE", "COMP_TYPE", "COMP_KEY", "COMP_POINT"}
		for _, rem := range remove {
			assert.Nil(t, autocomplete.Prepare(args, copyWithoutKey(env, rem)), fmt.Sprintf("%s was removed", rem))
		}
	})
	t.Run("when all components are available", func(t *testing.T) {
		ac := autocomplete.Prepare(args, env)

		require.NotNil(t, ac)
		assert.True(t, reflect.DeepEqual(&ExpAc, ac), "expected:%v\nactual: %v\n", ExpAc, ac)
	})
}

func copyWithoutKey(src []string, key string) []string {
	var cpy []string
	for _, v := range src {
		if strings.HasPrefix(v, key+"=") {
			continue
		}
		cpy = append(cpy, v)
	}
	return cpy
}
