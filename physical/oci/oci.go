// Copyright © 2019, Oracle and/or its affiliates.
package oci

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"golang.org/x/net/context"
)

// Verify Backend satisfies the correct interfaces
var _ physical.Backend = (*Backend)(nil)

const (
	// Limits maximum outstanding requests
	MaxNumberOfPermits = 256
)

var (
	metricDelete     = []string{"oci", "delete"}
	metricGet        = []string{"oci", "get"}
	metricList       = []string{"oci", "list"}
	metricPut        = []string{"oci", "put"}
	metricDeleteFull = []string{"oci", "deleteFull"}
	metricGetFull    = []string{"oci", "getFull"}
	metricListFull   = []string{"oci", "listFull"}
	metricPutFull    = []string{"oci", "putFull"}

	metricDeleteHa = []string{"oci", "deleteHa"}
	metricGetHa    = []string{"oci", "getHa"}
	metricPutHa    = []string{"oci", "putHa"}

	metricDeleteAcquirePool = []string{"oci", "deleteAcquirePool"}
	metricGetAcquirePool    = []string{"oci", "getAcquirePool"}
	metricListAcquirePool   = []string{"oci", "listAcquirePool"}
	metricPutAcquirePool    = []string{"oci", "putAcquirePool"}

	metricDeleteFailed         = []string{"oci", "deleteFailed"}
	metricGetFailed            = []string{"oci", "getFailed"}
	metricListFailed           = []string{"oci", "listFailed"}
	metricPutFailed            = []string{"oci", "putFailed"}
	metricHaWatchLockRetriable = []string{"oci", "haWatchLockRetriable"}
	metricPermitsUsed          = []string{"oci", "permitsUsed"}

	metric5xx = []string{"oci", "5xx"}
)

type Backend struct {
	client         *objectstorage.ObjectStorageClient
	bucketName     string
	logger         log.Logger
	permitPool     *permitpool.Pool
	namespaceName  string
	haEnabled      bool
	lockBucketName string
}

func NewBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	bucketName := conf["bucket_name"]
	if bucketName == "" {
		return nil, errors.New("missing bucket name")
	}

	namespaceName := conf["namespace_name"]
	if bucketName == "" {
		return nil, errors.New("missing namespace name")
	}

	lockBucketName := ""
	haEnabled := false
	var err error
	haEnabledStr := conf["ha_enabled"]
	if haEnabledStr != "" {
		haEnabled, err = strconv.ParseBool(haEnabledStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HA enabled: %w", err)
		}

		if haEnabled {
			lockBucketName = conf["lock_bucket_name"]
			if lockBucketName == "" {
				return nil, errors.New("missing lock bucket name")
			}
		}
	}

	authTypeAPIKeyBool := false
	authTypeAPIKeyStr := conf["auth_type_api_key"]
	if authTypeAPIKeyStr != "" {
		authTypeAPIKeyBool, err = strconv.ParseBool(authTypeAPIKeyStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing auth_type_api_key parameter: %w", err)
		}
	}

	var cp common.ConfigurationProvider
	if authTypeAPIKeyBool {
		cp = common.DefaultConfigProvider()
	} else {
		cp, err = auth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, fmt.Errorf("failed creating InstancePrincipalConfigurationProvider: %w", err)
		}
	}

	objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(cp)
	if err != nil {
		return nil, fmt.Errorf("failed creating NewObjectStorageClientWithConfigurationProvider: %w", err)
	}

	region := conf["region"]
	if region != "" {
		objectStorageClient.SetRegion(region)
	}

	logger.Debug("configuration",
		"bucket_name", bucketName,
		"region", region,
		"namespace_name", namespaceName,
		"ha_enabled", haEnabled,
		"lock_bucket_name", lockBucketName,
		"auth_type_api_key", authTypeAPIKeyBool,
	)

	return &Backend{
		client:         &objectStorageClient,
		bucketName:     bucketName,
		logger:         logger,
		permitPool:     permitpool.New(MaxNumberOfPermits),
		namespaceName:  namespaceName,
		haEnabled:      haEnabled,
		lockBucketName: lockBucketName,
	}, nil
}

func (o *Backend) Put(ctx context.Context, entry *physical.Entry) error {
	o.logger.Debug("PUT started")
	defer metrics.MeasureSince(metricPutFull, time.Now())
	startAcquirePool := time.Now()
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	if err := o.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer o.permitPool.Release()
	metrics.MeasureSince(metricPutAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricPut, time.Now())
	size := int64(len(entry.Value))
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		metrics.IncrCounter(metricPutFailed, 1)
		o.logger.Error("failed to generate UUID")
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	o.logger.Debug("PUT", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.PutObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &entry.Key,
		ContentLength:      &size,
		PutObjectBody:      io.NopCloser(bytes.NewReader(entry.Value)),
		OpcMeta:            nil,
		OpcClientRequestId: &opcClientRequestId,
	}

	resp, err := o.client.PutObject(ctx, request)
	if resp.RawResponse != nil && resp.RawResponse.Body != nil {
		defer resp.RawResponse.Body.Close()
	}

	if err != nil {
		metrics.IncrCounter(metricPutFailed, 1)
		return fmt.Errorf("failed to put data: %w", err)
	}

	o.logRequest("PUT", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)
	o.logger.Debug("PUT completed")

	return nil
}

