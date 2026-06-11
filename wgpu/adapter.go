//go:build !js

package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_request_device_callback_c(WGPURequestDeviceStatus status, WGPUDevice device, WGPUStringView message, void* userdata1, void* userdata2);
extern void gowebgpu_device_lost_callback_c(WGPUDevice const* device, WGPUDeviceLostReason reason, WGPUStringView message, void* userdata1, void* userdata2);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type Adapter struct {
	ref         C.WGPUAdapter
	instanceRef C.WGPUInstance
}

func (p *Adapter) EnumerateFeatures() []FeatureName {
	var supported C.WGPUSupportedFeatures
	C.wgpuAdapterGetFeatures(p.ref, &supported)
	defer C.wgpuSupportedFeaturesFreeMembers(supported)
	if supported.featureCount == 0 {
		return nil
	}
	features := make([]FeatureName, supported.featureCount)
	slice := unsafe.Slice((*C.WGPUFeatureName)(unsafe.Pointer(supported.features)), supported.featureCount)
	for i, f := range slice {
		features[i] = FeatureName(f)
	}
	return features
}

func limitsFromC(limits C.WGPULimits) Limits {
	return Limits{
		MaxTextureDimension1D:                     uint32(limits.maxTextureDimension1D),
		MaxTextureDimension2D:                     uint32(limits.maxTextureDimension2D),
		MaxTextureDimension3D:                     uint32(limits.maxTextureDimension3D),
		MaxTextureArrayLayers:                     uint32(limits.maxTextureArrayLayers),
		MaxBindGroups:                             uint32(limits.maxBindGroups),
		MaxBindingsPerBindGroup:                   uint32(limits.maxBindingsPerBindGroup),
		MaxDynamicUniformBuffersPerPipelineLayout: uint32(limits.maxDynamicUniformBuffersPerPipelineLayout),
		MaxDynamicStorageBuffersPerPipelineLayout: uint32(limits.maxDynamicStorageBuffersPerPipelineLayout),
		MaxSampledTexturesPerShaderStage:          uint32(limits.maxSampledTexturesPerShaderStage),
		MaxSamplersPerShaderStage:                 uint32(limits.maxSamplersPerShaderStage),
		MaxStorageBuffersPerShaderStage:           uint32(limits.maxStorageBuffersPerShaderStage),
		MaxStorageTexturesPerShaderStage:          uint32(limits.maxStorageTexturesPerShaderStage),
		MaxUniformBuffersPerShaderStage:           uint32(limits.maxUniformBuffersPerShaderStage),
		MaxUniformBufferBindingSize:               uint64(limits.maxUniformBufferBindingSize),
		MaxStorageBufferBindingSize:               uint64(limits.maxStorageBufferBindingSize),
		MinUniformBufferOffsetAlignment:           uint32(limits.minUniformBufferOffsetAlignment),
		MinStorageBufferOffsetAlignment:           uint32(limits.minStorageBufferOffsetAlignment),
		MaxVertexBuffers:                          uint32(limits.maxVertexBuffers),
		MaxBufferSize:                             uint64(limits.maxBufferSize),
		MaxVertexAttributes:                       uint32(limits.maxVertexAttributes),
		MaxVertexBufferArrayStride:                uint32(limits.maxVertexBufferArrayStride),
		MaxInterStageShaderComponents:             uint32(limits.maxInterStageShaderVariables),
		MaxInterStageShaderVariables:              uint32(limits.maxInterStageShaderVariables),
		MaxColorAttachments:                       uint32(limits.maxColorAttachments),
		MaxColorAttachmentBytesPerSample:          uint32(limits.maxColorAttachmentBytesPerSample),
		MaxComputeWorkgroupStorageSize:            uint32(limits.maxComputeWorkgroupStorageSize),
		MaxComputeInvocationsPerWorkgroup:         uint32(limits.maxComputeInvocationsPerWorkgroup),
		MaxComputeWorkgroupSizeX:                  uint32(limits.maxComputeWorkgroupSizeX),
		MaxComputeWorkgroupSizeY:                  uint32(limits.maxComputeWorkgroupSizeY),
		MaxComputeWorkgroupSizeZ:                  uint32(limits.maxComputeWorkgroupSizeZ),
		MaxComputeWorkgroupsPerDimension:          uint32(limits.maxComputeWorkgroupsPerDimension),
		MaxPushConstantSize:                       uint32(limits.maxImmediateSize),
	}
}

