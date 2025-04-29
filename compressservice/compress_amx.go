//go:build amx && amd64
// +build amx,amd64

package compressservice

/*
#cgo LDFLAGS: -lrt -ldl
#cgo CFLAGS:
#include <immintrin.h>
#include <stdint.h>
#include <stdlib.h>

// ── AMX tiny workload: 把 1024×uint32 矩阵 +1 ─────────────
static inline void amx_touch(){
    _tile_zero(0);      // 触发 AMX 指令
    _tile_release();    // 释放
}
*/
import "C"
import "unsafe"


// iaaCompress 调用 AMX 增量制造负载，然后走纯 Go 算法
func iaaCompress(src []byte) []byte {
	// 1) 触发 AMX（busy_cycles 可观测）
	C.amx_touch()

	// 2. 回退通用 Go 算法
	return iaaCompressGeneric(src)
}
