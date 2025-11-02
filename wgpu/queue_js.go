//go:build js

package wgpu

import (
	"runtime"
	"syscall/js"
	"unsafe"

	"github.com/cogentcore/webgpu/jsx"
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
	defer runtime.KeepAlive(data)

	address := uintptr(unsafe.Pointer(&data[0]))
	queueWriteBuffer.Invoke(g.jsValue, pointerToJS(buffer), offset, address, uint64(0), len(data))
	return
}

// WriteTexture as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writetexture
func (g Queue) WriteTexture(destination *TexelCopyTextureInfo, data []byte, dataLayout *TexelCopyBufferLayout, writeSize *Extent3D) (err error) {
	defer runtime.KeepAlive(data)

	address := uintptr(unsafe.Pointer(&data[0]))
	queueWriteTexture.Invoke(g.jsValue, pointerToJS(destination), address, len(data), pointerToJS(dataLayout), pointerToJS(writeSize))
	return
}

// OnSubmittedWorkDone as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-onsubmittedworkdone
func (g Queue) OnSubmittedWorkDone(callback QueueWorkDoneCallback) {
	jsx.Await(g.jsValue.Call("onSubmittedWorkDone")) // TODO(kai): is this correct?
	callback(QueueWorkDoneStatusSuccess)
}

func (g Queue) Release() {} // no-op

var queueWriteBuffer js.Value
var queueWriteTexture js.Value

func init() {
	queueWriteBuffer = js.Global().Call("eval", `
		(queue, buf, offset, addr, x, n) => {
			const mem = wasm.instance.exports.mem.buffer;
			const data = new Uint8ClampedArray(mem, addr, n);
			return queue.writeBuffer(buf, offset, data, x, n);
		} 
	`)

	queueWriteTexture = js.Global().Call("eval", `
		(queue, tex, addr, n, layout, writeSize) => {
			const mem = wasm.instance.exports.mem.buffer;
			const data = new Uint8ClampedArray(mem, addr, n);
			return queue.writeTexture(tex, data, layout, writeSize);
		} 
	`)
}
