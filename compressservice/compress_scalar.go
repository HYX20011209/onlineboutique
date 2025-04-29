package compressservice

// compressScalar 实现 0 阶整数自适应算术编码。
// 代码与你原先 generic 版完全一致，仅函数名改为 compressScalar。
func compressScalar(src []byte) []byte {
	const (
		top    = uint32(1 << 24)
		bottom = uint32(1 << 16)
	)

	var low uint32
	var range_ uint32 = 0xFFFFFFFF
	var freq [257]uint32
	for i := range freq {
		freq[i] = uint32(i)
	}

	out := make([]byte, 0, len(src)/2)
	shiftOut := func() {
		out = append(out, byte(low>>24))
		low <<= 8
		range_ <<= 8
	}

	for _, b := range src {
		total := freq[256]
		r := range_ / total
		highsym := freq[b+1]
		lowsym := freq[b]
		low += r * lowsym
		range_ = r * (highsym - lowsym)

		for range_ <= bottom {
			shiftOut()
			range_ = (range_ << 8) | 0xFF
		}

		for i := int(b) + 1; i <= 256; i++ {
			freq[i]++
		}
		if freq[256] > 1<<23 {
			var cum uint32
			for i := 0; i <= 256; i++ {
				delta := freq[i+1] - freq[i]
				cum += delta>>1 + 1
				freq[i+1] = cum
			}
		}
	}
	out = append(out,
		byte(low>>24), byte(low>>16), byte(low>>8), byte(low))
	return out
}