package alicloudoss

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

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
)

// Verify AliCloudOSSBackend satisfies the correct interfaces
var _ physical.Backend = (*AliCloudOSSBackend)(nil)

// AliCloudOSSBackend is a physical backend that stores data
// within an Alibaba OSS bucket.
type AliCloudOSSBackend struct {
	bucket     string
	client     *oss.Client
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewAliCloudOSSBackend constructs an OSS backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment.
func NewAliCloudOSSBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	endpoint := os.Getenv("ALICLOUD_OSS_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
		if endpoint == "" {
			return nil, fmt.Errorf("'endpoint' must be set")
		}
	}

	bucket := os.Getenv("ALICLOUD_OSS_BUCKET")
	if bucket == "" {
		bucket = conf["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("'bucket' must be set")
		}
	}

	accessKeyID := os.Getenv("ALICLOUD_ACCESS_KEY")
	if accessKeyID == "" {
		accessKeyID = conf["access_key"]
		if accessKeyID == "" {
			return nil, fmt.Errorf("'access_key' must be set")
		}
	}

	accessKeySecret := os.Getenv("ALICLOUD_SECRET_KEY")
	if accessKeySecret == "" {
		accessKeySecret = conf["secret_key"]
		if accessKeySecret == "" {
			return nil, fmt.Errorf("'secret_key' must be set")
		}
	}

	options := func(c *oss.Client) {
		c.Config.Timeout = 30
	}

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret, options)
	if err != nil {
		return nil, err
	}

	bucketObj, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	_, err = bucketObj.ListObjects()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to access bucket %q at endpoint %q: {{err}}", bucket, endpoint), err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	a := &AliCloudOSSBackend{
		client:     client,
		bucket:     bucket,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return a, nil
}

// Put is used to insert or update an entry
func (a *AliCloudOSSBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"alibaba", "put"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	bucket, err := a.client.Bucket(a.bucket)
	if err != nil {
		return err
	}

	return bucket.PutObject(entry.Key, bytes.NewReader(entry.Value))

}

// Get is used to fetch an entry
func (a *AliCloudOSSBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"alibaba", "get"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	bucket, err := a.client.Bucket(a.bucket)
	if err != nil {
		return nil, err
	}

	object, err := bucket.GetObject(key)
	if err != nil {
		switch err := err.(type) {
		case oss.ServiceError:
			if err.StatusCode == 404 && err.Code == "NoSuchKey" {
				return nil, nil
			}
		}
		return nil, err
	}

	data := bytes.NewBuffer(nil)
	_, err = io.Copy(data, object)
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: data.Bytes(),
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (a *AliCloudOSSBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"alibaba", "delete"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	bucket, err := a.client.Bucket(a.bucket)
	if err != nil {
		return err
	}

	return bucket.DeleteObject(key)
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AliCloudOSSBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"alibaba", "list"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	keys := []string{}

	bucket, err := a.client.Bucket(a.bucket)
	if err != nil {
		return nil, err
	}

	marker := oss.Marker("")

	for {
		result, err := bucket.ListObjects(oss.Prefix(prefix), oss.Delimiter("/"), marker)
		if err != nil {
			return nil, err
		}

		for _, commonPrefix := range result.CommonPrefixes {

			commonPrefix := strings.TrimPrefix(commonPrefix, prefix)
			keys = append(keys, commonPrefix)
		}

		for _, object := range result.Objects {
			// Add objects only from the current 'folder'
			key := strings.TrimPrefix(object.Key, prefix)
			keys = append(keys, key)
		}

		if !result.IsTruncated {
			break
		}

		marker = oss.Marker(result.NextMarker)
	}

	sort.Strings(keys)

	return keys, nil
}
