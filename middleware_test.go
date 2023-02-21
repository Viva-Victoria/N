package n

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMiddlewareFunc_Handle(t *testing.T) {
	assert.NotPanics(t, func() {
		var middlewareFund MiddlewareFunc
		middlewareFund.Handle(nil, nil)
	})
	assert.NotPanics(t, func() {
		var called bool
		middlewareFunc := MiddlewareFunc(func(ctx Context, handler Handler) Handler {
			called = true
			return nil
		})
		middlewareFunc.Handle(nil, nil)
		assert.True(t, called)
	})
}
