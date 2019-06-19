package run_test

import (
	"strconv"
	"testing"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/run"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

type (
	FakeLoader struct {
		list func(search string) []string
		load func(filename string) (*builder.APISet, error)
	}
)

func (f FakeLoader) List(search string) []string {
	return f.list(search)
}

func (f FakeLoader) Load(filename string) (*builder.APISet, error) {
	return f.load(filename)
}

func Test_tab_complete(t *testing.T) {
	t.Run("when no API file can be found, return a list for each registered extension", func(t *testing.T) {
		var actualSearch string
		loader := FakeLoader{
			list: func(search string) []string {
				actualSearch = search
				return []string{"any-item"}
			},
			load: func(filename string) (*builder.APISet, error) {
				return nil, xerrors.New("any failure")
			},
		}
		params := autocomplete.Mock("capi partial-api", "")
		args := []string{
			"capi", //the name of the complete'er
			params.Command,
			params.Word,
			params.PrevWord,
		}

		options := run.AutoComplete(loader, args, []string{
			"COMP_LINE=" + params.Line,
			"COMP_POINT=" + strconv.Itoa(params.Point),
			"COMP_KEY=" + params.Key,
			"COMP_TYPE=" + strconv.Itoa(int(params.Type)),
		})

		require.Equal(t, 1, len(options))
		assert.Equal(t, "any-item", options[0])
		assert.Equal(t, "partial-api", actualSearch)
	})
}
