// Copyright (c) 2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"context"
	"errors"
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

type snowflakeS3Util struct {
}

type s3Location struct {
	bucketName string
	s3Path     string
}

func (util *snowflakeS3Util) createClient(info *execResponseStageInfo, useAccelerateEndpoint bool) cloudClient {
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
	})
}

type s3HeaderAPI interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// cloudUtil implementation
func (util *snowflakeS3Util) getFileHeader(meta *fileMetadata, filename string) *fileHeader {
	headObjInput := util.getS3Object(meta, filename)
	var s3Cli s3HeaderAPI
	s3Cli, ok := meta.client.(*s3.Client)
	if !ok {
		return nil
	}
	if meta.mockHeader != nil {
		s3Cli = meta.mockHeader
	}
	out, err := s3Cli.HeadObject(context.Background(), headObjInput)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == notFound {
				meta.resStatus = notFoundFile
				return &fileHeader{
					digest:             "",
					contentLength:      0,
					encryptionMetadata: nil,
				}
			} else if ae.ErrorCode() == expiredToken {
				meta.resStatus = renewToken
				return nil
			}
			meta.resStatus = errStatus
			meta.lastError = err
			return nil
		}
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
	return &fileHeader{
		out.Metadata[sfcDigest],
		out.ContentLength,
		&encMeta,
	}
}

type s3UploadAPI interface {
	Upload(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

// cloudUtil implementation
func (util *snowflakeS3Util) uploadFile(
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

	s3loc := util.extractBucketNameAndPath(meta.stageInfo.Location)
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
	if meta.mockUploader != nil {
		uploader = meta.mockUploader
	}

	var err error
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
		file, _ := os.Open(dataFile)
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

// cloudUtil implementation
func (util *snowflakeS3Util) nativeDownloadFile(
	meta *fileMetadata,
	fullDstFileName string,
	maxConcurrency int64) error {
	s3loc := util.extractBucketNameAndPath(meta.stageInfo.Location)
	s3path := s3loc.s3Path + strings.TrimLeft(meta.dstFileName, "/")
	client, ok := meta.client.(*s3.Client)
	if !ok {
		return &SnowflakeError{
			Message: "failed to cast to s3 client",
		}
	}

	f, err := os.OpenFile(fullDstFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	downloader := manager.NewDownloader(client, func(u *manager.Downloader) {
		u.Concurrency = int(maxConcurrency)
	})
	if _, err = downloader.Download(context.Background(), f, &s3.GetObjectInput{
		Bucket: &s3loc.bucketName,
		Key:    &s3path,
	}); err != nil {
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

func (util *snowflakeS3Util) extractBucketNameAndPath(location string) *s3Location {
	stageLocation := expandUser(location)
	bucketName := stageLocation
	s3Path := ""

	if idx := strings.Index(stageLocation, "/"); idx >= 0 {
		bucketName = stageLocation[0:idx]
		s3Path = stageLocation[idx+1:]
		if s3Path != "" && !strings.HasSuffix(s3Path, "/") {
			s3Path += "/"
		}
	}
	return &s3Location{bucketName, s3Path}
}

func (util *snowflakeS3Util) getS3Object(meta *fileMetadata, filename string) *s3.HeadObjectInput {
	s3loc := util.extractBucketNameAndPath(meta.stageInfo.Location)
	s3path := s3loc.s3Path + strings.TrimLeft(filename, "/")
	return &s3.HeadObjectInput{
		Bucket: &s3loc.bucketName,
		Key:    &s3path,
	}
}
