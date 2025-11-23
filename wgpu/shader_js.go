//go:build js

package wgpu

import "syscall/js"

// ShaderModuleDescriptor as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpushadermoduledescriptor
type ShaderModuleDescriptor struct {
	Label      string
	WGSLSource *ShaderSourceWGSL
}

func (g ShaderModuleDescriptor) toJS() any {
	return map[string]any{
		"code": g.WGSLSource.Code,
	}
}

// ShaderModule as described:
// https://gpuweb.github.io/gpuweb/#gpushadermodule
type ShaderModule struct {
	jsValue js.Value
}

func (g ShaderModule) toJS() any {
	return g.jsValue
}

func (g ShaderModule) Release() {} // no-op
