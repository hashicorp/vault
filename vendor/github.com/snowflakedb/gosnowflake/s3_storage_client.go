// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

const (
	sfcDigest  = "sfc-digest"
	amzMatdesc = "x-amz-matdesc"
	amzKey     = "x-amz-key"
	amzIv      = "x-amz-iv"

	notFound             = "NotFound"
	expiredToken         = "ExpiredToken"
	errNoWsaeconnaborted = "10053"
)

type snowflakeS3Client struct {
}

type s3Location struct {
	bucketName string
	s3Path     string
}

func (util *snowflakeS3Client) createClient(info *execResponseStageInfo, useAccelerateEndpoint bool) (cloudClient, error) {
	stageCredentials := info.Creds
	var resolver s3.EndpointResolver
	if info.EndPoint != "" {
		resolver = s3.EndpointResolverFromURL("https://" + info.EndPoint) // FIPS endpoint
	}

	return s3.New(s3.Options{
		Region: info.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			stageCredentials.AwsKeyID,
			stageCredentials.AwsSecretKey,
			stageCredentials.AwsToken)),
		EndpointResolver: resolver,
		UseAccelerate:    useAccelerateEndpoint,
		HTTPClient: &http.Client{
			Transport: SnowflakeTransport,
		},
	}), nil
}

type s3HeaderAPI interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// cloudUtil implementation
func (util *snowflakeS3Client) getFileHeader(meta *fileMetadata, filename string) (*fileHeader, error) {
	headObjInput, err := util.getS3Object(meta, filename)
	if err != nil {
		return nil, err
	}
	var s3Cli s3HeaderAPI
	s3Cli, ok := meta.client.(*s3.Client)
	if !ok {
		return nil, fmt.Errorf("could not parse client to s3.Client")
	}
	// for testing only
	if meta.mockHeader != nil {
		s3Cli = meta.mockHeader
	}
	out, err := s3Cli.HeadObject(context.Background(), headObjInput)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == notFound {
				meta.resStatus = notFoundFile
				return nil, fmt.Errorf("could not find file")
			} else if ae.ErrorCode() == expiredToken {
				meta.resStatus = renewToken
				return nil, fmt.Errorf("received expired token. renewing")
			}
			meta.resStatus = errStatus
			meta.lastError = err
			return nil, fmt.Errorf("error while retrieving header")
		}
		meta.resStatus = errStatus
		meta.lastError = err
		return nil, fmt.Errorf("unexpected error while retrieving header: %v", err)
	}

	meta.resStatus = uploaded
	var encMeta encryptMetadata
	if out.Metadata[amzKey] != "" {
		encMeta = encryptMetadata{
			out.Metadata[amzKey],
			out.Metadata[amzIv],
			out.Metadata[amzMatdesc],
		}
	}
	contentLength := convertContentLength(out.ContentLength)
	return &fileHeader{
		out.Metadata[sfcDigest],
		contentLength,
		&encMeta,
	}, nil
}

// SNOW-974548 remove this function after upgrading AWS SDK
func convertContentLength(contentLength any) int64 {
	switch t := contentLength.(type) {
	case int64:
		return t
	case *int64:
		if t != nil {
			return *t
		}
	}
	return 0
}

