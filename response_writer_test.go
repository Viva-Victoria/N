package n

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type writer struct {
	http.ResponseWriter
}

type writerFlusher struct {
	writer
	http.Flusher
	flushFn func()
}

func (w writerFlusher) Flush() {
	if w.flushFn == nil {
		return
	}
	w.flushFn()
}

type writerCloser struct {
	writer
	io.Closer
	closeFn func() error
}

func (w writerCloser) Close() error {
	if w.closeFn == nil {
		return nil
	}
	return w.closeFn()
}

type writerFlusherCloser struct {
	writer
	writerFlusher
	writerCloser
}

func TestNewResponseWriter(t *testing.T) {
	t.Run("writer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			w := writer{}
			rsw := NewResponseWriter(w)
			assert.Equal(t, w, rsw.rs)
			assert.Nil(t, rsw.flusher)
			assert.Nil(t, rsw.closer)
		})
	})
	t.Run("flusher", func(t *testing.T) {
		assert.NotPanics(t, func() {
			w := writerFlusher{}
			rsw := NewResponseWriter(w)
			assert.Equal(t, w, rsw.rs)
			assert.Equal(t, w, rsw.flusher)
			assert.Nil(t, rsw.closer)
		})
	})
	t.Run("closer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			w := writerCloser{}
			rsw := NewResponseWriter(w)
			assert.Equal(t, w, rsw.rs)
			assert.Nil(t, rsw.flusher)
			assert.Equal(t, w, rsw.closer)
		})
	})
	t.Run("flusher-closer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			w := writerFlusherCloser{}
			rsw := NewResponseWriter(w)
			assert.Equal(t, w, rsw.rs)
			assert.Equal(t, w, rsw.flusher)
			assert.Equal(t, w, rsw.closer)
		})
	})
}

type writerMock struct {
	header func() http.Header
	write  func([]byte) (int, error)
	status func(statusCode int)
}

func (w writerMock) Header() http.Header {
	if w.header == nil {
		return nil
	}

	return w.header()
}

func (w writerMock) Write(bytes []byte) (int, error) {
	if w.write == nil {
		return 0, nil
	}

	return w.write(bytes)
}

func (w writerMock) WriteHeader(statusCode int) {
	if w.status == nil {
		return
	}

	w.status(statusCode)
}

func TestNResponseWriter_Header(t *testing.T) {
	for name, header := range map[string]http.Header{
		"nil": nil,
		"single": {
			"X-Test": []string{
				"A", "B", "C",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			mock := writerMock{
				header: func() http.Header {
					return header
				},
			}

			rs := NewResponseWriter(mock)
			if header == nil {
				assert.Nil(t, rs.Header())
			}

			assert.EqualValues(t, header, rs.Header())
		})
	}
}

func TestNResponseWriter_Write(t *testing.T) {
	t.Run("not-committed", func(t *testing.T) {
		rsw := NewResponseWriter(writerMock{})

		data := make([]byte, BufferLimit)
		_, _ = rand.Read(data)

		n, err := rsw.Write(data)
		require.Equal(t, len(data), n)
		require.NoError(t, err)

		require.EqualValues(t, data, rsw.body.Bytes())
		require.False(t, rsw.committed)
	})
	t.Run("committed", func(t *testing.T) {
		var (
			count  int
			buffer bytes.Buffer
		)

		rsw := NewResponseWriter(writerMock{
			write: func(bytes []byte) (int, error) {
				count = len(bytes)
				return buffer.Write(bytes)
			},
		})

		data := make([]byte, BufferLimit+1)
		_, _ = rand.Read(data)

		n, err := rsw.Write(data)
		require.Equal(t, len(data), n)
		require.NoError(t, err)

		require.NotNil(t, count)
		assert.Equal(t, count, n)

		require.EqualValues(t, data, buffer.Bytes())
		require.True(t, rsw.committed)
	})
}

func TestNResponseWriter_WriteHeader(t *testing.T) {
	var status int
	rs := NewResponseWriter(writerMock{
		status: func(statusCode int) {
			status = statusCode
		},
	})
	rs.WriteHeader(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, status)
}

func TestNResponseWriter_Flush(t *testing.T) {
	t.Run("writeBody", func(t *testing.T) {
		var count int
		rsw := NewResponseWriter(writerMock{
			write: func(i []byte) (int, error) {
				count = len(i)
				return count, nil
			},
		})

		data := make([]byte, 128)
		_, _ = rand.Read(data)
		_, _ = rsw.Write(data)
		rsw.Flush()

		require.Equal(t, len(data), count)
	})
	t.Run("flusher-nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			rsw := NewResponseWriter(writer{})
			rsw.Flush()
		})
	})
	t.Run("flusher", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var called bool
			rsw := NewResponseWriter(writerFlusher{
				flushFn: func() {
					called = true
				},
			})
			rsw.Flush()
			assert.True(t, called)
		})
	})
}

func TestNResponseWriter_Close(t *testing.T) {
	t.Run("writeBody", func(t *testing.T) {
		var count int
		rsw := NewResponseWriter(writerMock{
			write: func(i []byte) (int, error) {
				count = len(i)
				return count, nil
			},
		})

		data := make([]byte, 128)
		_, _ = rand.Read(data)
		_, _ = rsw.Write(data)
		rsw.Close()

		require.Equal(t, len(data), count)
	})
	t.Run("writeBody-error", func(t *testing.T) {
		rsw := NewResponseWriter(writerMock{
			write: func(i []byte) (int, error) {
				return 0, errors.New("mock")
			},
		})

		data := make([]byte, 128)
		_, _ = rand.Read(data)
		_, _ = rsw.Write(data)
		require.NotNil(t, rsw.Close())
	})
	t.Run("closer-nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			rsw := NewResponseWriter(writer{})
			require.NoError(t, rsw.Close())
		})
	})
	t.Run("closer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var called bool
			rsw := NewResponseWriter(writerCloser{
				closeFn: func() error {
					called = true
					return nil
				},
			})
			require.NoError(t, rsw.Close())
			assert.True(t, called)
		})
	})
}

func TestNResponseWriter_IsCommitted(t *testing.T) {
	rsw := NewResponseWriter(writerMock{
		write: func(i []byte) (int, error) {
			return 0, nil
		},
	})

	data := make([]byte, BufferLimit+1)
	_, _ = rand.Read(data)
	_, _ = rsw.Write(data)

	assert.True(t, rsw.IsCommitted())
}
