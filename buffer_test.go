package n

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_init(t *testing.T) {
	require.NotNil(t, _bufferPool)

	a := _bufferPool.Get()
	require.NotNil(t, a)

	buf := a.(*bytes.Buffer)
	require.Equal(t, BufferSize, buf.Cap())
}

func Test_getBuffer(t *testing.T) {
	buf := getBuffer()
	require.NotNil(t, buf)
	assert.Equal(t, BufferSize, buf.Cap())
}

func Test_putBuffer(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		putBuffer(nil)
	})
	t.Run("notNil", func(t *testing.T) {
		putBuffer(bytes.NewBuffer(make([]byte, BufferSize)))
	})
}
