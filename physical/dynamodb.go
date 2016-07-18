package physical

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

// DynamoDBBackend is a physical backend that stores data in
// a DynamoDB table. It can be run in high-availability mode
// as DynamoDB has locking capabilities.
type DynamoDBBackend struct {
	table     string
	client    *dynamodb.DynamoDB
	recovery  bool
	logger    *log.Logger
	haEnabled bool
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
	held       bool
	lock       sync.Mutex
	recovery   bool
}

// newDynamoDBBackend constructs a DynamoDB backend. If the
// configured DynamoDB table does not exist, it creates it.
func newDynamoDBBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
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
		return nil, fmt.Errorf("invalid read capacity: %s", readCapacityString)
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
		return nil, fmt.Errorf("invalid write capacity: %s", writeCapacityString)
	}
	if writeCapacity == 0 {
		writeCapacity = DefaultDynamoDBWriteCapacity
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		accessKey = conf["access_key"]
	}
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		secretKey = conf["secret_key"]
	}
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	if sessionToken == "" {
		sessionToken = conf["session_token"]
	}

	endpoint := os.Getenv("AWS_DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = conf["region"]
		if region == "" {
			region = DefaultDynamoDBRegion
		}
	}

	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			SessionToken:    sessionToken,
		}},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
		&ec2rolecreds.EC2RoleProvider{Client: ec2metadata.New(session.New())},
	})

	awsConf := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(region).
		WithEndpoint(endpoint)
	client := dynamodb.New(session.New(awsConf))

	if err := ensureTableExists(client, table, readCapacity, writeCapacity); err != nil {
		return nil, err
	}

	haEnabled := os.Getenv("DYNAMODB_HA_ENABLED")
	if haEnabled == "" {
		haEnabled = conf["ha_enabled"]
	}
	haEnabledBool, _ := strconv.ParseBool(haEnabled)

	recoveryMode := os.Getenv("RECOVERY_MODE")
	if recoveryMode == "" {
		recoveryMode = conf["recovery_mode"]
	}
	recoveryModeBool, _ := strconv.ParseBool(recoveryMode)

	return &DynamoDBBackend{
		table:     table,
		client:    client,
		recovery:  recoveryModeBool,
		haEnabled: haEnabledBool,
		logger:    logger,
	}, nil
}

// Put is used to insert or update an entry
func (d *DynamoDBBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"dynamodb", "put"}, time.Now())

	record := DynamoDBRecord{
		Path:  recordPathForVaultKey(entry.Key),
		Key:   recordKeyForVaultKey(entry.Key),
		Value: entry.Value,
	}
	item, err := dynamodbattribute.ConvertToMap(record)
	if err != nil {
		return fmt.Errorf("could not convert prefix record to DynamoDB item: %s", err)
	}
	requests := []*dynamodb.WriteRequest{{
		PutRequest: &dynamodb.PutRequest{
			Item: item,
		},
	}}

	for _, prefix := range prefixes(entry.Key) {
		record = DynamoDBRecord{
			Path: recordPathForVaultKey(prefix),
			Key:  fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)),
		}
		item, err := dynamodbattribute.ConvertToMap(record)
		if err != nil {
			return fmt.Errorf("could not convert prefix record to DynamoDB item: %s", err)
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
func (d *DynamoDBBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"dynamodb", "get"}, time.Now())

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
	if err := dynamodbattribute.ConvertFromMap(resp.Item, record); err != nil {
		return nil, err
	}

	return &Entry{
		Key:   vaultKey(record),
		Value: record.Value,
	}, nil
}

