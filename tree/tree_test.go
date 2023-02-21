package tree

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTestTree() Tree[int] {
	tree := NewTree[int]()
	tree.Add([]Matcher{
		NewFixedMatcher("a"),
		NewFixedMatcher("b"),
		NewFixedMatcher("c"),
	}, 1)
	tree.Add([]Matcher{
		NewFixedMatcher("a"),
		NewFixedMatcher("b"),
		NewFixedMatcher("c"),
		NewFixedMatcher("d"),
	}, 11)
	tree.Add([]Matcher{
		NewFixedMatcher("a"),
		NewFixedMatcher("a1"),
	}, 21)
	return tree
}

func TestTree_Add(t *testing.T) {
	tree := newTestTree()

	assert.Equal(t, 1, tree.root.MatchChild("a").MatchChild("b").MatchChild("c").GetValue())
	assert.Equal(t, 11, tree.root.MatchChild("a").MatchChild("b").MatchChild("c").MatchChild("d").GetValue())
	assert.Equal(t, 21, tree.root.MatchChild("a").MatchChild("a1").GetValue())
}

func TestTree_Find(t *testing.T) {
	tree := newTestTree()

	i, ok := tree.Find([]string{"a", "b", "c"})
	require.True(t, ok)
	assert.Equal(t, 1, i)

	i, ok = tree.Find([]string{"a", "b", "c", "d"})
	require.True(t, ok)
	assert.Equal(t, 11, i)

	i, ok = tree.Find([]string{"a", "a1"})
	require.True(t, ok)
	assert.Equal(t, 21, i)

	_, ok = tree.Find([]string{"a", "a2"})
	require.False(t, ok)
}

func TestTree_Weight(t *testing.T) {
	tree := NewTree[int]()

	tree.Add([]Matcher{
		NewFixedMatcher("users"),
		mustNewRegexpMatcher(`\d+`),
	}, 1)
	tree.Add([]Matcher{
		NewFixedMatcher("users"),
		NewFixedMatcher("list"),
	}, 2)

	assert.Equal(t, 1, tree.root.children[0].node.children[0].node.value)
	assert.Equal(t, 2, tree.root.children[0].node.children[1].node.value)

	tree.Weight()

	assert.Equal(t, 2, tree.root.children[0].node.children[0].node.value)
	assert.Equal(t, 1, tree.root.children[0].node.children[1].node.value)
}

func mustNewRegexpMatcher(regex string) RegexpMatcher {
	m, _ := NewRegexpMatcher(regex)
	return m
}
