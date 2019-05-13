package mongodb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
)

func PrepareTestContainer(t *testing.T, version string) (cleanup func(), retURL string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		t.Fatalf("Could not start local mongo docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	addr := fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp"))
	retURL = "mongodb://" + addr

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:   []string{addr},
			Timeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}
		defer session.Close()

		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
		return session.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to mongo docker container: %s", err)
	}

	return
}
