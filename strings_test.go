package n

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_splitByComma(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, splitByComma(``))
	})
	t.Run("one", func(t *testing.T) {
		assert.EqualValues(t, []string{"a"}, splitByComma(`[a]`))
	})
	t.Run("two", func(t *testing.T) {
		assert.EqualValues(t, []string{"a", " b"}, splitByComma(`[a, b]`))
	})
	t.Run("embed", func(t *testing.T) {
		assert.EqualValues(t, []string{`[1, 2]`, ` [3, 4]`}, splitByComma(`[[1, 2], [3, 4]]`))
	})
}
