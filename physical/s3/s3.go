// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify S3Backend satisfies the correct interfaces
var _ physical.Backend = (*S3Backend)(nil)

// S3Backend is a physical backend that stores data
// within an S3 bucket.
type S3Backend struct {
	bucket     string
	path       string
	kmsKeyId   string
	client     *s3.S3
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewS3Backend constructs a S3 backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment, AWS credential files or by IAM role.
func NewS3Backend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		bucket = conf["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("'bucket' must be set")
		}
	}

	path := conf["path"]

	accessKey, ok := conf["access_key"]
	if !ok {
		accessKey = ""
	}
	secretKey, ok := conf["secret_key"]
	if !ok {
		secretKey = ""
	}
	sessionToken, ok := conf["session_token"]
	if !ok {
		sessionToken = ""
	}
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = conf["region"]
			if region == "" {
				region = "us-east-1"
			}
		}
	}
	s3ForcePathStyleStr, ok := conf["s3_force_path_style"]
	if !ok {
		s3ForcePathStyleStr = "false"
	}
	s3ForcePathStyleBool, err := parseutil.ParseBool(s3ForcePathStyleStr)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean set for s3_force_path_style: %q", s3ForcePathStyleStr)
	}
	disableSSLStr, ok := conf["disable_ssl"]
	if !ok {
		disableSSLStr = "false"
	}
	disableSSLBool, err := parseutil.ParseBool(disableSSLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean set for disable_ssl: %q", disableSSLStr)
	}

	credsConfig := &awsutil.CredentialsConfig{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
		Logger:       logger,
	}
	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	pooledTransport := cleanhttp.DefaultPooledTransport()
	pooledTransport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		HTTPClient: &http.Client{
			Transport: pooledTransport,
		},
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(s3ForcePathStyleBool),
		DisableSSL:       aws.Bool(disableSSLBool),
	})
	if err != nil {
		return nil, err
	}
	s3conn := s3.New(sess)

	_, err = s3conn.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
	if err != nil {
		return nil, fmt.Errorf("unable to access bucket %q in region %q: %w", bucket, region, err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	kmsKeyId, ok := conf["kms_key_id"]
	if !ok {
		kmsKeyId = ""
	}

	s := &S3Backend{
		client:     s3conn,
		bucket:     bucket,
		path:       path,
		kmsKeyId:   kmsKeyId,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return s, nil
}

// Put is used to insert or update an entry
func (s *S3Backend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"s3", "put"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	// Setup key
	key := path.Join(s.path, entry.Key)

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(entry.Value),
	}

	if s.kmsKeyId != "" {
		putObjectInput.ServerSideEncryption = aws.String("aws:kms")
		putObjectInput.SSEKMSKeyId = aws.String(s.kmsKeyId)
	}

	_, err := s.client.PutObjectWithContext(ctx, putObjectInput)
	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch an entry
func (s *S3Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"s3", "get"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	// Setup key
	key = path.Join(s.path, key)

	resp, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if awsErr, ok := err.(awserr.RequestFailure); ok {
		// Return nil on 404s, error on anything else
		if awsErr.StatusCode() == 404 {
			return nil, nil
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("got nil response from S3 but no error")
	}

	data := bytes.NewBuffer(nil)
	if resp.ContentLength != nil {
		data = bytes.NewBuffer(make([]byte, 0, *resp.ContentLength))
	}
	_, err = io.Copy(data, resp.Body)
	if err != nil {
		return nil, err
	}

	// Strip path prefix
	if s.path != "" {
		key = strings.TrimPrefix(key, s.path+"/")
	}

	ent := &physical.Entry{
		Key:   key,
		Value: data.Bytes(),
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (s *S3Backend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"s3", "delete"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	// Setup key
	key = path.Join(s.path, key)

	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (s *S3Backend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"s3", "list"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	// Setup prefix
	prefix = path.Join(s.path, prefix)

	// Validate prefix (if present) is ending with a "/"
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	params := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	keys := []string{}

	err := s.client.ListObjectsV2PagesWithContext(ctx, params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			if page != nil {
				// Add truncated 'folder' paths
				for _, commonPrefix := range page.CommonPrefixes {
					// Avoid panic
					if commonPrefix == nil {
						continue
					}

					commonPrefix := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
					keys = append(keys, commonPrefix)
				}
				// Add objects only from the current 'folder'
				for _, key := range page.Contents {
					// Avoid panic
					if key == nil {
						continue
					}

					key := strings.TrimPrefix(*key.Key, prefix)
					keys = append(keys, key)
				}
			}
			return true
		})
	if err != nil {
		return nil, err
	}

	sort.Strings(keys)

	return keys, nil
}
