//go:build cuda
// +build cuda

package compressservice

/*
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -L${SRCDIR} -lcuda_xor -lcudart
#cgo CFLAGS:  -I/usr/local/cuda/include
#include <cuda_runtime.h>
#include <stdlib.h>

// kernel 原型声明
extern size_t dummy_gpu_xor(const void *src, size_t len, char **dst);
*/
import "C"
import "unsafe"

//go:generate nvcc -O3 -shared -Xcompiler -fPIC -o libcuda_xor.so ./cuda_xor.cu

func iaaCompress(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}
	var dst *C.char
	n := C.dummy_gpu_xor(unsafe.Pointer(&src[0]), C.size_t(len(src)), &dst)
	defer C.free(unsafe.Pointer(dst))
	return C.GoBytes(unsafe.Pointer(dst), C.int(n))
}
