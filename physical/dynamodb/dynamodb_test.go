// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package dynamodb

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/require"
)

func TestDynamoDBBackend(t *testing.T) {
	cleanup, svccfg := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := svccfg.Credentials.Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: svccfg.Credentials,
		Endpoint:    aws.String(svccfg.URL().String()),
		Region:      aws.String(region),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	conn := dynamodb.New(awsSession)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      svccfg.URL().String(),
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)

	t.Run("Marshalling upgrade", func(t *testing.T) {
		path := "test_key"

		// Manually write to DynamoDB using the old ConvertTo function
		// for marshalling data
		inputEntry := &physical.Entry{
			Key:   path,
			Value: []byte{0x0f, 0xcf, 0x4a, 0x0f, 0xba, 0x2b, 0x15, 0xf0, 0xaa, 0x75, 0x09},
		}

		record := DynamoDBRecord{
			Path:  recordPathForVaultKey(inputEntry.Key),
			Key:   recordKeyForVaultKey(inputEntry.Key),
			Value: inputEntry.Value,
		}

		item, err := dynamodbattribute.ConvertToMap(record)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		request := &dynamodb.PutItemInput{
			Item:      item,
			TableName: &table,
		}
		conn.PutItem(request)

		// Read back the data using the normal interface which should
		// handle the old marshalling format gracefully
		entry, err := b.Get(context.Background(), path)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if diff := deep.Equal(inputEntry, entry); diff != nil {
			t.Fatal(diff)
		}
	})
}

func TestDynamoDBHABackend(t *testing.T) {
	cleanup, svccfg := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := svccfg.Credentials.Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: svccfg.Credentials,
		Endpoint:    aws.String(svccfg.URL().String()),
		Region:      aws.String(region),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	conn := dynamodb.New(awsSession)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      svccfg.URL().String(),
	}

	b, err := NewDynamoDBBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	b2, err := NewDynamoDBBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
	testDynamoDBLockTTL(t, b.(physical.HABackend))
	testDynamoDBLockRenewal(t, b.(physical.HABackend))
}

// TestDynamoDBBackendPayPerRequest tests the DynamoDB backend
// with the PAY_PER_REQUEST billing mode
func TestDynamoDBBackendPayPerRequest(t *testing.T) {
	cleanup, svccfg := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := svccfg.Credentials.Get()
	require.NoError(t, err)

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: svccfg.Credentials,
		Endpoint:    aws.String(svccfg.URL().String()),
		Region:      aws.String(region),
	})
	require.NoError(t, err)

	conn := dynamodb.New(awsSession)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      svccfg.URL().String(),
		"billing_mode":  "PAY_PER_REQUEST",
	}, logger)
	require.NoError(t, err)

	dynamoTable, err := conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	require.NoError(t, err)
	billingMode := *(dynamoTable.Table.BillingModeSummary.BillingMode)
	require.Equal(t, "PAY_PER_REQUEST", billingMode)

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

// TestDynamoDBBackendUpdateBillingMode tests the DynamoDB backend
// and updating the billing mode
func TestDynamoDBBackendUpdateBillingMode(t *testing.T) {
	cleanup, svccfg := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := svccfg.Credentials.Get()
	require.NoError(t, err)

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: svccfg.Credentials,
		Endpoint:    aws.String(svccfg.URL().String()),
		Region:      aws.String(region),
	})
	require.NoError(t, err)

	conn := dynamodb.New(awsSession)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      svccfg.URL().String(),
	}, logger)
	require.NoError(t, err)

	dynamoTable, err := conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	require.NoError(t, err)
	billingMode := dynamoTable.Table.BillingModeSummary
	require.Nil(t, billingMode)

	// now run again, with the same table name but a different billing mode
	// and setting allow_update
	b, err = NewDynamoDBBackend(map[string]string{
		"access_key":             creds.AccessKeyID,
		"secret_key":             creds.SecretAccessKey,
		"session_token":          creds.SessionToken,
		"table":                  table,
		"region":                 region,
		"endpoint":               svccfg.URL().String(),
		"billing_mode":           "PAY_PER_REQUEST",
		"dynamodb_allow_updates": "true",
	}, logger)
	require.NoError(t, err)

	dynamoTable, err = conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	require.NoError(t, err)
	newBillingMode := *(dynamoTable.Table.BillingModeSummary.BillingMode)
	require.Equal(t, "PAY_PER_REQUEST", newBillingMode)

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

// TestDynamoDBBackendUpdateReadWriteCapacity tests the DynamoDB backend
// and updating the provisioned read and write capacity
func TestDynamoDBBackendUpdateReadWriteCapacity(t *testing.T) {
	cleanup, svccfg := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := svccfg.Credentials.Get()
	require.NoError(t, err)

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: svccfg.Credentials,
		Endpoint:    aws.String(svccfg.URL().String()),
		Region:      aws.String(region),
	})
	require.NoError(t, err)

	conn := dynamodb.New(awsSession)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      svccfg.URL().String(),
	}, logger)
	require.NoError(t, err)

	dynamoTable, err := conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	require.NoError(t, err)

	provisionedThroughput := dynamoTable.Table.ProvisionedThroughput
	require.NotNil(t, provisionedThroughput)
	require.Equal(t, int64(5), *(provisionedThroughput.ReadCapacityUnits))
	require.Equal(t, int64(5), *(provisionedThroughput.WriteCapacityUnits))

	// now run again, with the same table name but a capacity of 20
	// and setting allow_update
	b, err = NewDynamoDBBackend(map[string]string{
		"access_key":             creds.AccessKeyID,
		"secret_key":             creds.SecretAccessKey,
		"session_token":          creds.SessionToken,
		"table":                  table,
		"region":                 region,
		"endpoint":               svccfg.URL().String(),
		"read_capacity":          "20",
		"write_capacity":         "20",
		"dynamodb_allow_updates": "true",
	}, logger)
	require.NoError(t, err)

	dynamoTable, err = conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	require.NoError(t, err)

	provisionedThroughput = dynamoTable.Table.ProvisionedThroughput
	require.NotNil(t, provisionedThroughput)
	require.Equal(t, int64(20), *(provisionedThroughput.ReadCapacityUnits))
	require.Equal(t, int64(20), *(provisionedThroughput.WriteCapacityUnits))

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

