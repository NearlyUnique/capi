package autocomplete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_creating_mocked_complete_params(t *testing.T) {
	testData := []struct {
		before, after, name string
		expected            Params
	}{
		//0123456789 123456789
		{"acmd f", "", "single char", Params{Line: "acmd f", Point: 6, Word: "f", PrevWord: "acmd", Command: "acmd"}},
		{"acmd any ", "", "trailing char", Params{Line: "acmd any ", Point: 9, Word: "", PrevWord: "any", Command: "acmd"}},
		{"acmd any  ", "", "multiple trailing char", Params{Line: "acmd any  ", Point: 10, Word: "", PrevWord: "any", Command: "acmd"}},
	}
	for _, td := range testData {
		actual := Mock(td.before, td.after)

		//force all expected params to be type 9, key 9 as per MockComplete spec
		td.expected.Key = "9"
		td.expected.Type = TypeNormal

		assert.Equal(t, td.expected, actual, td.name)
	}
}
