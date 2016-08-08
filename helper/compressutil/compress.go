package compressutil

import (
	"bytes"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"io"
	"log"
)

const (
	// A byte value used as a canary prefix for the compressed information
	// which is used to distinguish if a JSON input is compressed or not.
	// The value of this constant should not be a first character of any
	// valid JSON string.
	CompressionCanary byte = 'Z'

	CompressionTypeLzw = "lzw"

	CompressionTypeGzip = "gzip"
)

// CompressionConfig is used to select a compression type to be performed by
// Compress and Decompress utilities.
// Supported types are:
// * CompressionTypeLzw
// * CompressionTypeGzip
//
// When using CompressionTypeGzip, the compression levels can also be chosen:
// * gzip.DefaultCompression
// * gzip.BestSpeed
// * gzip.BestCompression
type CompressionConfig struct {
	// Type of the compression algorithm to be used
	Type string

	// When using Gzip format, the compression level to employ
	GzipCompressionLevel int
}

// Compress places the canary byte in a buffer and uses the same buffer to fill
// in the compressed information of the given input. The configuration supports
// two type of compression: LZW and Gzip. When using Gzip compression format,
// if GzipCompressionLevel is not specified, the 'gzip.DefaultCompression' will
// be assumed.
func Compress(data []byte, config *CompressionConfig) ([]byte, error) {
	var buf bytes.Buffer
	var writer io.WriteCloser
	var err error

	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// Write the canary into the buffer first
	buf.Write([]byte{CompressionCanary})

	// Create writer to compress the input data based on the configured type
	switch config.Type {
	case CompressionTypeLzw:
		writer = lzw.NewWriter(&buf, lzw.LSB, 8)
	case CompressionTypeGzip:
		level := gzip.DefaultCompression
		if config.GzipCompressionLevel != gzip.NoCompression {
			level = config.GzipCompressionLevel
		}
		writer, err = gzip.NewWriterLevel(&buf, level)
	default:
		return nil, fmt.Errorf("unsupported compression type")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create a compression writer; err: %v", err)
	}

	if writer == nil {
		return nil, fmt.Errorf("failed to create a compression writer")
	}

	// Compress the input and place it in the same buffer containing the
	// canary byte.
	if _, err = writer.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compress JSON string; err: %v", err)
	}

	// Close the io.WriteCloser
	if err = writer.Close(); err != nil {
		return nil, err
	}

	// Return the compressed bytes with canary byte at the start
	return buf.Bytes(), nil
}

// Decompress checks if the first byte in the input matches the canary byte.
// If the first byte is a canary byte, then the input past the canary byte
// will be decompressed using the method specified in the given configuration.
// If the first byte isn't a canary byte, then the utility returns a boolean
// return value indicating that the input was not compressed.
func Decompress(data []byte, config *CompressionConfig) ([]byte, bool, error) {
	var err error
	var reader io.ReadCloser
	if data == nil || len(data) == 0 {
		return nil, false, fmt.Errorf("'data' being decompressed is invalid")
	}

	if config == nil {
		return nil, false, fmt.Errorf("config is nil")
	}

	// Read the first byte
	bytesReader := bytes.NewReader(data)
	firstByte, err := bytesReader.ReadByte()
	if err != nil {
		return nil, false, fmt.Errorf("failed to read the first byte from the input")
	}

	// If the first byte doesn't match the canary byte, it means that the
	// content was not compressed in the first place.
	if CompressionCanary != firstByte {
		// Indicate the caller that the input was not compressed
		return nil, true, nil
	} else {
		// If the first byte matches the canary byte, remove the canary
		// byte and try to decompress the data before JSON decoding it.
		if len(data) < 2 {
			return nil, false, fmt.Errorf("invalid 'data' after the canary")
		}
		data = data[1:]
	}

	// Create a reader to read the compressed data based on the configured
	// compression type
	switch config.Type {
	case CompressionTypeLzw:
		reader = lzw.NewReader(bytes.NewReader(data), lzw.LSB, 8)
	case CompressionTypeGzip:
		reader, err = gzip.NewReader(bytes.NewReader(data))
	default:
		return nil, false, fmt.Errorf("invalid 'data' being decompressed is invalid")
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to create a compression reader; err: %v", err)
	}

	if reader == nil {
		return nil, false, fmt.Errorf("failed to create a compression reader")
	}

	// Close the io.ReadCloser
	defer reader.Close()

	// Read all the compressed data into a buffer
	var jsonBuf bytes.Buffer
	if _, err = io.Copy(&jsonBuf, reader); err != nil {
		return nil, false, err
	}

	log.Printf("decompressedBytes: %s\n", string(jsonBuf.Bytes()))
	return jsonBuf.Bytes(), false, nil
}
