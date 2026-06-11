//go:build !js

package wgpu

/*
#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_queue_work_done_callback_c(WGPUQueueWorkDoneStatus status, WGPUStringView message, void* userdata1, void* userdata2);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

type Queue struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref         C.WGPUQueue
}

//export gowebgpu_queue_work_done_callback_go
func gowebgpu_queue_work_done_callback_go(status C.WGPUQueueWorkDoneStatus, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	_ = messageData
	_ = messageLen
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()

	cb, ok := handle.Value().(QueueWorkDoneCallback)
	if ok {
		cb(QueueWorkDoneStatus(status))
	}
}

func (p *Queue) OnSubmittedWorkDone(callback QueueWorkDoneCallback) {
	handle := cgo.NewHandle(callback)
	handlePtr := cgoHandlePtr(handle)
	cbInfo := C.WGPUQueueWorkDoneCallbackInfo{
		mode:      C.WGPUCallbackMode_WaitAnyOnly,
		callback:  C.WGPUQueueWorkDoneCallback(C.gowebgpu_queue_work_done_callback_c),
		userdata2: handlePtr,
	}
	C.wgpuQueueOnSubmittedWorkDone(p.ref, cbInfo)
}

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r := C.wgpuQueueSubmitForIndex(p.ref, 0, nil)
		return SubmissionIndex(r)
	}

	commandRefs := C.malloc(C.size_t(commandCount) * C.size_t(unsafe.Sizeof(C.WGPUCommandBuffer(nil))))
	defer C.free(commandRefs)

	commandRefsSlice := unsafe.Slice((*C.WGPUCommandBuffer)(commandRefs), commandCount)
	for i, v := range commands {
		commandRefsSlice[i] = v.ref
	}

	r := C.wgpuQueueSubmitForIndex(
		p.ref,
		C.size_t(commandCount),
		(*C.WGPUCommandBuffer)(commandRefs),
	)
	return SubmissionIndex(r)
}

func (p *Queue) WriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) (err error) {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*Queue).WriteBuffer(): ", func() {
		size := len(data)
		if size == 0 {
			C.wgpuQueueWriteBuffer(p.ref, buffer.ref, C.uint64_t(bufferOffset), nil, 0)
			return
		}
		C.wgpuQueueWriteBuffer(p.ref, buffer.ref, C.uint64_t(bufferOffset), unsafe.Pointer(&data[0]), C.size_t(size))
	})
}

func (p *Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) (err error) {
	var dst C.WGPUTexelCopyTextureInfo
	if destination != nil {
		dst = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var layout C.WGPUTexelCopyBufferLayout
	if dataLayout != nil {
		layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(dataLayout.Offset),
			bytesPerRow:  C.uint32_t(dataLayout.BytesPerRow),
			rowsPerImage: C.uint32_t(dataLayout.RowsPerImage),
		}
	}

	var writeExtent C.WGPUExtent3D
	if writeSize != nil {
		writeExtent = C.WGPUExtent3D{
			width:              C.uint32_t(writeSize.Width),
			height:             C.uint32_t(writeSize.Height),
			depthOrArrayLayers: C.uint32_t(writeSize.DepthOrArrayLayers),
		}
	}

	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*Queue).WriteTexture(): ", func() {
		size := len(data)
		if size == 0 {
			C.wgpuQueueWriteTexture(p.ref, &dst, nil, 0, &layout, &writeExtent)
			return
		}
		C.wgpuQueueWriteTexture(p.ref, &dst, unsafe.Pointer(&data[0]), C.size_t(size), &layout, &writeExtent)
	})
}

func (p *Queue) Release() {
	C.wgpuDeviceRelease(p.deviceRef)
	C.wgpuQueueRelease(p.ref)
}