package n

import (
	"gitea.voopsen/n/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRouter(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var (
			router = NewRouter("", log.LoggerMock{})
			called = false
		)

		router.Handle("/users/all", HandlerFunc(func(ctx Context) error {
			called = true

			require.NotNil(t, ctx.Request())
			require.NotNil(t, ctx.Response())

			ctx.Response().WriteHeader(http.StatusAccepted)
			return nil
		}))

		rs := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, "/users/all", nil)
		router.ServeHTTP(rs, rq)

		time.Sleep(time.Millisecond * 3)

		require.True(t, called)
		require.Equal(t, http.StatusAccepted, rs.Code)
	})
	t.Run("not-found", func(t *testing.T) {
		var (
			router = NewRouter("", log.LoggerMock{})
		)

		rs := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, "/users/all", nil)
		router.ServeHTTP(rs, rq)

		time.Sleep(time.Millisecond * 3)

		require.Equal(t, http.StatusNotFound, rs.Code)
	})
	t.Run("bad-route", func(t *testing.T) {
		var (
			actualError error
			router      = NewRouter("", log.LoggerMock{
				Error: func(err error) {
					actualError = err
				},
			})
		)

		route := router.Handle("/users/  /", HandlerFunc(func(ctx Context) error {
			return nil
		}))
		require.Nil(t, route)
		require.Error(t, actualError)

		rs := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, "/users//", nil)
		router.ServeHTTP(rs, rq)

		time.Sleep(time.Millisecond * 5)

		require.Equal(t, http.StatusNotFound, rs.Code)
	})
}

func TestNRouter_ServeHTTP(t *testing.T) {
	var (
		r      = NewRouter("", log.LoggerMock{})
		called = false
	)

	r.Get("/user/{id:\\d+}", HandlerFunc(func(ctx Context) error {
		called = true

		vars := ctx.Vars()
		require.NotEmpty(t, vars)

		var id int
		err := vars.Get("id", &id)
		require.NoError(t, err)
		assert.Equal(t, 15, id)
		return nil
	}))

	rs := httptest.NewRecorder()
	rq := httptest.NewRequest(http.MethodPost, "/user/15", nil)
	r.ServeHTTP(rs, rq)

	time.Sleep(time.Millisecond * 4)

	require.True(t, called)
}