func (p *Adapter) GetLimits() SupportedLimits {
	var limits C.WGPULimits
	C.wgpuAdapterGetLimits(p.ref, &limits)
	return SupportedLimits{limitsFromC(limits)}
}

func (p *Adapter) GetInfo() AdapterInfo {
	var info C.WGPUAdapterInfo
	C.wgpuAdapterGetInfo(p.ref, &info)
	defer C.wgpuAdapterInfoFreeMembers(info)

	return AdapterInfo{
		VendorId:          uint32(info.vendorID),
		VendorName:        goStringView(info.vendor),
		Architecture:      goStringView(info.architecture),
		DeviceId:          uint32(info.deviceID),
		Name:              goStringView(info.device),
		DriverDescription: goStringView(info.description),
		AdapterType:       AdapterType(info.adapterType),
		BackendType:       BackendType(info.backendType),
	}
}

func (p *Adapter) HasFeature(feature FeatureName) bool {
	hasFeature := C.wgpuAdapterHasFeature(p.ref, C.WGPUFeatureName(feature))
	return goBool(hasFeature)
}

type requestDeviceCb func(status RequestDeviceStatus, device *Device, message string)

//export gowebgpu_request_device_callback_go
func gowebgpu_request_device_callback_go(status C.WGPURequestDeviceStatus, device C.WGPUDevice, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()

	cb, ok := handle.Value().(requestDeviceCb)
	if ok {
		cb(RequestDeviceStatus(status), &Device{ref: device}, goStringViewFromParts(messageData, messageLen))
	}
}

//export gowebgpu_device_lost_callback_go
func gowebgpu_device_lost_callback_go(reason C.WGPUDeviceLostReason, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()

	cb, ok := handle.Value().(DeviceLostCallback)
	if ok {
		cb(DeviceLostReason(reason), goStringViewFromParts(messageData, messageLen))
	}
}

func requiredLimitsFromGo(l Limits) C.WGPULimits {
	return C.WGPULimits{
		maxTextureDimension1D:                     C.uint32_t(l.MaxTextureDimension1D),
		maxTextureDimension2D:                     C.uint32_t(l.MaxTextureDimension2D),
		maxTextureDimension3D:                     C.uint32_t(l.MaxTextureDimension3D),
		maxTextureArrayLayers:                     C.uint32_t(l.MaxTextureArrayLayers),
		maxBindGroups:                             C.uint32_t(l.MaxBindGroups),
		maxBindingsPerBindGroup:                   C.uint32_t(l.MaxBindingsPerBindGroup),
		maxDynamicUniformBuffersPerPipelineLayout: C.uint32_t(l.MaxDynamicUniformBuffersPerPipelineLayout),
		maxDynamicStorageBuffersPerPipelineLayout: C.uint32_t(l.MaxDynamicStorageBuffersPerPipelineLayout),
		maxSampledTexturesPerShaderStage:          C.uint32_t(l.MaxSampledTexturesPerShaderStage),
		maxSamplersPerShaderStage:                 C.uint32_t(l.MaxSamplersPerShaderStage),
		maxStorageBuffersPerShaderStage:           C.uint32_t(l.MaxStorageBuffersPerShaderStage),
		maxStorageTexturesPerShaderStage:          C.uint32_t(l.MaxStorageTexturesPerShaderStage),
		maxUniformBuffersPerShaderStage:           C.uint32_t(l.MaxUniformBuffersPerShaderStage),
		maxUniformBufferBindingSize:               C.uint64_t(l.MaxUniformBufferBindingSize),
		maxStorageBufferBindingSize:               C.uint64_t(l.MaxStorageBufferBindingSize),
		minUniformBufferOffsetAlignment:           C.uint32_t(l.MinUniformBufferOffsetAlignment),
		minStorageBufferOffsetAlignment:           C.uint32_t(l.MinStorageBufferOffsetAlignment),
		maxVertexBuffers:                          C.uint32_t(l.MaxVertexBuffers),
		maxBufferSize:                             C.uint64_t(l.MaxBufferSize),
		maxVertexAttributes:                       C.uint32_t(l.MaxVertexAttributes),
		maxVertexBufferArrayStride:                C.uint32_t(l.MaxVertexBufferArrayStride),
		maxInterStageShaderVariables:              C.uint32_t(l.MaxInterStageShaderVariables),
		maxColorAttachments:                       C.uint32_t(l.MaxColorAttachments),
		maxColorAttachmentBytesPerSample:          C.uint32_t(l.MaxColorAttachmentBytesPerSample),
		maxComputeWorkgroupStorageSize:            C.uint32_t(l.MaxComputeWorkgroupStorageSize),
		maxComputeInvocationsPerWorkgroup:         C.uint32_t(l.MaxComputeInvocationsPerWorkgroup),
		maxComputeWorkgroupSizeX:                  C.uint32_t(l.MaxComputeWorkgroupSizeX),
		maxComputeWorkgroupSizeY:                  C.uint32_t(l.MaxComputeWorkgroupSizeY),
		maxComputeWorkgroupSizeZ:                  C.uint32_t(l.MaxComputeWorkgroupSizeZ),
		maxComputeWorkgroupsPerDimension:          C.uint32_t(l.MaxComputeWorkgroupsPerDimension),
		maxImmediateSize:                          C.uint32_t(l.MaxPushConstantSize),
	}
}

