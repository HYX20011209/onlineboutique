package compressservice

import (
	"math/rand"
	"testing"
)

func BenchmarkIAA(b *testing.B) {
	data := make([]byte, 4*1024*1024) // 4 MiB
	rand.Read(data)
	c := impl{} // 直接调用实现
	for i := 0; i < b.N; i++ {
		res, _ := c.Compress(nil, data)
		if len(res) == 0 {
			b.Fatal("bad result")
		}
	}
}
