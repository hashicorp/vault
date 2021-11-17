// Copyright (c) 2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	gcsMetadataPrefix             = "x-goog-meta-"
	gcsMetadataSfcDigest          = gcsMetadataPrefix + sfcDigest
	gcsMetadataMatdescKey         = gcsMetadataPrefix + "matdesc"
	gcsMetadataEncryptionDataProp = gcsMetadataPrefix + "encryptiondata"
	gcsFileHeaderDigest           = "gcs-file-header-digest"
)

type snowflakeGcsUtil struct {
}

type gcsLocation struct {
	bucketName string
	path       string
}

func (util *snowflakeGcsUtil) createClient(info *execResponseStageInfo, _ bool) cloudClient {
	if info.Creds.GcsAccessToken != "" {
		return info.Creds.GcsAccessToken
	}
	return ""
}

// cloudUtil implementation
func (util *snowflakeGcsUtil) getFileHeader(meta *fileMetadata, filename string) *fileHeader {
	if meta.resStatus == uploaded || meta.resStatus == downloaded {
		return &fileHeader{
			digest:             meta.gcsFileHeaderDigest,
			contentLength:      meta.gcsFileHeaderContentLength,
			encryptionMetadata: meta.gcsFileHeaderEncryptionMeta,
		}
	}
	if meta.presignedURL != nil {
		meta.resStatus = notFoundFile
	} else {
		URL := util.generateFileURL(meta.stageInfo.Location, strings.TrimLeft(filename, "/"))
		accessToken := meta.client.(string)
		gcsHeaders := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		req, err := http.NewRequest("HEAD", URL.String(), nil)
		if err != nil {
			return nil
		}
		for k, v := range gcsHeaders {
			req.Header.Add(k, v)
		}
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 403 || resp.StatusCode == 408 || resp.StatusCode == 429 || resp.StatusCode == 500 || resp.StatusCode == 503 {
				meta.lastError = fmt.Errorf(resp.Status)
				meta.resStatus = needRetry
				return nil
			}
			if resp.StatusCode == 404 {
				meta.resStatus = notFoundFile
			} else if util.isTokenExpired(resp) {
				meta.lastError = fmt.Errorf(resp.Status)
				meta.resStatus = renewToken
			}
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = errStatus
			return nil
		}

		digest := resp.Header.Get(gcsMetadataSfcDigest)
		contentLength, _ := strconv.Atoi(resp.Header.Get("content-length"))
		var encryptionMeta *encryptMetadata
		if resp.Header.Get(gcsMetadataEncryptionDataProp) != "" {
			var encryptData *encryptionData
			json.Unmarshal([]byte(resp.Header.Get(gcsMetadataEncryptionDataProp)), encryptData)
			if encryptData != nil {
				encryptionMeta = &encryptMetadata{
					key: encryptData.WrappedContentKey.EncryptionKey,
					iv:  encryptData.ContentEncryptionIV,
				}
				if resp.Header.Get(gcsMetadataMatdescKey) != "" {
					encryptionMeta.matdesc = resp.Header.Get(gcsMetadataMatdescKey)
				}
			}
		}
		meta.resStatus = uploaded
		return &fileHeader{
			digest:             digest,
			contentLength:      int64(contentLength),
			encryptionMetadata: encryptionMeta,
		}
	}
	return nil
}

// cloudUtil implementation
func (util *snowflakeGcsUtil) uploadFile(dataFile string, meta *fileMetadata, encryptMeta *encryptMetadata, maxConcurrency int, multiPartThreshold int64) error {
	uploadURL := meta.presignedURL
	var accessToken string

	if uploadURL == nil {
		uploadURL = util.generateFileURL(meta.stageInfo.Location, strings.TrimLeft(meta.dstFileName, "/"))
		accessToken = meta.client.(string)
	}

	var contentEncoding string
	if meta.dstCompressionType != nil {
		contentEncoding = strings.ToLower(meta.dstCompressionType.name)
	}

	if contentEncoding == "gzip" {
		contentEncoding = ""
	}

	gcsHeaders := make(map[string]string)
	gcsHeaders[httpHeaderContentEncoding] = contentEncoding
	gcsHeaders[gcsMetadataSfcDigest] = meta.sha256Digest
	if accessToken != "" {
		gcsHeaders["Authorization"] = "Bearer " + accessToken
	}

	if encryptMeta != nil {
		encryptData := encryptionData{
			"FullBlob",
			contentKey{
				"symmKey1",
				encryptMeta.key,
				"AES_CBC_256",
			},
			encryptionAgent{
				"1.0",
				"AES_CBC_256",
			},
			encryptMeta.iv,
			keyMetadata{
				"Java 5.3.0",
			},
		}
		b, _ := json.Marshal(&encryptData)
		gcsHeaders[gcsMetadataEncryptionDataProp] = string(b)
		gcsHeaders[gcsMetadataMatdescKey] = encryptMeta.matdesc
	}

	var uploadSrc io.Reader
	if meta.srcStream != nil {
		uploadSrc = meta.srcStream
		if meta.realSrcStream != nil {
			uploadSrc = meta.realSrcStream
		}
	} else {
		uploadSrc, _ = os.OpenFile(dataFile, os.O_RDONLY, os.ModePerm)
	}

	req, err := http.NewRequest("PUT", uploadURL.String(), uploadSrc)
	if err != nil {
		return err
	}
	for k, v := range gcsHeaders {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		meta.lastError = fmt.Errorf(resp.Status)
		if resp.StatusCode == 403 || resp.StatusCode == 408 || resp.StatusCode == 429 || resp.StatusCode == 500 || resp.StatusCode == 503 {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = needRetry
		} else if accessToken == "" && resp.StatusCode == 400 && meta.lastError == nil {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = renewPresignedURL
		} else if accessToken != "" && util.isTokenExpired(resp) {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = renewToken
		}
		return meta.lastError
	}

	if meta.options.putCallback != nil {
		meta.options.putCallback = &snowflakeProgressPercentage{
			filename:        dataFile,
			fileSize:        float64(meta.srcFileSize),
			outputStream:    meta.options.putCallbackOutputStream,
			showProgressBar: meta.options.showProgressBar,
		}
	}

	meta.dstFileSize = meta.uploadSize
	meta.resStatus = uploaded

	meta.gcsFileHeaderDigest = gcsHeaders[gcsFileHeaderDigest]
	meta.gcsFileHeaderContentLength = meta.uploadSize
	if err = json.Unmarshal([]byte(gcsHeaders[gcsMetadataEncryptionDataProp]), &encryptMeta); err != nil {
		return err
	}
	meta.gcsFileHeaderEncryptionMeta = encryptMeta
	return nil
}

