// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"io"
	"os"

	"github.com/oracle/oci-go-sdk/common"
)

// multipartManifest provides thread-safe access to an ongoing manifest upload.
type multipartManifest struct {
	// key is UploadID, define it as map since user can upload multiple times
	// second key is part number
	parts map[string]map[int]uploadPart
}

type uploadPart struct {
	size     int64
	offset   int64
	partBody []byte
	partNum  int
	hash     *string
	opcMD5   *string
	etag     *string
	err      error
}

// splitFileToParts starts a goroutine to read a file and break down to parts and send the parts to
// uploadPart channel. It sends the error to error chanel. If done is closed, splitFileToParts
// abandones its works.
func (manifest *multipartManifest) splitFileToParts(done <-chan struct{}, partSize int64, file *os.File, fileSize int64) <-chan uploadPart {
	parts := make(chan uploadPart)

	// Number of parts of the file
	numberOfParts := int(fileSize / partSize)

	go func() {
		// close the channel after splitFile returns
		defer func() {
			common.Debugln("closing parts channel from splitFileParts")
			close(parts)
		}()

		// All buffer sizes are the same in the normal case. Offsets depend on the index.
		// Second go routine should start at 100, for example, given our
		// buffer size of 100.
		for i := 0; i < numberOfParts; i++ {
			offset := partSize * int64(i) // offset of the file, start with 0

			buffer := make([]byte, partSize)
			_, err := file.ReadAt(buffer, offset)

			part := uploadPart{
				partNum:  i + 1,
				size:     partSize,
				offset:   offset,
				err:      err,
				partBody: buffer,
			}

			select {
			case parts <- part:
			case <-done:
				return
			}
		}

		// check for any left over bytes. Add the residual number of bytes as the
		// the last chunk size.
		if remainder := fileSize % int64(partSize); remainder != 0 {
			part := uploadPart{offset: (int64(numberOfParts) * partSize), partNum: numberOfParts + 1}

			part.partBody = make([]byte, remainder)
			_, err := file.ReadAt(part.partBody, part.offset)

			part.size = remainder
			part.err = err

			select {
			case parts <- part:
			case <-done:
				return
			}
		}
	}()

	return parts
}

// splitStreamToParts starts a goroutine to read a stream and break down to parts and send the parts to
// uploadPart channel. It sends the error to error chanel. If done is closed, splitStreamToParts
// abandones its works.
func (manifest *multipartManifest) splitStreamToParts(done <-chan struct{}, partSize int64, reader io.Reader) <-chan uploadPart {
	parts := make(chan uploadPart)

	go func() {
		defer close(parts)
		partNum := 1
		for {
			buffer := make([]byte, partSize)
			_, err := reader.Read(buffer)

			if err == io.EOF {
				break
			}

			part := uploadPart{
				partNum:  partNum,
				size:     partSize,
				err:      err,
				partBody: buffer,
			}

			partNum++
			select {
			case parts <- part:
			case <-done:
				return
			}
		}
	}()

	return parts
}

// update the result in manifest
func (manifest *multipartManifest) updateManifest(result <-chan uploadPart, uploadID string) {
	if manifest.parts[uploadID] == nil {
		manifest.parts[uploadID] = make(map[int]uploadPart)
	}
	for r := range result {
		manifest.parts[uploadID][r.partNum] = r
	}
}
