package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupZooKeeperContainer sets up a real ZooKeeper node for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupZooKeeperContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 2181)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}

	c, ip, err = SetupContainer(ZooKeeperImageName, port, 10*time.Second, func() (string, error) {
		return run("--name", GenerateContainerID(), "-d", "-p", forward, ZooKeeperImageName)
	})
	return
}

// ConnectToZooKeeper starts a ZooKeeper image and passes the nodes connection string to the connector callback function.
// The connection string will match the ip:port pattern.
func ConnectToZooKeeper(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupZooKeeperContainer()
	if err != nil {
		return c, fmt.Errorf("Could not setup ZooKeeper container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not setup ZooKeeper container.")
}
