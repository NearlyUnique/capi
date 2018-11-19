package capi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_urls_with_query_string_are_correctly_encoded(t *testing.T) {
	actual := makeURL("http://example.com", "/any?arg=value")

	assert.Equal(t, "http://example.com/any?arg=value", actual)
}
