package gocb

import "github.com/golang/snappy"

type compressor struct {
	CompressionEnabled  bool
	CompressionMinSize  uint32
	CompressionMinRatio float64
}

type possiblyCompressedResponse interface {
	GetContentUncompressed() []byte
	GetContentCompressed() []byte
}

func (c *compressor) Compress(val []byte) ([]byte, bool) {
	if !c.CompressionEnabled {
		return val, false
	}

	valueLen := len(val)
	if valueLen > int(c.CompressionMinSize) {
		compressedValue := snappy.Encode(nil, val)
		if float64(len(compressedValue))/float64(valueLen) <= c.CompressionMinRatio {
			return compressedValue, true
		}
	}

	return val, false
}

func (c *compressor) Decompress(val possiblyCompressedResponse) ([]byte, error) {
	if val.GetContentUncompressed() != nil {
		return val.GetContentUncompressed(), nil
	}

	newValue, err := snappy.Decode(nil, val.GetContentCompressed())
	if err != nil {
		return nil, wrapError(err, "failed to decode content")
	}

	return newValue, nil
}
