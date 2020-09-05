package sarama

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"

	snappy "github.com/eapache/go-xerial-snappy"
	"github.com/pierrec/lz4"
)

var (
	lz4WriterPool = sync.Pool{
		New: func() interface{} {
			return lz4.NewWriter(nil)
		},
	}

	gzipWriterPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
	gzipWriterPoolForCompressionLevel1 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 1)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel2 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 2)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel3 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 3)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel4 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 4)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel5 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 5)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel6 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 6)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel7 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 7)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel8 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 8)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	gzipWriterPoolForCompressionLevel9 = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(nil, 9)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
)

func compress(cc CompressionCodec, level int, data []byte) ([]byte, error) {
	switch cc {
	case CompressionNone:
		return data, nil
	case CompressionGZIP:
		var (
			err    error
			buf    bytes.Buffer
			writer *gzip.Writer
		)

		switch level {
		case CompressionLevelDefault:
			writer = gzipWriterPool.Get().(*gzip.Writer)
			defer gzipWriterPool.Put(writer)
			writer.Reset(&buf)
		case 1:
			writer = gzipWriterPoolForCompressionLevel1.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel1.Put(writer)
			writer.Reset(&buf)
		case 2:
			writer = gzipWriterPoolForCompressionLevel2.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel2.Put(writer)
			writer.Reset(&buf)
		case 3:
			writer = gzipWriterPoolForCompressionLevel3.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel3.Put(writer)
			writer.Reset(&buf)
		case 4:
			writer = gzipWriterPoolForCompressionLevel4.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel4.Put(writer)
			writer.Reset(&buf)
		case 5:
			writer = gzipWriterPoolForCompressionLevel5.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel5.Put(writer)
			writer.Reset(&buf)
		case 6:
			writer = gzipWriterPoolForCompressionLevel6.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel6.Put(writer)
			writer.Reset(&buf)
		case 7:
			writer = gzipWriterPoolForCompressionLevel7.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel7.Put(writer)
			writer.Reset(&buf)
		case 8:
			writer = gzipWriterPoolForCompressionLevel8.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel8.Put(writer)
			writer.Reset(&buf)
		case 9:
			writer = gzipWriterPoolForCompressionLevel9.Get().(*gzip.Writer)
			defer gzipWriterPoolForCompressionLevel9.Put(writer)
			writer.Reset(&buf)
		default:
			writer, err = gzip.NewWriterLevel(&buf, level)
			if err != nil {
				return nil, err
			}
		}
		if _, err := writer.Write(data); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case CompressionSnappy:
		return snappy.Encode(data), nil
	case CompressionLZ4:
		writer := lz4WriterPool.Get().(*lz4.Writer)
		defer lz4WriterPool.Put(writer)

		var buf bytes.Buffer
		writer.Reset(&buf)

		if _, err := writer.Write(data); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case CompressionZSTD:
		return zstdCompress(nil, data)
	default:
		return nil, PacketEncodingError{fmt.Sprintf("unsupported compression codec (%d)", cc)}
	}
}
