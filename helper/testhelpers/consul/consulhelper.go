package consul

import (
	"context"
	"os"
	"strings"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

type Config struct {
	docker.ServiceHostPort
	Token string
}

func (c *Config) APIConfig() *consulapi.Config {
	apiConfig := consulapi.DefaultConfig()
	apiConfig.Address = c.Address()
	apiConfig.Token = c.Token
	return apiConfig
}

// PrepareTestContainer creates a Consul docker container.  If version is empty,
// the Consul version used will be given by the environment variable
// CONSUL_DOCKER_VERSION, or if that's empty, whatever we've hardcoded as the
// the latest Consul version.
func PrepareTestContainer(t *testing.T, version string) (func(), *Config) {
	if retAddress := os.Getenv("CONSUL_HTTP_ADDR"); retAddress != "" {
		shp, err := docker.NewServiceHostPortParse(retAddress)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{ServiceHostPort: *shp, Token: os.Getenv("CONSUL_HTTP_TOKEN")}
	}

	config := `acl { enabled = true default_policy = "deny" }`
	if version == "" {
		consulVersion := os.Getenv("CONSUL_DOCKER_VERSION")
		if consulVersion != "" {
			version = consulVersion
		} else {
			version = "1.7.2" // Latest Consul version, update as new releases come out
		}
	}
	if strings.HasPrefix(version, "1.3") {
		config = `datacenter = "test" acl_default_policy = "deny" acl_datacenter = "test" acl_master_token = "test"`
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "consul",
		ImageTag:  version,
		Cmd:       []string{"agent", "-dev", "-client", "0.0.0.0", "-hcl", config},
		Ports:     []string{"8500/tcp"},
		// TODO auth untested on this branch
		Repository:   os.Getenv("CONSUL_DOCKER_REPO"),
		AuthUsername: os.Getenv("CONSUL_DOCKER_USERNAME"),
		AuthPassword: os.Getenv("CONSUL_DOCKER_PASSWORD"),
	})

	if err != nil {
		t.Fatalf("Could not start docker Consul: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		shp := docker.NewServiceHostPort(host, port)
		apiConfig := consulapi.DefaultNonPooledConfig()
		apiConfig.Address = shp.Address()
		consul, err := consulapi.NewClient(apiConfig)
		if err != nil {
			return nil, err
		}

		// For version of Consul < 1.4
		if strings.HasPrefix(version, "1.3") {
			consulToken := "test"
			_, err = consul.KV().Put(&consulapi.KVPair{
				Key:   "setuptest",
				Value: []byte("setuptest"),
			}, &consulapi.WriteOptions{
				Token: consulToken,
			})
			if err != nil {
				return nil, err
			}
			return &Config{
				ServiceHostPort: *shp,
				Token:           consulToken,
			}, nil
		}

		// New default behavior
		aclbootstrap, _, err := consul.ACL().Bootstrap()
		if err != nil {
			return nil, err
		}
		consulToken := aclbootstrap.SecretID
		policy := &consulapi.ACLPolicy{
			Name:        "test",
			Description: "test",
			Rules: `node_prefix "" {
                policy = "write"
              }

              service_prefix "" {
                policy = "read"
              }
      `,
		}
		q := &consulapi.WriteOptions{
			Token: consulToken,
		}
		_, _, err = consul.ACL().PolicyCreate(policy, q)
		if err != nil {
			return nil, err
		}
		return &Config{
			ServiceHostPort: *shp,
			Token:           consulToken,
		}, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker Consul: %s", err)
	}
	return svc.Cleanup, svc.Config.(*Config)
}
