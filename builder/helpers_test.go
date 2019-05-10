package builder_test

import (
	"strings"
	"testing"

	"github.com/NearlyUnique/capi/builder"
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
