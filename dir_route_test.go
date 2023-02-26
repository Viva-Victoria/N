package n

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_pathFixer(t *testing.T) {
	require.NotNil(t, _pathFixer)
	assert.Equal(t, "/base/path", _pathFixer.Replace("/base/path"))
	assert.Equal(t, "/base/path", _pathFixer.Replace("\\base\\path"))
}

func TestNewDirRoute(t *testing.T) {
	var (
		fullPath string
	)

	r := NewDirRoute("/base/path/", func(path string, _ Handler) Route {
		fullPath = path
		return nil
	})

	r.Handle("/route", HandlerFunc(func(ctx Context) error {
		return nil
	}))

	assert.Equal(t, "/base/path/route", fullPath)
}

func TestNDirRoute_Handle(t *testing.T) {
	var (
		actualHandler Handler
		expectedRoute Route
		expectedError = errors.New("mock")
	)

	r := NewDirRoute("/base/path/", func(path string, h Handler) Route {
		assert.Equal(t, "/base/path/a", path)
		actualHandler = h
		expectedRoute = NewRoute(h)
		return expectedRoute
	})
	route := r.Handle("/a", HandlerFunc(func(ctx Context) error {
		return expectedError
	}))

	require.NotNil(t, route)
	assert.Equal(t, expectedRoute, route)

	require.NotNil(t, actualHandler)
	assert.Error(t, expectedError, actualHandler.Handle(nil))
}

func TestNDirRoute_Dir(t *testing.T) {
	var (
		actualHandler Handler
		expectedRoute Route
		expectedError = errors.New("mock")
	)

	r := NewDirRoute("/base/path/", func(path string, h Handler) Route {
		assert.Equal(t, "/base/path/sub/route", path)
		actualHandler = h
		expectedRoute = NewRoute(h)
		return expectedRoute
	})
	subR := r.Dir("/sub")
	require.NotNil(t, subR)

	actualRoute := subR.Handle("/route", HandlerFunc(func(ctx Context) error {
		return expectedError
	}))
	require.NotNil(t, actualRoute)
	assert.Equal(t, expectedRoute, actualRoute)

	require.NotNil(t, actualHandler)
	assert.Error(t, expectedError, actualHandler.Handle(nil))
}

func TestNDirRoute_Method(t *testing.T) {
	t.Parallel()
	t.Run("get", func(t *testing.T) {
		testHandler(t, func(d *NDirRoute) func(string, Handler) Route {
			return d.Get
		}, http.MethodGet)
	})
	t.Run("post", func(t *testing.T) {
		testHandler(t, func(d *NDirRoute) func(string, Handler) Route {
			return d.Post
		}, http.MethodPost)
	})
	t.Run("put", func(t *testing.T) {
		testHandler(t, func(d *NDirRoute) func(string, Handler) Route {
			return d.Put
		}, http.MethodPut)
	})
	t.Run("patch", func(t *testing.T) {
		testHandler(t, func(d *NDirRoute) func(string, Handler) Route {
			return d.Patch
		}, http.MethodPatch)
	})
	t.Run("delete", func(t *testing.T) {
		testHandler(t, func(d *NDirRoute) func(string, Handler) Route {
			return d.Delete
		}, http.MethodDelete)
	})
}

func testHandler(t *testing.T, createAddRouteFunc func(d *NDirRoute) func(string, Handler) Route, method string) {
	r := NewDirRoute("/base", func(path string, h Handler) Route {
		assert.Equal(t, "/base/route", path)
		return NewRoute(h)
	})
	route := createAddRouteFunc(r)("route", HandlerFunc(func(ctx Context) error {
		return nil
	}))

	testHandlerMethod(t, route, method, http.StatusOK)

	var secondMethod string
	for _, httpMethod := range []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete} {
		if httpMethod == method {
			continue
		}
		if len(secondMethod) == 0 {
			secondMethod = method
		}

		testHandlerMethod(t, route, httpMethod, http.StatusMethodNotAllowed)
	}

	route.Methods(secondMethod)
	testHandlerMethod(t, route, method, http.StatusOK)
	testHandlerMethod(t, route, secondMethod, http.StatusOK)
}

func testHandlerMethod(t *testing.T, route Route, method string, status int) {
	recorder := httptest.NewRecorder()
	err := route.Handler().Handle(NewContext(make(MapValues), httptest.NewRequest(method, "/base/route", nil), NewResponseWriter(recorder)))
	require.NoError(t, err)
	assert.Equal(t, status, recorder.Code)
}
