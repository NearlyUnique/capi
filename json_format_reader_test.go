package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_json_format_reader(t *testing.T) {
	t.Run("errors if content not json", func(t *testing.T) {
		content := []byte(`not json`)

		set, err := JSONFormatReader(content)

		require.Error(t, err)
		assert.Nil(t, set)
	})
	t.Run("can read well formed apiset json", func(t *testing.T) {
		content := []byte(`{"apis":[{"name":"first"}]}`)

		set, err := JSONFormatReader(content)

		require.NoError(t, err)
		require.Equal(t, 1, len(set.APIs))
		assert.Equal(t, "first", set.APIs[0].Name)
	})
}
