//go:build js

package wgpu

import (
	"syscall/js"
	"unsafe"
)

// BytesToJS converts the given bytes to a js Uint8ClampedArray
// by using the global wasm memory bytes. This avoids the
// copying present in [js.CopyBytesToJS].
func BytesToJS(b []byte) js.Value {
	ptr := uintptr(unsafe.Pointer(&b[0]))
	memoryBytes := js.Global().Get("Uint8ClampedArray").New(js.Global().Get("wasm").Get("instance").Get("exports").Get("mem").Get("buffer"))
	// using subarray instead of slice gives a 5x performance improvement due to no copying
	return memoryBytes.Call("subarray", ptr, ptr+uintptr(len(b)))
}
