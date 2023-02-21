package n

import (
	"bytes"
	"io"
	"net/http"
)

var (
	BufferLimit = 1024 * 1024 * 16
)

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	io.Closer
	IsCommitted() bool
}

type NResponseWriter struct {
	body      *bytes.Buffer
	rs        http.ResponseWriter
	flusher   http.Flusher
	closer    io.Closer
	committed bool
}

func NewResponseWriter(rs http.ResponseWriter) *NResponseWriter {
	flusher, _ := rs.(http.Flusher)
	closer, _ := rs.(io.Closer)

	return &NResponseWriter{
		rs:      rs,
		flusher: flusher,
		closer:  closer,
	}
}

func (N *NResponseWriter) Header() http.Header {
	return N.rs.Header()
}

func (N *NResponseWriter) Write(bytes []byte) (int, error) {
	if len(bytes) <= BufferLimit {
		N.body = getBuffer()
		return N.body.Write(bytes)
	}

	N.committed = true
	return N.rs.Write(bytes)
}

func (N *NResponseWriter) WriteHeader(statusCode int) {
	N.rs.WriteHeader(statusCode)
}

func (N *NResponseWriter) Flush() {
	_, _ = N.writeBody()
	if N.flusher == nil {
		return
	}

	N.flusher.Flush()
	return
}

func (N *NResponseWriter) Close() error {
	if _, err := N.writeBody(); err != nil {
		return err
	}

	if N.closer == nil {
		return nil
	}

	return N.closer.Close()
}

func (N *NResponseWriter) IsCommitted() bool {
	return N.committed
}

func (N *NResponseWriter) writeBody() (int64, error) {
	if N.body == nil {
		return 0, nil
	}

	defer func() {
		putBuffer(N.body)
		N.body = nil
		N.committed = true
	}()
	return io.Copy(N.rs, N.body)
}
