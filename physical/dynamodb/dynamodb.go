// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	pkgPath "path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/hashicorp/go-hclog"

	metrics "github.com/armon/go-metrics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/physical"

	"github.com/cenkalti/backoff/v3"
)

const (
	// DefaultDynamoDBRegion is used when no region is configured
	// explicitly.
	DefaultDynamoDBRegion = "us-east-1"
	// DefaultDynamoDBTableName is used when no table name
	// is configured explicitly.
	DefaultDynamoDBTableName = "vault-dynamodb-backend"

	// DefaultDynamoDBReadCapacity is the default read capacity
	// that is used when none is configured explicitly.
	DefaultDynamoDBReadCapacity = 5
	// DefaultDynamoDBWriteCapacity is the default write capacity
	// that is used when none is configured explicitly.
	DefaultDynamoDBWriteCapacity = 5

	// DynamoDBEmptyPath is the string that is used instead of
	// empty strings when stored in DynamoDB.
	DynamoDBEmptyPath = " "
	// DynamoDBLockPrefix is the prefix used to mark DynamoDB records
	// as locks. This prefix causes them not to be returned by
	// List operations.
	DynamoDBLockPrefix = "_"

	// The lock TTL matches the default that Consul API uses, 15 seconds.
	DynamoDBLockTTL = 15 * time.Second

	// The amount of time to wait between the lock renewals
	DynamoDBLockRenewInterval = 5 * time.Second

	// DynamoDBLockRetryInterval is the amount of time to wait
	// if a lock fails before trying again.
	DynamoDBLockRetryInterval = time.Second
	// DynamoDBWatchRetryMax is the number of times to re-try a
	// failed watch before signaling that leadership is lost.
	DynamoDBWatchRetryMax = 5
	// DynamoDBWatchRetryInterval is the amount of time to wait
	// if a watch fails before trying again.
	DynamoDBWatchRetryInterval = 5 * time.Second
)

// Verify DynamoDBBackend satisfies the correct interfaces
var (
	_ physical.Backend   = (*DynamoDBBackend)(nil)
	_ physical.HABackend = (*DynamoDBBackend)(nil)
	_ physical.Lock      = (*DynamoDBLock)(nil)
)

// DynamoDBBackend is a physical backend that stores data in
// a DynamoDB table. It can be run in high-availability mode
// as DynamoDB has locking capabilities.
type DynamoDBBackend struct {
	table      string
	client     *dynamodb.DynamoDB
	logger     log.Logger
	haEnabled  bool
	permitPool *PermitPoolWithMetrics
}

// DynamoDBRecord is the representation of a vault entry in
// DynamoDB. The vault key is split up into two components
// (Path and Key) in order to allow more efficient listings.
type DynamoDBRecord struct {
	Path  string
	Key   string
	Value []byte
}

// DynamoDBLock implements a lock using an DynamoDB client.
type DynamoDBLock struct {
	backend    *DynamoDBBackend
	value, key string
	identity   string
	held       bool
	lock       sync.Mutex
	// Allow modifying the Lock durations for ease of unit testing.
	renewInterval      time.Duration
	ttl                time.Duration
	watchRetryInterval time.Duration
}

type DynamoDBLockRecord struct {
	Path     string
	Key      string
	Value    []byte
	Identity []byte
	Expires  int64
}

type PermitPoolWithMetrics struct {
	physical.PermitPool
	pendingPermits int32
	poolSize       int
}

