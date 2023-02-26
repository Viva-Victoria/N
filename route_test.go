package n

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNRoute_Handler(t *testing.T) {
	var expectedError = errors.New("mock")
	r := NewRoute(HandlerFunc(func(ctx Context) error {
		return expectedError
	}))
	assert.Error(t, expectedError, r.Handler().Handle(nil))
}

func TestNRoute_Methods(t *testing.T) {
	r := NewRoute(HandlerFunc(func(ctx Context) error {
		return nil
	}))

	recorder := httptest.NewRecorder()
	assert.NoError(t, r.Handler().Handle(NewContext(nil, httptest.NewRequest(http.MethodGet, "/base", nil), NewResponseWriter(recorder))))
	assert.Equal(t, http.StatusOK, recorder.Code)

	r.Methods(http.MethodGet)
	recorder = httptest.NewRecorder()
	assert.NoError(t, r.Handler().Handle(NewContext(nil, httptest.NewRequest(http.MethodPost, "/base", nil), NewResponseWriter(recorder))))
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestNRoute_Use(t *testing.T) {
	var (
		originCalled bool
		firstCalled  bool
		secondCalled bool
	)

	r := NewRoute(HandlerFunc(func(ctx Context) error {
		originCalled = true
		return nil
	}))
	r.Use(MiddlewareFunc(func(ctx Context, handler Handler) Handler {
		firstCalled = true
		return handler
	}))
	r.Use(MiddlewareFunc(func(ctx Context, handler Handler) Handler {
		secondCalled = true
		return handler
	}))

	assert.NoError(t, r.Handler().Handle(nil))
	assert.True(t, originCalled)
	assert.True(t, firstCalled)
	assert.True(t, secondCalled)
}
