package kocto_test

import (
	"encoding/json"
	"testing"

	"kamaesoft.visualstudio.com/kocto/_git/kocto"
)

type G map[string]any

func (g G) Get(field string) any {
	return g[field]
}

var groupTestData = []G{
	{"A": 1, "B": "1", "C": 1, "D": 1},
	{"A": 1, "B": "1", "C": 2, "D": 1},
	{"A": 1, "B": "1", "C": 3, "D": 1},
	{"A": 1, "B": "2", "C": 1, "D": 1},
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
	groups := []string{"A", "B", "D", "C"}
	data := make([]kocto.Indexable, len(groupTestData))
	for i := range groupTestData {
		data[i] = groupTestData[i]
	}

	tree := kocto.Group(groups, data)

	j, err := json.Marshal(tree)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(string(j))
	t.Fail()
}