// NewDynamoDBBackend constructs a DynamoDB backend. If the
// configured DynamoDB table does not exist, it creates it.
func NewDynamoDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	table := os.Getenv("AWS_DYNAMODB_TABLE")
	if table == "" {
		table = conf["table"]
		if table == "" {
			table = DefaultDynamoDBTableName
		}
	}
	readCapacityString := os.Getenv("AWS_DYNAMODB_READ_CAPACITY")
	if readCapacityString == "" {
		readCapacityString = conf["read_capacity"]
		if readCapacityString == "" {
			readCapacityString = "0"
		}
	}
	readCapacity, err := strconv.Atoi(readCapacityString)
	if err != nil {
		return nil, fmt.Errorf("invalid read capacity: %q", readCapacityString)
	}
	if readCapacity == 0 {
		readCapacity = DefaultDynamoDBReadCapacity
	}

	writeCapacityString := os.Getenv("AWS_DYNAMODB_WRITE_CAPACITY")
	if writeCapacityString == "" {
		writeCapacityString = conf["write_capacity"]
		if writeCapacityString == "" {
			writeCapacityString = "0"
		}
	}
	writeCapacity, err := strconv.Atoi(writeCapacityString)
	if err != nil {
		return nil, fmt.Errorf("invalid write capacity: %q", writeCapacityString)
	}
	if writeCapacity == 0 {
		writeCapacity = DefaultDynamoDBWriteCapacity
	}

	endpoint := os.Getenv("AWS_DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}
	region := os.Getenv("AWS_DYNAMODB_REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			region = os.Getenv("AWS_DEFAULT_REGION")
			if region == "" {
				region = conf["region"]
				if region == "" {
					region = DefaultDynamoDBRegion
				}
			}
		}
	}

	dynamodbMaxRetryString := os.Getenv("AWS_DYNAMODB_MAX_RETRIES")
	if dynamodbMaxRetryString == "" {
		dynamodbMaxRetryString = conf["dynamodb_max_retries"]
	}
	dynamodbMaxRetry := aws.UseServiceDefaultRetries
	if dynamodbMaxRetryString != "" {
		var err error
		dynamodbMaxRetry, err = strconv.Atoi(dynamodbMaxRetryString)
		if err != nil {
			return nil, fmt.Errorf("invalid max retry: %q", dynamodbMaxRetryString)
		}
	}

	credsConfig := &awsutil.CredentialsConfig{
		AccessKey:    conf["access_key"],
		SecretKey:    conf["secret_key"],
		SessionToken: conf["session_token"],
		Logger:       logger,
	}
	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	pooledTransport := cleanhttp.DefaultPooledTransport()
	pooledTransport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	awsConf := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithHTTPClient(&http.Client{
			Transport: pooledTransport,
		}).
		WithMaxRetries(dynamodbMaxRetry)

	awsSession, err := session.NewSession(awsConf)
	if err != nil {
		return nil, fmt.Errorf("Could not establish AWS session: %w", err)
	}

	client := dynamodb.New(awsSession)

	if err := ensureTableExists(client, table, readCapacity, writeCapacity); err != nil {
		return nil, err
	}

	haEnabled := os.Getenv("DYNAMODB_HA_ENABLED")
	if haEnabled == "" {
		haEnabled = conf["ha_enabled"]
	}
	haEnabledBool, _ := strconv.ParseBool(haEnabled)

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

	return &DynamoDBBackend{
		table:      table,
		client:     client,
		permitPool: NewPermitPoolWithMetrics(maxParInt),
		haEnabled:  haEnabledBool,
		logger:     logger,
	}, nil
}

// Put is used to insert or update an entry
func (d *DynamoDBBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"dynamodb", "put"}, time.Now())

	record := DynamoDBRecord{
		Path:  recordPathForVaultKey(entry.Key),
		Key:   recordKeyForVaultKey(entry.Key),
		Value: entry.Value,
	}
	item, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("could not convert prefix record to DynamoDB item: %w", err)
	}
	requests := []*dynamodb.WriteRequest{{
		PutRequest: &dynamodb.PutRequest{
			Item: item,
		},
	}}

	for _, prefix := range physical.Prefixes(entry.Key) {
		record = DynamoDBRecord{
			Path: recordPathForVaultKey(prefix),
			Key:  fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)),
		}
		item, err := dynamodbattribute.MarshalMap(record)
		if err != nil {
			return fmt.Errorf("could not convert prefix record to DynamoDB item: %w", err)
		}
		requests = append(requests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		})
	}

	return d.batchWriteRequests(requests)
}

