package cassandra

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	dockertest "gopkg.in/ory-am/dockertest.v2"
)

var (
	testImagePull sync.Once
)

func prepareTestContainer(t *testing.T, s logical.Storage, b logical.Backend) (cid dockertest.ContainerID, retURL string) {
	if os.Getenv("CASSANDRA_HOST") != "" {
		return "", os.Getenv("CASSANDRA_HOST")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testImagePull.Do(func() {
		dockertest.Pull("cassandra")
	})

	cwd, _ := os.Getwd()

	cid, connErr := dockertest.ConnectToCassandra("latest", 60, 1000*time.Millisecond, func(connURL string) bool {
		// This will cause a validation to run
		resp, err := b.HandleRequest(&logical.Request{
			Storage:   s,
			Operation: logical.UpdateOperation,
			Path:      "config/connection",
			Data: map[string]interface{}{
				"hosts":            connURL,
				"username":         "cassandra",
				"password":         "cassandra",
				"protocol_version": 3,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			// It's likely not up and running yet, so return false and try again
			return false
		}

		retURL = connURL
		return true
	}, []string{"-v", cwd + "/test-fixtures/:/etc/cassandra/"}...)

	if connErr != nil {
		if cid != "" {
			cid.KillRemove()
		}
		t.Fatalf("could not connect to database: %v", connErr)
	}

	return
}

func cleanupTestContainer(t *testing.T, cid dockertest.ContainerID) {
	err := cid.KillRemove()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_basic(t *testing.T) {
	if os.Getenv("TRAVIS") != "true" {
		t.SkipNow()
	}
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, hostname := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, hostname),
			testAccStepRole(t),
			testAccStepReadCreds(t, "test"),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	if os.Getenv("TRAVIS") != "true" {
		t.SkipNow()
	}
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, hostname := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, hostname),
			testAccStepRole(t),
			testAccStepRoleWithOptions(t),
			testAccStepReadRole(t, "test", testRole),
			testAccStepReadRole(t, "test2", testRole),
			testAccStepDeleteRole(t, "test"),
			testAccStepDeleteRole(t, "test2"),
			testAccStepReadRole(t, "test", ""),
			testAccStepReadRole(t, "test2", ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CASSANDRA_HOST"); v == "" {
		t.Fatal("CASSANDRA_HOST must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T, hostname string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"hosts":            hostname,
			"username":         "cassandra",
			"password":         "cassandra",
			"protocol_version": 3,
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Data: map[string]interface{}{
			"creation_cql": testRole,
		},
	}
}

func testAccStepRoleWithOptions(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test2",
		Data: map[string]interface{}{
			"creation_cql": testRole,
			"lease":        "30s",
			"consistency":  "All",
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, cql string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if cql == "" {
					return nil
				}

				return fmt.Errorf("response is nil")
			}

			var d struct {
				CreationCQL string `mapstructure:"creation_cql"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.CreationCQL != cql {
				return fmt.Errorf("bad: %#v\n%#v\n%#v\n", resp, cql, d.CreationCQL)
			}

			return nil
		},
	}
}

const testRole = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
