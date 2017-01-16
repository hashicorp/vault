package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupEtcdContainer sets up a real etcd instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupEtcdContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 2379)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(EtcdImageName, port, 10*time.Second, func() (string, error) {
		return run(
			"--name", GenerateContainerID(),
			"-d",
			"-p", forward,
			EtcdImageName,
			"etcd",
			"-name", "etcd-test",
			"-advertise-client-urls", "http://127.0.0.1:2379",
			"-listen-client-urls", "http://0.0.0.0:2379", // Allow clients from any IP
		)
	})

	return c, ip, port, err
}

// ConnectToEtcd starts a etcd image and passes the address to the
// connector callback function.
func ConnectToEtcd(tries int, delay time.Duration, connector func(address string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupEtcdContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up etcd container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		address := fmt.Sprintf("%s:%d", ip, port)
		if connector(address) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up etcd container.")
}