// Delete is used to permanently delete an entry
func (d *DynamoDBBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"dynamodb", "delete"}, time.Now())

	requests := []*dynamodb.WriteRequest{{
		DeleteRequest: &dynamodb.DeleteRequest{
			Key: map[string]*dynamodb.AttributeValue{
				"Path": {S: aws.String(recordPathForVaultKey(key))},
				"Key":  {S: aws.String(recordKeyForVaultKey(key))},
			},
		},
	}}

	// clean up now empty 'folders'
	prefixes := prefixes(key)
	sort.Sort(sort.Reverse(sort.StringSlice(prefixes)))
	for _, prefix := range prefixes {
		items, err := d.List(prefix)
		if err != nil {
			return err
		}
		if len(items) == 1 {
			requests = append(requests, &dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						"Path": {S: aws.String(recordPathForVaultKey(prefix))},
						"Key":  {S: aws.String(fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)))},
					},
				},
			})
		}
	}

	return d.batchWriteRequests(requests)
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (d *DynamoDBBackend) List(prefix string) ([]string, error) {
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
	err := d.client.QueryPages(queryInput, func(out *dynamodb.QueryOutput, lastPage bool) bool {
		var record DynamoDBRecord
		for _, item := range out.Items {
			dynamodbattribute.ConvertFromMap(item, &record)
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

// LockWith is used for mutual exclusion based on the given key.
func (d *DynamoDBBackend) LockWith(key, value string) (Lock, error) {
	return &DynamoDBLock{
		backend:  d,
		key:      filepath.Join(filepath.Dir(key), DynamoDBLockPrefix+filepath.Base(key)),
		value:    value,
		recovery: d.recovery,
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
		batch := requests[:batchSize]
		requests = requests[batchSize:]

		_, err := d.client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				d.table: batch,
			},
		})
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
		// after acquiring it successfully, we must watch
		// the lock in order to close the leader channel
		// once it is lost.
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
	if err := l.backend.Delete(l.key); err != nil {
		return err
	}
	return nil
}

// Value checks whether or not the lock is held by any instance of DynamoDBLock,
// including this one, and returns the current value.
func (l *DynamoDBLock) Value() (bool, string, error) {
	entry, err := l.backend.Get(l.key)
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

	record := DynamoDBRecord{
		Path:  recordPathForVaultKey(l.key),
		Key:   recordKeyForVaultKey(l.key),
		Value: []byte(l.value),
	}
	item, err := dynamodbattribute.ConvertToMap(record)
	if err != nil {
		errors <- err
		return
	}

	for {
		select {
		case <-stop:
			ticker.Stop()
		case <-ticker.C:
			_, err := l.backend.client.PutItem(&dynamodb.PutItemInput{
				TableName:           aws.String(l.backend.table),
				Item:                item,
				ConditionExpression: aws.String("attribute_not_exists(#p) or attribute_not_exists(#k)"),
				ExpressionAttributeNames: map[string]*string{
					"#p": aws.String("Path"),
					"#k": aws.String("Key"),
				},
			})
			if err != nil {
				if err, ok := err.(awserr.Error); ok && err.Code() != "ConditionalCheckFailedException" {
					errors <- err
				}
				if l.recovery {
					_, err := l.backend.client.DeleteItem(&dynamodb.DeleteItemInput{
						TableName: aws.String(l.backend.table),
						Key: map[string]*dynamodb.AttributeValue{
							"Path": {S: aws.String(record.Path)},
							"Key":  {S: aws.String(record.Key)},
						},
					})
					if err != nil {
						errors <- fmt.Errorf("could not delete lock record: %s", err)
					} else {
						l.recovery = false
					}
				}
			} else {
				ticker.Stop()
				close(success)
			}
		}
	}
}

// watch checks whether the lock has changed in the
// DynamoDB table and closes the leader channel if so.
// The interval is set by `DynamoDBWatchRetryInterval`.
// If an error occurs during the check, watch will retry
// the operation for `DynamoDBWatchRetryMax` times and
// close the leader channel if it can't succeed.
func (l *DynamoDBLock) watch(lost chan struct{}) {
	retries := DynamoDBWatchRetryMax

	ticker := time.NewTicker(DynamoDBWatchRetryInterval)
WatchLoop:
	for {
		select {
		case <-ticker.C:
			item, err := l.backend.Get(l.key)
			if err != nil {
				retries--
				if retries == 0 {
					break WatchLoop
				}
				continue
			}

			if item == nil || string(item.Value) != l.value {
				break WatchLoop
			}
		}
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
	if awserr, ok := err.(awserr.Error); ok {
		if awserr.Code() == "ResourceNotFoundException" {
			_, err = client.CreateTable(&dynamodb.CreateTableInput{
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
		}
	}
	if err != nil {
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
		return filepath.Dir(key)
	}
	return DynamoDBEmptyPath
}

// recordKeyForVaultKey transforms a vault key into
// a value suitable for the `DynamoDBRecord`'s `Key`
// property. This path equals the the vault key's
// last component.
func recordKeyForVaultKey(key string) string {
	return filepath.Base(key)
}

// vaultKey returns the vault key for a given record
// from the DynamoDB table. This is the combination of
// the records Path and Key.
func vaultKey(record *DynamoDBRecord) string {
	path := unescapeEmptyPath(record.Path)
	if path == "" {
		return record.Key
	}
	return filepath.Join(record.Path, record.Key)
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

// prefixes returns all parent 'folders' for a given
// vault key.
// e.g. for 'foo/bar/baz', it returns ['foo', 'foo/bar']
func prefixes(s string) []string {
	components := strings.Split(s, "/")
	result := []string{}
	for i := 1; i < len(components); i++ {
		result = append(result, strings.Join(components[:i], "/"))
	}
	return result
}
