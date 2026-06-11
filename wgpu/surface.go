//go:build !js

package wgpu

/*
#include <stdlib.h>
#include "./lib/wgpu.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

type Surface struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref         C.WGPUSurface
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
	p.instanceRef = device.instanceRef

	var cfg *C.WGPUSurfaceConfiguration
	if config != nil {
		cfg = &C.WGPUSurfaceConfiguration{
			device:      p.deviceRef,
			format:      C.WGPUTextureFormat(config.Format),
			usage:       C.WGPUTextureUsage(config.Usage),
			alphaMode:   C.WGPUCompositeAlphaMode(config.AlphaMode),
			width:       C.uint32_t(config.Width),
			height:      C.uint32_t(config.Height),
			presentMode: C.WGPUPresentMode(config.PresentMode),
		}
		viewFormatCount := len(config.ViewFormats)
		if viewFormatCount > 0 {
			viewFormats := C.malloc(C.size_t(unsafe.Sizeof(C.WGPUTextureFormat(0))) * C.size_t(viewFormatCount))
			defer C.free(viewFormats)

			viewFormatsSlice := unsafe.Slice((*TextureFormat)(viewFormats), viewFormatCount)
			copy(viewFormatsSlice, config.ViewFormats)

			cfg.viewFormatCount = C.size_t(viewFormatCount)
			cfg.viewFormats = (*C.WGPUTextureFormat)(viewFormats)
		}
	}

	C.wgpuSurfaceConfigure(p.ref, cfg)
}

// NOTE: you should typically not call [Texture.Release] on the returned texture.
// Instead, you should call [TextureView.Release] on any [TextureView] you create from it.
func (p *Surface) GetCurrentTexture() (*Texture, error) {
	var surfaceTexture C.WGPUSurfaceTexture
	C.wgpuSurfaceGetCurrentTexture(p.ref, &surfaceTexture)

	switch surfaceTexture.status {
	case C.WGPUSurfaceGetCurrentTextureStatus_SuccessOptimal,
		C.WGPUSurfaceGetCurrentTextureStatus_SuccessSuboptimal:
		return &Texture{
			deviceRef:   p.deviceRef,
			instanceRef: p.instanceRef,
			ref:         surfaceTexture.texture,
		}, nil
	default:
		if surfaceTexture.texture != nil {
			C.wgpuTextureRelease(surfaceTexture.texture)
		}
		return nil, errors.New("wgpu.(*Surface).GetCurrentTexture(): failed to acquire texture")
	}
}

func (p *Surface) Present() {
	C.wgpuSurfacePresent(p.ref)
}

func (p *Surface) Release() {
	C.wgpuSurfaceRelease(p.ref)
}
