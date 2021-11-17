// Copyright (c) 2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	usr "os/user"
	"path/filepath"
	"strings"
)

type snowflakeFileUtil struct {
}

func (util *snowflakeFileUtil) compressFileWithGzipFromStream(srcStream **bytes.Buffer) (*bytes.Buffer, int) {
	r := getReaderFromBuffer(srcStream)
	buf, _ := ioutil.ReadAll(r)
	var c bytes.Buffer
	w := gzip.NewWriter(&c)
	w.Write(buf) // write buf to gzip writer
	w.Close()
	return &c, c.Len()
}

func (util *snowflakeFileUtil) compressFileWithGzip(fileName string, tmpDir string) (string, int64) {
	basename := baseName(fileName)
	gzipFileName := filepath.Join(tmpDir, basename+"_c.gz")

	fr, _ := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	defer fr.Close()
	fw, _ := os.OpenFile(gzipFileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	gzw := gzip.NewWriter(fw)
	defer gzw.Close()
	io.Copy(gzw, fr)

	stat, _ := os.Stat(gzipFileName)
	return gzipFileName, stat.Size()
}

func (util *snowflakeFileUtil) getDigestAndSize(src **bytes.Buffer) (string, int64) {
	chunkSize := 16 * 4 * 1024
	m := sha256.New()
	r := getReaderFromBuffer(src)
	for {
		chunk := make([]byte, chunkSize)
		n, err := r.Read(chunk)
		if n == 0 || err != nil {
			break
		}
		m.Write(chunk[:n])
	}
	return base64.StdEncoding.EncodeToString(m.Sum(nil)), int64((*src).Len())
}

func (util *snowflakeFileUtil) getDigestAndSizeForStream(stream **bytes.Buffer) (string, int64) {
	return util.getDigestAndSize(stream)
}

func (util *snowflakeFileUtil) getDigestAndSizeForFile(fileName string) (string, int64, error) {
	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", 0, err
	}
	buf := bytes.NewBuffer(src)
	digest, size := util.getDigestAndSize(&buf)
	return digest, size, err
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

	/* GCS */
	presignedURL                *url.URL
	gcsFileHeaderDigest         string
	gcsFileHeaderContentLength  int64
	gcsFileHeaderEncryptionMeta *encryptMetadata

	/* mock */
	mockUploader s3UploadAPI
	mockHeader   s3HeaderAPI
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
func expandUser(path string) string {
	usr, _ := usr.Current()
	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

// getDirectory retrieves the current working directory
func getDirectory() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
