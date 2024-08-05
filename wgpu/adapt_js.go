//go:build js

package wgpu

import (
	"runtime"
	"syscall/js"
	"unsafe"
)

var (
	// NOTE: We use a global ArrayBuffer and a few TypedArray views on top
	// of it for WebGPU calls that require such instead of allocating new ones
	// for each call.

	bufferSize  int
	arrayBuffer js.Value
	uint8Array  js.Value
)

// ensureBufferSize ensures that the global ArrayBuffer has a size
// that is equal or larger to the specified size.
func ensureBufferSize(size int) {
	if size <= bufferSize {
		return
	}

	// Grow to the smallest power of two that will satisfy the requested size.
	if bufferSize == 0 {
		bufferSize = 256 * 256
	}
	for bufferSize < size {
		bufferSize *= 2
	}

	arrayBuffer = js.Global().Get("ArrayBuffer").New(bufferSize)
	uint8Array = js.Global().Get("Uint8Array").New(arrayBuffer)
}

// DataTypes represents allowed data slice types.
type DataTypes interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~float32 | ~float64
}

// stageBufferData inserts the specified data into the global
// ArrayBuffer to be used in WebGPU calls.
//
// This function will panic if data cannot be converted to []byte
// through the asByteSlice function.
func stageBufferData[T DataTypes](data []T) uint64 {
	byteData := asByteSlice(data)
	ensureBufferSize(len(byteData))
	js.CopyBytesToJS(uint8Array, byteData)
	runtime.KeepAlive(data)
	return uint64(len(byteData))
}

// asByteSlice returns a []byte representation for the
// specified arbitrary slice type.
//
// This utility function is related to the following issues:
// https://github.com/golang/go/issues/32402
// https://github.com/golang/go/issues/31980
func asByteSlice[T DataTypes](data []T) []byte {
	if len(data) == 0 {
		return nil
	}
	dataSize := byteSize(data)
	return unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), dataSize)
}

// byteSize returns the number of bytes that would be
// needed to represent data once it is converted to a
// byte slice through asByteSlice.
func byteSize[T DataTypes](data []T) int {
	if len(data) == 0 {
		return 0
	}
	return len(data) * int(unsafe.Sizeof(data[0]))
}
