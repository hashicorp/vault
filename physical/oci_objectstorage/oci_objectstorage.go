// Copyright © 2019, Oracle and/or its affiliates. All rights reserved.
package oci_objectStorage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Verify Backend satisfies the correct interfaces
var _ physical.Backend = (*Backend)(nil)

const (
	// Limits maximum outstanding requests
	MaxNumberOfPermits = 256
)

var (
	metricDelete     = []string{"oci_objectstorage", "delete"}
	metricGet        = []string{"oci_objectstorage", "get"}
	metricList       = []string{"oci_objectstorage", "list"}
	metricPut        = []string{"oci_objectstorage", "put"}
	metricDeleteFull = []string{"oci_objectstorage", "deleteFull"}
	metricGetFull    = []string{"oci_objectstorage", "getFull"}
	metricListFull   = []string{"oci_objectstorage", "listFull"}
	metricPutFull    = []string{"oci_objectstorage", "putFull"}

	metricDeleteHa = []string{"oci_objectstorage", "deleteHa"}
	metricGetHa    = []string{"oci_objectstorage", "getHa"}
	metricPutHa    = []string{"oci_objectstorage", "putHa"}

	metricDeleteAcquirePool = []string{"oci_objectstorage", "deleteAcquirePool"}
	metricGetAcquirePool    = []string{"oci_objectstorage", "getAcquirePool"}
	metricListAcquirePool   = []string{"oci_objectstorage", "listAcquirePool"}
	metricPutAcquirePool    = []string{"oci_objectstorage", "putAcquirePool"}

	metricDeleteFailed         = []string{"oci_objectstorage", "deleteFailed"}
	metricGetFailed            = []string{"oci_objectstorage", "getFailed"}
	metricListFailed           = []string{"oci_objectstorage", "listFailed"}
	metricPutFailed            = []string{"oci_objectstorage", "putFailed"}
	metricHaWatchLockRetriable = []string{"oci_objectstorage", "haWatchLockRetriable"}
	metricPermitsUsed          = []string{"oci_objectstorage", "permitsUsed"}

	metric5xx = []string{"oci_objectstorage", "5xx"}
)

type Backend struct {
	client         *objectstorage.ObjectStorageClient
	bucketName     string
	logger         log.Logger
	permitPool     *physical.PermitPool
	namespaceName  string
	haEnabled      bool
	lockBucketName string
}

func NewBackendConstructor(
	client *objectstorage.ObjectStorageClient,
	bucketName string,
	logger log.Logger,
	permitPool *physical.PermitPool,
	namespaceName string,
	haEnabled bool,
	lockBucketName string) *Backend {
	return &Backend{
		client:         client,
		bucketName:     bucketName,
		logger:         logger,
		permitPool:     permitPool,
		namespaceName:  namespaceName,
		haEnabled:      haEnabled,
		lockBucketName: lockBucketName,
	}
}

func NewBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	bucketName := conf["bucketName"]
	if bucketName == "" {
		return nil, errors.New("missing bucket name")
	}

	namespaceName := conf["namespaceName"]
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
			return nil, errwrap.Wrapf("failed to parse HA enabled: {{err}}", err)
		}

		if haEnabled {
			lockBucketName = conf["lockBucketName"]
			if lockBucketName == "" {
				return nil, errors.New("missing lock bucket name")
			}
		}
	}

	authTypeAPIKeyBool := false
	authTypeAPIKeyStr := conf["authTypeAPIKey"]
	if authTypeAPIKeyStr != "" {
		authTypeAPIKeyBool, err = strconv.ParseBool(authTypeAPIKeyStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing authTypeAPIKey parameter: {{err}}", err)
		}
	}

	var cp common.ConfigurationProvider
	if authTypeAPIKeyBool {
		cp = common.DefaultConfigProvider()
	} else {
		cp, err = auth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, errwrap.Wrapf("failed creating InstancePrincipalConfigurationProvider: {{err}}", err)
		}
	}

	objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(cp)
	if err != nil {
		return nil, errwrap.Wrapf("failed creating NewObjectStorageClientWithConfigurationProvider: {{err}}", err)
	}

	logger.Debug("physical/oci: configuration",
		"bucketName", bucketName,
		"namespaceName", namespaceName,
		"haEnabled", haEnabled,
		"lockBucketName", lockBucketName,
		"authTypeAPIKey", authTypeAPIKeyBool,
	)

	return &Backend{
		client:         &objectStorageClient,
		bucketName:     bucketName,
		logger:         logger,
		permitPool:     physical.NewPermitPool(MaxNumberOfPermits),
		namespaceName:  namespaceName,
		haEnabled:      haEnabled,
		lockBucketName: lockBucketName,
	}, nil
}

