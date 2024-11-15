// Copyright (c) 2021-2023 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type snowflakeAzureClient struct {
}

type azureLocation struct {
	containerName string
	path          string
}

type azureAPI interface {
	UploadStream(ctx context.Context, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	UploadFile(ctx context.Context, file *os.File, o *azblob.UploadFileOptions) (azblob.UploadFileResponse, error)
	DownloadFile(ctx context.Context, file *os.File, o *blob.DownloadFileOptions) (int64, error)
	DownloadStream(ctx context.Context, o *blob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error)
	GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error)
}

func (util *snowflakeAzureClient) createClient(info *execResponseStageInfo, _ bool) (cloudClient, error) {
	sasToken := info.Creds.AzureSasToken
	u, err := url.Parse(fmt.Sprintf("https://%s.%s/%s%s", info.StorageAccount, info.EndPoint, info.Path, sasToken))
	if err != nil {
		return nil, err
	}
	client, err := azblob.NewClientWithNoCredential(u.String(), &azblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				MaxRetries: 60,
				RetryDelay: 2 * time.Second,
			},
			Transport: &http.Client{
				Transport: SnowflakeTransport,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// cloudUtil implementation
func (util *snowflakeAzureClient) getFileHeader(meta *fileMetadata, filename string) (*fileHeader, error) {
	client, ok := meta.client.(*azblob.Client)
	if !ok {
		return nil, fmt.Errorf("failed to parse client to azblob.Client")
	}

	azureLoc, err := util.extractContainerNameAndPath(meta.stageInfo.Location)
	if err != nil {
		return nil, err
	}
	path := azureLoc.path + strings.TrimLeft(filename, "/")
	containerClient, err := createContainerClient(client.URL())
	if err != nil {
		return nil, &SnowflakeError{
			Message: "failed to create container client",
		}
	}
	var blobClient azureAPI
	blobClient = containerClient.NewBlockBlobClient(path)
	// for testing only
	if meta.mockAzureClient != nil {
		blobClient = meta.mockAzureClient
	}
	resp, err := blobClient.GetProperties(context.Background(), &blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{},
		CPKInfo:          &blob.CPKInfo{},
	})
	if err != nil {
		var se *azcore.ResponseError
		if errors.As(err, &se) {
			if se.ErrorCode == string(bloberror.BlobNotFound) {
				meta.resStatus = notFoundFile
				return nil, fmt.Errorf("could not find file")
			} else if se.StatusCode == 403 {
				meta.resStatus = renewToken
				return nil, fmt.Errorf("received 403, attempting to renew")
			}
		}
		meta.resStatus = errStatus
		return nil, err
	}

	meta.resStatus = uploaded
	metadata := resp.Metadata
	var encData encryptionData

	_, ok = metadata["Encryptiondata"]
	if ok {
		if err = json.Unmarshal([]byte(*metadata["Encryptiondata"]), &encData); err != nil {
			return nil, err
		}
	}

	matdesc, ok := metadata["Matdesc"]
	if !ok {
		// matdesc is not in response, use empty string
		matdesc = new(string)
	}
	encryptionMetadata := encryptMetadata{
		encData.WrappedContentKey.EncryptionKey,
		encData.ContentEncryptionIV,
		*matdesc,
	}

	digest, ok := metadata["Sfcdigest"]
	if !ok {
		// sfcdigest is not in response, use empty string
		digest = new(string)
	}
	return &fileHeader{
		*digest,
		int64(len(metadata)),
		&encryptionMetadata,
	}, nil
}

// cloudUtil implementation
func (util *snowflakeAzureClient) uploadFile(
	dataFile string,
	meta *fileMetadata,
	encryptMeta *encryptMetadata,
	maxConcurrency int,
	multiPartThreshold int64) error {
	azureMeta := map[string]*string{
		"sfcdigest": &meta.sha256Digest,
	}
	if encryptMeta != nil {
		ed := &encryptionData{
			EncryptionMode: "FullBlob",
			WrappedContentKey: contentKey{
				"symmKey1",
				encryptMeta.key,
				"AES_CBC_256",
			},
			EncryptionAgent: encryptionAgent{
				"1.0",
				"AES_CBC_128",
			},
			ContentEncryptionIV: encryptMeta.iv,
			KeyWrappingMetadata: keyMetadata{
				"Java 5.3.0",
			},
		}
		metadata, err := json.Marshal(ed)
		if err != nil {
			return err
		}
		encryptionMetadata := string(metadata)
		azureMeta["encryptiondata"] = &encryptionMetadata
		azureMeta["matdesc"] = &encryptMeta.matdesc
	}

	azureLoc, err := util.extractContainerNameAndPath(meta.stageInfo.Location)
	if err != nil {
		return err
	}
	path := azureLoc.path + strings.TrimLeft(meta.dstFileName, "/")
	client, ok := meta.client.(*azblob.Client)
	if !ok {
		return &SnowflakeError{
			Message: "failed to cast to azure client",
		}
	}
	containerClient, err := createContainerClient(client.URL())

	if err != nil {
		return &SnowflakeError{
			Message: "failed to create container client",
		}
	}
	var blobClient azureAPI
	blobClient = containerClient.NewBlockBlobClient(path)
	// for testing only
	if meta.mockAzureClient != nil {
		blobClient = meta.mockAzureClient
	}
	if meta.srcStream != nil {
		uploadSrc := meta.srcStream
		if meta.realSrcStream != nil {
			uploadSrc = meta.realSrcStream
		}
		_, err = blobClient.UploadStream(context.Background(), uploadSrc, &azblob.UploadStreamOptions{
			BlockSize: int64(uploadSrc.Len()),
			Metadata:  azureMeta,
		})
	} else {
		var f *os.File
		f, err = os.Open(dataFile)
		if err != nil {
			return err
		}
		defer f.Close()

		contentType := "application/octet-stream"
		contentEncoding := "utf-8"
		blobOptions := &azblob.UploadFileOptions{
			HTTPHeaders: &blob.HTTPHeaders{
				BlobContentType:     &contentType,
				BlobContentEncoding: &contentEncoding,
			},
			Metadata:    azureMeta,
			Concurrency: uint16(maxConcurrency),
		}
		if meta.options.putAzureCallback != nil {
			blobOptions.Progress = meta.options.putAzureCallback.call
		}
		_, err = blobClient.UploadFile(context.Background(), f, blobOptions)
	}
	if err != nil {
		var se *azcore.ResponseError
		if errors.As(err, &se) {
			if se.StatusCode == 403 && util.detectAzureTokenExpireError(se.RawResponse) {
				meta.resStatus = renewToken
			} else {
				meta.resStatus = needRetry
				meta.lastError = err
			}
			return err
		}
		meta.resStatus = errStatus
		return err
	}

	meta.dstFileSize = meta.uploadSize
	meta.resStatus = uploaded
	return nil
}

// cloudUtil implementation
func (util *snowflakeAzureClient) nativeDownloadFile(
	meta *fileMetadata,
	fullDstFileName string,
	maxConcurrency int64) error {
	azureLoc, err := util.extractContainerNameAndPath(meta.stageInfo.Location)
	if err != nil {
		return err
	}
	path := azureLoc.path + strings.TrimLeft(meta.srcFileName, "/")
	client, ok := meta.client.(*azblob.Client)
	if !ok {
		return &SnowflakeError{
			Message: "failed to cast to azure client",
		}
	}
	containerClient, err := createContainerClient(client.URL())
	if err != nil {
		return &SnowflakeError{
			Message: "failed to create container client",
		}
	}
	var blobClient azureAPI
	blobClient = containerClient.NewBlockBlobClient(path)
	// for testing only
	if meta.mockAzureClient != nil {
		blobClient = meta.mockAzureClient
	}
	if meta.options.GetFileToStream {
		blobDownloadResponse, err := blobClient.DownloadStream(context.Background(), &azblob.DownloadStreamOptions{})
		if err != nil {
			return err
		}
		retryReader := blobDownloadResponse.NewRetryReader(context.Background(), &azblob.RetryReaderOptions{})
		defer retryReader.Close()
		_, err = meta.dstStream.ReadFrom(retryReader)
		if err != nil {
			return err
		}
	} else {
		f, err := os.OpenFile(fullDstFileName, os.O_CREATE|os.O_WRONLY, readWriteFileMode)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = blobClient.DownloadFile(
			context.Background(), f, &azblob.DownloadFileOptions{
				Concurrency: uint16(maxConcurrency)})
		if err != nil {
			return err
		}
	}
	meta.resStatus = downloaded
	return nil
}

func (util *snowflakeAzureClient) extractContainerNameAndPath(location string) (*azureLocation, error) {
	stageLocation, err := expandUser(location)
	if err != nil {
		return nil, err
	}
	containerName := stageLocation
	path := ""

	if strings.Contains(stageLocation, "/") {
		containerName = stageLocation[:strings.Index(stageLocation, "/")]
		path = stageLocation[strings.Index(stageLocation, "/")+1:]
		if path != "" && !strings.HasSuffix(path, "/") {
			path += "/"
		}
	}
	return &azureLocation{containerName, path}, nil
}

func (util *snowflakeAzureClient) detectAzureTokenExpireError(resp *http.Response) bool {
	if resp.StatusCode != 403 {
		return false
	}
	azureErr, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	errStr := string(azureErr)
	return strings.Contains(errStr, "Signature not valid in the specified time frame") ||
		strings.Contains(errStr, "Server failed to authenticate the request")
}

func createContainerClient(clientURL string) (*container.Client, error) {
	return container.NewClientWithNoCredential(clientURL, &container.ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: &http.Client{
			Transport: SnowflakeTransport,
		},
	}})
}
