package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_indexOrEmpty_safely_returns_a_slice_element_or_the_default_value(t *testing.T) {
	testData := []struct {
		source   []string
		index    int
		expected string
	}{
		{[]string{}, -1, ""},
		{[]string{}, 0, ""},
		{[]string{}, 1, ""},
		{[]string{"first"}, -1, ""},
		{[]string{"first"}, 0, "first"},
		{[]string{"first"}, 1, ""},
		{[]string{"first", "second", "third"}, -1, ""},
		{[]string{"first", "second", "third"}, 0, "first"},
		{[]string{"first", "second", "third"}, 1, "second"},
		{[]string{"first", "second", "third"}, 3, ""},
		{[]string{"first", "second", "third"}, 99, ""},
	}

	for _, td := range testData {
		actual := indexOrEmpty(td.source, td.index)
		assert.Equal(t, td.expected, actual)
	}
}
