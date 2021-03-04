package sarama

import (
	"sync"

	"github.com/klauspost/compress/zstd"
)

var (
	zstdDec *zstd.Decoder
	zstdEnc *zstd.Encoder

	zstdEncOnce, zstdDecOnce sync.Once
)

func zstdDecompress(dst, src []byte) ([]byte, error) {
	zstdDecOnce.Do(func() {
		zstdDec, _ = zstd.NewReader(nil)
	})
	return zstdDec.DecodeAll(src, dst)
}

func zstdCompress(dst, src []byte) ([]byte, error) {
	zstdEncOnce.Do(func() {
		zstdEnc, _ = zstd.NewWriter(nil, zstd.WithZeroFrames(true))
	})
	return zstdEnc.EncodeAll(src, dst), nil
}
