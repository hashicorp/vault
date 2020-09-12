package tencentcloudcos

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	PROVIDER_METRIC_KEY     = "tencentcloudcos"
	PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION         = "TENCENTCLOUD_REGION"
	PROVIDER_COS_BUCKET     = "TENCENTCLOUD_COS_BUCKET"
)

// Verify TencentCloudCOSBackend satisfies the correct interfaces
var _ physical.Backend = (*TencentCloudCOSBackend)(nil)

// TencentCloudCOSBackend is a physical backend that stores data
// with a TencentCloud COS bucket.
type TencentCloudCOSBackend struct {
	bucket     string
	client     *cos.Client
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewTencentCloudCOSBackend constructs a cos backend using a pre-existing bucket.
// Credentials can be provided to the backend, sourced from the environment.
func NewTencentCloudCOSBackend(config map[string]string, logger log.Logger) (physical.Backend, error) {
	if config == nil {
		config = map[string]string{}
	}

	var (
		bucket       string
		region       string
		accessKey    string
		secretKey    string
		sessionToken string
	)

	switch {
	case os.Getenv(PROVIDER_COS_BUCKET) != "":
		bucket = os.Getenv(PROVIDER_COS_BUCKET)
	case config["bucket"] != "":
		bucket = config["bucket"]
	default:
		return nil, fmt.Errorf("'bucket' must be set")
	}

	switch {
	case os.Getenv(PROVIDER_REGION) != "":
		region = os.Getenv(PROVIDER_REGION)
	case config["region"] != "":
		region = config["region"]
	default:
		region = "ap-guangzhou"
	}

	switch {
	case os.Getenv(PROVIDER_SECRET_ID) != "":
		accessKey = os.Getenv(PROVIDER_SECRET_ID)
	case config["access_key"] != "":
		accessKey = config["access_key"]
	default:
		return nil, fmt.Errorf("'access_key' must be set")
	}

	switch {
	case os.Getenv(PROVIDER_SECRET_KEY) != "":
		secretKey = os.Getenv(PROVIDER_SECRET_KEY)
	case config["secret_key"] != "":
		secretKey = config["secret_key"]
	default:
		return nil, fmt.Errorf("'secret_key' must be set")
	}

	switch {
	case os.Getenv(PROVIDER_SECURITY_TOKEN) != "":
		sessionToken = os.Getenv(PROVIDER_SECURITY_TOKEN)
	case config["session_token"] != "":
		sessionToken = config["session_token"]
	default:
		sessionToken = ""
	}

	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	if err != nil {
		return nil, err
	}

	client := cos.NewClient(
		&cos.BaseURL{BucketURL: u},
		&http.Client{
			Timeout: 60 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:     accessKey,
				SecretKey:    secretKey,
				SessionToken: sessionToken,
			},
		},
	)

	_, rsp, err := client.Bucket.Get(context.Background(), nil)
	if rsp == nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to access bucket %q: {{err}}", bucket), fmt.Errorf("no response"))
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 404 {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to access bucket %q: {{err}}", bucket), fmt.Errorf("bucket not exists"))
	}

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to access bucket %q: {{err}}", bucket), err)
	}

	var maxParallel int
	if v, ok := config["max_parallel"]; ok {
		maxParallel, err = strconv.Atoi(v)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParallel)
		}
	}

	a := &TencentCloudCOSBackend{
		client:     client,
		bucket:     bucket,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParallel),
	}

	return a, nil
}

// Put is used to insert or update an entry
func (a *TencentCloudCOSBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{PROVIDER_METRIC_KEY, "put"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosMetaXXX: &http.Header{
				"X-Cos-Meta-Md5": []string{fmt.Sprintf("%x", md5.Sum(entry.Value))},
			},
		},
	}

	r := bytes.NewReader(entry.Value)
	rsp, err := a.client.Object.Put(context.Background(), entry.Key, r, opt)
	if rsp == nil {
		return fmt.Errorf("failed to save file to %v: %v", entry.Key, fmt.Errorf("no response"))
	}
	defer rsp.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to save file to %v: %v", entry.Key, err)
	}

	return nil
}

// Get is used to fetch an entry
func (a *TencentCloudCOSBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{PROVIDER_METRIC_KEY, "get"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	rsp, err := a.client.Object.Get(context.Background(), key, nil)
	if rsp == nil {
		return nil, fmt.Errorf("failed to save file to %v: %v", key, fmt.Errorf("no response"))
	}
	defer rsp.Body.Close()

	if err != nil {
		if rsp.StatusCode == 404 {
			err = nil
		} else {
			err = fmt.Errorf("failed to open file at %v: %v", key, err)
		}
		return nil, err
	}

	checksum := rsp.Header.Get("X-Cos-Meta-Md5")
	if len(checksum) != 32 {
		return nil, fmt.Errorf("failed to open file at %v: checksum %s invalid", key, checksum)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at %v: %v", key, err)
	}

	check := fmt.Sprintf("%x", md5.Sum(data))
	if check != checksum {
		return nil, fmt.Errorf("failed to open file at %v: checksum mismatch, %s != %s", key, check, checksum)
	}

	ent := &physical.Entry{
		Key:   key,
		Value: data,
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (a *TencentCloudCOSBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{PROVIDER_METRIC_KEY, "delete"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	rsp, err := a.client.Object.Delete(context.Background(), key)
	if rsp == nil {
		return fmt.Errorf("failed to delete file %v: %v", key, fmt.Errorf("no response"))
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 404 {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to delete file %v: %v", key, err)
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *TencentCloudCOSBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{PROVIDER_METRIC_KEY, "list"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	marker := ""
	keys := []string{}

	for {
		fs, rsp, err := a.client.Bucket.Get(context.Background(), &cos.BucketGetOptions{Prefix: prefix, Delimiter: "/", Marker: marker})
		if rsp == nil {
			return nil, fmt.Errorf("failed to list bucket %v(%v:%v) : %v", a.bucket, prefix, marker, fmt.Errorf("no response"))
		}
		defer rsp.Body.Close()

		if rsp.StatusCode == 404 {
			return nil, fmt.Errorf("failed to list bucket %v(%v:%v) : %v", a.bucket, prefix, marker, fmt.Errorf("bucket not exists"))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list bucket %v(%v:%v) : %v", a.bucket, prefix, marker, err)
		}

		for _, commonPrefix := range fs.CommonPrefixes {
			commonPrefix = strings.TrimPrefix(commonPrefix, prefix)
			keys = append(keys, commonPrefix)
		}

		for _, object := range fs.Contents {
			key := strings.TrimPrefix(object.Key, prefix)
			keys = append(keys, key)
		}

		if !fs.IsTruncated {
			break
		}

		marker = fs.NextMarker
	}

	sort.Strings(keys)

	return keys, nil
}
