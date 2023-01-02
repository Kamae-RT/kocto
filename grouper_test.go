package kocto_test

import (
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

	if len(tree.Nodes) != 2 {
		t.Log("should have 2 A nodes")
		t.FailNow()
	} else if tree.Nodes[0].Group != "A" || tree.Nodes[0].Key != 1 {
		t.Fail()
	}

	if len(tree.Nodes[0].Nodes) != 2 {
		t.Log("A: 1 node should and 2 subnodes")
		t.FailNow()
	}
	if tree.Nodes[0].Nodes[0].Group != "B" || tree.Nodes[0].Nodes[0].Key != "1" ||
		tree.Nodes[0].Nodes[1].Group != "B" || tree.Nodes[0].Nodes[1].Key != "2" {

		t.Log("incorrect subnodes")
		t.Fail()
	}
}
