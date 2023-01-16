package kocto_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"kamaesoft.visualstudio.com/kocto/_git/kocto"
)

type G map[string]any

func (g G) Get(field string) any {
	return g[field]
}

var groupTestData = []G{
	{"A": 1, "B": "1", "C": 1},
	{"A": 1, "B": "1", "C": 2},
	{"A": 1, "B": "1", "C": 3},
	{"A": 1, "B": "2", "C": 1},
	{"A": 1, "B": "2", "C": 2},
	{"A": 1, "B": "2", "C": 3},
	{"A": 2, "B": "1", "C": 1},
	{"A": 2, "B": "1", "C": 2},
	{"A": 2, "B": "1", "C": 3},
	{"A": 2, "B": "2", "C": 1},
	{"A": 2, "B": "2", "C": 2},
	{"A": 2, "B": "2", "C": 3},
}

func TestGrouper(t *testing.T) {
	tests := []struct {
		groups    []string
		groupSize int
	}{
		{[]string{"A"}, 2},
		{[]string{"B"}, 2},
		{[]string{"A", "B"}, 4},
		{[]string{"B", "A"}, 4},
		{[]string{"A", "B", "C"}, 12},
		{[]string{"A", "C", "B"}, 12},
		{[]string{"C", "B", "A"}, 12},
		{[]string{"C", "A", "B"}, 12},
	}

	is := is.NewRelaxed(t)

	data := make([]kocto.Indexable, len(groupTestData))
	for i := range groupTestData {
		data[i] = groupTestData[i]
	}

	for _, tt := range tests {
		grouppedData := kocto.Group(tt.groups, data)

		// check if the number of groups is correct
		if len(grouppedData) != tt.groupSize {
			is.Equal(tt.groupSize, len(grouppedData)) // size of groups should be equal
		}

		// check if any child is misplaced
		for _, gd := range grouppedData {
			for _, d := range gd.Data {
				for _, g := range tt.groups {
					is.True(strings.Contains(gd.Key, fmt.Sprint(d.Get(g)))) // group is incorrect
				}
			}
		}
	}
}
