//go:build !js

package wgpu

/*
#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_pop_error_scope_callback_c(WGPUPopErrorScopeStatus status, WGPUErrorType type, WGPUStringView message, void* userdata1, void* userdata2);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

const wgpuStrlen = ^C.size_t(0)

func emptyStringView() C.WGPUStringView {
	return C.WGPUStringView{data: nil, length: wgpuStrlen}
}

func stringViewOf(s string) (C.WGPUStringView, func()) {
	if s == "" {
		return emptyStringView(), func() {}
	}
	cstr := C.CString(s)
	sv := C.WGPUStringView{
		data:   cstr,
		length: C.size_t(len(s)),
	}
	return sv, func() { C.free(unsafe.Pointer(cstr)) }
}

func goStringView(sv C.WGPUStringView) string {
	return goStringViewFromParts(uintptr(unsafe.Pointer(sv.data)), uintptr(sv.length))
}

func goStringViewFromParts(data uintptr, length uintptr) string {
	if data == 0 {
		return ""
	}
	if length == uintptr(wgpuStrlen) {
		return C.GoString((*C.char)(unsafe.Pointer(data)))
	}
	return C.GoStringN((*C.char)(unsafe.Pointer(data)), C.int(length))
}

func (p *Instance) waitFuture(future C.WGPUFuture) {
	if future.id == 0 || p.ref == nil {
		return
	}
	wait := C.WGPUFutureWaitInfo{future: future}
	C.wgpuInstanceWaitAny(p.ref, 1, &wait, ^C.uint64_t(0))
}

func pushValidationScope(device C.WGPUDevice) {
	C.wgpuDevicePushErrorScope(device, C.WGPUErrorFilter_Validation)
}

type popErrorScopeState struct {
	errPtr *error
	prefix string
}

//export gowebgpu_pop_error_scope_callback_go
func gowebgpu_pop_error_scope_callback_go(status C.WGPUPopErrorScopeStatus, typ C.WGPUErrorType, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()
	state, ok := handle.Value().(popErrorScopeState)
	if !ok {
		return
	}
	if status == C.WGPUPopErrorScopeStatus_Success && typ != C.WGPUErrorType_NoError && state.errPtr != nil {
		*state.errPtr = errors.New(state.prefix + goStringViewFromParts(messageData, messageLen))
	}
}

func popValidationScope(device C.WGPUDevice, instance C.WGPUInstance, prefix string, errPtr *error) {
	handle := cgo.NewHandle(popErrorScopeState{errPtr: errPtr, prefix: prefix})
	handlePtr := cgoHandlePtr(handle)
	cbInfo := C.WGPUPopErrorScopeCallbackInfo{
		mode:      C.WGPUCallbackMode_WaitAnyOnly,
		callback:  C.WGPUPopErrorScopeCallback(C.gowebgpu_pop_error_scope_callback_c),
		userdata2: handlePtr,
	}
	future := C.wgpuDevicePopErrorScope(device, cbInfo)
	(&Instance{ref: instance}).waitFuture(future)
}

func (p *Device) withValidation(prefix string, fn func()) error {
	pushValidationScope(p.ref)
	fn()
	var err error
	popValidationScope(p.ref, p.instanceRef, prefix, &err)
	return err
}

func withDeviceValidation(device C.WGPUDevice, instance C.WGPUInstance, prefix string, fn func()) error {
	pushValidationScope(device)
	fn()
	var err error
	popValidationScope(device, instance, prefix, &err)
	return err
}

func cgoHandlePtr(h cgo.Handle) unsafe.Pointer {
	p := C.malloc(C.size_t(unsafe.Sizeof(uintptr(0))))
	*(*cgo.Handle)(p) = h
	return p
}

func freeCgoHandlePtr(p unsafe.Pointer) {
	if p != nil {
		C.free(p)
	}
}

func optionalBool(b bool) C.WGPUOptionalBool {
	if b {
		return C.WGPUOptionalBool_True
	}
	return C.WGPUOptionalBool_False
}

func mapAsyncStatusFromC(status C.WGPUMapAsyncStatus) BufferMapAsyncStatus {
	switch status {
	case C.WGPUMapAsyncStatus_Success:
		return BufferMapAsyncStatusSuccess
	case C.WGPUMapAsyncStatus_Error:
		return BufferMapAsyncStatusValidationError
	case C.WGPUMapAsyncStatus_Aborted:
		return BufferMapAsyncStatusDeviceLost
	case C.WGPUMapAsyncStatus_CallbackCancelled:
		return BufferMapAsyncStatusDestroyedBeforeCallback
	default:
		return BufferMapAsyncStatusUnknown
	}
}
