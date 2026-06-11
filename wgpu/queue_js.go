//go:build js

package wgpu

import (
	"syscall/js"

	"github.com/openfluke/webgpu/jsx"
)

// Queue as described:
// https://gpuweb.github.io/gpuweb/#gpuqueue
type Queue struct {
	jsValue js.Value
}

func (g Queue) toJS() any {
	return g.jsValue
}

// Submit as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-submit
func (g Queue) Submit(commandBuffers ...*CommandBuffer) {
	jsSequence := mapSlice(commandBuffers, func(buffer *CommandBuffer) any {
		return pointerToJS(buffer)
	})
	g.jsValue.Call("submit", jsSequence)
}

// WriteBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writebuffer
func (g Queue) WriteBuffer(buffer *Buffer, offset uint64, data []byte) (err error) {
	g.jsValue.Call("writeBuffer", pointerToJS(buffer), offset, jsx.BytesToJS(data), uint64(0), len(data))
	return
}

// WriteTexture as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writetexture
func (g Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) (err error) {
	g.jsValue.Call("writeTexture", pointerToJS(destination), jsx.BytesToJS(data), pointerToJS(dataLayout), pointerToJS(writeSize))
	return
}

// OnSubmittedWorkDone as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-onsubmittedworkdone
func (g Queue) OnSubmittedWorkDone(callback QueueWorkDoneCallback) {
	// Don't use jsx.Await - handle promise with callbacks

	promise := g.jsValue.Call("onSubmittedWorkDone")

	// Set up success handler
	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(QueueWorkDoneStatusSuccess)
		return nil
	})

	// Set up error handler
	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(QueueWorkDoneStatusError) // or whatever error status you have
		return nil
	})

	// Handle the promise
	promise.Call("then", successCallback).Call("catch", errorCallback)
}

func (g Queue) Release() {} // no-op
