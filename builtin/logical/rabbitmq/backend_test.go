package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/dockertest"
)

const (
	envRabbitMQConnectionURI = "RABBITMQ_CONNECTION_URI"
	envRabbitMQUsername      = "RABBITMQ_USERNAME"
	envRabbitMQPassword      = "RABBITMQ_PASSWORD"
)

func prepareRabbitMQTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv(envRabbitMQConnectionURI) != "" {
		return func() {}, os.Getenv(envRabbitMQConnectionURI), 0
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	runOpts := &dockertest.RunOptions{
		Repository: "rabbitmq",
		Tag:        "3-management",
	}
	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		t.Fatalf("Could not start local rabbitmq docker container: %s", err)
	}

	cleanup := func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	port, _ := strconv.Atoi(resource.GetPort("15672/tcp"))
	address := fmt.Sprintf("http://127.0.0.1:%d", port)

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		rmqc, err := rabbithole.NewClient(address, "guest", "guest")
		if err != nil {
			return err
		}

		_, err = rmqc.Overview()
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to rabbitmq docker container: %s", err)
	}
	return cleanup, address, port
}

func TestBackend_basic(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}
	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	cleanup, uri, _ := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri),
			testAccStepRole(t),
			testAccStepReadCreds(t, b, uri, "web"),
		},
	})

}

func TestBackend_roleCrud(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}
	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	cleanup, uri, _ := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri),
			testAccStepRole(t),
			testAccStepReadRole(t, "web", "administrator", `{"/": {"configure": ".*", "write": ".*", "read": ".*"}}`),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", "", ""),
		},
	})
}

func testAccPreCheckFunc(t *testing.T, uri string) func() {
	return func() {
		if uri == "" {
			t.Fatal("RabbitMQ URI must be set for acceptance tests")
		}
	}
}

func testAccStepConfig(t *testing.T, uri string) logicaltest.TestStep {
	username := os.Getenv(envRabbitMQUsername)
	if len(username) == 0 {
		username = "guest"
	}
	password := os.Getenv(envRabbitMQPassword)
	if len(password) == 0 {
		password = "guest"
	}

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"connection_uri": uri,
			"username":       username,
			"password":       password,
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"tags":   "administrator",
			"vhosts": `{"/": {"configure": ".*", "write": ".*", "read": ".*"}}`,
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, uri, name string) logicaltest.TestStep {
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

			client, err := rabbithole.NewClient(uri, d.Username, d.Password)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.ListVhosts()
			if err != nil {
				t.Fatalf("unable to list vhosts with generated credentials: %s", err)
			}

			resp, err = b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.RevokeOperation,
				Secret: &logical.Secret{
					InternalData: map[string]interface{}{
						"secret_type": "creds",
						"username":    d.Username,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp != nil {
				if resp.IsError() {
					return fmt.Errorf("error on resp: %#v", *resp)
				}
			}

			client, err = rabbithole.NewClient(uri, d.Username, d.Password)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.ListVhosts()
			if err == nil {
				t.Fatalf("expected to fail listing vhosts: %s", err)
			}

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name, tags, rawVHosts string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if tags == "" && rawVHosts == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Tags   string                     `mapstructure:"tags"`
				VHosts map[string]vhostPermission `mapstructure:"vhosts"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Tags != tags {
				return fmt.Errorf("bad: %#v", resp)
			}

			var vhosts map[string]vhostPermission
			if err := jsonutil.DecodeJSON([]byte(rawVHosts), &vhosts); err != nil {
				return fmt.Errorf("bad expected vhosts %#v: %s", vhosts, err)
			}

			for host, permission := range vhosts {
				actualPermission, ok := d.VHosts[host]
				if !ok {
					return fmt.Errorf("expected vhost: %s", host)
				}

				if actualPermission.Configure != permission.Configure {
					return fmt.Errorf("expected permission %s to be %s, got %s", "configure", permission.Configure, actualPermission.Configure)
				}

				if actualPermission.Write != permission.Write {
					return fmt.Errorf("expected permission %s to be %s, got %s", "write", permission.Write, actualPermission.Write)
				}

				if actualPermission.Read != permission.Read {
					return fmt.Errorf("expected permission %s to be %s, got %s", "read", permission.Read, actualPermission.Read)
				}
			}

			return nil
		},
	}
}
