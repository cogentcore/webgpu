//go:build !js

package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView message, void* userdata1, void* userdata2);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type Instance struct {
	ref C.WGPUInstance
}

func CreateInstance(descriptor *InstanceDescriptor) *Instance {
	desc := C.WGPUInstanceDescriptor{}

	if descriptor != nil {
		instanceExtras := (*C.WGPUInstanceExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUInstanceExtras{}))))
		defer C.free(unsafe.Pointer(instanceExtras))
		*instanceExtras = C.WGPUInstanceExtras{}

		instanceExtras.chain.next = nil
		instanceExtras.chain.sType = C.WGPUSType_InstanceExtras
		instanceExtras.backends = C.WGPUInstanceBackend(descriptor.Backends)
		instanceExtras.dx12ShaderCompiler = C.WGPUDx12Compiler(descriptor.Dx12ShaderCompiler)

		if descriptor.DxcPath != "" {
			dxcPath, freeDxcPath := stringViewOf(descriptor.DxcPath)
			defer freeDxcPath()
			instanceExtras.dxcPath = dxcPath
		}

		desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(instanceExtras))
	}

	ref := C.wgpuCreateInstance(&desc)
	if ref == nil {
		panic("Failed to acquire Instance")
	}

	return &Instance{ref}
}

type SurfaceDescriptorFromWindowsHWND struct {
	Hinstance unsafe.Pointer
	Hwnd      unsafe.Pointer
}

type SurfaceDescriptorFromXcbWindow struct {
	Connection unsafe.Pointer
	Window     uint32
}

type SurfaceDescriptorFromXlibWindow struct {
	Display unsafe.Pointer
	Window  uint32
}

type SurfaceDescriptorFromMetalLayer struct {
	Layer unsafe.Pointer
}

type SurfaceDescriptorFromWaylandSurface struct {
	Display unsafe.Pointer
	Surface unsafe.Pointer
}

type SurfaceDescriptorFromAndroidNativeWindow struct {
	Window unsafe.Pointer
}

type SurfaceDescriptor struct {
	Label string

	WindowsHWND         *SurfaceDescriptorFromWindowsHWND
	XcbWindow           *SurfaceDescriptorFromXcbWindow
	XlibWindow          *SurfaceDescriptorFromXlibWindow
	MetalLayer          *SurfaceDescriptorFromMetalLayer
	WaylandSurface      *SurfaceDescriptorFromWaylandSurface
	AndroidNativeWindow *SurfaceDescriptorFromAndroidNativeWindow
}

