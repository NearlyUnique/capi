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

func Test_can_detect_logical_position_on_line(t *testing.T) {
	line := "some-app an-api-name a-command-name"
	//       0123456789 123456789 123456789 1234
	//                 10        20        30
	testData := []struct {
		name     string
		line     string
		pos      int
		expected int
	}{
		{"impossible low", line, 7, 0},
		{"impossible high", line, len(line) + 1, -1},
		{"before", line, 8, 0},
		{"start api", line, 9, 1},
		{"mid api", line, 13, 1},
		{"end api", line, 19, 1},
		{"between api and command", line, 20, 1},
		{"on command", line, 21, 2},
	}

	for _, td := range testData {
		pos := logicalPosition(td.line, td.pos)

		assert.Equal(t, td.expected, pos, td.name)
	}
}
