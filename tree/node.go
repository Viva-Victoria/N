package tree

import "sort"

type NodeChild[V any] struct {
	matcher Matcher
	node    Node[V]
}

type Node[V any] struct {
	children []NodeChild[V]
	value    V
}

func (n *Node[V]) Sort(less func(a, b Matcher) bool) {
	sort.Slice(n.children, func(i, j int) bool {
		return less(n.children[i].matcher, n.children[j].matcher)
	})
}

func (n *Node[V]) FindChild(r Matcher) *Node[V] {
	for i, m := range n.children {
		if m.matcher.String() == r.String() {
			return &n.children[i].node
		}
	}

	return nil
}

func (n *Node[V]) MatchChild(s string) *Node[V] {
	for i, m := range n.children {
		if m.matcher.MatchString(s) {
			return &n.children[i].node
		}
	}

	return nil
}

func (n *Node[V]) AddChild(m Matcher) *Node[V] {
	n.children = append(n.children, NodeChild[V]{
		matcher: m,
		node:    Node[V]{},
	})

	return &n.children[len(n.children)-1].node
}

func (n *Node[V]) GetValue() V {
	return n.value
}

func (n *Node[V]) SetValue(v V) {
	n.value = v
}
