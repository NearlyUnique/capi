package builder_test

import (
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
)

func Test_default_flag_source_will_find_no_values(t *testing.T) {
	var actualErr error

	fn := builder.NewFlagSource(nil, func(e error) { actualErr = e })

	assert.NoError(t, actualErr)
	assert.Equal(t, "", fn("any"))
}

func Test_when_single_command_is_found(t *testing.T) {
	var actualErr error
	collectError := func(e error) { actualErr = e }

	t.Run("api base path is included", func(t *testing.T) {
		fn := builder.NewFlagSource([]string{"api", "cmd", "--anykey=value", "--other", "thing"}, collectError)

		assert.NoError(t, actualErr)
		assert.Equal(t, "value", fn("anykey"))
		assert.Equal(t, "thing", fn("other"))
	})
}