func (p *Instance) CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	var desc C.WGPUSurfaceDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label, freeLabel := stringViewOf(descriptor.Label)
			defer freeLabel()
			desc.label = label
		}

		if descriptor.WindowsHWND != nil {
			windowsHWND := (*C.WGPUSurfaceSourceWindowsHWND)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceWindowsHWND{}))))
			defer C.free(unsafe.Pointer(windowsHWND))

			windowsHWND.chain.next = nil
			windowsHWND.chain.sType = C.WGPUSType_SurfaceSourceWindowsHWND
			windowsHWND.hinstance = descriptor.WindowsHWND.Hinstance
			windowsHWND.hwnd = descriptor.WindowsHWND.Hwnd

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(windowsHWND))
		}

		if descriptor.XcbWindow != nil {
			xcbWindow := (*C.WGPUSurfaceSourceXCBWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceXCBWindow{}))))
			defer C.free(unsafe.Pointer(xcbWindow))

			xcbWindow.chain.next = nil
			xcbWindow.chain.sType = C.WGPUSType_SurfaceSourceXCBWindow
			xcbWindow.connection = descriptor.XcbWindow.Connection
			xcbWindow.window = C.uint32_t(descriptor.XcbWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xcbWindow))
		}

		if descriptor.XlibWindow != nil {
			xlibWindow := (*C.WGPUSurfaceSourceXlibWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceXlibWindow{}))))
			defer C.free(unsafe.Pointer(xlibWindow))

			xlibWindow.chain.next = nil
			xlibWindow.chain.sType = C.WGPUSType_SurfaceSourceXlibWindow
			xlibWindow.display = descriptor.XlibWindow.Display
			xlibWindow.window = C.uint64_t(descriptor.XlibWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xlibWindow))
		}

		if descriptor.MetalLayer != nil {
			metalLayer := (*C.WGPUSurfaceSourceMetalLayer)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceMetalLayer{}))))
			defer C.free(unsafe.Pointer(metalLayer))

			metalLayer.chain.next = nil
			metalLayer.chain.sType = C.WGPUSType_SurfaceSourceMetalLayer
			metalLayer.layer = descriptor.MetalLayer.Layer

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(metalLayer))
		}

		if descriptor.WaylandSurface != nil {
			waylandSurface := (*C.WGPUSurfaceSourceWaylandSurface)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceWaylandSurface{}))))
			defer C.free(unsafe.Pointer(waylandSurface))

			waylandSurface.chain.next = nil
			waylandSurface.chain.sType = C.WGPUSType_SurfaceSourceWaylandSurface
			waylandSurface.display = descriptor.WaylandSurface.Display
			waylandSurface.surface = descriptor.WaylandSurface.Surface

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(waylandSurface))
		}

		if descriptor.AndroidNativeWindow != nil {
			androidNativeWindow := (*C.WGPUSurfaceSourceAndroidNativeWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceSourceAndroidNativeWindow{}))))
			defer C.free(unsafe.Pointer(androidNativeWindow))

			androidNativeWindow.chain.next = nil
			androidNativeWindow.chain.sType = C.WGPUSType_SurfaceSourceAndroidNativeWindow
			androidNativeWindow.window = descriptor.AndroidNativeWindow.Window

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(androidNativeWindow))
		}
	}

	ref := C.wgpuInstanceCreateSurface(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire Surface")
	}
	return &Surface{ref: ref}
}

type requestAdapterCb func(status RequestAdapterStatus, adapter *Adapter, message string)

//export gowebgpu_request_adapter_callback_go
func gowebgpu_request_adapter_callback_go(status C.WGPURequestAdapterStatus, adapter C.WGPUAdapter, messageData uintptr, messageLen uintptr, userdata2 uintptr) {
	ptr := unsafe.Pointer(userdata2)
	handle := *(*cgo.Handle)(ptr)
	defer freeCgoHandlePtr(ptr)
	defer handle.Delete()

	cb, ok := handle.Value().(requestAdapterCb)
	if ok {
		cb(RequestAdapterStatus(status), &Adapter{ref: adapter}, goStringViewFromParts(messageData, messageLen))
	}
}

func (p *Instance) RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	var opts *C.WGPURequestAdapterOptions

	if options != nil {
		opts = &C.WGPURequestAdapterOptions{}

		if options.CompatibleSurface != nil {
			opts.compatibleSurface = options.CompatibleSurface.ref
		}
		opts.powerPreference = C.WGPUPowerPreference(options.PowerPreference)
		opts.forceFallbackAdapter = cBool(options.ForceFallbackAdapter)
		opts.backendType = C.WGPUBackendType(options.BackendType)
	}

	var status RequestAdapterStatus
	var adapter *Adapter

	done := make(chan struct{})
	var cb requestAdapterCb = func(s RequestAdapterStatus, a *Adapter, _ string) {
		status = s
		adapter = a
		close(done)
	}
	handle := cgo.NewHandle(cb)
	handlePtr := cgoHandlePtr(handle)
	cbInfo := C.WGPURequestAdapterCallbackInfo{
		mode:      C.WGPUCallbackMode_WaitAnyOnly,
		callback:  C.WGPURequestAdapterCallback(C.gowebgpu_request_adapter_callback_c),
		userdata2: handlePtr,
	}
	future := C.wgpuInstanceRequestAdapter(p.ref, opts, cbInfo)
	p.waitFuture(future)
	<-done

	if status != RequestAdapterStatusSuccess {
		return nil, errors.New("failed to request adapter")
	}
	adapter.instanceRef = p.ref
	return adapter, nil
}

