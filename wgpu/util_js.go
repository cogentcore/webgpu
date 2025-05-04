//go:build js

package wgpu

import (
	"log/slog"
	"syscall/js"
	"unsafe"
)

// BytesToJS converts the given bytes to a js Uint8ClampedArray
// by using the global wasm memory bytes. This avoids the
// copying present in [js.CopyBytesToJS].
func BytesToJS(b []byte) js.Value {
	ptr := uintptr(unsafe.Pointer(&b[0]))
	// We directly pass the offset and length to the constructor to avoid calling subarray or slice,
	// thereby improving performance and safety (this fixes a detached array buffer crash).
	return js.Global().Get("Uint8ClampedArray").New(js.Global().Get("wasm").Get("instance").Get("exports").Get("mem").Get("buffer"), ptr, len(b))
}

// mapSlice can be used to transform one slice into another by providing a
// function to do the mapping.
func mapSlice[S, T any](slice []S, fn func(S) T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// AwaitJS is a helper function equivalent to await in JS.
// It is copied from https://go-review.googlesource.com/c/go/+/150917/
func AwaitJS(promise js.Value) (result js.Value, ok bool) {
	if promise.Type() != js.TypeObject || promise.Get("then").Type() != js.TypeFunction {
		return promise, true
	}

	done := make(chan struct{})

	onResolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		result = args[0]
		ok = true
		close(done)
		return nil
	})
	defer onResolve.Release()

	onReject := js.FuncOf(func(this js.Value, args []js.Value) any {
		result = args[0]
		ok = false
		slog.Error("wgpu.AwaitJS: promise rejected", "reason", result)
		close(done)
		return nil
	})
	defer onReject.Release()

	promise.Call("then", onResolve, onReject)
	<-done
	return
}

// no-ops
func SetLogLevel(level LogLevel) {}
func GetVersion() Version        { return 0 }
