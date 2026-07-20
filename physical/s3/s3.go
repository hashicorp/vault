// Copyright IBM Corp. 2016, 2026
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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	awsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
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
	bucket     string
	path       string
	kmsKeyId   string
	client     *s3.Client
	logger     log.Logger
	permitPool *permitpool.Pool
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
	region := resolveS3Region(conf)
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

	pooledTransport := cleanhttp.DefaultPooledTransport()
	pooledTransport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	cfg, err := buildS3CredentialChain(accessKey, secretKey, sessionToken, region, logger, &http.Client{Transport: pooledTransport})
	if err != nil {
		return nil, err
	}

	// Build the endpoint URL. disable_ssl only makes sense alongside a
	// custom endpoint (e.g. MinIO); real AWS S3 requires HTTPS.
	endpointURL := endpoint
	if disableSSLBool {
		if endpoint == "" {
			logger.Warn("disable_ssl is set without a custom endpoint; ignoring because real AWS S3 requires HTTPS")
		} else if strings.HasPrefix(endpoint, "https://") {
			// Endpoint scheme is https:// but disable_ssl is true; switch to http://.
			endpointURL = "http://" + strings.TrimPrefix(endpoint, "https://")
		} else if !strings.HasPrefix(endpoint, "http://") {
			endpointURL = "http://" + endpoint
		}
	} else if endpoint != "" && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpointURL = "https://" + endpoint
	}

	s3conn := s3.NewFromConfig(*cfg, func(o *s3.Options) {
		if endpointURL != "" {
			o.BaseEndpoint = aws.String(endpointURL)
		}
		o.UsePathStyle = s3ForcePathStyleBool
	})

	_, err = s3conn.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{Bucket: &bucket})
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
		permitPool: permitpool.New(maxParInt),
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
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			// Return nil on 404s
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

	params := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	keys := []string{}

	paginator := s3.NewListObjectsV2Paginator(s.client, params)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		// Add truncated 'folder' paths
		for _, commonPrefix := range page.CommonPrefixes {
			keys = append(keys, strings.TrimPrefix(aws.ToString(commonPrefix.Prefix), prefix))
		}
		// Add objects only from the current 'folder'
		for _, obj := range page.Contents {
			keys = append(keys, strings.TrimPrefix(aws.ToString(obj.Key), prefix))
		}
	}

	sort.Strings(keys)

	return keys, nil
}

// resolveS3Region returns the AWS region for the S3 backend following the documented
// priority order from configuration/storage/s3.mdx:
// AWS_REGION env var → AWS_DEFAULT_REGION env var → region conf key → "us-east-1".
func resolveS3Region(conf map[string]string) string {
	if r := os.Getenv("AWS_REGION"); r != "" {
		return r
	}
	if r := os.Getenv("AWS_DEFAULT_REGION"); r != "" {
		return r
	}
	if r := conf["region"]; r != "" {
		return r
	}
	return "us-east-1"
}

// buildS3CredentialChain constructs the AWS credential chain for the S3 backend.
// It mirrors the credential resolution order documented at configuration/storage/s3.mdx:
// static config → environment variables → credential files → instance role.
// WithSharedCredentials(false) prevents awsutil from injecting an empty credentials-file
// path that would cause LoadDefaultConfig to fail when a [default] profile exists in
// ~/.aws/config alongside explicit static keys.
func buildS3CredentialChain(accessKey, secretKey, sessionToken, region string, logger log.Logger, httpClient *http.Client) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
		Logger:       logger,
		Region:       region,
		HTTPClient:   httpClient,
	}
	return credsConfig.GenerateCredentialChain(context.Background(), awsutil.WithSharedCredentials(false))
}
