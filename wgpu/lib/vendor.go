package lib

// NOTE: this is for `go mod vendor` and does not affect binary size.

import (
	_ "github.com/cogentcore/webgpu/wgpu/lib/android/386"
	_ "github.com/cogentcore/webgpu/wgpu/lib/android/amd64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/android/arm"
	_ "github.com/cogentcore/webgpu/wgpu/lib/android/arm64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/darwin/amd64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/darwin/arm64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/ios/amd64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/ios/arm64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/linux/amd64"
	_ "github.com/cogentcore/webgpu/wgpu/lib/windows/amd64"
)
