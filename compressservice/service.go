package compressservice

import (
	"context"
	"errors"

	"github.com/ServiceWeaver/weaver"
)

//go:generate weaver generate .

// Compressor 定义统一接口，便于 adservice 调用。
type Compressor interface {
	// Compress 返回 dst = IAACompress(src)
	Compress(ctx context.Context, src []byte) ([]byte, error)
}

type impl struct {
	weaver.Implements[Compressor]
}

func (c *impl) Compress(ctx context.Context, src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("empty input")
	}
	return iaaCompress(src), nil // iaaCompress 在不同文件里实现
}
