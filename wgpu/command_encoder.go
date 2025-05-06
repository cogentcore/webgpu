//go:build !js

package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_error_callback_c(enum WGPUPopErrorScopeStatus status, WGPUErrorType type, WGPUStringView message, void * userdata, void * userdata2);

static inline void gowebgpu_command_encoder_clear_buffer(WGPUCommandEncoder commandEncoder, WGPUBuffer buffer, uint64_t offset, uint64_t size, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderClearBuffer(commandEncoder, buffer, offset, size);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_copy_buffer_to_buffer(WGPUCommandEncoder commandEncoder, WGPUBuffer source, uint64_t sourceOffset, WGPUBuffer destination, uint64_t destinationOffset, uint64_t size, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderCopyBufferToBuffer(commandEncoder, source, sourceOffset, destination, destinationOffset, size);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_copy_buffer_to_texture(WGPUCommandEncoder commandEncoder, WGPUTexelCopyBufferInfo const * source, WGPUTexelCopyTextureInfo const * destination, WGPUExtent3D const * copySize, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderCopyBufferToTexture(commandEncoder, source, destination, copySize);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_copy_texture_to_buffer(WGPUCommandEncoder commandEncoder, WGPUTexelCopyTextureInfo const * source, WGPUTexelCopyBufferInfo const * destination, WGPUExtent3D const * copySize, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderCopyTextureToBuffer(commandEncoder, source, destination, copySize);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_copy_texture_to_texture(WGPUCommandEncoder commandEncoder, WGPUTexelCopyTextureInfo const * source, WGPUTexelCopyTextureInfo const * destination, WGPUExtent3D const * copySize, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderCopyTextureToTexture(commandEncoder, source, destination, copySize);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline WGPUCommandBuffer gowebgpu_command_encoder_finish(WGPUCommandEncoder commandEncoder, WGPUCommandBufferDescriptor const * descriptor, WGPUDevice device, void * error_userdata) {
	WGPUCommandBuffer ref = NULL;
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	ref = wgpuCommandEncoderFinish(commandEncoder, descriptor);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);

	return ref;
}

static inline void gowebgpu_command_encoder_insert_debug_marker(WGPUCommandEncoder commandEncoder, char const * markerLabel, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderInsertDebugMarker(commandEncoder, (WGPUStringView) {markerLabel, WGPU_STRLEN});

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_pop_debug_group(WGPUCommandEncoder commandEncoder, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderPopDebugGroup(commandEncoder);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_push_debug_group(WGPUCommandEncoder commandEncoder, char const * groupLabel, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderPushDebugGroup(commandEncoder, (WGPUStringView) {groupLabel, WGPU_STRLEN});

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_resolve_query_set(WGPUCommandEncoder commandEncoder, WGPUQuerySet querySet, uint32_t firstQuery, uint32_t queryCount, WGPUBuffer destination, uint64_t destinationOffset, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderResolveQuerySet(commandEncoder, querySet, firstQuery, queryCount, destination, destinationOffset);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_write_timestamp(WGPUCommandEncoder commandEncoder, WGPUQuerySet querySet, uint32_t queryIndex, WGPUDevice device, void * error_userdata) {
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	wgpuCommandEncoderWriteTimestamp(commandEncoder, querySet, queryIndex);

	WGPUPopErrorScopeCallbackInfo const err_cb = {
		.callback = gowebgpu_error_callback_c,
		.userdata1 = error_userdata,
	};

	wgpuDevicePopErrorScope(device, err_cb);
}

static inline void gowebgpu_command_encoder_release(WGPUCommandEncoder commandEncoder, WGPUDevice device) {
	wgpuDeviceRelease(device);
	wgpuCommandEncoderRelease(commandEncoder);
}

*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type CommandEncoder struct {
	deviceRef C.WGPUDevice
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
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc = &C.WGPUComputePassDescriptor{
			label: C.WGPUStringView{data: label, length: C.WGPU_STRLEN},
		}
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, desc)
	if ref == nil {
		panic("Failed to acquire ComputePassEncoder")
	}

	C.wgpuDeviceAddRef(p.deviceRef)
	return &ComputePassEncoder{deviceRef: p.deviceRef, ref: ref}
}

func (p *CommandEncoder) BeginRenderPass(descriptor *RenderPassDescriptor) *RenderPassEncoder {
	var desc C.WGPURenderPassDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label.data = label
			desc.label.length = C.WGPU_STRLEN
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
	return &RenderPassEncoder{deviceRef: p.deviceRef, ref: ref}
}

func (p *CommandEncoder) ClearBuffer(buffer *Buffer, offset uint64, size uint64) (err error) {
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).ClearBuffer(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_clear_buffer(
		p.ref,
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) CopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinatonOffset uint64, size uint64) (err error) {
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).CopyBufferToBuffer(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_copy_buffer_to_buffer(
		p.ref,
		source.ref,
		C.uint64_t(sourceOffset),
		destination.ref,
		C.uint64_t(destinatonOffset),
		C.uint64_t(size),
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) CopyBufferToTexture(source *TexelCopyBufferInfo, destination *TexelCopyTextureInfo, copySize *Extent3D) (err error) {
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

	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).CopyBufferToTexture(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_copy_buffer_to_texture(
		p.ref,
		&src,
		&dst,
		&cpySize,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) CopyTextureToBuffer(source *TexelCopyTextureInfo, destination *TexelCopyBufferInfo, copySize *Extent3D) (err error) {
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

	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).CopyTextureToBuffer(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_copy_texture_to_buffer(
		p.ref,
		&src,
		&dst,
		&cpySize,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) CopyTextureToTexture(source *TexelCopyTextureInfo, destination *TexelCopyTextureInfo, copySize *Extent3D) (err error) {
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

	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).CopyTextureToTexture(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_copy_texture_to_texture(
		p.ref,
		&src,
		&dst,
		&cpySize,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) (*CommandBuffer, error) {
	var desc *C.WGPUCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc = &C.WGPUCommandBufferDescriptor{
			label: C.WGPUStringView{data: label, length: C.WGPU_STRLEN},
		}
	}

	var err error = nil
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).Finish(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	ref := C.gowebgpu_command_encoder_finish(
		p.ref,
		desc,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	if err != nil {
		C.wgpuCommandBufferRelease(ref)
		return nil, err
	}

	return &CommandBuffer{ref}, nil
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) (err error) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).InsertDebugMarker(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_insert_debug_marker(
		p.ref,
		markerLabelStr,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) PopDebugGroup() (err error) {
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).PopDebugGroup(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_pop_debug_group(
		p.ref,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) (err error) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).PushDebugGroup(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_push_debug_group(
		p.ref,
		groupLabelStr,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) ResolveQuerySet(querySet *QuerySet, firstQuery uint32, queryCount uint32, destination *Buffer, destinationOffset uint64) (err error) {
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).ResolveQuerySet(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_resolve_query_set(
		p.ref,
		querySet.ref,
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
		destination.ref,
		C.uint64_t(destinationOffset),
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) WriteTimestamp(querySet *QuerySet, queryIndex uint32) (err error) {
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*CommandEncoder).WriteTimestamp(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	C.gowebgpu_command_encoder_write_timestamp(
		p.ref,
		querySet.ref,
		C.uint32_t(queryIndex),
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	return
}

func (p *CommandEncoder) Release() {
	C.gowebgpu_command_encoder_release(p.ref, p.deviceRef)
}
