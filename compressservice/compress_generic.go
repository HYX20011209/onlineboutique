//go:build !avx2 && !cuda && !amx
// +build !avx2,!cuda,!amx

package compressservice

// iaaCompress：一个极简 0 阶整数自适应算术编码。
// 仅演示用，无完整边界检测与异常处理。
func iaaCompress(src []byte) []byte {
	return compressScalar(src) // 直接调纯 Go 实现
}