func (p *Instance) EnumerateAdapters(options *InstanceEnumerateAdapterOptons) []*Adapter {
	var opts *C.WGPUInstanceEnumerateAdapterOptions
	if options != nil {
		opts = &C.WGPUInstanceEnumerateAdapterOptions{
			backends: C.WGPUInstanceBackend(options.Backends),
		}
	}

	size := C.wgpuInstanceEnumerateAdapters(p.ref, opts, nil)
	if size == 0 {
		return nil
	}

	adapterRefs := make([]C.WGPUAdapter, size)
	C.wgpuInstanceEnumerateAdapters(p.ref, opts, (*C.WGPUAdapter)(unsafe.Pointer(&adapterRefs[0])))

	adapters := make([]*Adapter, size)
	for i, ref := range adapterRefs {
		adapters[i] = &Adapter{ref: ref, instanceRef: p.ref}
	}
	return adapters
}

type RegistryReport struct {
	NumAllocated        uint64
	NumKeptFromUser     uint64
	NumReleasedFromUser uint64
	ElementSize         uint64
}

type HubReport struct {
	Adapters         RegistryReport
	Devices          RegistryReport
	Queues           RegistryReport
	PipelineLayouts  RegistryReport
	ShaderModules    RegistryReport
	BindGroupLayouts RegistryReport
	BindGroups       RegistryReport
	CommandBuffers   RegistryReport
	RenderBundles    RegistryReport
	RenderPipelines  RegistryReport
	ComputePipelines RegistryReport
	PipelineCaches   RegistryReport
	QuerySets        RegistryReport
	Buffers          RegistryReport
	Textures         RegistryReport
	TextureViews     RegistryReport
	Samplers         RegistryReport
}

type GlobalReport struct {
	Surfaces RegistryReport
	Hub      HubReport
}

func (p *Instance) GenerateReport() GlobalReport {
	var r C.WGPUGlobalReport
	C.wgpuGenerateReport(p.ref, &r)

	mapRegistryReport := func(creport C.WGPURegistryReport) RegistryReport {
		return RegistryReport{
			NumAllocated:        uint64(creport.numAllocated),
			NumKeptFromUser:     uint64(creport.numKeptFromUser),
			NumReleasedFromUser: uint64(creport.numReleasedFromUser),
			ElementSize:         uint64(creport.elementSize),
		}
	}

	hub := r.hub
	return GlobalReport{
		Surfaces: mapRegistryReport(r.surfaces),
		Hub: HubReport{
			Adapters:         mapRegistryReport(hub.adapters),
			Devices:          mapRegistryReport(hub.devices),
			Queues:           mapRegistryReport(hub.queues),
			PipelineLayouts:  mapRegistryReport(hub.pipelineLayouts),
			ShaderModules:    mapRegistryReport(hub.shaderModules),
			BindGroupLayouts: mapRegistryReport(hub.bindGroupLayouts),
			BindGroups:       mapRegistryReport(hub.bindGroups),
			CommandBuffers:   mapRegistryReport(hub.commandBuffers),
			RenderBundles:    mapRegistryReport(hub.renderBundles),
			RenderPipelines:  mapRegistryReport(hub.renderPipelines),
			ComputePipelines: mapRegistryReport(hub.computePipelines),
			PipelineCaches:   mapRegistryReport(hub.pipelineCaches),
			QuerySets:        mapRegistryReport(hub.querySets),
			Buffers:          mapRegistryReport(hub.buffers),
			Textures:         mapRegistryReport(hub.textures),
			TextureViews:     mapRegistryReport(hub.textureViews),
			Samplers:         mapRegistryReport(hub.samplers),
		},
	}
}

func (p *Instance) Release() {
	C.wgpuInstanceRelease(p.ref)
}
