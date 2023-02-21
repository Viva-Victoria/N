package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixedMatcher(t *testing.T) {
	f := NewFixedMatcher("test")
	assert.Equal(t, "test", f.String())
	assert.True(t, f.MatchString("test"))
	assert.False(t, f.MatchString("toast"))
}

func TestRegexpMatcher(t *testing.T) {
	r, _ := NewRegexpMatcher(`\d+`)
	assert.Equal(t, `\d+`, r.String())
	assert.True(t, r.MatchString("15"))
	assert.False(t, r.MatchString("aF"))
}

func TestRegexpMatcher_Fail(t *testing.T) {
	_, err := NewRegexpMatcher(`?!`)
	assert.NotNil(t, err)
}
