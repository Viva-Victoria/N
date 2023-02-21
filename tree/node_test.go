package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNode_GetValue(t *testing.T) {
	var node Node[int]
	node.value = 15
	assert.Equal(t, 15, node.GetValue())
}

func TestNode_SetValue(t *testing.T) {
	var node Node[int]
	node.SetValue(15)
	assert.Equal(t, 15, node.value)
}

func TestNode_AddChild(t *testing.T) {
	var node Node[int]
	node.AddChild(NewFixedMatcher("a"))
	assert.Len(t, node.children, 1)
}

func TestNode_FindChild(t *testing.T) {
	var node Node[int]
	a := node.AddChild(NewFixedMatcher("a"))
	a.SetValue(10)
	a.AddChild(NewFixedMatcher("a1")).SetValue(11)

	b := node.AddChild(NewFixedMatcher("b"))
	b.SetValue(20)
	b.AddChild(NewFixedMatcher("b1")).SetValue(21)

	assert.Equal(t, 10, node.FindChild(NewFixedMatcher("a")).GetValue())
	assert.Equal(t, 20, node.FindChild(NewFixedMatcher("b")).GetValue())
	assert.Nil(t, node.FindChild(NewFixedMatcher("a1")))
}

func TestNode_MatchChild(t *testing.T) {
	m, _ := NewRegexpMatcher(`^\d+$`)

	var node Node[int]
	a := node.AddChild(m)
	a.SetValue(1)

	assert.Equal(t, 1, node.MatchChild("15").GetValue())
	assert.Nil(t, node.MatchChild("a15"))
}
