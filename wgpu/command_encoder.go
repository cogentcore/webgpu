//go:build !js

package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

*/
import "C"
import (
	"unsafe"
)

type CommandEncoder struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref       C.WGPUCommandEncoder
}

type ComputePassDescriptor struct {
	Label string

	// unused in wgpu
	// TimestampWrites []ComputePassTimestampWrite
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc *C.WGPUComputePassDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label, freeLabel := stringViewOf(descriptor.Label)
		defer freeLabel()
		desc = &C.WGPUComputePassDescriptor{label: label}
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, desc)
	if ref == nil {
		panic("Failed to acquire ComputePassEncoder")
	}

	C.wgpuDeviceAddRef(p.deviceRef)
	return &ComputePassEncoder{deviceRef: p.deviceRef, instanceRef: p.instanceRef, ref: ref}
}

func (p *CommandEncoder) BeginRenderPass(descriptor *RenderPassDescriptor) *RenderPassEncoder {
	var desc C.WGPURenderPassDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label, freeLabel := stringViewOf(descriptor.Label)
			defer freeLabel()
			desc.label = label
		}

		colorAttachmentCount := len(descriptor.ColorAttachments)
		if colorAttachmentCount > 0 {
			colorAttachments := C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassColorAttachment{})) * C.size_t(colorAttachmentCount))
			defer C.free(colorAttachments)

			colorAttachmentsSlice := unsafe.Slice((*C.WGPURenderPassColorAttachment)(colorAttachments), colorAttachmentCount)

			for i, v := range descriptor.ColorAttachments {
				colorAttachment := C.WGPURenderPassColorAttachment{
					loadOp:     C.WGPULoadOp(v.LoadOp),
					storeOp:    C.WGPUStoreOp(v.StoreOp),
					depthSlice: C.WGPU_DEPTH_SLICE_UNDEFINED,
					clearValue: C.WGPUColor{
						r: C.double(v.ClearValue.R),
						g: C.double(v.ClearValue.G),
						b: C.double(v.ClearValue.B),
						a: C.double(v.ClearValue.A),
					},
				}
				if v.View != nil {
					colorAttachment.view = v.View.ref
				}
				if v.ResolveTarget != nil {
					colorAttachment.resolveTarget = v.ResolveTarget.ref
				}

				colorAttachmentsSlice[i] = colorAttachment
			}

			desc.colorAttachmentCount = C.size_t(colorAttachmentCount)
			desc.colorAttachments = (*C.WGPURenderPassColorAttachment)(colorAttachments)
		}

		if descriptor.DepthStencilAttachment != nil {
			depthStencilAttachment := (*C.WGPURenderPassDepthStencilAttachment)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassDepthStencilAttachment{}))))
			defer C.free(unsafe.Pointer(depthStencilAttachment))

			if descriptor.DepthStencilAttachment.View != nil {
				depthStencilAttachment.view = descriptor.DepthStencilAttachment.View.ref
			}
			depthStencilAttachment.depthLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.DepthLoadOp)
			depthStencilAttachment.depthStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.DepthStoreOp)
			depthStencilAttachment.depthClearValue = C.float(descriptor.DepthStencilAttachment.DepthClearValue)
			depthStencilAttachment.depthReadOnly = cBool(descriptor.DepthStencilAttachment.DepthReadOnly)
			depthStencilAttachment.stencilLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.StencilLoadOp)
			depthStencilAttachment.stencilStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.StencilStoreOp)
			depthStencilAttachment.stencilClearValue = C.uint32_t(descriptor.DepthStencilAttachment.StencilClearValue)
			depthStencilAttachment.stencilReadOnly = cBool(descriptor.DepthStencilAttachment.DepthReadOnly)

			desc.depthStencilAttachment = depthStencilAttachment
		}
	}

	ref := C.wgpuCommandEncoderBeginRenderPass(p.ref, &desc)
	C.wgpuDeviceAddRef(p.deviceRef)
	return &RenderPassEncoder{deviceRef: p.deviceRef, instanceRef: p.instanceRef, ref: ref}
}

func (p *CommandEncoder) ClearBuffer(buffer *Buffer, offset uint64, size uint64) error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).ClearBuffer(): ", func() {
		C.wgpuCommandEncoderClearBuffer(p.ref, buffer.ref, C.uint64_t(offset), C.uint64_t(size))
	})
}

