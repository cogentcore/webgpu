//go:build !js

package wgpu

/*
#include <stdlib.h>
#include "./lib/wgpu.h"
*/
import "C"

type Texture struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref         C.WGPUTexture
}

func (p *Texture) CreateView(descriptor *TextureViewDescriptor) (*TextureView, error) {
	var desc *C.WGPUTextureViewDescriptor

	if descriptor != nil {
		desc = &C.WGPUTextureViewDescriptor{
			format:          C.WGPUTextureFormat(descriptor.Format),
			dimension:       C.WGPUTextureViewDimension(descriptor.Dimension),
			baseMipLevel:    C.uint32_t(descriptor.BaseMipLevel),
			mipLevelCount:   C.uint32_t(descriptor.MipLevelCount),
			baseArrayLayer:  C.uint32_t(descriptor.BaseArrayLayer),
			arrayLayerCount: C.uint32_t(descriptor.ArrayLayerCount),
			aspect:          C.WGPUTextureAspect(descriptor.Aspect),
		}

		if descriptor.Label != "" {
			label, freeLabel := stringViewOf(descriptor.Label)
			defer freeLabel()
			desc.label = label
		}
	}

	var err error
	pushValidationScope(p.deviceRef)
	ref := C.wgpuTextureCreateView(p.ref, desc)
	popValidationScope(p.deviceRef, p.instanceRef, "wgpu.(*Texture).CreateView(): ", &err)
	if err != nil {
		C.wgpuTextureViewRelease(ref)
		return nil, err
	}

	return &TextureView{ref}, nil
}

func (p *Texture) Destroy() {
	C.wgpuTextureDestroy(p.ref)
}

func (p *Texture) GetDepthOrArrayLayers() uint32 {
	return uint32(C.wgpuTextureGetDepthOrArrayLayers(p.ref))
}

func (p *Texture) GetDimension() TextureDimension {
	return TextureDimension(C.wgpuTextureGetDimension(p.ref))
}

func (p *Texture) GetFormat() TextureFormat {
	return TextureFormat(C.wgpuTextureGetFormat(p.ref))
}

func (p *Texture) GetHeight() uint32 {
	return uint32(C.wgpuTextureGetHeight(p.ref))
}

func (p *Texture) GetMipLevelCount() uint32 {
	return uint32(C.wgpuTextureGetMipLevelCount(p.ref))
}

func (p *Texture) GetSampleCount() uint32 {
	return uint32(C.wgpuTextureGetSampleCount(p.ref))
}

func (p *Texture) GetUsage() TextureUsage {
	return TextureUsage(C.wgpuTextureGetUsage(p.ref))
}

func (p *Texture) GetWidth() uint32 {
	return uint32(C.wgpuTextureGetWidth(p.ref))
}

func (p *Texture) Release() {
	C.wgpuDeviceRelease(p.deviceRef)
	C.wgpuTextureRelease(p.ref)
}
