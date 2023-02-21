package n

import (
	"bytes"
	"sync"
)

var (
	BufferSize     = 1024 * 8 // 8 KB
	BufferPoolSize = 16

	_bufferPool = &sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, BufferSize))
		},
	}
)

func init() {
	for i := 0; i < BufferPoolSize; i++ {
		_bufferPool.Put(_bufferPool.Get())
	}
}

func getBuffer() *bytes.Buffer {
	buffer := _bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

func putBuffer(b *bytes.Buffer) {
	if b == nil {
		return
	}

	_bufferPool.Put(b)
}
