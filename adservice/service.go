// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adservice

import (
	"context"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	maxAdsToServe = 2
)

// Ad represents an advertisement.
type Ad struct {
	weaver.AutoMarshal
	RedirectURL string // URL to redirect to when an ad is clicked.
	Text        string // Short advertisement text to display.
}

type AdService interface {
	GetAds(ctx context.Context, keywords []string) ([]Ad, error)
}

type impl struct {
	weaver.Implements[AdService]
	ads map[string]Ad
}

func (s *impl) Init(ctx context.Context) error {
	s.Logger(ctx).Info("Ad Service started")
	s.ads = createAdsMap()
	return nil
}

func burnCPU(iter int) {
	total := 0
	for i := 0; i < iter; i++ {
		// 做点没意义的计算，例如
		total += i * i % 99999
	}
}

// naiveMultiply 朴素矩阵乘法 C = A * B (A, B, C 均为 size x size)
func naiveMultiply(A, B [][]float64, size int) [][]float64 {
	C := make([][]float64, size)
	for i := 0; i < size; i++ {
		C[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			sum := 0.0
			for k := 0; k < size; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}
	return C
}

// gonumMultiply 使用 gonum/mat 提供的 Mul() 函数做矩阵乘法
func gonumMultiply(A, B mat.Matrix) mat.Matrix {
	rA, cA := A.Dims()
	rB, cB := B.Dims()
	if cA != rB {
		panic("dimension mismatch")
	}
	C := mat.NewDense(rA, cB, nil)
	// 内部调用 BLAS，可利用SIMD、多线程等
	C.Mul(A, B)
	return C
}

func burnMatrixIfEnabled() {
	// 通过 MATRIX_SIZE 控制矩阵大小
	//sizeStr := os.Getenv("MATRIX_SIZE")
	sizeStr := "1024"
	//if sizeStr == "" {
	//	// 未设置就不做额外运算
	//	return
	//}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 512 // 默认512x512
	}

	// 通过 USE_GONUM 判断用Gonum还是朴素方法
	//useGonum := (os.Getenv("USE_GONUM") != "")
	useGonum := true

	// 生成随机矩阵 A, B (size x size)
	A2D := make([][]float64, size)
	B2D := make([][]float64, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		A2D[i] = make([]float64, size)
		B2D[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			A2D[i][j] = rand.Float64()
			B2D[i][j] = rand.Float64()
		}
	}

	start := time.Now()
	if useGonum {
		// 将 A2D, B2D 转为 gonum 的 Dense
		A_gonum := mat.NewDense(size, size, nil)
		B_gonum := mat.NewDense(size, size, nil)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				A_gonum.Set(i, j, A2D[i][j])
				B_gonum.Set(i, j, B2D[i][j])
			}
		}
		// gonum 乘法
		C := gonumMultiply(A_gonum, B_gonum)
		val := C.At(0, 0) // 避免编译器优化
		elapsed := time.Since(start)
		// 可换成 s.Logger(ctx).Info(...)，这里只示范简单打印
		println("[burnMatrixIfEnabled] gonumMultiply done, C[0,0] =", val, "elapsed=", elapsed)
	} else {
		// 朴素乘法
		C2D := naiveMultiply(A2D, B2D, size)
		val := C2D[0][0]
		elapsed := time.Since(start)
		println("[burnMatrixIfEnabled] naiveMultiply done, C[0,0] =", val, "elapsed=", elapsed)
	}
}

// GetAds returns a list of ads that best match the given context keywords.
func (s *impl) GetAds(ctx context.Context, keywords []string) ([]Ad, error) {
	burnMatrixIfEnabled() // <-- 新增：执行可选的矩阵乘法烧CPU

	//burnCPU(100000000) // 计算量大小可调
	//burnCPU(10000000000)
	//burnCPU(10000000000)
	//burnCPU(10000000000)
	//burnCPU(10000000000)
	s.Logger(ctx).Info("received ad request", "keywords", keywords)
	span := trace.SpanFromContext(ctx)
	var allAds []Ad
	if len(keywords) > 0 {
		span.AddEvent("Constructing Ads using context", trace.WithAttributes(
			attribute.String("Context Keys", strings.Join(keywords, ",")),
			attribute.Int("Context Keys length", len(keywords)),
		))
		for _, kw := range keywords {
			allAds = append(allAds, s.getAdsByCategory(kw)...)
		}
		if allAds == nil {
			// Serve random ads.
			span.AddEvent("No Ads found based on context. Constructing random Ads.")
			allAds = s.getRandomAds()
		}
	} else {
		span.AddEvent("No Context provided. Constructing random Ads.")
		allAds = s.getRandomAds()
	}
	return allAds, nil
}

func (s *impl) getAdsByCategory(category string) []Ad {
	return []Ad{s.ads[category]}
}

func (s *impl) getRandomAds() []Ad {
	ads := make([]Ad, maxAdsToServe)
	vals := maps.Values(s.ads)
	for i := 0; i < maxAdsToServe; i++ {
		ads[i] = vals[rand.Intn(len(vals))]
	}
	return ads
}

func createAdsMap() map[string]Ad {
	return map[string]Ad{
		"hair": {
			RedirectURL: "/product/2ZYFJ3GM2N",
			Text:        "Hairdryer for sale. 50% off.",
		},
		"clothing": {
			RedirectURL: "/product/66VCHSJNUP",
			Text:        "Tank top for sale. 20% off.",
		},
		"accessories": {
			RedirectURL: "/product/1YMWWN1N4O",
			Text:        "Watch for sale. Buy one, get second kit for free",
		},
		"footwear": {
			RedirectURL: "/product/L9ECAV7KIM",
			Text:        "Loafers for sale. Buy one, get second one for free",
		},
		"decor": {
			RedirectURL: "/product/0PUK6V6EV0",
			Text:        "Candle holder for sale. 30% off.",
		},
		"kitchen": {
			RedirectURL: "/product/9SIQT8TOJO",
			Text:        "Bamboo glass jar for sale. 10% off.",
		},
	}
}
