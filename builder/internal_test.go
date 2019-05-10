package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_urls_are_joined_regardless_of_edge_slash_placement(t *testing.T) {
	testData := []struct {
		base, path, expected string
	}{
		{"http://example.com/", "", "http://example.com/"},
		{"http://example.com", "", "http://example.com"},
		{"", "/path", "/path"},
		{"http://example.com/", "/", "http://example.com/"},
		{"http://example.com/", "/path", "http://example.com/path"},
		{"http://example.com", "/path", "http://example.com/path"},
		{"http://example.com", "path", "http://example.com/path"},
	}
	for _, data := range testData {
		actual := joinUrlFragments(data.base, data.path)

		assert.Equal(t, data.expected, actual)
	}
}
