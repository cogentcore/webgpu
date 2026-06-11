//go:build js

package wgpu

import (
	"fmt"
	"syscall/js"
)

// Adapter as described:
// https://gpuweb.github.io/gpuweb/#gpuadapter
type Adapter struct {
	jsValue js.Value
}

func (g Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {

	device := js.Global().Get("webgpuDevice")
	if device.IsUndefined() {
		return nil, fmt.Errorf("WebGPU device not pre-initialized. Call setupWebGPU() in JavaScript first")
	}

	if !device.Truthy() {
		return nil, fmt.Errorf("no WebGPU device available")
	}

	// Also get the queue since Device will need it
	//queue := js.Global().Get("webgpuQueue")

	return &Device{
		jsValue: device,
	}, nil
}

func (g Adapter) GetInfo() AdapterInfo {
	return AdapterInfo{} // TODO(kai): implement?
}

func (g Adapter) GetLimits() SupportedLimits {
	return SupportedLimits{limitsFromJS(g.jsValue.Get("limits"))}
}

func (g Adapter) Release() {} // no-op