func (o *Backend) Put(ctx context.Context, entry *physical.Entry) error {
	o.logger.Debug("physical/oci: put")
	defer metrics.MeasureSince(metricPutFull, time.Now())
	startAcquirePool := time.Now()
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	o.permitPool.Acquire()
	defer o.permitPool.Release()
	metrics.MeasureSince(metricPutAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricPut, time.Now())
	size := int64(binary.Size(entry.Value))
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		o.logger.Debug("physical/oci: Put: error generating UUID")
		opcClientRequestId = ""
	}
	o.logger.Debug("physical/oci: put", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.PutObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &entry.Key,
		ContentLength:      &size,
		PutObjectBody:      ioutil.NopCloser(bytes.NewReader(entry.Value)),
		OpcMeta:            nil,
		OpcClientRequestId: &opcClientRequestId,
	}

	resp, err := o.client.PutObject(ctx, request)
	o.logRequest("put", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)

	if err != nil {
		metrics.SetGauge(metricPutFailed, 1)

		if resp.RawResponse != nil {
			defer resp.RawResponse.Body.Close()
		}

		return errwrap.Wrapf("failed to put data: {{err}}", err)
	}
	o.logger.Debug("physical/oci: put end")

	defer resp.RawResponse.Body.Close()

	return nil
}

func (o *Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	o.logger.Debug("physical/oci: get")
	defer metrics.MeasureSince(metricGetFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	o.permitPool.Acquire()
	defer o.permitPool.Release()
	metrics.MeasureSince(metricGetAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricGet, time.Now())
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		o.logger.Debug("physical/oci: Get: error generating UUID")
		opcClientRequestId = ""
	}
	o.logger.Debug("physical/oci: Get", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.GetObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &key,
		OpcClientRequestId: &opcClientRequestId,
	}

	response, err := o.client.GetObject(ctx, request)
	o.logRequest("get", response.RawResponse, response.OpcClientRequestId, response.OpcRequestId, err)

	if err != nil {

		if response.RawResponse != nil && response.RawResponse.StatusCode == http.StatusNotFound {
			defer response.RawResponse.Body.Close()
			return nil, nil
		}

		metrics.SetGauge(metricGetFailed, 1)
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read Value for %q: {{err}}", key), err)
	}

	defer response.RawResponse.Body.Close()

	body, err := ioutil.ReadAll(response.Content)
	if err != nil {
		metrics.SetGauge(metricGetFailed, 1)
		return nil, errwrap.Wrapf("failed to decode Value into bytes: {{err}}", err)
	}

	o.logger.Debug("physical/oci: get end")

	return &physical.Entry{
		Key:   key,
		Value: body,
	}, nil
}

