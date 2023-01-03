package kocto

import "fmt"

type Indexable interface {
	Get(field string) any
}

type GroupTree struct {
	groupNode `json:"data"`
}

type FlatGroup struct {
	Key  string
	Data []Indexable
}

func Group(groups []string, data []Indexable) []FlatGroup {
	tree := GroupTree{groupNode: newNode("", nil)}

	for _, d := range data {
		tree.Insert(groups, d)
	}

	return tree.groupNode.Flatten("")
}

func (t *GroupTree) Insert(groups []string, d Indexable) {
	t.groupNode.Insert(groups[0], groups[1:], d)
}

type groupNode struct {
	Group string `json:"field"`
	Key   any    `json:"value"`

	Nodes []groupNode `json:"nodes"`
	Data  []Indexable `json:"data"`
}

func newNode(group string, key any) groupNode {
	return groupNode{
		Group: group,
		Key:   key,
		Nodes: make([]groupNode, 0),
		Data:  make([]Indexable, 0),
	}
}

func (n *groupNode) Insert(group string, groups []string, d Indexable) {
	if group == "" {
		n.Data = append(n.Data, d)
		return
	}

	val := d.Get(group)

	nIdx := findNode(val, n.Nodes)
	if nIdx < 0 {
		n.Nodes = append(n.Nodes, newNode(group, val))
		nIdx = len(n.Nodes) - 1
	}

	nextGroup := ""
	if len(groups) > 0 {
		nextGroup = groups[0]
	}

	nextGroups := []string{}
	if len(groups) > 1 {
		nextGroups = groups[1:]
	}

	n.Nodes[nIdx].Insert(nextGroup, nextGroups, d)
}

func findNode(val any, nodes []groupNode) int {
	idx := -1
	for i, n := range nodes {
		if n.Key == val {
			return i
		}
	}

	return idx
}

func (n groupNode) Flatten(key string) []FlatGroup {
	groups := make([]FlatGroup, 0)

	if len(n.Nodes) > 0 {
		for _, c := range n.Nodes {
			k := gKey(key, c.Key)
			groups = append(groups, c.Flatten(k)...)
		}
	} else {
		g := FlatGroup{Key: key, Data: n.Data}
		groups = append(groups, g)
	}

	return groups
}

func gKey(s1 string, s2 any) string {
	s := ""
	if s2 != nil {
		s = fmt.Sprint(s2)
	}

	if s1 == "" {
		return s
	}

	return fmt.Sprint(s1, ",", s)
}