func (o *Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	o.logger.Debug("GET started")
	defer metrics.MeasureSince(metricGetFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	if err := o.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer o.permitPool.Release()
	metrics.MeasureSince(metricGetAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricGet, time.Now())
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		o.logger.Error("failed to generate UUID")
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}
	o.logger.Debug("GET", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.GetObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &key,
		OpcClientRequestId: &opcClientRequestId,
	}

	resp, err := o.client.GetObject(ctx, request)
	if resp.RawResponse != nil && resp.RawResponse.Body != nil {
		defer resp.RawResponse.Body.Close()
	}
	o.logRequest("GET", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)

	if err != nil {
		if resp.RawResponse != nil && resp.RawResponse.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		metrics.IncrCounter(metricGetFailed, 1)
		return nil, fmt.Errorf("failed to read Value: %w", err)
	}

	body, err := io.ReadAll(resp.Content)
	if err != nil {
		metrics.IncrCounter(metricGetFailed, 1)
		return nil, fmt.Errorf("failed to decode Value into bytes: %w", err)
	}

	o.logger.Debug("GET completed")

	return &physical.Entry{
		Key:   key,
		Value: body,
	}, nil
}

func (o *Backend) Delete(ctx context.Context, key string) error {
	o.logger.Debug("DELETE started")
	defer metrics.MeasureSince(metricDeleteFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	if err := o.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer o.permitPool.Release()
	metrics.MeasureSince(metricDeleteAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricDelete, time.Now())
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		o.logger.Error("Delete: error generating UUID")
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	o.logger.Debug("Delete", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.DeleteObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &key,
		OpcClientRequestId: &opcClientRequestId,
	}

	resp, err := o.client.DeleteObject(ctx, request)
	if resp.RawResponse != nil && resp.RawResponse.Body != nil {
		defer resp.RawResponse.Body.Close()
	}

	o.logRequest("DELETE", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)

	if err != nil {
		if resp.RawResponse != nil && resp.RawResponse.StatusCode == http.StatusNotFound {
			return nil
		}
		metrics.IncrCounter(metricDeleteFailed, 1)
		return fmt.Errorf("failed to delete Key: %w", err)
	}
	o.logger.Debug("DELETE completed")

	return nil
}

func (o *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	o.logger.Debug("LIST started")
	defer metrics.MeasureSince(metricListFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	if err := o.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer o.permitPool.Release()

	metrics.MeasureSince(metricListAcquirePool, startAcquirePool)
	defer metrics.MeasureSince(metricList, time.Now())
	var keys []string
	delimiter := "/"
	var start *string

	for {
		opcClientRequestId, err := uuid.GenerateUUID()
		if err != nil {
			o.logger.Error("List: error generating UUID")
			return nil, fmt.Errorf("failed to generate UUID %w", err)
		}
		o.logger.Debug("LIST", "opc-client-request-id", opcClientRequestId)
		request := objectstorage.ListObjectsRequest{
			NamespaceName:      &o.namespaceName,
			BucketName:         &o.bucketName,
			Prefix:             &prefix,
			Delimiter:          &delimiter,
			Start:              start,
			OpcClientRequestId: &opcClientRequestId,
		}

		resp, err := o.client.ListObjects(ctx, request)
		o.logRequest("LIST", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)

		if err != nil {
			metrics.IncrCounter(metricListFailed, 1)
			return nil, fmt.Errorf("failed to list using prefix: %w", err)
		}

		for _, commonPrefix := range resp.Prefixes {
			commonPrefix := strings.TrimPrefix(commonPrefix, prefix)
			keys = append(keys, commonPrefix)
		}

		for _, object := range resp.Objects {
			key := strings.TrimPrefix(*object.Name, prefix)
			keys = append(keys, key)
		}

		// Duplicate keys are not expected
		keys = strutil.RemoveDuplicates(keys, false)

		if resp.NextStartWith == nil {
			resp.RawResponse.Body.Close()
			break
		}

		start = resp.NextStartWith
		resp.RawResponse.Body.Close()
	}

	sort.Strings(keys)
	o.logger.Debug("LIST completed")
	return keys, nil
}

func (o *Backend) logRequest(operation string, response *http.Response, clientOpcRequestIdPtr *string, opcRequestIdPtr *string, err error) {
	statusCode := 0
	clientOpcRequestId := " "
	opcRequestId := " "

	if response != nil {
		statusCode = response.StatusCode
		if statusCode/100 == 5 {
			metrics.IncrCounter(metric5xx, 1)
		}
	}

	if clientOpcRequestIdPtr != nil {
		clientOpcRequestId = *clientOpcRequestIdPtr
	}

	if opcRequestIdPtr != nil {
		opcRequestId = *opcRequestIdPtr
	}

	statusCodeStr := "No response"
	if statusCode != 0 {
		statusCodeStr = strconv.Itoa(statusCode)
	}

	logLine := fmt.Sprintf("%s client:opc-request-id %s opc-request-id: %s status-code: %s",
		operation, clientOpcRequestId, opcRequestId, statusCodeStr)
	if err != nil && statusCode/100 == 5 {
		o.logger.Error(logLine, "error", err)
	}
}