func (o *Backend) Delete(ctx context.Context, key string) error {
	o.logger.Debug("physical/oci: delete")
	defer metrics.MeasureSince(metricDeleteFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	o.permitPool.Acquire()
	defer o.permitPool.Release()
	metrics.MeasureSince(metricDeleteAcquirePool, startAcquirePool)

	defer metrics.MeasureSince(metricDelete, time.Now())
	opcClientRequestId, err := uuid.GenerateUUID()
	if err != nil {
		o.logger.Debug("physical/oci: Delete: error generating UUID")
		opcClientRequestId = ""
	}
	o.logger.Debug("physical/oci: Delete", "opc-client-request-id", opcClientRequestId)
	request := objectstorage.DeleteObjectRequest{
		NamespaceName:      &o.namespaceName,
		BucketName:         &o.bucketName,
		ObjectName:         &key,
		OpcClientRequestId: &opcClientRequestId,
	}

	resp, err := o.client.DeleteObject(ctx, request)
	o.logRequest("delete", resp.RawResponse, resp.OpcClientRequestId, resp.OpcRequestId, err)

	if err != nil {

		if resp.RawResponse != nil && resp.RawResponse.StatusCode == http.StatusNotFound {
			defer resp.RawResponse.Body.Close()
			return nil
		}

		metrics.SetGauge(metricDeleteFailed, 1)
		return errwrap.Wrapf("failed to delete Key: {{err}}", err)
	}
	o.logger.Debug("physical/oci: delete end")
	defer resp.RawResponse.Body.Close()

	return nil
}

func (o *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	o.logger.Debug("physical/oci: list")
	defer metrics.MeasureSince(metricListFull, time.Now())
	metrics.SetGauge(metricPermitsUsed, float32(o.permitPool.CurrentPermits()))
	startAcquirePool := time.Now()
	o.permitPool.Acquire()
	defer o.permitPool.Release()

	metrics.MeasureSince(metricListAcquirePool, startAcquirePool)
	defer metrics.MeasureSince(metricList, time.Now())
	var keys []string
	delimiter := "/"
	var start *string

	for {

		opcClientRequestId, err := uuid.GenerateUUID()
		if err != nil {
			o.logger.Debug("physical/oci: List: error generating UUID")
			opcClientRequestId = ""
		}
		o.logger.Debug("physical/oci: List", "opc-client-request-id", opcClientRequestId)
		request := objectstorage.ListObjectsRequest{
			NamespaceName:      &o.namespaceName,
			BucketName:         &o.bucketName,
			Prefix:             &prefix,
			Delimiter:          &delimiter,
			Start:              start,
			OpcClientRequestId: &opcClientRequestId,
		}

		response, err := o.client.ListObjects(ctx, request)
		o.logRequest("list", response.RawResponse, response.OpcClientRequestId, response.OpcRequestId, err)

		if err != nil {
			metrics.SetGauge(metricListFailed, 1)

			return nil, errwrap.Wrapf("failed to list using prefix: {{err}}", err)
		}

		for _, commonPrefix := range response.Prefixes {
			commonPrefix := strings.TrimPrefix(commonPrefix, prefix)
			keys = strutil.AppendIfMissing(keys, commonPrefix)
		}

		for _, object := range response.Objects {
			key := strings.TrimPrefix(*object.Name, prefix)
			keys = strutil.AppendIfMissing(keys, key)
		}

		if response.NextStartWith == nil {
			response.RawResponse.Body.Close()
			break
		}

		start = response.NextStartWith
		response.RawResponse.Body.Close()
	}

	sort.Strings(keys)
	o.logger.Debug("physical/oci: list end")
	return keys, nil
}

func (o *Backend) logRequest(operation string, response *http.Response, clientOpcRequestIdPtr *string, opcRequestIdPtr *string, err error) {

	statusCode := 0
	clientOpcRequestId := " "
	opcRequestId := " "

	if response != nil {
		statusCode = response.StatusCode
		if statusCode/100 == 5 {
			metrics.SetGauge(metric5xx, 1)
		}
	}

	if clientOpcRequestIdPtr != nil {
		clientOpcRequestId = *clientOpcRequestIdPtr
	}

	if opcRequestIdPtr != nil {
		opcRequestId = *opcRequestIdPtr
	}

	logLine := fmt.Sprintf("physical/oci: %s client:opc-request-id %s opc-request-id: %s status-code: %d",
		operation, clientOpcRequestId, opcRequestId, statusCode)
	if err == nil {
		o.logger.Debug(logLine)
	} else {
		o.logger.Error(logLine, "error", err)
	}
}
