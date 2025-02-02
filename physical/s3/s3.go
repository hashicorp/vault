// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package s3

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify S3Backend satisfies the correct interfaces
var _ physical.Backend = (*S3Backend)(nil)

// S3Backend is a physical backend that stores data
// within an S3 bucket.
type S3Backend struct {
	context    context.Context
	bucket     string
	path       string
	kmsKeyId   string
	client     *s3.Client
	haEnabled  bool
	logger     log.Logger
	permitPool *permitpool.Pool
}

// NewS3Backend constructs a S3 backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment, AWS credential files or by IAM role.
func NewS3Backend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	ctx := context.TODO()
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

	s3ForcePathStyleBool, err := parseutil.ParseBool(conf["s3_force_path_style"])
	if err != nil {
		return nil, fmt.Errorf("invalid boolean set for s3_force_path_style: %q", conf["s3_force_path_style"])
	}
	disableSSLBool, err := parseutil.ParseBool(conf["disable_ssl"])
	if err != nil {
		return nil, fmt.Errorf("invalid boolean set for disable_ssl: %q", conf["disable_ssl"])
	}

	// Set custom transport
	pooledTransport := cleanhttp.DefaultPooledTransport()
	pooledTransport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithHTTPClient(&http.Client{Transport: pooledTransport}),
	)
	if accessKey != "" && secretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Set S3-specific configurations
	s3Options := func(o *s3.Options) {
		o.UsePathStyle = s3ForcePathStyleBool
		if disableSSLBool {
			o.EndpointOptions.DisableHTTPS = true
		}
		if endpoint != "" {
			o.EndpointResolver = s3.EndpointResolverFromURL(endpoint)
		}
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg, s3Options)

	_, err = s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to access bucket %q in region %q: %w", bucket, region, err)
	}

	maxParInt, err := strconv.Atoi(conf["max_parallel"])
	if err != nil && conf["max_parallel"] != "" {
		return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
	}
	if logger.IsDebug() {
		logger.Debug("max_parallel set", "max_parallel", maxParInt)
	}

	kmsKeyId := conf["kms_key_id"]

	// envHAEnabled is the name of the environment variable to search for the
	// boolean indicating if HA is enabled.
	haEnabled := os.Getenv("S3_STORAGE_HA_ENABLED")
	if haEnabled == "" {
		haEnabled = conf["ha_enabled"]
	}
	haEnabledBool, err := strconv.ParseBool(haEnabled)
	if err != nil && conf["ha_enabled"] != "" {
		return nil, fmt.Errorf("failed to parse ha_enabled value: %w", err)
	}

	s := &S3Backend{
		client:     s3Client,
		context:    ctx,
		bucket:     bucket,
		path:       path,
		kmsKeyId:   kmsKeyId,
		logger:     logger,
		permitPool: permitpool.New(maxParInt),
		haEnabled:  haEnabledBool,
	}
	return s, nil
}

// Put is used to insert or update an entry
func (s *S3Backend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"s3", "put"}, time.Now())

	if err := s.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer s.permitPool.Release()

	// Setup key
	key := path.Join(s.path, entry.Key)

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(entry.Value),
	}

	if s.kmsKeyId != "" {
		putObjectInput.ServerSideEncryption = types.ServerSideEncryptionAwsKms
		putObjectInput.SSEKMSKeyId = aws.String(s.kmsKeyId)
	}

	_, err := s.client.PutObject(ctx, putObjectInput)
	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch an entry
func (s *S3Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"s3", "get"}, time.Now())

	if err := s.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer s.permitPool.Release()

	// Setup key
	key = path.Join(s.path, key)

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	var nsk *types.NoSuchKey
	if err != nil {
		if errors.As(err, &nsk) {
			return nil, nil
		}
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

	if err := s.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer s.permitPool.Release()

	// Setup key
	key = path.Join(s.path, key)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
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

	if err := s.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer s.permitPool.Release()

	// Setup prefix
	prefix = path.Join(s.path, prefix)

	// Validate prefix (if present) is ending with a "/"
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	var keys []string

	// Create paginator
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	// Iterate over all the pages
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		// Add truncated 'folder' paths
		for _, commonPrefix := range page.CommonPrefixes {
			// Avoid panic
			if commonPrefix.Prefix == nil {
				continue
			}

			commonPrefix := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
			keys = append(keys, commonPrefix)
		}

		// Add objects only from the current 'folder'
		for _, key := range page.Contents {
			// Avoid panic
			if key.Key == nil {
				continue
			}

			key := strings.TrimPrefix(*key.Key, prefix)
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys, nil
}