// Get is used to fetch an entry
func (d *DynamoDBBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"dynamodb", "get"}, time.Now())

	d.permitPool.Acquire()
	defer d.permitPool.Release()

	resp, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName:      aws.String(d.table),
		ConsistentRead: aws.Bool(true),
		Key: map[string]*dynamodb.AttributeValue{
			"Path": {S: aws.String(recordPathForVaultKey(key))},
			"Key":  {S: aws.String(recordKeyForVaultKey(key))},
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, nil
	}

	record := &DynamoDBRecord{}
	if err := dynamodbattribute.UnmarshalMap(resp.Item, record); err != nil {
		return nil, err
	}

	return &physical.Entry{
		Key:   vaultKey(record),
		Value: record.Value,
	}, nil
}

// Delete is used to permanently delete an entry
func (d *DynamoDBBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"dynamodb", "delete"}, time.Now())

	requests := []*dynamodb.WriteRequest{{
		DeleteRequest: &dynamodb.DeleteRequest{
			Key: map[string]*dynamodb.AttributeValue{
				"Path": {S: aws.String(recordPathForVaultKey(key))},
				"Key":  {S: aws.String(recordKeyForVaultKey(key))},
			},
		},
	}}

	// Clean up empty "folders" by looping through all levels of the path to the item being deleted looking for
	// children. Loop from deepest path to shallowest, and only consider items children if they are not going to be
	// deleted by our batch delete request. If a path has no valid children, then it should be considered an empty
	// "folder" and be deleted along with the original item in our batch job. Because we loop from deepest path to
	// shallowest, once we find a path level that contains valid children we can stop the cleanup operation.
	prefixes := physical.Prefixes(key)
	sort.Sort(sort.Reverse(sort.StringSlice(prefixes)))
	for index, prefix := range prefixes {
		// Because delete batches its requests, we need to pass keys we know are going to be deleted into
		// hasChildren so it can exclude those when it determines if there WILL be any children left after
		// the delete operations have completed.
		var excluded []string
		if index == 0 {
			// This is the value we know for sure is being deleted
			excluded = append(excluded, recordKeyForVaultKey(key))
		} else {
			// The previous path doesn't count as a child, since if we're still looping, we've found no children
			excluded = append(excluded, recordKeyForVaultKey(prefixes[index-1]))
		}

		hasChildren, err := d.hasChildren(prefix, excluded)
		if err != nil {
			return err
		}

		if !hasChildren {
			// If there are no children other than ones we know are being deleted then cleanup empty "folder" pointers
			requests = append(requests, &dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						"Path": {S: aws.String(recordPathForVaultKey(prefix))},
						"Key":  {S: aws.String(fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)))},
					},
				},
			})
		} else {
			// This loop starts at the deepest path and works backwards looking for children
			// once a deeper level of the path has been found to have children there is no
			// more cleanup that needs to happen, otherwise we might remove folder pointers
			// to that deeper path making it "undiscoverable" with the list operation
			break
		}
	}

	return d.batchWriteRequests(requests)
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (d *DynamoDBBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"dynamodb", "list"}, time.Now())

	prefix = strings.TrimSuffix(prefix, "/")

	keys := []string{}
	prefix = escapeEmptyPath(prefix)
	queryInput := &dynamodb.QueryInput{
		TableName:      aws.String(d.table),
		ConsistentRead: aws.Bool(true),
		KeyConditions: map[string]*dynamodb.Condition{
			"Path": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{
					S: aws.String(prefix),
				}},
			},
		},
	}

	d.permitPool.Acquire()
	defer d.permitPool.Release()

	err := d.client.QueryPages(queryInput, func(out *dynamodb.QueryOutput, lastPage bool) bool {
		var record DynamoDBRecord
		for _, item := range out.Items {
			dynamodbattribute.UnmarshalMap(item, &record)
			if !strings.HasPrefix(record.Key, DynamoDBLockPrefix) {
				keys = append(keys, record.Key)
			}
		}
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// hasChildren returns true if there exist items below a certain path prefix.
// To do so, the method fetches such items from DynamoDB. This method is primarily
// used by Delete. Because DynamoDB requests are batched this method is being called
// before any deletes take place. To account for that hasChildren accepts a slice of
// strings representing values we expect to find that should NOT be counted as children
// because they are going to be deleted.
func (d *DynamoDBBackend) hasChildren(prefix string, exclude []string) (bool, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	prefix = escapeEmptyPath(prefix)

	queryInput := &dynamodb.QueryInput{
		TableName:      aws.String(d.table),
		ConsistentRead: aws.Bool(true),
		KeyConditions: map[string]*dynamodb.Condition{
			"Path": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{
					S: aws.String(prefix),
				}},
			},
		},
		// Avoid fetching too many items from DynamoDB for performance reasons.
		// We want to know if there are any children we don't expect to see.
		// Answering that question requires fetching a minimum of one more item
		// than the number we expect. In most cases this value will be 2
		Limit: aws.Int64(int64(len(exclude) + 1)),
	}

	d.permitPool.Acquire()
	defer d.permitPool.Release()

	out, err := d.client.Query(queryInput)
	if err != nil {
		return false, err
	}
	var childrenExist bool
	for _, item := range out.Items {
		for _, excluded := range exclude {
			// Check if we've found an item we didn't expect to. Look for "folder" pointer keys (trailing slash)
			// and regular value keys (no trailing slash)
			if *item["Key"].S != excluded && *item["Key"].S != fmt.Sprintf("%s/", excluded) {
				childrenExist = true
				break
			}
		}
		if childrenExist {
			// We only need to find ONE child we didn't expect to.
			break
		}
	}

	return childrenExist, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (d *DynamoDBBackend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &DynamoDBLock{
		backend:            d,
		key:                pkgPath.Join(pkgPath.Dir(key), DynamoDBLockPrefix+pkgPath.Base(key)),
		value:              value,
		identity:           identity,
		renewInterval:      DynamoDBLockRenewInterval,
		ttl:                DynamoDBLockTTL,
		watchRetryInterval: DynamoDBWatchRetryInterval,
	}, nil
}

func (d *DynamoDBBackend) HAEnabled() bool {
	return d.haEnabled
}

// batchWriteRequests takes a list of write requests and executes them in badges
// with a maximum size of 25 (which is the limit of BatchWriteItem requests).
func (d *DynamoDBBackend) batchWriteRequests(requests []*dynamodb.WriteRequest) error {
	for len(requests) > 0 {
		batchSize := int(math.Min(float64(len(requests)), 25))
		batch := map[string][]*dynamodb.WriteRequest{d.table: requests[:batchSize]}
		requests = requests[batchSize:]

		var err error

		d.permitPool.Acquire()

		boff := backoff.NewExponentialBackOff()
		boff.MaxElapsedTime = 600 * time.Second

		for len(batch) > 0 {
			var output *dynamodb.BatchWriteItemOutput
			output, err = d.client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: batch,
			})

			if err != nil {
				break
			}

			if len(output.UnprocessedItems) == 0 {
				break
			} else {
				duration := boff.NextBackOff()
				if duration != backoff.Stop {
					batch = output.UnprocessedItems
					time.Sleep(duration)
				} else {
					err = errors.New("dynamodb: timeout handling UnproccessedItems")
					break
				}
			}
		}

		d.permitPool.Release()
		if err != nil {
			return err
		}
	}
	return nil
}

