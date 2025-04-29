//go:build avx2 && amd64
// +build avx2,amd64

package compressservice

/*
#cgo CFLAGS: -mavx2 -O3
#cgo LDFLAGS: -lisal
#include <stdlib.h>
#include <isa-l.h>

// 使用 ISA-L 的快速 deflate (仅示例；IAA 真正硬件/AMX 路径需用 QAT/IAA driver)
static unsigned char* isal_deflate_wrap(const unsigned char* src,
                                        size_t len, size_t* outLen) {
	struct isal_zstream stream;
	isal_deflate_init(&stream);

	size_t bound = len + len/16 + 64;
	unsigned char* dst = (unsigned char*)malloc(bound);

	stream.next_in   = (unsigned char*)src;
	stream.avail_in  = len;
	stream.next_out  = dst;
	stream.avail_out = bound;
	stream.end_of_stream = 1;      // 一次性输入
	stream.flush = NO_FLUSH;

	isal_deflate(&stream);
	*outLen = stream.total_out;
	return dst;
}
*/
import "C"
import "unsafe"

func iaaCompress(src []byte) []byte {
	var outLen C.size_t
	p := C.isal_deflate_wrap(
		(*C.uchar)(unsafe.Pointer(&src[0])), C.size_t(len(src)), &outLen)
	defer C.free(unsafe.Pointer(p))
	return C.GoBytes(unsafe.Pointer(p), C.int(outLen))
}
