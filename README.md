# WebGPU

Current upstream version: v25.0.2.1

Go bindings for WebGPU, a cross-platform, safe graphics API. It runs natively using [wgpu-native](https://github.com/gfx-rs/wgpu-native) on Vulkan, Metal, D3D12, and OpenGL ES based on https://github.com/rajveermalviya/go-webgpu. It also comes with web (JS) support based on https://github.com/mokiat/wasmgpu.

For more information, see:

- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu-native](https://github.com/webgpu-native/webgpu-headers)

The included static libraries are built via [GitHub Actions](.github/workflows/build-wgpu.yml).

## Examples

|[boids][b]|[cube][c]|[triangle][t]|
:-:|:-:|:-:
| [![b-i]][b] | [![c-i]][c] | [![t-i]][t] |

[b-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/boids/image-msaa.png
[b]: examples/boids
[c-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/cube/image-msaa.png
[c]: examples/cube
[t-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/triangle/image-msaa.png
[t]: examples/triangle
