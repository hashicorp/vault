package consul

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	dockertest "gopkg.in/ory-am/dockertest.v2"
)

var (
	testImagePull sync.Once
)

func PrepareConsulTestContainer(t *testing.T) (cid dockertest.ContainerID, retAddress string) {
	if os.Getenv("CONSUL_HTTP_ADDR") != "" {
		return "", os.Getenv("CONSUL_HTTP_ADDR")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testImagePull.Do(func() {
		dockertest.Pull(dockertest.ConsulImageName)
	})

	try := 0
	cid, connErr := dockertest.ConnectToConsul(60, 500*time.Millisecond, func(connAddress string) bool {
		try += 1
		// Build a client and verify that the credentials work
		config := api.DefaultConfig()
		config.Address = connAddress
		config.Token = dockertest.ConsulACLMasterToken
		client, err := api.NewClient(config)
		if err != nil {
			if try > 50 {
				panic(err)
			}
			return false
		}

		_, err = client.KV().Put(&api.KVPair{
			Key:   "setuptest",
			Value: []byte("setuptest"),
		}, nil)
		if err != nil {
			if try > 50 {
				panic(err)
			}
			return false
		}

		retAddress = connAddress
		return true
	})

	if connErr != nil {
		t.Fatalf("could not connect to consul: %v", connErr)
	}

	return
}

func CleanupConsulTestContainer(t *testing.T, cid dockertest.ContainerID) {
	err := cid.KillRemove()
	if err != nil {
		t.Fatal(err)
	}
}
