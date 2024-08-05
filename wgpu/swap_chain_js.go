//go:build js

package wgpu

import (
	"strings"
	"syscall/js"
)

// SwapChain as described:
// https://gpuweb.github.io/gpuweb/#gpucanvascontext
// (CanvasContext is the closest equivalent to SwapChain in js)
type SwapChain struct {
	jsValue js.Value
}

func (g SwapChain) GetCurrentTextureView() (*TextureView, error) {
	texture := g.jsValue.Call("getCurrentTexture")
	// We can just use the properties of the texture as the descriptor.
	descriptor := map[string]any{
		"dimension":       texture.Get("dimension"),
		"mipLevelCount":   texture.Get("mipLevelCount"),
		"arrayLayerCount": texture.Get("depthOrArrayLayers"),
	}
	// We ensure the format is srgb, which must be done here
	// (see https://gpuweb.github.io/gpuweb/#canvas-configuration).
	format := texture.Get("format").String()
	if !strings.HasSuffix(format, "-srgb") {
		format += "-srgb" // TODO(kai): we probably shouldn't always do this
	}
	descriptor["format"] = format
	return &TextureView{jsValue: texture.Call("createView", descriptor)}, nil
}

func (g SwapChain) Present() {} // no-op

func (g SwapChain) Release() {} // no-op
