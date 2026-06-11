//go:build js

package wgpu

import (
	"fmt"
	"log"
	"syscall/js"
)

// Instance as described:
// https://gpuweb.github.io/gpuweb/#gpu-interface
// (Instance is called GPU in js)
type Instance struct {
	jsValue js.Value
}

func CreateInstance(descriptor *InstanceDescriptor) *Instance {
	gpu := js.Global().Get("navigator").Get("gpu")
	if !gpu.Truthy() {
		log.Println("WebGPU not supported")
		return nil
	}
	return &Instance{jsValue: gpu}
}

func (g Instance) RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	// Simple WASM fix - require pre-initialization
	preAdapter := js.Global().Get("webgpuAdapter")
	if preAdapter.IsUndefined() {
		return nil, fmt.Errorf("WebGPU adapter not pre-initialized. Call setupWebGPU() in JavaScript first")
	}

	if !preAdapter.Truthy() {
		return nil, fmt.Errorf("no WebGPU adapter available")
	}

	return &Adapter{jsValue: preAdapter}, nil
}

func (g Instance) EnumerateAdapters(options *InstanceEnumerateAdapterOptons) []*Adapter {
	a, err := g.RequestAdapter(&RequestAdapterOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return []*Adapter{a}
}

// SurfaceDescriptor must contain a valid HTML canvas element on web.
type SurfaceDescriptor struct {
	// Canvas must be specified.
	Canvas js.Value

	Label string
}

func (g Instance) CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	if descriptor.Canvas.IsUndefined() {
		panic("wgpu.Instance.CreateSurface: descriptor.Canvas must be specified")
	}
	jsContext := descriptor.Canvas.Call("getContext", "webgpu")
	return &Surface{jsContext}
}

func (g Instance) GenerateReport() any { return nil } // no-op

func (g Instance) Release() {} // no-op
