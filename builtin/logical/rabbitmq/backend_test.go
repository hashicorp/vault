// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/mitchellh/mapstructure"
)

const (
	envRabbitMQConnectionURI = "RABBITMQ_CONNECTION_URI"
	envRabbitMQUsername      = "RABBITMQ_USERNAME"
	envRabbitMQPassword      = "RABBITMQ_PASSWORD"
)

const (
	testTags        = "administrator"
	testVHosts      = `{"/": {"configure": ".*", "write": ".*", "read": ".*"}}`
	testVHostTopics = `{"/": {"amq.topic": {"write": ".*", "read": ".*"}}}`

	roleName = "web"
)

func prepareRabbitMQTestContainer(t *testing.T) (func(), string) {
	if os.Getenv(envRabbitMQConnectionURI) != "" {
		return func() {}, os.Getenv(envRabbitMQConnectionURI)
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/library/rabbitmq",
		ImageTag:      "3-management",
		ContainerName: "rabbitmq",
		Ports:         []string{"15672/tcp"},
	})
	if err != nil {
		t.Fatalf("could not start docker rabbitmq: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		connURL := fmt.Sprintf("http://%s:%d", host, port)
		rmqc, err := rabbithole.NewClient(connURL, "guest", "guest")
		if err != nil {
			return nil, err
		}

		_, err = rmqc.Overview()
		if err != nil {
			return nil, err
		}

		return docker.NewServiceURLParse(connURL)
	})
	if err != nil {
		t.Fatalf("could not start docker rabbitmq: %s", err)
	}
	return svc.Cleanup, svc.Config.URL().String()
}

func TestBackend_basic(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	cleanup, uri := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri, ""),
			testAccStepRole(t),
			testAccStepReadCreds(t, b, uri, roleName),
		},
	})
}

func TestBackend_returnsErrs(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	cleanup, uri := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri, ""),
			{
				Operation: logical.CreateOperation,
				Path:      fmt.Sprintf("roles/%s", roleName),
				Data: map[string]interface{}{
					"tags":         testTags,
					"vhosts":       `{"invalid":{"write": ".*", "read": ".*"}}`,
					"vhost_topics": testVHostTopics,
				},
			},
			{
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("creds/%s", roleName),
				ErrorOk:   true,
			},
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	cleanup, uri := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri, ""),
			testAccStepRole(t),
			testAccStepReadRole(t, roleName, testTags, testVHosts, testVHostTopics),
			testAccStepDeleteRole(t, roleName),
			testAccStepReadRole(t, roleName, "", "", ""),
		},
	})
}

func TestBackend_roleWithPasswordPolicy(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env %q set", logicaltest.TestEnvVar))
		return
	}

	backendConfig := logical.TestBackendConfig()
	passGen := func() (password string, err error) {
		return base62.Random(30)
	}
	backendConfig.System.(*logical.StaticSystemView).SetPasswordPolicy("testpolicy", passGen)
	b, _ := Factory(context.Background(), backendConfig)

	cleanup, uri := prepareRabbitMQTestContainer(t)
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck:       testAccPreCheckFunc(t, uri),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, uri, "testpolicy"),
			testAccStepRole(t),
			testAccStepReadCreds(t, b, uri, roleName),
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

func testAccStepConfig(t *testing.T, uri string, passwordPolicy string) logicaltest.TestStep {
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
			"connection_uri":  uri,
			"username":        username,
			"password":        password,
			"password_policy": passwordPolicy,
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("roles/%s", roleName),
		Data: map[string]interface{}{
			"tags":         testTags,
			"vhosts":       testVHosts,
			"vhost_topics": testVHostTopics,
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

func testAccStepReadRole(t *testing.T, name, tags, rawVHosts string, rawVHostTopics string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if tags == "" && rawVHosts == "" && rawVHostTopics == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Tags        string                                     `mapstructure:"tags"`
				VHosts      map[string]vhostPermission                 `mapstructure:"vhosts"`
				VHostTopics map[string]map[string]vhostTopicPermission `mapstructure:"vhost_topics"`
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

			var vhostTopics map[string]map[string]vhostTopicPermission
			if err := jsonutil.DecodeJSON([]byte(rawVHostTopics), &vhostTopics); err != nil {
				return fmt.Errorf("bad expected vhostTopics %#v: %s", vhostTopics, err)
			}

			for host, permissions := range vhostTopics {
				for exchange, permission := range permissions {
					actualPermissions, ok := d.VHostTopics[host]
					if !ok {
						return fmt.Errorf("expected vhost topics: %s", host)
					}

					actualPermission, ok := actualPermissions[exchange]
					if !ok {
						return fmt.Errorf("expected vhost topic exchange: %s", exchange)
					}

					if actualPermission.Write != permission.Write {
						return fmt.Errorf("expected permission %s to be %s, got %s", "write", permission.Write, actualPermission.Write)
					}

					if actualPermission.Read != permission.Read {
						return fmt.Errorf("expected permission %s to be %s, got %s", "read", permission.Read, actualPermission.Read)
					}
				}
			}

			return nil
		},
	}
}
