package kocto

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestHeaderString(t *testing.T) {
	is := is.New(t)

	t.Run("simple headers", func(t *testing.T) {
		is := is.New(t)

		headers := map[string][]string{
			"Header1": {"Value1"},
			"Header2": {"Value2"},
		}

		// might change based on map order
		expected1 := "Header1: Value1, Header2: Value2"
		expected2 := "Header2: Value1, Header1: Value1"
		str := headersString(headers)

		is.True(str == expected1 || str == expected2)
	})

	t.Run("multiple values in headers", func(t *testing.T) {
		is := is.New(t)

		headers := map[string][]string{
			"Header1": {"Value11", "Value12"},
			"Header2": {"Value2"},
			"Header3": {"Value31", "Value32"},
		}

		str := headersString(headers)

		expectedStrs := []string{
			"Header1: Value11",
			"Header1: Value12",
			"Header2: Value2",
			"Header3: Value31",
			"Header3: Value32",
		}

		for _, estr := range expectedStrs {
			is.True(strings.Contains(str, estr))
		}
	})
}
