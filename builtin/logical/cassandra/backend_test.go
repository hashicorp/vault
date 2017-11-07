package cassandra

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var (
	testImagePull sync.Once
)

func prepareCassandraTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv("CASSANDRA_HOST") != "" {
		return func() {}, os.Getenv("CASSANDRA_HOST"), 0
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	cwd, _ := os.Getwd()
	cassandraMountPath := fmt.Sprintf("%s/test-fixtures/:/etc/cassandra/", cwd)

	ro := &dockertest.RunOptions{
		Repository: "cassandra",
		Tag:        "latest",
		Env:        []string{"CASSANDRA_BROADCAST_ADDRESS=127.0.0.1"},
		Mounts:     []string{cassandraMountPath},
	}
	resource, err := pool.RunWithOptions(ro)
	if err != nil {
		t.Fatalf("Could not start local cassandra docker container: %s", err)
	}

	cleanup := func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	port, _ := strconv.Atoi(resource.GetPort("9042/tcp"))
	address := fmt.Sprintf("127.0.0.1:%d", port)

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		clusterConfig := gocql.NewCluster(address)
		clusterConfig.Authenticator = gocql.PasswordAuthenticator{
			Username: "cassandra",
			Password: "cassandra",
		}
		clusterConfig.ProtoVersion = 4
		clusterConfig.Port = port

		session, err := clusterConfig.CreateSession()
		if err != nil {
			return fmt.Errorf("error creating session: %s", err)
		}
		defer session.Close()
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to cassandra docker container: %s", err)
	}
	return cleanup, address, port
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

	cleanup, hostname, _ := prepareCassandraTestContainer(t)
	defer cleanup()

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

	cleanup, hostname, _ := prepareCassandraTestContainer(t)
	defer cleanup()

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
