package array_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/array"
	"github.com/stretchr/testify/assert"
)

func TestStringUnique(t *testing.T) {
	arr := []string{
		"aaa",
		"bbb",
		"ccc",
		"aaa",
		"ddd",
		"ccc",
	}

	assert.EqualValues(t, 4, len(array.StringUnique(arr)))
}
