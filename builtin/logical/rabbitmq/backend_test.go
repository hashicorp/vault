package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/mitchellh/mapstructure"
)

// Set the following env vars for the below test case to work.
//
// RABBITMQ_CONNECTION_URI
// RABBITMQ_USERNAME
// RABBITMQ_PASSWORD
func TestBackend_basic(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadCreds(t, b, "web"),
		},
	})

}

func TestBackend_roleCrud(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadRole(t, "web", "administrator", `{"/": {"configure": ".*", "write": ".*", "read": ".*"}}`),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", "", ""),
		},
	})
}

const (
	envRabbitMQConnectionURI = "RABBITMQ_CONNECTION_URI"
	envRabbitMQUsername      = "RABBITMQ_USERNAME"
	envRabbitMQPassword      = "RABBITMQ_PASSWORD"
)

func testAccPreCheck(t *testing.T) {
	if uri := os.Getenv(envRabbitMQConnectionURI); uri == "" {
		t.Fatalf(fmt.Sprintf("%s must be set for acceptance tests", envRabbitMQConnectionURI))
	}
	if username := os.Getenv(envRabbitMQUsername); username == "" {
		t.Fatalf(fmt.Sprintf("%s must be set for acceptance tests", envRabbitMQUsername))
	}
	if password := os.Getenv(envRabbitMQPassword); password == "" {
		t.Fatalf(fmt.Sprintf("%s must be set for acceptance tests", envRabbitMQPassword))
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"connection_uri": os.Getenv(envRabbitMQConnectionURI),
			"username":       os.Getenv(envRabbitMQUsername),
			"password":       os.Getenv(envRabbitMQPassword),
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

func testAccStepReadCreds(t *testing.T, b logical.Backend, name string) logicaltest.TestStep {
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

			uri := os.Getenv(envRabbitMQConnectionURI)

			client, err := rabbithole.NewClient(uri, d.Username, d.Password)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.ListVhosts()
			if err != nil {
				t.Fatalf("unable to list vhosts with generated credentials: %s", err)
			}

			resp, err = b.HandleRequest(&logical.Request{
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
					return fmt.Errorf("Error on resp: %#v", *resp)
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