type s3UploadAPI interface {
	Upload(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

// cloudUtil implementation
func (util *snowflakeS3Client) uploadFile(
	dataFile string,
	meta *fileMetadata,
	encryptMeta *encryptMetadata,
	maxConcurrency int,
	multiPartThreshold int64) error {
	s3Meta := map[string]string{
		httpHeaderContentType: httpHeaderValueOctetStream,
		sfcDigest:             meta.sha256Digest,
	}
	if encryptMeta != nil {
		s3Meta[amzIv] = encryptMeta.iv
		s3Meta[amzKey] = encryptMeta.key
		s3Meta[amzMatdesc] = encryptMeta.matdesc
	}

	s3loc, err := util.extractBucketNameAndPath(meta.stageInfo.Location)
	if err != nil {
		return err
	}
	s3path := s3loc.s3Path + strings.TrimLeft(meta.dstFileName, "/")

	client, ok := meta.client.(*s3.Client)
	if !ok {
		return &SnowflakeError{
			Message: "failed to cast to s3 client",
		}
	}
	var uploader s3UploadAPI
	uploader = manager.NewUploader(client, func(u *manager.Uploader) {
		u.Concurrency = maxConcurrency
		u.PartSize = int64Max(multiPartThreshold, manager.DefaultUploadPartSize)
	})
	// for testing only
	if meta.mockUploader != nil {
		uploader = meta.mockUploader
	}

	if meta.srcStream != nil {
		uploadStream := meta.srcStream
		if meta.realSrcStream != nil {
			uploadStream = meta.realSrcStream
		}
		_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
			Bucket:   &s3loc.bucketName,
			Key:      &s3path,
			Body:     bytes.NewBuffer(uploadStream.Bytes()),
			Metadata: s3Meta,
		})
	} else {
		var file *os.File
		file, err = os.Open(dataFile)
		if err != nil {
			return err
		}
		_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
			Bucket:   &s3loc.bucketName,
			Key:      &s3path,
			Body:     file,
			Metadata: s3Meta,
		})
	}

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == expiredToken {
				meta.resStatus = renewToken
				return err
			} else if strings.Contains(ae.ErrorCode(), errNoWsaeconnaborted) {
				meta.lastError = err
				meta.resStatus = needRetryWithLowerConcurrency
				return err
			}
		}
		meta.lastError = err
		meta.resStatus = needRetry
		return err
	}
	meta.dstFileSize = meta.uploadSize
	meta.resStatus = uploaded
	return nil
}

type s3DownloadAPI interface {
	Download(ctx context.Context, w io.WriterAt, params *s3.GetObjectInput, optFns ...func(*manager.Downloader)) (int64, error)
}

// cloudUtil implementation
func (util *snowflakeS3Client) nativeDownloadFile(
	meta *fileMetadata,
	fullDstFileName string,
	maxConcurrency int64) error {
	s3Obj, _ := util.getS3Object(meta, meta.srcFileName)
	client, ok := meta.client.(*s3.Client)
	if !ok {
		return &SnowflakeError{
			Message: "failed to cast to s3 client",
		}
	}

	f, err := os.OpenFile(fullDstFileName, os.O_CREATE|os.O_WRONLY, readWriteFileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	var downloader s3DownloadAPI
	downloader = manager.NewDownloader(client, func(u *manager.Downloader) {
		u.Concurrency = int(maxConcurrency)
	})
	// for testing only
	if meta.mockDownloader != nil {
		downloader = meta.mockDownloader
	}

	if meta.options.GetFileToStream {
		buf := manager.NewWriteAtBuffer([]byte{})
		_, err = downloader.Download(context.Background(), buf, &s3.GetObjectInput{
			Bucket: s3Obj.Bucket,
			Key:    s3Obj.Key,
		})
		meta.dstStream = bytes.NewBuffer(buf.Bytes())
	} else {
		_, err = downloader.Download(context.Background(), f, &s3.GetObjectInput{
			Bucket: s3Obj.Bucket,
			Key:    s3Obj.Key,
		})
	}

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == expiredToken {
				meta.resStatus = renewToken
				return err
			} else if strings.Contains(ae.ErrorCode(), errNoWsaeconnaborted) {
				meta.lastError = err
				meta.resStatus = needRetryWithLowerConcurrency
				return err
			}
			meta.lastError = err
			meta.resStatus = errStatus
			return err
		}
		meta.lastError = err
		meta.resStatus = needRetry
		return err
	}
	meta.resStatus = downloaded
	return nil
}

func (util *snowflakeS3Client) extractBucketNameAndPath(location string) (*s3Location, error) {
	stageLocation, err := expandUser(location)
	if err != nil {
		return nil, err
	}
	bucketName := stageLocation
	s3Path := ""

	if idx := strings.Index(stageLocation, "/"); idx >= 0 {
		bucketName = stageLocation[0:idx]
		s3Path = stageLocation[idx+1:]
		if s3Path != "" && !strings.HasSuffix(s3Path, "/") {
			s3Path += "/"
		}
	}
	return &s3Location{bucketName, s3Path}, nil
}

func (util *snowflakeS3Client) getS3Object(meta *fileMetadata, filename string) (*s3.HeadObjectInput, error) {
	s3loc, err := util.extractBucketNameAndPath(meta.stageInfo.Location)
	if err != nil {
		return nil, err
	}
	s3path := s3loc.s3Path + strings.TrimLeft(filename, "/")
	return &s3.HeadObjectInput{
		Bucket: &s3loc.bucketName,
		Key:    &s3path,
	}, nil
}
