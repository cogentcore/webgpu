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

type ComputePassEncoder struct {
	deviceRef   C.WGPUDevice
	instanceRef C.WGPUInstance
	ref         C.WGPUComputePassEncoder
}

func (p *ComputePassEncoder) BeginPipelineStatisticsQuery(querySet *QuerySet, queryIndex uint32) {
	C.wgpuComputePassEncoderBeginPipelineStatisticsQuery(p.ref, querySet.ref, C.uint32_t(queryIndex))
}

func (p *ComputePassEncoder) DispatchWorkgroups(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	C.wgpuComputePassEncoderDispatchWorkgroups(p.ref, C.uint32_t(workgroupCountX), C.uint32_t(workgroupCountY), C.uint32_t(workgroupCountZ))
}

func (p *ComputePassEncoder) DispatchWorkgroupsIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuComputePassEncoderDispatchWorkgroupsIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
}

func (p *ComputePassEncoder) End() error {
	return withDeviceValidation(p.deviceRef, p.instanceRef, "wgpu.(*ComputePassEncoder).End(): ", func() {
		C.wgpuComputePassEncoderEnd(p.ref)
	})
}

func (p *ComputePassEncoder) EndPipelineStatisticsQuery() {
	C.wgpuComputePassEncoderEndPipelineStatisticsQuery(p.ref)
}

func (p *ComputePassEncoder) InsertDebugMarker(markerLabel string) {
	label, freeLabel := stringViewOf(markerLabel)
	defer freeLabel()
	C.wgpuComputePassEncoderInsertDebugMarker(p.ref, label)
}

func (p *ComputePassEncoder) PopDebugGroup() {
	C.wgpuComputePassEncoderPopDebugGroup(p.ref)
}

func (p *ComputePassEncoder) PushDebugGroup(groupLabel string) {
	label, freeLabel := stringViewOf(groupLabel)
	defer freeLabel()
	C.wgpuComputePassEncoderPushDebugGroup(p.ref, label)
}

func (p *ComputePassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		C.wgpuComputePassEncoderSetBindGroup(p.ref, C.uint32_t(groupIndex), group.ref, 0, nil)
	} else {
		C.wgpuComputePassEncoderSetBindGroup(
			p.ref, C.uint32_t(groupIndex), group.ref,
			C.size_t(dynamicOffsetCount), (*C.uint32_t)(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
}

func (p *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	C.wgpuComputePassEncoderSetPipeline(p.ref, pipeline.ref)
}

func (p *ComputePassEncoder) Release() {
	C.wgpuDeviceRelease(p.deviceRef)
	C.wgpuComputePassEncoderRelease(p.ref)
}
