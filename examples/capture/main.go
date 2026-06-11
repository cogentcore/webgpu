package main

import (
	"image"
	"image/png"
	"os"
	"unsafe"

	"github.com/openfluke/webgpu/wgpu"
)

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	switch os.Getenv("WGPU_LOG_LEVEL") {
	case "OFF":
		wgpu.SetLogLevel(wgpu.LogLevelOff)
	case "ERROR":
		wgpu.SetLogLevel(wgpu.LogLevelError)
	case "WARN":
		wgpu.SetLogLevel(wgpu.LogLevelWarn)
	case "INFO":
		wgpu.SetLogLevel(wgpu.LogLevelInfo)
	case "DEBUG":
		wgpu.SetLogLevel(wgpu.LogLevelDebug)
	case "TRACE":
		wgpu.SetLogLevel(wgpu.LogLevelTrace)
	}
}

type BufferDimensions struct {
	width               uint64
	height              uint64
	unpaddedBytesPerRow uint64
	paddedBytesPerRow   uint64
}

func newBufferDimensions(width uint64, height uint64) BufferDimensions {
	const bytesPerPixel = unsafe.Sizeof(uint32(0))
	unpaddedBytesPerRow := width * uint64(bytesPerPixel)
	align := uint64(wgpu.CopyBytesPerRowAlignment)
	paddedBytesPerRowPadding := (align - unpaddedBytesPerRow%align) % align
	paddedBytesPerRow := unpaddedBytesPerRow + uint64(paddedBytesPerRowPadding)
	return BufferDimensions{
		width,
		height,
		unpaddedBytesPerRow,
		paddedBytesPerRow,
	}
}

func main() {
	width := 100
	height := 200

	instance := wgpu.CreateInstance(nil)
	defer instance.Release()

	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
	})
	if err != nil {
		panic(err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		panic(err)
	}
	defer device.Release()
	queue := device.GetQueue()
	defer queue.Release()

	bufferDimensions := newBufferDimensions(uint64(width), uint64(height))

	bufferSize := bufferDimensions.paddedBytesPerRow * bufferDimensions.height
	// The output buffer lets us retrieve the data as an array
	outputBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  bufferSize,
		Usage: wgpu.BufferUsageMapRead | wgpu.BufferUsageCopyDst,
	})
	if err != nil {
		panic(err)
	}
	defer outputBuffer.Release()

	textureExtent := wgpu.Extent3D{
		Width:              uint32(bufferDimensions.width),
		Height:             uint32(bufferDimensions.height),
		DepthOrArrayLayers: 1,
	}

	// The render pipeline renders data into this texture
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension2D,
		Format:        wgpu.TextureFormatRGBA8UnormSrgb,
		Usage:         wgpu.TextureUsageRenderAttachment | wgpu.TextureUsageCopySrc,
	})
	if err != nil {
		panic(err)
	}
	defer texture.Release()

	// Set the background to be red
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		panic(err)
	}
	defer encoder.Release()

	textureView, err := texture.CreateView(nil)
	if err != nil {
		panic(err)
	}
	defer textureView.Release()

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:       textureView,
			LoadOp:     wgpu.LoadOpClear,
			StoreOp:    wgpu.StoreOpStore,
			ClearValue: wgpu.ColorRed,
		}},
	})
	defer renderPass.Release()
	renderPass.End()

	// Copy the data from the texture to the buffer
	encoder.CopyTextureToBuffer(
		texture.AsImageCopy(),
		&wgpu.ImageCopyBuffer{
			Buffer: outputBuffer,
			Layout: wgpu.TextureDataLayout{
				Offset:       0,
				BytesPerRow:  uint32(bufferDimensions.paddedBytesPerRow),
				RowsPerImage: wgpu.CopyStrideUndefined,
			},
		},
		&textureExtent,
	)

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		panic(err)
	}
	defer cmdBuffer.Release()

	queue.Submit(cmdBuffer)

	outputBuffer.MapAsync(wgpu.MapModeRead, 0, bufferSize, func(status wgpu.BufferMapAsyncStatus) {
		if status != wgpu.BufferMapAsyncStatusSuccess {
			panic("failed to map buffer")
		}
	})
	defer outputBuffer.Unmap()

	device.Poll(true, nil)

	data := outputBuffer.GetMappedRange(0, uint(bufferSize))

	// Code to print the image data on JS, which does not support os.Create:
	// u := js.Global().Get("Uint8Array").New(len(data))
	// js.CopyBytesToJS(u, data)
	// js.Global().Get("console").Call("log", u)
	// return

	// Save png
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	imageEncoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = imageEncoder.Encode(f, &image.NRGBA{
		Pix:    data,
		Stride: int(bufferDimensions.paddedBytesPerRow),
		Rect:   image.Rect(0, 0, width, height),
	})
	if err != nil {
		panic(err)
	}
}
