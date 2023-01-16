package kocto

import (
	"testing"

	"github.com/matryer/is"
)

func TestBatcher(t *testing.T) {
	b := NewBatcher[int](2)

	is := is.NewRelaxed(t)

	is.True(!b.Add(1)) // batch should not be full
	is.True(b.Add(2))  // batch should be full

	b.Flush(func(t []int) {})

	is.True(!b.Add(3)) // batch shoud be able to full again
	is.True(b.Add(4)) // batch should be full again
}
