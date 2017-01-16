package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupCassandraContainer sets up a real Cassandra node for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupCassandraContainer(versionTag string, optionalParams ...string) (c ContainerID, ip string, port int, err error) {
	port = RandomPort()

	// Forward for the CQL port.
	forward := fmt.Sprintf("%d:%d", port, 9042)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}

	imageName := fmt.Sprintf("%s:%s", CassandraImageName, versionTag)

	c, ip, err = SetupContainer(imageName, port, 10*time.Second, func() (string, error) {
		return run(append(optionalParams, "--name", GenerateContainerID(), "-d", "-p", forward, imageName)...)
	})
	return
}

// ConnectToCassandra starts a Cassandra image and passes the nodes connection string to the connector callback function.
// The connection string will match the ip:port pattern, where port is the mapped CQL port.
func ConnectToCassandra(versionTag string, tries int, delay time.Duration, connector func(url string) bool, optionalParams ...string) (c ContainerID, err error) {
	c, ip, port, err := SetupCassandraContainer(versionTag, optionalParams...)
	if err != nil {
		return c, fmt.Errorf("Could not setup Cassandra container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not setup Cassandra container.")
}
