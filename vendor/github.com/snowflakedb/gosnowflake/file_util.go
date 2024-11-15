// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
	"os"
	usr "os/user"
	"path/filepath"
	"strings"
)

type snowflakeFileUtil struct {
}

const (
	fileChunkSize                 = 16 * 4 * 1024
	readWriteFileMode os.FileMode = 0666
)

func (util *snowflakeFileUtil) compressFileWithGzipFromStream(srcStream **bytes.Buffer) (*bytes.Buffer, int, error) {
	r := getReaderFromBuffer(srcStream)
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, -1, err
	}
	var c bytes.Buffer
	w := gzip.NewWriter(&c)
	w.Write(buf) // write buf to gzip writer
	w.Close()
	return &c, c.Len(), nil
}

func (util *snowflakeFileUtil) compressFileWithGzip(fileName string, tmpDir string) (string, int64, error) {
	basename := baseName(fileName)
	gzipFileName := filepath.Join(tmpDir, basename+"_c.gz")

	fr, err := os.Open(fileName)
	if err != nil {
		return "", -1, err
	}
	defer fr.Close()
	fw, err := os.OpenFile(gzipFileName, os.O_WRONLY|os.O_CREATE, readWriteFileMode)
	if err != nil {
		return "", -1, err
	}
	gzw := gzip.NewWriter(fw)
	defer gzw.Close()
	io.Copy(gzw, fr)

	stat, err := os.Stat(gzipFileName)
	if err != nil {
		return "", -1, err
	}
	return gzipFileName, stat.Size(), nil
}

func (util *snowflakeFileUtil) getDigestAndSizeForStream(stream **bytes.Buffer) (string, int64, error) {
	m := sha256.New()
	r := getReaderFromBuffer(stream)
	chunk := make([]byte, fileChunkSize)

	for {
		n, err := r.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", 0, err
		}
		m.Write(chunk[:n])
	}
	return base64.StdEncoding.EncodeToString(m.Sum(nil)), int64((*stream).Len()), nil
}

func (util *snowflakeFileUtil) getDigestAndSizeForFile(fileName string) (string, int64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	var total int64
	m := sha256.New()
	chunk := make([]byte, fileChunkSize)

	for {
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", 0, err
		}
		total += int64(n)
		m.Write(chunk[:n])
	}
	f.Seek(0, io.SeekStart)
	return base64.StdEncoding.EncodeToString(m.Sum(nil)), total, nil
}

// file metadata for PUT/GET
type fileMetadata struct {
	name               string
	sfa                *snowflakeFileTransferAgent
	stageLocationType  cloudType
	resStatus          resultStatus
	stageInfo          *execResponseStageInfo
	encryptionMaterial *snowflakeFileEncryption

	srcFileName        string
	realSrcFileName    string
	srcFileSize        int64
	srcCompressionType *compressionType
	uploadSize         int64
	dstFileSize        int64
	dstFileName        string
	dstCompressionType *compressionType

	client             cloudClient // *s3.Client (S3), *azblob.ContainerURL (Azure), string (GCS)
	requireCompress    bool
	parallel           int64
	sha256Digest       string
	overwrite          bool
	tmpDir             string
	errorDetails       error
	lastError          error
	noSleepingTime     bool
	lastMaxConcurrency int
	localLocation      string
	options            *SnowflakeFileTransferOptions

	/* streaming PUT */
	srcStream     *bytes.Buffer
	realSrcStream *bytes.Buffer

	/* streaming GET */
	dstStream *bytes.Buffer

	/* GCS */
	presignedURL                *url.URL
	gcsFileHeaderDigest         string
	gcsFileHeaderContentLength  int64
	gcsFileHeaderEncryptionMeta *encryptMetadata

	/* mock */
	mockUploader    s3UploadAPI
	mockDownloader  s3DownloadAPI
	mockHeader      s3HeaderAPI
	mockGcsClient   gcsAPI
	mockAzureClient azureAPI
}

type fileTransferResultType struct {
	name               string
	srcFileName        string
	dstFileName        string
	srcFileSize        int64
	dstFileSize        int64
	srcCompressionType *compressionType
	dstCompressionType *compressionType
	resStatus          resultStatus
	errorDetails       error
}

type fileHeader struct {
	digest             string
	contentLength      int64
	encryptionMetadata *encryptMetadata
}

func getReaderFromBuffer(src **bytes.Buffer) io.Reader {
	var b bytes.Buffer
	tee := io.TeeReader(*src, &b) // read src to buf
	*src = &b                     // revert pointer back
	return tee
}

// baseName returns the pathname of the path provided
func baseName(path string) string {
	base := filepath.Base(path)
	if base == "." || base == "/" {
		return ""
	}
	if len(base) > 1 && (path[len(path)-1:] == "." || path[len(path)-1:] == "/") {
		return ""
	}
	return base
}

// expandUser returns the argument with an initial component of ~
func expandUser(path string) (string, error) {
	usr, err := usr.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path, nil
}

// getDirectory retrieves the current working directory
func getDirectory() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