// Lock tries to acquire the lock by repeatedly trying to create
// a record in the DynamoDB table. It will block until either the
// stop channel is closed or the lock could be acquired successfully.
// The returned channel will be closed once the lock is deleted or
// changed in the DynamoDB table.
func (l *DynamoDBLock) Lock(stopCh <-chan struct{}) (doneCh <-chan struct{}, retErr error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.held {
		return nil, fmt.Errorf("lock already held")
	}

	done := make(chan struct{})
	// close done channel even in case of error
	defer func() {
		if retErr != nil {
			close(done)
		}
	}()

	var (
		stop    = make(chan struct{})
		success = make(chan struct{})
		errors  = make(chan error)
		leader  = make(chan struct{})
	)
	// try to acquire the lock asynchronously
	go l.tryToLock(stop, success, errors)

	select {
	case <-success:
		l.held = true
		// after acquiring it successfully, we must renew the lock periodically,
		// and watch the lock in order to close the leader channel
		// once it is lost.
		go l.periodicallyRenewLock(leader)
		go l.watch(leader)
	case retErr = <-errors:
		close(stop)
		return nil, retErr
	case <-stopCh:
		close(stop)
		return nil, nil
	}

	return leader, retErr
}

// Unlock releases the lock by deleting the lock record from the
// DynamoDB table.
func (l *DynamoDBLock) Unlock() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.held {
		return nil
	}

	l.held = false

	// Conditionally delete after check that the key is actually this Vault's and
	// not been already claimed by another leader
	condition := "#identity = :identity"
	deleteMyLock := &dynamodb.DeleteItemInput{
		TableName:           &l.backend.table,
		ConditionExpression: &condition,
		Key: map[string]*dynamodb.AttributeValue{
			"Path": {S: aws.String(recordPathForVaultKey(l.key))},
			"Key":  {S: aws.String(recordKeyForVaultKey(l.key))},
		},
		ExpressionAttributeNames: map[string]*string{
			"#identity": aws.String("Identity"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":identity": {B: []byte(l.identity)},
		},
	}

	_, err := l.backend.client.DeleteItem(deleteMyLock)
	if isConditionCheckFailed(err) {
		err = nil
	}

	return err
}

