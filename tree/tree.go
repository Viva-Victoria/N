package tree

type Tree[V any] struct {
	root  Node[V]
	count int
}

func NewTree[V any]() Tree[V] {
	return Tree[V]{}
}

func (t *Tree[V]) Add(path []Matcher, v V) {
	currentNode := &t.root
	for _, matcher := range path {
		if existLeaf := currentNode.FindChild(matcher); existLeaf != nil {
			currentNode = existLeaf
			continue
		}

		currentNode = currentNode.AddChild(matcher)
	}

	currentNode.SetValue(v)
	t.count++
}

func (t *Tree[V]) Get(path []Matcher) (V, bool) {
	var v V

	currentNode := &t.root
	for _, matcher := range path {
		existLeaf := currentNode.FindChild(matcher)
		if existLeaf == nil {
			return v, false
		}

		currentNode = existLeaf
	}

	return currentNode.GetValue(), true
}

func (t *Tree[V]) Find(path []string) (V, bool) {
	var (
		v           V
		currentNode = &t.root
	)

	for _, pathSegment := range path {
		currentNode = currentNode.MatchChild(pathSegment)
		if currentNode == nil {
			return v, false
		}
	}

	return currentNode.GetValue(), true
}

func (t *Tree[V]) Weight() {
	nodes := make([]Node[V], 0, t.count)

	i := 0
	currentNode := &t.root
	for len(nodes) < t.count {
		for _, meta := range currentNode.children {
			nodes = append(nodes, meta.node)
		}

		currentNode = &currentNode.children[i].node
		i++
	}

	for _, node := range nodes {
		node.Sort(func(a, b Matcher) bool {
			_, ok := a.(FixedMatcher)
			return ok
		})
	}
}
