package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/physical"
)

// S3Backend is a physical backend that stores data
// within an S3 bucket.
type S3Backend struct {
	bucket     string
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

	credsConfig := &awsutil.CredentialsConfig{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
	}
	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	pooledTransport := cleanhttp.DefaultPooledTransport()
	pooledTransport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	s3conn := s3.New(session.New(&aws.Config{
		Credentials: creds,
		HTTPClient: &http.Client{
			Transport: pooledTransport,
		},
		Endpoint: aws.String(endpoint),
		Region:   aws.String(region),
	}))

	_, err = s3conn.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
	if err != nil {
		return nil, fmt.Errorf("unable to access bucket '%s' in region %s: %v", bucket, region, err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("s3: max_parallel set", "max_parallel", maxParInt)
		}
	}

	s := &S3Backend{
		client:     s3conn,
		bucket:     bucket,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return s, nil
}

// Put is used to insert or update an entry
func (s *S3Backend) Put(entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"s3", "put"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(entry.Key),
		Body:   bytes.NewReader(entry.Value),
	})

	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch an entry
func (s *S3Backend) Get(key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"s3", "get"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
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

	data := make([]byte, *resp.ContentLength)
	_, err = io.ReadFull(resp.Body, data)
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: data,
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (s *S3Backend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"s3", "delete"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
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
func (s *S3Backend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"s3", "list"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	params := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	keys := []string{}

	err := s.client.ListObjectsV2Pages(params,
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
