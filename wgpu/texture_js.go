//go:build js

package wgpu

import (
	"fmt"
	"syscall/js"
)

// TextureView as described:
// https://gpuweb.github.io/gpuweb/#gputextureview
type TextureView struct {
	jsValue js.Value
}

func (g TextureView) toJS() any {
	return g.jsValue
}

func (g TextureView) Release() {} // no-op

// Texture as described:
// https://gpuweb.github.io/gpuweb/#gputexture
type Texture struct {
	jsValue js.Value
}

func (g Texture) toJS() any {
	return g.jsValue
}

// GetFormat as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-format
func (g Texture) GetFormat() TextureFormat {
	jsFormat := g.jsValue.Get("format").String()

	for tf := TextureFormatUndefined + 1; tf.String() != ""; tf++ {
		if tf.String() == jsFormat {
			return tf
		}
	}

	panic(fmt.Sprintf("unknown texture format %q", jsFormat))
}

// GetDepthOrArrayLayers as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-depthorarraylayers
func (g Texture) GetDepthOrArrayLayers() uint32 {
	return uint32(g.jsValue.Get("depthOrArrayLayers").Int())
}

// GetMipLevelCount as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-miplevelcount
func (g Texture) GetMipLevelCount() uint32 {
	return uint32(g.jsValue.Get("mipLevelCount").Int())
}

// GetWidth as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-width
func (g Texture) GetWidth() uint32 {
	return uint32(g.jsValue.Get("width").Int())
}

// GetWidth as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-height
func (g Texture) GetHeight() uint32 {
	return uint32(g.jsValue.Get("height").Int())
}

// GetWidth as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-samplecount
func (g Texture) GetSampleCount() uint32 {
	return uint32(g.jsValue.Get("sampleCount").Int())
}

// CreateView as described:
// https://gpuweb.github.io/gpuweb/#dom-gputexture-createview
func (g Texture) CreateView(descriptor *TextureViewDescriptor) (*TextureView, error) {
	jsView := g.jsValue.Call("createView", pointerToJS(descriptor))
	return &TextureView{
		jsValue: jsView,
	}, nil
}

func (g Texture) Present() {} // no-op

func (g Texture) Release() {} // no-op
