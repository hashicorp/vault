package compressutil

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"io"
	"log"
)

const (
	CompressionCanaryJSON byte = 'Z'
)

func Compress(data []byte, canary byte) ([]byte, error) {
	// Create a buffer and place the canary as its first byte
	var buf bytes.Buffer
	buf.Write([]byte{canary})

	// Create writer to compress the JSON encoded bytes
	writer := lzw.NewWriter(&buf, lzw.LSB, 8)

	// Compress the JSON bytes
	if _, err := writer.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compress JSON string; err: %v", err)
	}

	// Close the io.WriteCloser
	if err := writer.Close(); err != nil {
		return nil, err
	}

	log.Printf("compressutil.Compress: len(compressedBytes): %d\n", len(buf.Bytes()))

	// Return the compressed bytes with canary byte at the start
	return buf.Bytes(), nil
}

func Decompress(data []byte, canary byte) ([]byte, bool, error) {
	if data == nil || len(data) < 2 {
		return nil, false, fmt.Errorf("'data' being decompressed is invalid")
	}

	// Read the first byte
	bytesReader := bytes.NewReader(data)
	firstByte, err := bytesReader.ReadByte()
	if err != nil {
		return nil, false, fmt.Errorf("failed to read the first byte from the input")
	}

	// If the first byte doesn't match the canaryByte, it means that the
	// content was not compressed in the first place. Try JSON decoding it.
	if canary != firstByte {
		return nil, true, nil
	} else {
		// If the first byte matches the canaryByte, remove the canary
		// byte and try to decompress the data before JSON decoding it.
		data = data[1:]
	}

	// Create a reader to read the compressed data
	reader := lzw.NewReader(bytes.NewReader(data), lzw.LSB, 8)

	// Close the io.ReadCloser
	defer reader.Close()

	// Read all the compressed data into a buffer
	var jsonBuf bytes.Buffer
	if _, err := io.Copy(&jsonBuf, reader); err != nil {
		return nil, false, err
	}

	return jsonBuf.Bytes(), false, nil
}