// Similar to testHABackend, but using internal implementation details to
// trigger the lock failure scenario by setting the lock renew period for one
// of the locks to a higher value than the lock TTL.
func testDynamoDBLockTTL(t *testing.T, ha physical.HABackend) {
	// Set much smaller lock times to speed up the test.
	lockTTL := time.Second * 3
	renewInterval := time.Second * 1
	watchInterval := time.Second * 1

	// Get the lock
	origLock, err := ha.LockWith("dynamodbttl", "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// set the first lock renew period to double the expected TTL.
	lock := origLock.(*DynamoDBLock)
	lock.renewInterval = lockTTL * 2
	lock.ttl = lockTTL
	lock.watchRetryInterval = watchInterval

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("failed to get leader ch")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "bar" {
		t.Fatalf("bad value: %v", err)
	}

	// Second acquisition should succeed because the first lock should
	// not renew within the 3 sec TTL.
	origLock2, err := ha.LockWith("dynamodbttl", "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	lock2 := origLock2.(*DynamoDBLock)
	lock2.renewInterval = renewInterval
	lock2.ttl = lockTTL
	lock2.watchRetryInterval = watchInterval

	// Cancel attempt eventually so as not to block unit tests forever
	stopCh := make(chan struct{})
	time.AfterFunc(lockTTL*10, func() {
		close(stopCh)
	})

	// Attempt to lock should work
	leaderCh2, err := lock2.Lock(stopCh)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh2 == nil {
		t.Fatalf("should get leader ch")
	}

	// Check the value
	held, val, err = lock2.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "baz" {
		t.Fatalf("bad value: %v", err)
	}

	// The first lock should have lost the leader channel
	leaderChClosed := false
	blocking := make(chan struct{})
	// Attempt to read from the leader or the blocking channel, which ever one
	// happens first.
	go func() {
		select {
		case <-time.After(watchInterval * 3):
			return
		case <-leaderCh:
			leaderChClosed = true
			close(blocking)
		case <-blocking:
			return
		}
	}()

	<-blocking
	if !leaderChClosed {
		t.Fatalf("original lock did not have its leader channel closed.")
	}

	// Cleanup
	lock2.Unlock()
}

// Similar to testHABackend, but using internal implementation details to
// trigger a renewal before a "watch" check, which has been a source of
// race conditions.
func testDynamoDBLockRenewal(t *testing.T, ha physical.HABackend) {
	renewInterval := time.Second * 1
	watchInterval := time.Second * 5

	// Get the lock
	origLock, err := ha.LockWith("dynamodbrenewal", "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// customize the renewal and watch intervals
	lock := origLock.(*DynamoDBLock)
	lock.renewInterval = renewInterval
	lock.watchRetryInterval = watchInterval

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("failed to get leader ch")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "bar" {
		t.Fatalf("bad value: %v", err)
	}

	// Release the lock, which will delete the stored item
	if err := lock.Unlock(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait longer than the renewal time, but less than the watch time
	time.Sleep(1500 * time.Millisecond)

	// Attempt to lock with new lock
	newLock, err := ha.LockWith("dynamodbrenewal", "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Cancel attempt in 6 sec so as not to block unit tests forever
	stopCh := make(chan struct{})
	time.AfterFunc(6*time.Second, func() {
		close(stopCh)
	})

	// Attempt to lock should work
	leaderCh2, err := newLock.Lock(stopCh)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh2 == nil {
		t.Fatalf("should get leader ch")
	}

	// Check the value
	held, val, err = newLock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "baz" {
		t.Fatalf("bad value: %v", err)
	}

	// Cleanup
	newLock.Unlock()
}

type Config struct {
	docker.ServiceURL
	Credentials *credentials.Credentials
}

var _ docker.ServiceConfig = &Config{}

func prepareDynamoDBTestContainer(t *testing.T) (func(), *Config) {
	// Skipping on ARM, as this image can't run on ARM architecture
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	// If environment variable is set, assume caller wants to target a real
	// DynamoDB.
	if endpoint := os.Getenv("AWS_DYNAMODB_ENDPOINT"); endpoint != "" {
		s, err := docker.NewServiceURLParse(endpoint)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{*s, credentials.NewEnvCredentials()}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/cnadiminti/dynamodb-local",
		ImageTag:      "latest",
		ContainerName: "dynamodb",
		Ports:         []string{"8000/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start local DynamoDB: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectDynamoDB)
	if err != nil {
		t.Fatalf("Could not start local DynamoDB: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config)
}

func connectDynamoDB(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 400 {
		return nil, err
	}

	return &Config{
		ServiceURL:  *docker.NewServiceURL(u),
		Credentials: credentials.NewStaticCredentials("fake", "fake", ""),
	}, nil
}
