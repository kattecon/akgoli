package utils

import (
	"bytes"
)

// Based on https://github.com/oxtoacart/bpool/blob/master/sizedbufferpool.go
// with the addition of WithBuffer, which borrows a buffer, runs a callback,
// and returns the buffer to the pool via the existing get and put helpers.

// SizedBufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel. Buffers are created lazily with the requested alloc capacity.
//
// The size parameter bounds only the number of idle buffers retained in the
// channel. When a caller requests a buffer and the channel is empty, a new
// buffer is allocated immediately without waiting. The total number of live
// buffers across concurrent callers is therefore not bounded by size.
//
// After each use, any buffer whose capacity exceeds alloc is discarded and
// replaced with a fresh buffer of exactly alloc capacity.
type SizedBufferPool struct {
	c chan *bytes.Buffer
	a int
}

// NewSizedBufferPool creates a new SizedBufferPool bounded to the given size.
// size defines the number of buffers to be retained in the pool and alloc sets
// the initial capacity of new buffers to minimize calls to make().
//
// The value of alloc should seek to provide a buffer that is representative of
// most data written to the buffer (i.e. 95th percentile) without being
// overly large (which will increase static memory consumption). You may wish to
// track the capacity of your last N buffers (i.e. using an []int) prior to
// returning them to the pool as input into calculating a suitable alloc value.
//
// Both size and alloc must be non-negative. A negative size panics during
// construction. A negative alloc panics on the first cache miss when a new
// buffer is created.
func NewSizedBufferPool(size int, alloc int) (bp *SizedBufferPool) {
	return &SizedBufferPool{
		c: make(chan *bytes.Buffer, size),
		a: alloc,
	}
}

// WithBuffer borrows a buffer from the pool, passes it to f, and returns it
// to the pool when f returns. The buffer is empty at the start of the call.
// The caller must not retain the buffer pointer or any slice obtained from it
// (such as b.Bytes()) after f returns. The deferred return runs even during
// panic unwinding, so another goroutine may receive and overwrite the same
// storage immediately.
func (bp *SizedBufferPool) WithBuffer(f func(b *bytes.Buffer)) {
	b := bp.get()
	defer bp.put(b)
	f(b)
}

// get returns a Buffer from the SizedBufferPool, or creates a new one if none are
// available in the pool. Buffers have a pre-allocated capacity.
func (bp *SizedBufferPool) get() (b *bytes.Buffer) {
	select {
	case b = <-bp.c:
	// reuse existing buffer.
	default:
		// create new buffer.
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}
	return
}

// put returns the given Buffer to the SizedBufferPool.
func (bp *SizedBufferPool) put(b *bytes.Buffer) {
	b.Reset()

	// Release buffers over our maximum capacity and re-create a pre-sized
	// buffer to replace it.
	// Note that the cap(b.Bytes()) provides the capacity from the read off-set
	// only, but as we've called b.Reset() the full capacity of the underlying
	// byte slice is returned.
	if cap(b.Bytes()) > bp.a {
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}

	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}
