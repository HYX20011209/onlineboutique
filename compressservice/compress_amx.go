//go:build amx && amd64
// +build amx,amd64

package compressservice

/*
#cgo LDFLAGS: -lrt -ldl
#include <immintrin.h>

// 轻量 AMX：清零 tile-0
static inline void amx_touch() {
    _tile_zero(0);
    _tile_release();
}
*/
import "C"
import (
	"log"
	"runtime"
	"syscall"
)

// Linux x86_64 constants
const (
	_ARCH_REQ_XCOMP_PERM = 0x1022 // from <asm/prctl.h>
	_XFEATURE_XTILEDATA  = 18     // tile data component
)

// tryEnableAMX issues arch_prctl to get AMX permission.
// Returns true on success, false otherwise.
func tryEnableAMX() bool {
	_, _, errno := syscall.RawSyscall(syscall.SYS_ARCH_PRCTL,
		uintptr(_ARCH_REQ_XCOMP_PERM),
		uintptr(_XFEATURE_XTILEDATA), 0)
	return errno == 0
}

func iaaCompress(src []byte) []byte {
	// 确保只请求一次
	once.Do(func() {
		if ok := tryEnableAMX(); !ok {
			log.Println("AMX not available; falling back to scalar")
			hasAMX = false
		}
	})

	if hasAMX {
		C.amx_touch() // 触发一次 AMX 指令
	}
	return compressScalar(src)
}

var (
	hasAMX = true
	once   runtime.Once
)