// Value checks whether or not the lock is held by any instance of DynamoDBLock,
// including this one, and returns the current value.
func (l *DynamoDBLock) Value() (bool, string, error) {
	entry, err := l.backend.Get(context.Background(), l.key)
	if err != nil {
		return false, "", err
	}
	if entry == nil {
		return false, "", nil
	}

	return true, string(entry.Value), nil
}

// tryToLock tries to create a new item in DynamoDB
// every `DynamoDBLockRetryInterval`. As long as the item
// cannot be created (because it already exists), it will
// be retried. If the operation fails due to an error, it
// is sent to the errors channel.
// When the lock could be acquired successfully, the success
// channel is closed.
func (l *DynamoDBLock) tryToLock(stop, success chan struct{}, errors chan error) {
	ticker := time.NewTicker(DynamoDBLockRetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			err := l.updateItem(true)
			if err != nil {
				if err, ok := err.(awserr.Error); ok {
					// Don't report a condition check failure, this means that the lock
					// is already being held.
					if !isConditionCheckFailed(err) {
						errors <- err
					}
				} else {
					// Its not an AWS error, and is probably not transient, bail out.
					errors <- err
					return
				}
			} else {
				close(success)
				return
			}
		}
	}
}

func (l *DynamoDBLock) periodicallyRenewLock(done chan struct{}) {
	ticker := time.NewTicker(l.renewInterval)
	for {
		select {
		case <-ticker.C:
			// This should not renew the lock if the lock was deleted from under you.
			err := l.updateItem(false)
			if err != nil {
				if !isConditionCheckFailed(err) {
					l.backend.logger.Error("error renewing leadership lock", "error", err)
				}
			}
		case <-done:
			ticker.Stop()
			return
		}
	}
}

