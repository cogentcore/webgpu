//go:build js

package wgpu

import (
	"syscall/js"
)

// ComputePipelineDescriptor as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpucomputepipelinedescriptor
type ComputePipelineDescriptor struct {
	Label   string
	Layout  *PipelineLayout
	Compute ProgrammableStageDescriptor
}

func (g ComputePipelineDescriptor) toJS() any {
	result := make(map[string]any)
	result["label"] = g.Label
	result["layout"] = pointerToJS(g.Layout)
	result["compute"] = g.Compute.toJS()
	return result
}

// ComputePipeline as described:
// https://gpuweb.github.io/gpuweb/#gpucomputepipeline
type ComputePipeline struct {
	jsValue js.Value
}

func (g ComputePipeline) toJS() any {
	return g.jsValue
}

func (g ComputePipeline) Release() {} // no-op

func (g ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	return &BindGroupLayout{
		jsValue: g.jsValue.Call("getBindGroupLayout", groupIndex),
	}
}
