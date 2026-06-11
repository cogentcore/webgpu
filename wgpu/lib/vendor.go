package lib

// NOTE: these imports allow `go mod vendor` to include the
// static libraries, and they do not affect binary sizes.

import (
	_ "github.com/openfluke/webgpu/wgpu/lib/android/386"
	_ "github.com/openfluke/webgpu/wgpu/lib/android/amd64"
	_ "github.com/openfluke/webgpu/wgpu/lib/android/arm"
	_ "github.com/openfluke/webgpu/wgpu/lib/android/arm64"
	_ "github.com/openfluke/webgpu/wgpu/lib/darwin/amd64"
	_ "github.com/openfluke/webgpu/wgpu/lib/darwin/arm64"
	_ "github.com/openfluke/webgpu/wgpu/lib/ios/amd64"
	_ "github.com/openfluke/webgpu/wgpu/lib/ios/arm64"
	_ "github.com/openfluke/webgpu/wgpu/lib/linux/amd64"
	_ "github.com/openfluke/webgpu/wgpu/lib/windows/amd64"
)
