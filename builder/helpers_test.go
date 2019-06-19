package builder_test

import (
	"strings"
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/require"
)

func fakeSource(k, s string) builder.SourceFn {
	return func(key string) string {
		if key == k {
			return s
		}
		return ""
	}
}

func assertStartsWith(t *testing.T, expected, actual string) {
	t.Helper()
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("'%s' does not start with '%s'", actual, expected)
	}
}

func assertEndsWith(t *testing.T, expected, actual string) {
	t.Helper()
	if !strings.HasSuffix(actual, expected) {
		t.Errorf("'%s' does not end with '%s'", actual, expected)
	}
}

func firstCmd(t *testing.T, set builder.APISet) *builder.Command {
	require.True(t, len(set.APIs) > 0)
	api, err := set.FindAPI(set.APIs[0].Name)
	require.NoError(t, err)
	require.True(t, len(api[0].Commands) > 0)
	cmd, err := api[0].FindCommand(api[0].Commands[0].Name)
	require.NoError(t, err)
	return cmd[0]
}
