// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aerospike

import (
	"context"
	"math/bits"
	"strings"
	"testing"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v8"
	"github.com/docker/docker/api/types/container"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestAerospikeBackend(t *testing.T) {
	if bits.UintSize == 32 {
		t.Skip("Aerospike storage is only supported on 64-bit architectures")
	}
	cleanup, config := prepareAerospikeContainer(t)
	defer cleanup()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewAerospikeBackend(map[string]string{
		"hostname":  config.hostname,
		"port":      config.port,
		"namespace": config.namespace,
		"set":       config.set,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

type aerospikeConfig struct {
	hostname  string
	port      string
	namespace string
	set       string
}

func prepareAerospikeContainer(t *testing.T) (func(), *aerospikeConfig) {
	containerLogs := new(strings.Builder)
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/aerospike/aerospike-server",
		ContainerName: "aerospikedb",
		ImageTag:      "6.4",
		Ports:         []string{"3000/tcp", "3001/tcp", "3002/tcp", "3003/tcp"},
		LogConsumer: func(s string) {
			containerLogs.Write([]byte(s + "\n"))
		},
		Resources: container.Resources{
			// 15,000 is the default in 6.4 and Docker >= 29 uses containerd >= v2.1.5,
			// which uses systemd's default LimitNOFILE for containers, changing the
			// open file descriptor limit (ulimit -n) from 1048576 to 1024. Here we
			// explicitly allow more even though it certainly won't use them.
			Ulimits: []*container.Ulimit{{Name: "nofile", Soft: 15_000, Hard: 15_000}},
		},
	})
	if err != nil {
		time.Sleep(1 * time.Second) // Allow our log consumer to get all container logs
		t.Fatalf("Could not start local Aerospike: %s, container logs: %s", err, containerLogs.String())
	}

	svc, err := runner.StartService(context.Background(),
		func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
			cfg := docker.NewServiceHostPort(host, port)

			time.Sleep(time.Second)
			client, err := aero.NewClient(host, port)
			if err != nil {
				return nil, err
			}

			node, err := client.Cluster().GetRandomNode()
			if err != nil {
				return nil, err
			}

			_, err = node.RequestInfo(aero.NewInfoPolicy(), "namespaces")
			if err != nil {
				return nil, err
			}

			return cfg, nil
		},
	)
	if err != nil {
		time.Sleep(1 * time.Second) // Allow our log consumer to get all container logs
		t.Fatalf("Could not start local Aerospike: %s, container logs: %s", err, containerLogs.String())
	}

	return svc.Cleanup, &aerospikeConfig{
		hostname:  svc.Config.URL().Hostname(),
		port:      svc.Config.URL().Port(),
		namespace: "test",
		set:       "vault",
	}
}
