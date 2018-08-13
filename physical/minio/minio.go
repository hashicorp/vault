package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/physical"
	minio "github.com/minio/minio-go"
)

// Verify MinioBackend satisfies the correct interfaces
var _ physical.Backend = (*MinioBackend)(nil)

type MinioBackend struct {
	bucket     string
	region     string
	client     *minio.Client
	logger     log.Logger
	permitPool *physical.PermitPool
}

func NewMinioBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = conf["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("'bucket' must be set")
		}
	}

	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = conf["region"]
		if region == "" {
			region = "us-east-1"
		}
	}

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	if accessKeyID == "" {
		accessKeyID = conf["access_key_id"]
	}

	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		secretAccessKey = conf["secret_access_key"]
	}

	useSSL := true
	disableSSLStr, ok := conf["disable_ssl"]
	if ok {
		var err error
		disableSSL, err := parseutil.ParseBool(disableSSLStr)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean set for 'disable_ssl':%q", disableSSLStr)
		}
		useSSL = !disableSSL
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		var err error
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, errwrap.Wrapf("unable to create the minio client: {{err}}", err)
	}

	m := &MinioBackend{
		bucket:     bucket,
		region:     region,
		client:     minioClient,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}

	err = m.createBucketIfNotPresent()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MinioBackend) createBucketIfNotPresent() error {
	exists, err := m.client.BucketExists(m.bucket)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to verify if bucket %q exists: {{err}}", m.bucket), err)
	}
	if exists {
		return nil
	}

	err = m.client.MakeBucket(m.bucket, m.region)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to make bucket %q in region %q: {{err}}", m.bucket, m.region), err)
	}
	return nil
}

func (m *MinioBackend) removeBucket() error {
	return m.client.RemoveBucket(m.bucket)
}

func (m *MinioBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"minio", "put"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, err := m.client.PutObject(m.bucket,
		entry.Key,
		bytes.NewReader(entry.Value),
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		errwrap.Wrapf(fmt.Sprintf("failed to put the object with key %q: {{err}}", entry.Key), err)
	}
	return nil
}

func (m *MinioBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"minio", "get"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	o, err := m.client.GetObject(m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to get the object with key %q: {{err}}", key), err)
	}

	data := bytes.NewBuffer(nil)
	_, err = io.Copy(data, o)
	if err != nil {
		return nil, nil
	}

	e := &physical.Entry{
		Key:   key,
		Value: data.Bytes(),
	}
	return e, nil
}

func (m *MinioBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"minio", "delete"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	err := m.client.RemoveObject(m.bucket, key)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to remove the object with key %q: {{err}}", key), err)
	}
	return nil
}

func (m *MinioBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"minio", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	keys := []string{}

	doneCh := make(chan struct{})
	defer close(doneCh)

	recursive := true
	objectCh := m.client.ListObjectsV2(m.bucket, prefix, recursive, doneCh)
	for o := range objectCh {
		if o.Err != nil {
			return nil, errwrap.Wrapf("failed to list the objects: {{err}}", o.Err)
		}
		key := strings.TrimPrefix(o.Key, prefix)
		keys = append(keys, key)
	}

	sort.Strings(keys)

	m.logger.Info(strings.Join(keys, " "))
	return keys, nil
}