// Attempts to put/update the dynamodb item using condition expressions to
// evaluate the TTL.
func (l *DynamoDBLock) updateItem(createIfMissing bool) error {
	now := time.Now()

	conditionExpression := ""
	if createIfMissing {
		conditionExpression += "attribute_not_exists(#path) or " +
			"attribute_not_exists(#key) or "
	} else {
		conditionExpression += "attribute_exists(#path) and " +
			"attribute_exists(#key) and "
	}

	// To work when upgrading from older versions that did not include the
	// Identity attribute, we first check if the attr doesn't exist, and if
	// it does, then we check if the identity is equal to our own.
	// We also write if the lock expired.
	conditionExpression += "(attribute_not_exists(#identity) or #identity = :identity or #expires <= :now)"

	_, err := l.backend.client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(l.backend.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Path": {S: aws.String(recordPathForVaultKey(l.key))},
			"Key":  {S: aws.String(recordKeyForVaultKey(l.key))},
		},
		UpdateExpression: aws.String("SET #value=:value, #identity=:identity, #expires=:expires"),
		// If both key and path already exist, we can only write if
		// A. identity is equal to our identity (or the identity doesn't exist)
		// or
		// B. The ttl on the item is <= to the current time
		ConditionExpression: aws.String(conditionExpression),
		ExpressionAttributeNames: map[string]*string{
			"#path":     aws.String("Path"),
			"#key":      aws.String("Key"),
			"#identity": aws.String("Identity"),
			"#expires":  aws.String("Expires"),
			"#value":    aws.String("Value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":identity": {B: []byte(l.identity)},
			":value":    {B: []byte(l.value)},
			":now":      {N: aws.String(strconv.FormatInt(now.UnixNano(), 10))},
			":expires":  {N: aws.String(strconv.FormatInt(now.Add(l.ttl).UnixNano(), 10))},
		},
	})

	return err
}

// watch checks whether the lock has changed in the
// DynamoDB table and closes the leader channel if so.
// The interval is set by `DynamoDBWatchRetryInterval`.
// If an error occurs during the check, watch will retry
// the operation for `DynamoDBWatchRetryMax` times and
// close the leader channel if it can't succeed.
func (l *DynamoDBLock) watch(lost chan struct{}) {
	retries := DynamoDBWatchRetryMax

	ticker := time.NewTicker(l.watchRetryInterval)
WatchLoop:
	for {
		select {
		case <-ticker.C:
			resp, err := l.backend.client.GetItem(&dynamodb.GetItemInput{
				TableName:      aws.String(l.backend.table),
				ConsistentRead: aws.Bool(true),
				Key: map[string]*dynamodb.AttributeValue{
					"Path": {S: aws.String(recordPathForVaultKey(l.key))},
					"Key":  {S: aws.String(recordKeyForVaultKey(l.key))},
				},
			})
			if err != nil {
				retries--
				if retries == 0 {
					break WatchLoop
				}
				continue
			}

			if resp == nil {
				break WatchLoop
			}
			record := &DynamoDBLockRecord{}
			err = dynamodbattribute.UnmarshalMap(resp.Item, record)
			if err != nil || string(record.Identity) != l.identity {
				break WatchLoop
			}
		}
		retries = DynamoDBWatchRetryMax
	}

	close(lost)
}

