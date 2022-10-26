package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// PrepareTestContainer calls PrepareTestContainerWithDatabase without a
// database name value, which results in configuring a database named "test"
func PrepareTestContainer(t *testing.T, version string) (cleanup func(), retURL string) {
	return PrepareTestContainerWithDatabase(t, version, "")
}

// PrepareTestContainerWithDatabase configures a test container with a given
// database name, to test non-test/admin database configurations
func PrepareTestContainerWithDatabase(t *testing.T, version, dbName string) (func(), string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "mongo",
		ImageTag:  version,
		Ports:     []string{"27017/tcp"},
	})
	if err != nil {
		t.Fatalf("could not start docker mongo: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		connURL := fmt.Sprintf("mongodb://%s:%d", host, port)
		if dbName != "" {
			connURL = fmt.Sprintf("%s/%s", connURL, dbName)
		}

		ctx, _ = context.WithTimeout(context.Background(), 1*time.Minute)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
		if err != nil {
			return nil, err
		}

		err = client.Ping(ctx, readpref.Primary())
		if err = client.Disconnect(ctx); err != nil {
			t.Fatal()
		}

		return docker.NewServiceURLParse(connURL)
	})
	if err != nil {
		t.Fatalf("could not start docker mongo: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}
