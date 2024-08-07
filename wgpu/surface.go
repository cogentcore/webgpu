//go:build !js

package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_error_callback_c(WGPUErrorType type, char const * message, void * userdata);

static inline WGPUTexture gowebgpu_surface_get_current_texture(WGPUSurface surface, WGPUDevice device, void * error_userdata) {
	WGPUSurfaceTexture ref = {};
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuSurfaceGetCurrentTexture(surface, &ref);
	wgpuDevicePopErrorScope(device, gowebgpu_error_callback_c, error_userdata);
	return ref.texture;
}

*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type Surface struct {
	deviceRef C.WGPUDevice
	ref       C.WGPUSurface
}

func (p *Surface) GetCapabilities(adapter *Adapter) (ret SurfaceCapabilities) {
	var caps C.WGPUSurfaceCapabilities
	C.wgpuSurfaceGetCapabilities(p.ref, adapter.ref, &caps)

	if caps.alphaModeCount == 0 && caps.formatCount == 0 && caps.presentModeCount == 0 {
		return
	}
	if caps.formatCount > 0 {
		caps.formats = (*C.WGPUTextureFormat)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUTextureFormat(0))) * caps.formatCount))
		defer C.free(unsafe.Pointer(caps.formats))
	}
	if caps.presentModeCount > 0 {
		caps.presentModes = (*C.WGPUPresentMode)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUPresentMode(0))) * caps.presentModeCount))
		defer C.free(unsafe.Pointer(caps.presentModes))
	}
	if caps.alphaModeCount > 0 {
		caps.alphaModes = (*C.WGPUCompositeAlphaMode)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUCompositeAlphaMode(0))) * caps.alphaModeCount))
		defer C.free(unsafe.Pointer(caps.alphaModes))
	}

	C.wgpuSurfaceGetCapabilities(p.ref, adapter.ref, &caps)

	if caps.formatCount > 0 {
		formatsTmp := unsafe.Slice((*TextureFormat)(caps.formats), caps.formatCount)
		ret.Formats = make([]TextureFormat, caps.formatCount)
		copy(ret.Formats, formatsTmp)
	}
	if caps.presentModeCount > 0 {
		presentModesTmp := unsafe.Slice((*PresentMode)(caps.presentModes), caps.presentModeCount)
		ret.PresentModes = make([]PresentMode, caps.presentModeCount)
		copy(ret.PresentModes, presentModesTmp)
	}
	if caps.alphaModeCount > 0 {
		alphaModesTmp := unsafe.Slice((*CompositeAlphaMode)(caps.alphaModes), caps.alphaModeCount)
		ret.AlphaModes = make([]CompositeAlphaMode, caps.alphaModeCount)
		copy(ret.AlphaModes, alphaModesTmp)
	}

	return
}

func (p *Surface) Configure(adapter *Adapter, device *Device, config *SurfaceConfiguration) {
	p.deviceRef = device.ref
	C.wgpuSurfaceConfigure(p.ref, config) // TODO(kai): convert to C struct
}

func (p *Surface) GetCurrentTexture() (*Texture, error) {
	var err error = nil
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*Surface).GetCurrentTexture(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	ref := C.gowebgpu_surface_get_current_texture(
		p.ref,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	if err != nil {
		if ref != nil {
			C.wgpuTextureRelease(ref)
		}
		return nil, err
	}

	return &Texture{p.deviceRef, ref}, nil
}

func (p *Surface) Present() {
	C.wgpuSurfacePresent(p.ref)
}

func (p *Surface) Release() {
	C.wgpuSurfaceRelease(p.ref)
}