// ensureTableExists creates a DynamoDB table with a given
// DynamoDB client. If the table already exists, it is not
// being reconfigured.
func ensureTableExists(client *dynamodb.DynamoDB, table string, readCapacity, writeCapacity int) error {
	_, err := client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			if awsError.Code() == "ResourceNotFoundException" {
				_, err := client.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(table),
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(int64(readCapacity)),
						WriteCapacityUnits: aws.Int64(int64(writeCapacity)),
					},
					KeySchema: []*dynamodb.KeySchemaElement{{
						AttributeName: aws.String("Path"),
						KeyType:       aws.String("HASH"),
					}, {
						AttributeName: aws.String("Key"),
						KeyType:       aws.String("RANGE"),
					}},
					AttributeDefinitions: []*dynamodb.AttributeDefinition{{
						AttributeName: aws.String("Path"),
						AttributeType: aws.String("S"),
					}, {
						AttributeName: aws.String("Key"),
						AttributeType: aws.String("S"),
					}},
				})
				if err != nil {
					return err
				}

				err = client.WaitUntilTableExists(&dynamodb.DescribeTableInput{
					TableName: aws.String(table),
				})
				if err != nil {
					return err
				}
				// table created successfully
				return nil
			}
		}
		return err
	}

	return nil
}

// recordPathForVaultKey transforms a vault key into
// a value suitable for the `DynamoDBRecord`'s `Path`
// property. This path equals the the vault key without
// its last component.
func recordPathForVaultKey(key string) string {
	if strings.Contains(key, "/") {
		return pkgPath.Dir(key)
	}
	return DynamoDBEmptyPath
}

// recordKeyForVaultKey transforms a vault key into
// a value suitable for the `DynamoDBRecord`'s `Key`
// property. This path equals the the vault key's
// last component.
func recordKeyForVaultKey(key string) string {
	return pkgPath.Base(key)
}

// vaultKey returns the vault key for a given record
// from the DynamoDB table. This is the combination of
// the records Path and Key.
func vaultKey(record *DynamoDBRecord) string {
	path := unescapeEmptyPath(record.Path)
	if path == "" {
		return record.Key
	}
	return pkgPath.Join(record.Path, record.Key)
}

// escapeEmptyPath is used to escape the root key's path
// with a value that can be stored in DynamoDB. DynamoDB
// does not allow values to be empty strings.
func escapeEmptyPath(s string) string {
	if s == "" {
		return DynamoDBEmptyPath
	}
	return s
}

// unescapeEmptyPath is the opposite of `escapeEmptyPath`.
func unescapeEmptyPath(s string) string {
	if s == DynamoDBEmptyPath {
		return ""
	}
	return s
}

// isConditionCheckFailed tests whether err is an ErrCodeConditionalCheckFailedException
// from the AWS SDK.
func isConditionCheckFailed(err error) bool {
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			return err.Code() == dynamodb.ErrCodeConditionalCheckFailedException
		}
	}

	return false
}

// NewPermitPoolWithMetrics returns a new permit pool with the provided
// number of permits which emits metrics
func NewPermitPoolWithMetrics(permits int) *PermitPoolWithMetrics {
	return &PermitPoolWithMetrics{
		PermitPool:     *physical.NewPermitPool(permits),
		pendingPermits: 0,
		poolSize:       permits,
	}
}

// Acquire returns when a permit has been acquired
func (c *PermitPoolWithMetrics) Acquire() {
	atomic.AddInt32(&c.pendingPermits, 1)
	c.emitPermitMetrics()
	c.PermitPool.Acquire()
	atomic.AddInt32(&c.pendingPermits, -1)
	c.emitPermitMetrics()
}

// Release returns a permit to the pool
func (c *PermitPoolWithMetrics) Release() {
	c.PermitPool.Release()
	c.emitPermitMetrics()
}

// Get the number of requests in the permit pool
func (c *PermitPoolWithMetrics) CurrentPermits() int {
	return c.PermitPool.CurrentPermits()
}

func (c *PermitPoolWithMetrics) emitPermitMetrics() {
	metrics.SetGauge([]string{"dynamodb", "permit_pool", "pending_permits"}, float32(c.pendingPermits))
	metrics.SetGauge([]string{"dynamodb", "permit_pool", "active_permits"}, float32(c.PermitPool.CurrentPermits()))
	metrics.SetGauge([]string{"dynamodb", "permit_pool", "pool_size"}, float32(c.poolSize))
}
