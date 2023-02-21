package n

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerFunc_Handle(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var handlerFunc HandlerFunc
			_ = handlerFunc.Handle(nil)
		})
	})
	t.Run("not-nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var called bool
			handlerFunc := HandlerFunc(func(ctx Context) error {
				called = true
				return nil
			})

			assert.NoError(t, handlerFunc.Handle(nil))
			assert.True(t, called)
		})
	})
}
