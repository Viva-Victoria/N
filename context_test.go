package n

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNContext(t *testing.T) {
	t.Run("request", func(t *testing.T) {
		expectedRQ := httptest.NewRequest(http.MethodGet, "/path", nil)
		actualRQ := NewContext(nil, expectedRQ, nil).Request()
		assert.Equal(t, expectedRQ, actualRQ)
	})
	t.Run("response", func(t *testing.T) {
		expectedRS := NewResponseWriter(httptest.NewRecorder())
		expectedRS.WriteHeader(http.StatusNotFound)
		actualRS := NewContext(nil, nil, expectedRS).Response()
		assert.Equal(t, expectedRS, actualRS)
	})
	t.Run("vars", func(t *testing.T) {
		expectedVars := make(MapValues)
		expectedVars.Set("test", []string{"1", "2", "3"})
		actualVars := NewContext(expectedVars, nil, nil).Vars()
		assert.Equal(t, expectedVars, actualVars)
	})
}

func TestNContext_Header(t *testing.T) {
	rq := httptest.NewRequest(http.MethodPost, "/path", nil)
	rq.Header = http.Header{
		"X-Data": []string{
			"Test",
			"A",
		},
	}

	ctx := NewContext(nil, rq, nil)
	header := ctx.Header()
	require.NotNil(t, header)
	assert.Equal(t, "Test", header.GetString("X-Data"))

	var array []string
	err := header.Get("X-Data", &array)
	require.NoError(t, err)
	assert.EqualValues(t, []string{"Test", "A"}, array)
}
