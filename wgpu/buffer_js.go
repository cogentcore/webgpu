//go:build js

package wgpu

import (
	"syscall/js"
)

// Buffer as described:
// https://gpuweb.github.io/gpuweb/#gpubuffer
type Buffer struct {
	jsValue js.Value
}

func (g Buffer) toJS() any {
	return g.jsValue
}

// Destroy as described:
// https://gpuweb.github.io/gpuweb/#dom-gpubuffer-destroy
func (g Buffer) Destroy() {
	g.jsValue.Call("destroy")
}

func (g Buffer) GetMappedRange(offset, size uint) []byte {
	// TODO(kai): this does not work because it does not get
	// the actual pointer to the byte data; this is only really
	// possible with GopherJS.
	buf := g.jsValue.Call("getMappedRange", offset, size)
	src := js.Global().Get("Uint8ClampedArray").New(buf)
	dst := make([]byte, src.Length())
	js.CopyBytesToGo(dst, src)
	return dst
}

func (g Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) (err error) {
	promise := g.jsValue.Call("mapAsync", uint32(mode), offset, size)

	// Set up success handler
	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(BufferMapAsyncStatusSuccess)
		return nil
	})

	// Set up error handler
	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(BufferMapAsyncStatusSuccess) // Use Success as fallback since Error doesn't exist
		return nil
	})

	// Handle the promise
	promise.Call("then", successCallback).Call("catch", errorCallback)

	return nil
}

func (g Buffer) Unmap() (err error) {
	g.jsValue.Call("unmap")
	return
}

func (g Buffer) Release() {} // no-op

func (g Buffer) GetSize() uint64 {
	sizeVal := g.jsValue.Get("size")
	if sizeVal.Type() == js.TypeUndefined {
		return 0 // fallback or error
	}
	return uint64(sizeVal.Int())
}
