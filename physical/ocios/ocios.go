package ocios

import (
	"bytes"
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/physical"
)

const objectNotFound = "ObjectNotFound"

// Verify OCIOSBackend satisfies the correct interfaces
var _ physical.Backend = (*OCIOSBackend)(nil)

// OCIOSBackend is a physical backend that stores data
// within an Oracle OS bucket.
type OCIOSBackend struct {
	namespace  string
	bucket     string
	client     objectstorage.ObjectStorageClient
	logger     log.Logger
	permitPool *physical.PermitPool
}

// staticConfigProvider allows a user to simply construct a configuration provider from raw values.
type staticConfigProvider struct {
	tenancy              string
	user                 string
	region               string
	fingerprint          string
	privateKey           string
	privateKeyPassphrase *string
}

func (p staticConfigProvider) PrivateRSAKey() (key *rsa.PrivateKey, err error) {
	return common.PrivateKeyFromBytes([]byte(p.privateKey), p.privateKeyPassphrase)
}

func (p staticConfigProvider) KeyID() (keyID string, err error) {
	tenancy, err := p.TenancyOCID()
	if err != nil {
		return
	}

	user, err := p.UserOCID()
	if err != nil {
		return
	}

	fingerprint, err := p.KeyFingerprint()
	if err != nil {
		return
	}

	return fmt.Sprintf("%s/%s/%s", tenancy, user, fingerprint), nil
}

func (p staticConfigProvider) TenancyOCID() (string, error) {
	if p.tenancy == "" {
		return "", fmt.Errorf("can not get tenancy")
	}
	return p.tenancy, nil
}

func (p staticConfigProvider) UserOCID() (string, error) {
	if p.user == "" {
		return "", fmt.Errorf("can not get user")
	}
	return p.user, nil
}

func (p staticConfigProvider) KeyFingerprint() (string, error) {
	if p.fingerprint == "" {
		return "", fmt.Errorf("can not get fingerprint")
	}
	return p.fingerprint, nil
}

func (p staticConfigProvider) Region() (string, error) {
	if p.region == "" {
		return "", fmt.Errorf("can not get region")
	}
	return p.region, nil
}

// NewOCIOSBackend constructs an Oracle Object Storage backend
// using a pre-existing bucket. Credentials can be provided to the backend,
// sourced from the environment.
func NewOCIOSBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {

	var staticConfigProvider staticConfigProvider

	tenancy := os.Getenv("OCI_TENANCY_OCID")
	if tenancy == "" {
		tenancy = conf["tenancy_ocid"]
	}
	staticConfigProvider.tenancy = tenancy

	user := os.Getenv("OCI_USER_OCID")
	if user == "" {
		user = conf["user_ocid"]
	}
	staticConfigProvider.user = user

	region := os.Getenv("OCI_REGION")
	if region == "" {
		region = conf["region"]
	}
	staticConfigProvider.region = region

	fingerprint := os.Getenv("OCI_FINGERPRINT")
	if fingerprint == "" {
		fingerprint = conf["fingerprint"]
	}
	staticConfigProvider.fingerprint = fingerprint

	var privateKey string
	privateKeyPath := os.Getenv("OCI_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = conf["private_key_path"]
	} else {
		privateKeyData, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}
		privateKey = string(privateKeyData)
	}
	staticConfigProvider.privateKey = privateKey

	privateKeyPassword := os.Getenv("OCI_PRIVATE_KEY_PASSWORD")
	if privateKeyPassword == "" {
		privateKeyPassword = conf["private_key_password"]
		if privateKeyPassword != "" {
			staticConfigProvider.privateKeyPassphrase = &privateKeyPassword
		}
	}

	bucket := os.Getenv("OCI_STORAGE_BUCKET")
	if bucket == "" {
		bucket = conf["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("'bucket' must be set")
		}
	}

	configProviders := []common.ConfigurationProvider{
		staticConfigProvider,
		common.DefaultConfigProvider(),
	}

	configProvider, err := common.ComposingConfigurationProvider(configProviders)
	if err != nil {
		return nil, err
	}

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	response, err := client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		return nil, err
	}

	namespace := *response.Value

	_, err = client.GetBucket(context.Background(), objectstorage.GetBucketRequest{
		NamespaceName: &namespace,
		BucketName:    &bucket,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to access bucket %q: {{err}}", bucket), err)
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

	a := &OCIOSBackend{
		client:     client,
		namespace:  namespace,
		bucket:     bucket,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return a, nil
}

// Put is used to insert or update an entry
func (a *OCIOSBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"opcobject", "put"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	contentLength := int64(len(entry.Value))

	request := objectstorage.PutObjectRequest{
		NamespaceName: &a.namespace,
		BucketName:    &a.bucket,
		ObjectName:    &entry.Key,
		ContentLength: &contentLength,
		PutObjectBody: &CloseableBuffer{bytes.NewBuffer(entry.Value)},
	}

	_, err := a.client.PutObject(ctx, request)

	return err

}

// Get is used to fetch an entry
func (a *OCIOSBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"opcobject", "get"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	request := objectstorage.GetObjectRequest{
		NamespaceName: &a.namespace,
		BucketName:    &a.bucket,
		ObjectName:    &key,
	}

	response, err := a.client.GetObject(ctx, request)
	if err != nil {
		if serviceError, ok := common.IsServiceError(err); ok && serviceError.GetCode() == objectNotFound {
			return nil, nil
		}
		return nil, err
	}

	data := bytes.NewBuffer(nil)
	_, err = io.Copy(data, response.Content)
	defer response.Content.Close()
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
func (a *OCIOSBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"opcobject", "delete"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	request := objectstorage.DeleteObjectRequest{
		NamespaceName: &a.namespace,
		BucketName:    &a.bucket,
		ObjectName:    &key,
	}

	_, err := a.client.DeleteObject(ctx, request)

	if err != nil {
		if serviceError, ok := common.IsServiceError(err); ok && serviceError.GetCode() == objectNotFound {
			return nil
		}
	}

	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *OCIOSBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"opcobject", "list"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	keys := []string{}

	delimiter := "/"
	var start *string

	for {
		request := objectstorage.ListObjectsRequest{
			NamespaceName: &a.namespace,
			BucketName:    &a.bucket,
			Prefix:        &prefix,
			Delimiter:     &delimiter,
			Start:         start,
		}

		response, err := a.client.ListObjects(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, commonPrefix := range response.Prefixes {
			commonPrefix := strings.TrimPrefix(commonPrefix, prefix)
			keys = append(keys, commonPrefix)
		}

		for _, object := range response.Objects {
			// Add objects only from the current 'folder'
			key := strings.TrimPrefix(*object.Name, prefix)
			keys = append(keys, key)
		}

		start = response.NextStartWith
		if start == nil {
			break
		}
	}

	sort.Strings(keys)

	return keys, nil
}

// CloseableBuffer is a bytes.Buffer which implements the io.ReadCloser interface for compatibility
type CloseableBuffer struct {
	*bytes.Buffer
}

// Close implements io.ReadCloser.Close() without doing anything, because of Buffer's in-memory essence
func (CloseableBuffer) Close() error { return nil }

// Make sure CloseableBuffer implements io.ReadCloser
var _ io.ReadCloser = &CloseableBuffer{}
