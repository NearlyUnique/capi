package autocomplete

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_slice_to_map(t *testing.T) {
	m := sliceToMap([]string{
		"some_value=simple",
		"another=with=equal",
	})

	assert.Equal(t, m["some_value"], "simple")
	assert.Equal(t, m["another"], "with=equal")
}

func Test_use_real_env(t *testing.T) {
	os.Setenv("some_real_env_var", "some_real_env_value")
	m := sliceToMap(os.Environ())

	assert.Equal(t, m["some_real_env_var"], "some_real_env_value")
}
