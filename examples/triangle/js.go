//go:build js

package main

import (
	"syscall/js"
	"time"

	"github.com/openfluke/webgpu/wgpu"
)

type window struct{}

func (w window) GetSize() (int, int) {
	vv := js.Global().Get("visualViewport")
	return vv.Get("width").Int(), vv.Get("height").Int()
}

func main() {
	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	w := &window{}
	width, height := w.GetSize()
	canvas.Set("width", width)
	canvas.Set("height", height)
	canvas.Set("style", "width:100vw; height:100vh")

	s, err := InitState(w, &wgpu.SurfaceDescriptor{Canvas: canvas})
	if err != nil {
		panic(err)
	}
	defer s.Destroy()
	ticker := time.NewTicker(time.Second / 30)
	for range ticker.C {
		s.Render()
	}
}
