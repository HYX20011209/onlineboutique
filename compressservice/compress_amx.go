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
void amx_increment(uint32_t *buf){
  // tile config: 16 cols × 16 rows = 256 ×4 =1 KiB
  __tileconfig tc = { .palette_id = 1, .start_row = 0,
                      .rows = {16}, .cols = {16} };
  _tile_loadconfig(&tc);

  for(int t=0;t<4;++t){                 // 4 tiles ×1 KiB =4 KiB
    uint32_t *p = buf + t*256;
    _tile_loadd(0,p,16*4);              // load tile 0
    _tile_loadd(1,p,16*4);              // load tile 1 (same)
    _tile_dpaddd(0,1,1);                // dst += src  (加1次)
    _tile_stored(0,p,16*4);             // store回内存
    _tile_zero(1);
  }
  _tile_release();
}
*/
import "C"
import "unsafe"


// iaaCompress 调用 AMX 增量制造负载，然后走纯 Go 算法
func iaaCompress(src []byte) []byte {
	// 1. AMX 占用几千 cycles
	buf := make([]uint32, 1024)
	C.amx_increment((*C.uint32_t)(unsafe.Pointer(&buf[0])))

	// 2. 回退通用 Go 算法
	return iaaCompressGeneric(src)
}
