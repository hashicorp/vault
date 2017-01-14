package dockertest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	// ConsulDatacenter must be defined when starting a Consul datacenter; this
	// value will be used for both the datacenter and the ACL datacenter
	ConsulDatacenter = "test"

	// ConsulACLDefaultPolicy defines the default policy to use with Consul ACLs
	ConsulACLDefaultPolicy = "deny"

	// ConsulACLMasterToken defines the master ACL token
	ConsulACLMasterToken = "test"

	// A function with no arguments that outputs a valid JSON string to be used
	// as the value of the environment variable CONSUL_LOCAL_CONFIG.
	ConsulLocalConfigGen = DefaultConsulLocalConfig
)

func DefaultConsulLocalConfig() (string, error) {
	type d struct {
		Datacenter       string `json:"datacenter,omitempty"`
		ACLDatacenter    string `json:"acl_datacenter,omitempty"`
		ACLDefaultPolicy string `json:"acl_default_policy,omitempty"`
		ACLMasterToken   string `json:"acl_master_token,omitempty"`
	}

	vals := &d{
		Datacenter:       ConsulDatacenter,
		ACLDatacenter:    ConsulDatacenter,
		ACLDefaultPolicy: ConsulACLDefaultPolicy,
		ACLMasterToken:   ConsulACLMasterToken,
	}

	ret, err := json.Marshal(vals)
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

// SetupConsulContainer sets up a real Consul instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupConsulContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 8500)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	localConfig, err := ConsulLocalConfigGen()
	if err != nil {
		return "", "", 0, err
	}
	c, ip, err = SetupContainer(ConsulImageName, port, 15*time.Second, func() (string, error) {
		return run(
			"--name", GenerateContainerID(),
			"-d",
			"-p", forward,
			"-e", fmt.Sprintf("CONSUL_LOCAL_CONFIG=%s", localConfig),
			ConsulImageName,
			"agent",
			"-dev",               // Run in dev mode
			"-client", "0.0.0.0", // Allow clients from any IP, otherwise the bridge IP will be where clients come from and it will be rejected
		)
	})
	return
}

// ConnectToConsul starts a Consul image and passes the address to the
// connector callback function.
func ConnectToConsul(tries int, delay time.Duration, connector func(address string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupConsulContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up Consul container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		address := fmt.Sprintf("%s:%d", ip, port)
		if connector(address) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up Consul container.")
}
