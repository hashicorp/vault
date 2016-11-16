package physical

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	log "github.com/mgutz/logxi/v1"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/vault/helper/logformat"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var ConsistencyDelays = delays{
	beforeList: 5 * time.Second,
	beforeGet:  0 * time.Second,
}

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

		time.Sleep(ConsistencyDelays.beforeList)
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
			// occassionally we get "storage: object doesn't exist" which is fine
			bucket.Object(objAttrs.Name).Delete(context.Background())
		}

		// not a list operation, but google lists to make sure the bucket is empty on delete
		time.Sleep(ConsistencyDelays.beforeList)
		err := bucket.Delete(context.Background())
		if err != nil {
			t.Fatalf("error deleting bucket '%s': '%v'", bucketName, err)
		}
	}()

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("gcs", logger, map[string]string{
		"bucket":           bucketName,
		"credentials_file": credentialsFile,
	})

	if err != nil {
		t.Fatalf("error creating google cloud storage backend: '%s'", err)
	}

	testEventuallyConsistentBackend(t, b, ConsistencyDelays)
	testEventuallyConsistentBackend_ListPrefix(t, b, ConsistencyDelays)

}
