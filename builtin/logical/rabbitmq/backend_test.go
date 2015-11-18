package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
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
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadRole(t, "web", "administrator", `{"/": {"configure": ".*", "write": ".*", "read": ".*"}}`),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", "", ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if uri := os.Getenv("RABBITMQ_MG_URI"); uri == "" {
		t.Fatal("RABBITMQ_MG_URI must be set for acceptance tests")
	}
	if username := os.Getenv("RABBITMQ_MG_USERNAME"); username == "" {
		t.Fatal("RABBITMQ_MG_USERNAME must be set for acceptance tests")
	}
	if password := os.Getenv("RABBITMQ_MG_PASSWORD"); password == "" {
		t.Fatal("RABBITMQ_MG_PASSWORD must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"uri":      os.Getenv("RABBITMQ_MG_URI"),
			"username": os.Getenv("RABBITMQ_MG_USERNAME"),
			"password": os.Getenv("RABBITMQ_MG_PASSWORD"),
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
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

			uri := os.Getenv("RABBITMQ_MG_URI")

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
			if err := json.Unmarshal([]byte(rawVHosts), &vhosts); err != nil {
				return fmt.Errorf("bad expected vhosts %#v: %s", vhosts, err)
			}

			for host, permission := range vhosts {
				actualPermission, ok := d.VHosts[host]
				if !ok {
					return fmt.Errorf("expected vhost: %s", host)
				}

				if actualPermission.Configure != permission.Configure {
					fmt.Errorf("expected permission %s to be %s, got %s", "configure", permission.Configure, actualPermission.Configure)
				}

				if actualPermission.Write != permission.Write {
					fmt.Errorf("expected permission %s to be %s, got %s", "write", permission.Write, actualPermission.Write)
				}

				if actualPermission.Read != permission.Read {
					fmt.Errorf("expected permission %s to be %s, got %s", "read", permission.Read, actualPermission.Read)
				}
			}

			return nil
		},
	}
}
