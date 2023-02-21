package n

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"testing"
)

func Test_HeaderValues(t *testing.T) {
	header := HeaderValues{
		header: make(http.Header),
	}

	header.Add("X-Test", "value1")
	header.Add("X-Test", "value2")
	assert.Equal(t, []string{"value1", "value2"}, header.header["X-Test"])
	assert.Equal(t, []string{"value1", "value2"}, header.Values("X-Test"))

	header.Set("X-Test", "value")
	assert.Equal(t, []string{"value"}, header.header["X-Test"])

	header.Set("X-Test-I", 63)
	assert.Equal(t, []string{"63"}, header.header["X-Test-I"])

	var s string
	err := header.Get("X-Test", &s)
	require.NoError(t, err)
	assert.Equal(t, []string{s}, header.header["X-Test"])

	var i int
	err = header.Get("X-Test-I", &i)
	require.NoError(t, err)
	assert.Equal(t, []string{strconv.Itoa(i)}, header.header["X-Test-I"])

	value := header.GetString("X-Test")
	assert.Equal(t, []string{value}, header.header["X-Test"])

	header.Delete("X-Test")
	header.Delete("X-Test-I")
	assert.Empty(t, header.header)
}

func Test_MapValues(t *testing.T) {
	header := MapValues(make(http.Header))

	header.Add("X-Test", "value1")
	header.Add("X-Test", "value2")
	assert.Equal(t, []string{"value1", "value2"}, header["X-Test"])
	assert.Equal(t, []string{"value1", "value2"}, header.Values("X-Test"))

	header.Set("X-Test", "value")
	assert.Equal(t, []string{"value"}, header["X-Test"])

	header.Set("X-Test-I", 63)
	assert.Equal(t, []string{"63"}, header["X-Test-I"])

	var s string
	err := header.Get("X-Test", &s)
	require.NoError(t, err)
	assert.Equal(t, []string{s}, header["X-Test"])

	var i int
	err = header.Get("X-Test-I", &i)
	require.NoError(t, err)
	assert.Equal(t, []string{strconv.Itoa(i)}, header["X-Test-I"])

	value := header.GetString("X-Test")
	assert.Equal(t, []string{value}, header["X-Test"])

	header.Delete("X-Test")
	header.Delete("X-Test-I")
	assert.Empty(t, header)
}
