//go:build js

package main

import (
	"syscall/js"
	"time"

	"github.com/cogentcore/webgpu/wgpu"
)

type window struct{}

func (w window) GetSize() (int, int) {
	vv := js.Global().Get("visualViewport")
	return vv.Get("width").Int(), vv.Get("height").Int()
}

func main() {
	s, err := InitState(&window{}, &wgpu.SurfaceDescriptor{})
	if err != nil {
		panic(err)
	}
	defer s.Destroy()
	ticker := time.NewTicker(time.Second / 30)
	for range ticker.C {
		s.Render()
	}
}