// cloudUtil implementation
func (util *snowflakeGcsUtil) nativeDownloadFile(
	meta *fileMetadata,
	fullDstFileName string,
	maxConcurrency int64) error {
	downloadURL := meta.presignedURL
	var accessToken string
	gcsHeaders := make(map[string]string)

	if downloadURL == nil {
		downloadURL = util.generateFileURL(meta.stageInfo.Location, strings.TrimLeft(meta.dstFileName, "/"))
		accessToken = meta.client.(string)
		if accessToken != "" {
			gcsHeaders["Authorization"] = "Bearer " + accessToken
		}
	}

	req, err := http.NewRequest("GET", downloadURL.String(), nil)
	if err != nil {
		return err
	}
	for k, v := range gcsHeaders {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		meta.lastError = fmt.Errorf(resp.Status)
		if resp.StatusCode == 403 || resp.StatusCode == 408 || resp.StatusCode == 429 || resp.StatusCode == 500 || resp.StatusCode == 503 {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = needRetry
		} else if accessToken == "" && resp.StatusCode == 400 && meta.lastError == nil {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = renewPresignedURL
		} else if accessToken != "" && util.isTokenExpired(resp) {
			meta.lastError = fmt.Errorf(resp.Status)
			meta.resStatus = renewToken
		}
		return meta.lastError
	}

	f, err := os.OpenFile(fullDstFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = io.Copy(f, resp.Body); err != nil {
		return err
	}

	var encryptMeta encryptMetadata
	if resp.Header.Get(gcsMetadataEncryptionDataProp) != "" {
		var encryptData *encryptionData
		if err = json.Unmarshal([]byte(resp.Header.Get(gcsMetadataEncryptionDataProp)), &encryptData); err != nil {
			return err
		}
		if encryptData != nil {
			encryptMeta = encryptMetadata{
				encryptData.WrappedContentKey.EncryptionKey,
				encryptData.ContentEncryptionIV,
				"",
			}
			if key := resp.Header.Get(gcsMetadataMatdescKey); key != "" {
				encryptMeta.matdesc = key
			}
		}
	}

	fi, err := os.Stat(fullDstFileName)
	if err != nil {
		return err
	}
	meta.srcFileSize = fi.Size()
	meta.resStatus = downloaded
	meta.gcsFileHeaderDigest = resp.Header.Get(gcsMetadataSfcDigest)
	meta.gcsFileHeaderContentLength = resp.ContentLength
	meta.gcsFileHeaderEncryptionMeta = &encryptMeta
	return nil
}

func (util *snowflakeGcsUtil) extractBucketNameAndPath(location string) *gcsLocation {
	containerName := location
	var path string
	if strings.Contains(location, "/") {
		containerName = location[:strings.Index(location, "/")]
		path = location[strings.Index(location, "/")+1:]
		if path != "" && strings.HasSuffix(location, "/") {
			path += "/"
		}
	}
	return &gcsLocation{containerName, path}
}

func (util *snowflakeGcsUtil) generateFileURL(stageLocation string, filename string) *url.URL {
	gcsLoc := util.extractBucketNameAndPath(stageLocation)
	fullFilePath := gcsLoc.path + filename
	URL, _ := url.Parse("https://storage.googleapis.com/" + gcsLoc.bucketName + "/" + url.QueryEscape(fullFilePath))
	return URL
}

func (util *snowflakeGcsUtil) isTokenExpired(resp *http.Response) bool {
	return resp.StatusCode == 401
}
