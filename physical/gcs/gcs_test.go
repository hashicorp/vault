package gcs

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func TestGCSBackend(t *testing.T) {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// projectID is only required for creating a bucket for this test
	projectID := os.Getenv("GOOGLE_PROJECT_ID")

	if credentialsFile == "" || projectID == "" {
		t.SkipNow()
	}

	client, err := storage.NewClient(
		context.Background(),
		option.WithServiceAccountFile(credentialsFile),
	)

	if err != nil {
		t.Fatalf("error creating storage client: '%v'", err)
	}

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucketName := fmt.Sprintf("vault-gcs-testacc-%d", randInt)

	bucket := client.Bucket(bucketName)
	err = bucket.Create(context.Background(), projectID, nil)

	if err != nil {
		t.Fatalf("error creating bucket '%v': '%v'", bucketName, err)
	}

	// test bucket teardown
	defer func() {
		objects_it := bucket.Objects(context.Background(), nil)

		// have to delete all objects before deleting bucket
		for {
			objAttrs, err := objects_it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("error listing bucket '%v' contents: '%v'", bucketName, err)
			}

			// ignore errors in deleting a single object, we only care about deleting the bucket
			// occasionally we get "storage: object doesn't exist" which is fine
			bucket.Object(objAttrs.Name).Delete(context.Background())
		}

		err := bucket.Delete(context.Background())
		if err != nil {
			t.Fatalf("error deleting bucket '%s': '%v'", bucketName, err)
		}
	}()

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewGCSBackend(map[string]string{
		"bucket":           bucketName,
		"credentials_file": credentialsFile,
	}, logger)

	if err != nil {
		t.Fatalf("error creating google cloud storage backend: '%s'", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestGCSHABackend(t *testing.T) {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// projectID is only required for creating a bucket for this test
	projectID := os.Getenv("GOOGLE_PROJECT_ID")

	if credentialsFile == "" || projectID == "" {
		t.SkipNow()
	}

	client, err := storage.NewClient(
		context.Background(),
		option.WithServiceAccountFile(credentialsFile),
	)

	if err != nil {
		t.Fatalf("error creating storage client: '%v'", err)
	}

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucketName := fmt.Sprintf("vault-gcs-testacc-%d", randInt)

	bucket := client.Bucket(bucketName)
	err = bucket.Create(context.Background(), projectID, nil)

	if err != nil {
		t.Fatalf("error creating bucket '%v': '%v'", bucketName, err)
	}

	// test bucket teardown
	defer func() {
		objects := bucket.Objects(context.Background(), nil)

		// have to delete all objects before deleting bucket
		for {
			objAttrs, err := objects.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("error listing bucket '%v' contents: '%v'", bucketName, err)
			}

			// ignore errors in deleting a single object, we only care about deleting the bucket
			// occassionally we get "storage: object doesn't exist" which is fine
			bucket.Object(objAttrs.Name).Delete(context.Background())
		}

		err = bucket.Delete(context.Background())
		if err != nil {
			t.Fatalf("error deleting bucket '%s': '%v'", bucketName, err)
		}
	}()

	logger := logformat.NewVaultLogger(log.LevelAll)

	b, err := NewGCSBackend(map[string]string{
		"bucket":           bucketName,
		"credentials_file": credentialsFile,
	}, logger)

	if err != nil {
		t.Fatalf("error creating google cloud storage backend: '%s'", err)
	}

	ha, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("dynamodb does not implement HABackend")
	}
	physical.ExerciseHABackend(t, ha, ha)
	testGCSLockTTL(t, ha)
}

// Similar to testHABackend, but using internal implementation details to
// trigger the lock failure scenario by setting the lock renew period for one
// of the locks to a higher value than the lock TTL.
func testGCSLockTTL(t *testing.T, ha physical.HABackend) {
	// Set much smaller lock times to speed up the test.
	lockTTL := time.Second * 3
	renewInterval := time.Second * 1
	watchInterval := time.Second * 1

	// Get the lock
	origLock, err := ha.LockWith("gcsttl", "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// set the first lock renew period to double the expected TTL.
	lock := origLock.(*GCSLock)
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
	origLock2, err := ha.LockWith("gcsttl", "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	lock2 := origLock2.(*GCSLock)
	lock2.renewInterval = renewInterval
	lock2.ttl = lockTTL
	lock2.watchRetryInterval = watchInterval

	// // Cancel attempt in 10 sec so as not to block unit tests forever
	stopCh := make(chan struct{})
	time.AfterFunc(time.Second*10, func() {
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
