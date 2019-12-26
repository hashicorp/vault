package docker

import (
	"strings"
	"time"

	"github.com/mitchellh/go-testing-interface"
	"github.com/ory/dockertest"
)

func CleanupResource(t testing.T, pool *dockertest.Pool, resource *dockertest.Resource) {
	var err error
	for i := 0; i < 10; i++ {
		err = pool.Purge(resource)
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}

	if strings.Contains(err.Error(), "No such container") {
		return
	}
	t.Fatalf("Failed to cleanup local container: %s", err)
}
