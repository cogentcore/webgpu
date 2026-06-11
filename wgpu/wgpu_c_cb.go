//go:build !js

package wgpu

/*

#include "./lib/wgpu.h"

void gowebgpu_pop_error_scope_callback_c(WGPUPopErrorScopeStatus status, WGPUErrorType type, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_pop_error_scope_callback_go(WGPUPopErrorScopeStatus status, WGPUErrorType type, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_pop_error_scope_callback_go(status, type, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

void gowebgpu_buffer_map_callback_c(WGPUMapAsyncStatus status, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_buffer_map_callback_go(WGPUMapAsyncStatus status, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_buffer_map_callback_go(status, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_request_adapter_callback_go(WGPURequestAdapterStatus status, WGPUAdapter adapter, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_request_adapter_callback_go(status, adapter, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

void gowebgpu_request_device_callback_c(WGPURequestDeviceStatus status, WGPUDevice device, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_request_device_callback_go(WGPURequestDeviceStatus status, WGPUDevice device, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_request_device_callback_go(status, device, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

void gowebgpu_device_lost_callback_c(WGPUDevice const* device, WGPUDeviceLostReason reason, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_device_lost_callback_go(WGPUDeviceLostReason reason, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_device_lost_callback_go(reason, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

void gowebgpu_queue_work_done_callback_c(WGPUQueueWorkDoneStatus status, WGPUStringView message, void* userdata1, void* userdata2) {
  extern void gowebgpu_queue_work_done_callback_go(WGPUQueueWorkDoneStatus status, uintptr_t messageData, uintptr_t messageLen, uintptr_t userdata2);
  gowebgpu_queue_work_done_callback_go(status, (uintptr_t)message.data, (uintptr_t)message.length, (uintptr_t)userdata2);
}

*/
import "C"
