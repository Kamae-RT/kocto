package kocto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatcher(t *testing.T) {
	b := NewBatcher[int](2)

	assert.False(t, b.Add(1))
	assert.True(t, b.Add(2))

	prev := make([]int, len(b.data))
	copy(prev, b.data)

	assert.ElementsMatch(t, prev, b.data)

	b.Flush(func(t []int) {})
	assert.False(t, b.Add(3))
	assert.True(t, b.Add(4))

	assert.NotContains(t, b.data, prev)
}