func (p *Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {
	var desc C.WGPUDeviceDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label, freeLabel := stringViewOf(descriptor.Label)
			defer freeLabel()
			desc.label = label
		}

		requiredFeatureCount := len(descriptor.RequiredFeatures)
		if requiredFeatureCount != 0 {
			requiredFeatures := C.malloc(C.size_t(requiredFeatureCount) * C.size_t(unsafe.Sizeof(C.WGPUFeatureName(0))))
			defer C.free(requiredFeatures)

			requiredFeaturesSlice := unsafe.Slice((*FeatureName)(requiredFeatures), requiredFeatureCount)
			copy(requiredFeaturesSlice, descriptor.RequiredFeatures)

			desc.requiredFeatures = (*C.WGPUFeatureName)(requiredFeatures)
			desc.requiredFeatureCount = C.size_t(requiredFeatureCount)
		}

		if descriptor.RequiredLimits != nil {
			limits := (*C.WGPULimits)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPULimits{}))))
			defer C.free(unsafe.Pointer(limits))
			*limits = requiredLimitsFromGo(descriptor.RequiredLimits.Limits)
			desc.requiredLimits = limits
		}

		if descriptor.DeviceLostCallback != nil {
			lostHandle := cgo.NewHandle(descriptor.DeviceLostCallback)
			lostHandlePtr := cgoHandlePtr(lostHandle)
			desc.deviceLostCallbackInfo = C.WGPUDeviceLostCallbackInfo{
				mode:      C.WGPUCallbackMode_AllowSpontaneous,
				callback:  C.WGPUDeviceLostCallback(C.gowebgpu_device_lost_callback_c),
				userdata2: lostHandlePtr,
			}
		}

		if descriptor.TracePath != "" {
			deviceExtras := (*C.WGPUDeviceExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUDeviceExtras{}))))
			defer C.free(unsafe.Pointer(deviceExtras))
			*deviceExtras = C.WGPUDeviceExtras{}

			deviceExtras.chain.next = nil
			deviceExtras.chain.sType = C.WGPUSType_DeviceExtras

			tracePath, freeTracePath := stringViewOf(descriptor.TracePath)
			defer freeTracePath()
			deviceExtras.tracePath = tracePath

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(deviceExtras))
		}
	}

	var status RequestDeviceStatus
	var device *Device

	done := make(chan struct{})
	var cb requestDeviceCb = func(s RequestDeviceStatus, d *Device, _ string) {
		status = s
		device = d
		close(done)
	}
	handle := cgo.NewHandle(cb)
	handlePtr := cgoHandlePtr(handle)
	cbInfo := C.WGPURequestDeviceCallbackInfo{
		mode:      C.WGPUCallbackMode_WaitAnyOnly,
		callback:  C.WGPURequestDeviceCallback(C.gowebgpu_request_device_callback_c),
		userdata2: handlePtr,
	}
	future := C.wgpuAdapterRequestDevice(p.ref, &desc, cbInfo)
	(&Instance{ref: p.instanceRef}).waitFuture(future)
	<-done

	if status != RequestDeviceStatusSuccess {
		return nil, errors.New("failed to request device")
	}

	device.instanceRef = p.instanceRef
	return device, nil
}

func (p *Adapter) Release() {
	C.wgpuAdapterRelease(p.ref)
}
