// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

// LatestConsulVersion is the most recent version of Consul which is used unless
// another version is specified in the test config or environment. This will
// probably go stale as we don't always update it on every release but we rarely
// rely on specific new Consul functionality so that's probably not a problem.
const LatestConsulVersion = "1.15.3"

type Config struct {
	docker.ServiceHostPort
	Token             string
	ContainerHTTPAddr string
}

func (c *Config) APIConfig() *consulapi.Config {
	apiConfig := consulapi.DefaultConfig()
	apiConfig.Address = c.Address()
	apiConfig.Token = c.Token
	return apiConfig
}

// PrepareTestContainer is a test helper that creates a Consul docker container
// or fails the test if unsuccessful. See RunContainer for more details on the
// configuration.
func PrepareTestContainer(t *testing.T, version string, isEnterprise bool, doBootstrapSetup bool) (func(), *Config) {
	t.Helper()

	cleanup, config, err := RunContainer(context.Background(), "", version, isEnterprise, doBootstrapSetup)
	if err != nil {
		t.Fatalf("failed starting consul: %s", err)
	}
	return cleanup, config
}

// RunContainer runs Consul in a Docker container unless CONSUL_HTTP_ADDR is
// already found in the environment. Consul version is determined by the version
// argument. If version is empty string, the CONSUL_DOCKER_VERSION environment
// variable is used and if that is empty too, LatestConsulVersion is used
// (defined above). If namePrefix is provided we assume you have chosen a unique
// enough prefix to avoid collision with other tests that may be running in
// parallel and so _do not_ add an additional unique ID suffix. We will also
// ensure previous instances are deleted and leave the container running for
// debugging. This is useful for using Consul as part of at testcluster (i.e.
// when Vault is in Docker too). If namePrefix is empty then a unique suffix is
// added since many older tests rely on a uniq instance of the container. This
// is used by `PrepareTestContainer` which is used typically in tests that rely
// on Consul but run tested code within the test process.
func RunContainer(ctx context.Context, namePrefix, version string, isEnterprise bool, doBootstrapSetup bool) (func(), *Config, error) {
	if retAddress := os.Getenv("CONSUL_HTTP_ADDR"); retAddress != "" {
		shp, err := docker.NewServiceHostPortParse(retAddress)
		if err != nil {
			return nil, nil, err
		}
		return func() {}, &Config{ServiceHostPort: *shp, Token: os.Getenv("CONSUL_HTTP_TOKEN")}, nil
	}

	config := `acl { enabled = true default_policy = "deny" }`
	if version == "" {
		consulVersion := os.Getenv("CONSUL_DOCKER_VERSION")
		if consulVersion != "" {
			version = consulVersion
		} else {
			version = LatestConsulVersion
		}
	}
	if strings.HasPrefix(version, "1.3") {
		config = `datacenter = "test" acl_default_policy = "deny" acl_datacenter = "test" acl_master_token = "test"`
	}

	name := "consul"
	repo := "docker.mirror.hashicorp.services/library/consul"
	var envVars []string
	// If running the enterprise container, set the appropriate values below.
	if isEnterprise {
		version += "-ent"
		name = "consul-enterprise"
		repo = "docker.mirror.hashicorp.services/hashicorp/consul-enterprise"
		license, hasLicense := os.LookupEnv("CONSUL_LICENSE")
		envVars = append(envVars, "CONSUL_LICENSE="+license)

		if !hasLicense {
			return nil, nil, fmt.Errorf("Failed to find enterprise license")
		}
	}
	if namePrefix != "" {
		name = namePrefix + name
	}

	if dockerRepo, hasEnvRepo := os.LookupEnv("CONSUL_DOCKER_REPO"); hasEnvRepo {
		repo = dockerRepo
	}

	dockerOpts := docker.RunOptions{
		ContainerName: name,
		ImageRepo:     repo,
		ImageTag:      version,
		Env:           envVars,
		Cmd:           []string{"agent", "-dev", "-client", "0.0.0.0", "-hcl", config},
		Ports:         []string{"8500/tcp"},
		AuthUsername:  os.Getenv("CONSUL_DOCKER_USERNAME"),
		AuthPassword:  os.Getenv("CONSUL_DOCKER_PASSWORD"),
	}

	// Add a unique suffix if there is no per-test prefix provided
	addSuffix := true
	if namePrefix != "" {
		// Don't add a suffix if the caller already provided a prefix
		addSuffix = false
		// Also enable predelete and non-removal to make debugging easier for test
		// cases with named containers).
		dockerOpts.PreDelete = true
		dockerOpts.DoNotAutoRemove = true
	}

	runner, err := docker.NewServiceRunner(dockerOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not start docker Consul: %s", err)
	}

	svc, _, err := runner.StartNewService(ctx, addSuffix, false, func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		shp := docker.NewServiceHostPort(host, port)
		apiConfig := consulapi.DefaultNonPooledConfig()
		apiConfig.Address = shp.Address()
		consul, err := consulapi.NewClient(apiConfig)
		if err != nil {
			return nil, err
		}

		// Make sure Consul is up
		if _, err = consul.Status().Leader(); err != nil {
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
		var consulToken string
		if doBootstrapSetup {
			aclbootstrap, _, err := consul.ACL().Bootstrap()
			if err != nil {
				return nil, err
			}
			consulToken = aclbootstrap.SecretID
			policy := &consulapi.ACLPolicy{
				Name:        "test",
				Description: "test",
				Rules: `node_prefix "" {
					policy = "write"
				}

				service_prefix "" {
					policy = "read"
				}`,
			}
			q := &consulapi.WriteOptions{
				Token: consulToken,
			}
			_, _, err = consul.ACL().PolicyCreate(policy, q)
			if err != nil {
				return nil, err
			}

			// Create a Consul role that contains the test policy, for Consul 1.5 and newer
			currVersion, _ := goversion.NewVersion(version)
			roleVersion, _ := goversion.NewVersion("1.5")
			if currVersion.GreaterThanOrEqual(roleVersion) {
				ACLList := []*consulapi.ACLTokenRoleLink{{Name: "test"}}

				role := &consulapi.ACLRole{
					Name:        "role-test",
					Description: "consul roles test",
					Policies:    ACLList,
				}

				_, _, err = consul.ACL().RoleCreate(role, q)
				if err != nil {
					return nil, err
				}
			}

			// Configure a namespace and partition if testing enterprise Consul
			if isEnterprise {
				// Namespaces require Consul 1.7 or newer
				namespaceVersion, _ := goversion.NewVersion("1.7")
				if currVersion.GreaterThanOrEqual(namespaceVersion) {
					namespace := &consulapi.Namespace{
						Name:        "ns1",
						Description: "ns1 test",
					}

					_, _, err = consul.Namespaces().Create(namespace, q)
					if err != nil {
						return nil, err
					}

					nsPolicy := &consulapi.ACLPolicy{
						Name:        "ns-test",
						Description: "namespace test",
						Namespace:   "ns1",
						Rules: `service_prefix "" {
							policy = "read"
						}`,
					}
					_, _, err = consul.ACL().PolicyCreate(nsPolicy, q)
					if err != nil {
						return nil, err
					}
				}

				// Partitions require Consul 1.11 or newer
				partitionVersion, _ := goversion.NewVersion("1.11")
				if currVersion.GreaterThanOrEqual(partitionVersion) {
					partition := &consulapi.Partition{
						Name:        "part1",
						Description: "part1 test",
					}

					_, _, err = consul.Partitions().Create(ctx, partition, q)
					if err != nil {
						return nil, err
					}

					partPolicy := &consulapi.ACLPolicy{
						Name:        "part-test",
						Description: "partition test",
						Partition:   "part1",
						Rules: `service_prefix "" {
							policy = "read"
						}`,
					}
					_, _, err = consul.ACL().PolicyCreate(partPolicy, q)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		return &Config{
			ServiceHostPort: *shp,
			Token:           consulToken,
		}, nil
	})
	if err != nil {
		return nil, nil, err
	}

	// Find the container network info.
	if len(svc.Container.NetworkSettings.Networks) < 1 {
		svc.Cleanup()
		return nil, nil, fmt.Errorf("failed to find any network settings for container")
	}
	cfg := svc.Config.(*Config)
	for _, eps := range svc.Container.NetworkSettings.Networks {
		// Just pick the first network, we assume only one for now.
		// Pull out the real container IP and set that up
		cfg.ContainerHTTPAddr = fmt.Sprintf("http://%s:8500", eps.IPAddress)
		break
	}
	return svc.Cleanup, cfg, nil
}