func (p *CommandEncoder) CopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinatonOffset uint64, size uint64) error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).CopyBufferToBuffer(): ", func() {
		C.wgpuCommandEncoderCopyBufferToBuffer(p.ref, source.ref, C.uint64_t(sourceOffset), destination.ref, C.uint64_t(destinatonOffset), C.uint64_t(size))
	})
}

func (p *CommandEncoder) CopyBufferToTexture(source *ImageCopyBuffer, destination *ImageCopyTexture, copySize *Extent3D) error {
	var src C.WGPUTexelCopyBufferInfo
	if source != nil {
		if source.Buffer != nil {
			src.buffer = source.Buffer.ref
		}
		src.layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(source.Layout.Offset),
			bytesPerRow:  C.uint32_t(source.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(source.Layout.RowsPerImage),
		}
	}

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

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).CopyBufferToTexture(): ", func() {
		C.wgpuCommandEncoderCopyBufferToTexture(p.ref, &src, &dst, &cpySize)
	})
}

func (p *CommandEncoder) CopyTextureToBuffer(source *ImageCopyTexture, destination *ImageCopyBuffer, copySize *Extent3D) error {
	var src C.WGPUTexelCopyTextureInfo
	if source != nil {
		src = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUTexelCopyBufferInfo
	if destination != nil {
		if destination.Buffer != nil {
			dst.buffer = destination.Buffer.ref
		}
		dst.layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(destination.Layout.Offset),
			bytesPerRow:  C.uint32_t(destination.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(destination.Layout.RowsPerImage),
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).CopyTextureToBuffer(): ", func() {
		C.wgpuCommandEncoderCopyTextureToBuffer(p.ref, &src, &dst, &cpySize)
	})
}

func (p *CommandEncoder) CopyTextureToTexture(source *ImageCopyTexture, destination *ImageCopyTexture, copySize *Extent3D) error {
	var src C.WGPUTexelCopyTextureInfo
	if source != nil {
		src = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

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

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).CopyTextureToTexture(): ", func() {
		C.wgpuCommandEncoderCopyTextureToTexture(p.ref, &src, &dst, &cpySize)
	})
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) (*CommandBuffer, error) {
	var desc *C.WGPUCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label, freeLabel := stringViewOf(descriptor.Label)
		defer freeLabel()
		desc = &C.WGPUCommandBufferDescriptor{label: label}
	}

	var err error
	pushValidationScope(p.deviceRef)
	ref := C.wgpuCommandEncoderFinish(p.ref, desc)
	popValidationScope(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).Finish(): ", &err)
	if err != nil {
		C.wgpuCommandBufferRelease(ref)
		return nil, err
	}

	return &CommandBuffer{ref}, nil
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) error {
	label, freeLabel := stringViewOf(markerLabel)
	defer freeLabel()
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).InsertDebugMarker(): ", func() {
		C.wgpuCommandEncoderInsertDebugMarker(p.ref, label)
	})
}

func (p *CommandEncoder) PopDebugGroup() error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).PopDebugGroup(): ", func() {
		C.wgpuCommandEncoderPopDebugGroup(p.ref)
	})
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) error {
	label, freeLabel := stringViewOf(groupLabel)
	defer freeLabel()
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).PushDebugGroup(): ", func() {
		C.wgpuCommandEncoderPushDebugGroup(p.ref, label)
	})
}

func (p *CommandEncoder) ResolveQuerySet(querySet *QuerySet, firstQuery uint32, queryCount uint32, destination *Buffer, destinationOffset uint64) error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).ResolveQuerySet(): ", func() {
		C.wgpuCommandEncoderResolveQuerySet(p.ref, querySet.ref, C.uint32_t(firstQuery), C.uint32_t(queryCount), destination.ref, C.uint64_t(destinationOffset))
	})
}

func (p *CommandEncoder) WriteTimestamp(querySet *QuerySet, queryIndex uint32) error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*CommandEncoder).WriteTimestamp(): ", func() {
		C.wgpuCommandEncoderWriteTimestamp(p.ref, querySet.ref, C.uint32_t(queryIndex))
	})
}

func (p *CommandEncoder) Release() {
	C.wgpuDeviceRelease(p.deviceRef)
	C.wgpuCommandEncoderRelease(p.ref)
}
