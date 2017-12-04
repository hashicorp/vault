package dynamodb

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoDBBackend(t *testing.T) {
	cleanup, endpoint, credsProvider := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := credsProvider.Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	conn := dynamodb.New(session.New(&aws.Config{
		Credentials: credsProvider,
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}))

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      endpoint,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestDynamoDBHABackend(t *testing.T) {
	cleanup, endpoint, credsProvider := prepareDynamoDBTestContainer(t)
	defer cleanup()

	creds, err := credsProvider.Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	conn := dynamodb.New(session.New(&aws.Config{
		Credentials: credsProvider,
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}))

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := logformat.NewVaultLogger(log.LevelTrace)
	b, err := NewDynamoDBBackend(map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
		"region":        region,
		"endpoint":      endpoint,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ha, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("dynamodb does not implement HABackend")
	}
	physical.ExerciseHABackend(t, ha, ha)
	testDynamoDBLockTTL(t, ha)
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

	// Cancel attempt in 6 sec so as not to block unit tests forever
	stopCh := make(chan struct{})
	time.AfterFunc(lockTTL*2, func() {
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

func prepareDynamoDBTestContainer(t *testing.T) (cleanup func(), retAddress string, creds *credentials.Credentials) {
	// If environment variable is set, assume caller wants to target a real
	// DynamoDB.
	if os.Getenv("AWS_DYNAMODB_ENDPOINT") != "" {
		return func() {}, os.Getenv("AWS_DYNAMODB_ENDPOINT"), credentials.NewEnvCredentials()
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("deangiberson/aws-dynamodb-local", "latest", []string{})
	if err != nil {
		t.Fatalf("Could not start local DynamoDB: %s", err)
	}

	retAddress = "http://localhost:" + resource.GetPort("8000/tcp")
	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local DynamoDB: %s", err)
		}
	}

	// exponential backoff-retry, because the DynamoDB may not be able to accept
	// connections yet
	if err := pool.Retry(func() error {
		var err error
		resp, err := http.Get(retAddress)
		if err != nil {
			return err
		}
		if resp.StatusCode != 400 {
			return fmt.Errorf("Expected DynamoDB to return status code 400, got (%s) instead.", resp.Status)
		}
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retAddress, credentials.NewStaticCredentials("fake", "fake", "")
}
