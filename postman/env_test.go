package postman_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NearlyUnique/capi/postman"
	"github.com/stretchr/testify/require"
)

func Test_load_postman_env_file_based_on_env_var(t *testing.T) {
	var actualFilename string

	opener := func(filename string) (reader io.Reader, e error) {
		actualFilename = filename
		return bytes.NewReader([]byte(`{
	"values": [
		{ "key": "a_key", "value": "a_value" }
	]
}`)), nil
	}
	t.Run("env file name only can be used", func(t *testing.T) {
		postie := postman.NewSource(
			[]string{"any"},
			[]string{"CAPI_ENV=a-filename"},
			opener,
			nil)

		assert.Equal(t, "a-filename"+postman.EnvironmentExtension, actualFilename)
		require.NotNil(t, postie)
		assert.Equal(t, "a_value", postie("a_key"))
	})
	t.Run("full env file name can be used", func(t *testing.T) {
		postie := postman.NewSource(
			[]string{"any"},
			[]string{"CAPI_ENV=a-filename" + postman.EnvironmentExtension},
			opener,
			nil)

		assert.Equal(t, "a-filename"+postman.EnvironmentExtension, actualFilename)
		require.NotNil(t, postie)
		assert.Equal(t, "a_value", postie("a_key"))
	})
}
