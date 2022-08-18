package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Modified code from here https://github.com/oxtoacart/bpool

func TestSizedBufferPool(t *testing.T) {

	var size int = 4
	var capacity int = 1024

	bufPool := NewSizedBufferPool(size, capacity)

	bufPool.WithBuffer(func(b *bytes.Buffer) {
		// Check the cap before we use the buffer.
		assert.Equal(t, cap(b.Bytes()), capacity)

		// Grow the buffer beyond our capacity and return it to the pool
		b.Grow(capacity * 3)
	})

	// Add some additional buffers to fill up the pool.
	for i := 0; i < size; i++ {
		bufPool.put(bytes.NewBuffer(make([]byte, 0, bufPool.a*2)))
	}

	// Check that oversized buffers are being replaced.
	assert.GreaterOrEqual(t, len(bufPool.c), size)

	// Close the channel so we can iterate over it.
	close(bufPool.c)

	// Check that there are buffers of the correct capacity in the pool.
	for buffer := range bufPool.c {
		assert.Equal(t, cap(buffer.Bytes()), bufPool.a)
	}

}
