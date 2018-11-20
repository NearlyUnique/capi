package capi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_urls_with_query_string_are_correctly_encoded(t *testing.T) {
	actual, err := makeURL("http://example.com", "/any?arg=value")

	require.NoError(t, err)
	assert.Equal(t, "http://example.com/any?arg=value", actual)
}
