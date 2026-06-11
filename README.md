# WebGPU (openfluke)

Go bindings for [WebGPU](https://gpuweb.github.io/gpuweb/), maintained by [openfluke](https://github.com/openfluke).

This fork tracks modern **wgpu-native v29** with Go-side bindings (no C compatibility shims). It runs natively on Vulkan, Metal, D3D12, and OpenGL ES, and supports the browser via WASM.

**Module:** `github.com/openfluke/webgpu`  
**Recommended release:** `v0.1.0` (first openfluke release; ships wgpu-native **v29.0.0.0**)

Forked from [rajveermalviya/go-webgpu](https://github.com/rajveermalviya/go-webgpu) (native) and [mokiat/wasmgpu](https://github.com/mokiat/wasmgpu) (JS/WASM). Upstream credit preserved; API and vendored binaries diverge from those projects.

## Why v0.1.0?

| Tag | What it was |
|-----|-------------|
| `v0.0.1` / `v0.0.2` | Initial openfluke fork of the old bindings (~wgpu-native v19 era) |
| **`v0.1.0`** | **Current line:** full Go migration to the v29 C API + v29 static libs on all platforms |

This is a **breaking** change vs `v0.0.x` for native (`!js`) code. Bump your `require` and retest GPU paths after upgrading.

## What's in v0.1.0

### wgpu-native v29 stack

- **Headers:** `wgpu/lib/webgpu.h`, `wgpu/lib/wgpu.h` (WebGPU C API + wgpu-native extensions)
- **Binaries:** prebuilt `libwgpu_native.a` per target under `wgpu/lib/`:

| Platform | Path |
|----------|------|
| macOS | `darwin/amd64`, `darwin/arm64` |
| Linux | `linux/amd64`, `linux/arm64` |
| Windows | `windows/amd64` (gnu), `windows/arm64` (msvc static) |
| iOS | `ios/arm64` (device), `ios/amd64` (x86_64 simulator) |
| Android | `android/arm64`, `android/arm`, `android/amd64`, `android/386` |

### Go binding changes (native)

Bindings were rewritten for the v29 API instead of using C redirect macros:

- `WGPUStringView` for labels, entry points, shader source, messages
- Async ops via `WGPUFuture` + `wgpuInstanceWaitAny` (`RequestAdapter`, `RequestDevice`, `MapAsync`, error scopes, queue work done)
- `wgpuDeviceAddRef`, `wgpuAdapterGetFeatures`, `wgpuDeviceGetLimits` with `maxImmediateSize`
- Shader chains: `WGPUShaderSourceWGSL` / `SPIRV` / `WGPUShaderSourceGLSL`
- Surfaces: `WGPUSurfaceSource*` descriptors
- Texel copy types: `WGPUTexelCopyTextureInfo`, etc.
- Push constants → pipeline `immediateDataSize` / `SetImmediates` (native extension)
- Validation error scopes implemented in Go (`v29.go`)

### Tested with

- **loom / lucy** on macOS arm64 + Metal: adapter → device → buffer → GPU forward parity (Apple M5)

Other platforms have matching v29 libs in-tree; smoke-test on each target you ship.

## Install

```bash
go get github.com/openfluke/webgpu@v0.1.0
```

Local development (e.g. from the endgame monorepo):

```go
replace github.com/openfluke/webgpu => ../webgpu
```

## Usage

Native and JS builds use the same import path; build tags select the backend:

```go
import "github.com/openfluke/webgpu/wgpu"
```

- `//go:build !js` — wgpu-native (this README)
- `//go:build js` — browser WebGPU via WASM

## Examples

| [boids](examples/boids) | [cube](examples/cube) | [triangle](examples/triangle) |
|:---:|:---:|:---:|
| ![boids](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/boids/image-msaa.png) | ![cube](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/cube/image-msaa.png) | ![triangle](https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/triangle/image-msaa.png) |

## Rebuilding native libs

Vendored artifacts came from official **wgpu-native v29** release zips. To refresh:

1. Download matching `wgpu-*-release.zip` from [wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases)
2. Copy `lib/libwgpu_native.a` into the corresponding `wgpu/lib/<platform>/<arch>/` directory
3. Ensure shared headers under `wgpu/lib/` stay in sync with that release

CI can also rebuild via [.github/workflows/build-wgpu.yml](.github/workflows/build-wgpu.yml) (`workflow_dispatch`).

## References

- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu.h](https://github.com/webgpu-native/webgpu-headers)
- [wgpu-native](https://github.com/gfx-rs/wgpu-native)

## License

See [LICENSE](LICENSE) (inherits upstream licensing).
