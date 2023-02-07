package aerospike

import (
	"context"
	"testing"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v5"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestAerospikeBackend(t *testing.T) {
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
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/aerospike/aerospike-server",
		ContainerName: "aerospikedb",
		ImageTag:      "5.6.0.5",
		Ports:         []string{"3000/tcp", "3001/tcp", "3002/tcp", "3003/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start local Aerospike: %s", err)
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
		t.Fatalf("Could not start local Aerospike: %s", err)
	}

	return svc.Cleanup, &aerospikeConfig{
		hostname:  svc.Config.URL().Hostname(),
		port:      svc.Config.URL().Port(),
		namespace: "test",
		set:       "vault",
	}
}
