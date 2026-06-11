//go:build !js

package wgpu

/*
#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_buffer_map_callback_c(WGPUMapAsyncStatus status, WGPUStringView message, void* userdata1, void* userdata2);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

type Buffer struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref         C.WGPUBuffer
}

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
}

func (p *Buffer) GetMappedRange(offset, size uint) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(buf), size)
}

func (p *Buffer) GetSize() uint64 {
	return uint64(C.wgpuBufferGetSize(p.ref))
}

func (p *Buffer) GetUsage() BufferUsage {
	return BufferUsage(C.wgpuBufferGetUsage(p.ref))
}

//export gowebgpu_buffer_map_callback_go
func gowebgpu_buffer_map_callback_go(status C.WGPUMapAsyncStatus, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	_ = messageData
	_ = messageLen
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()

	cb, ok := handle.Value().(BufferMapCallback)
	if ok {
		cb(mapAsyncStatusFromC(status))
	}
}

func (p *Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) (err error) {
	callbackHandle := cgo.NewHandle(callback)
	callbackHandlePtr := cgoHandlePtr(callbackHandle)
	cbInfo := C.WGPUBufferMapCallbackInfo{
		mode:      C.WGPUCallbackMode_WaitAnyOnly,
		callback:  C.WGPUBufferMapCallback(C.gowebgpu_buffer_map_callback_c),
		userdata2: callbackHandlePtr,
	}
	if err = withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*Buffer).MapAsync(): ", func() {
		future := C.wgpuBufferMapAsync(p.ref, C.WGPUMapMode(mode), C.size_t(offset), C.size_t(size), cbInfo)
		(&Instance{ref: p.instanceRef}).waitFuture(future)
	}); err != nil {
		return err
	}
	return nil
}

func (p *Buffer) Unmap() (err error) {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*Buffer).Unmap(): ", func() {
		C.wgpuBufferUnmap(p.ref)
	})
}

func (p *Buffer) Release() {
	C.wgpuDeviceRelease(p.deviceRef)
	C.wgpuBufferRelease(p.ref)
}
