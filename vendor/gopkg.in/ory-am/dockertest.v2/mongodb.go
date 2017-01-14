package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupMongoContainer sets up a real MongoDB instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupMongoContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 27017)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(MongoDBImageName, port, 10*time.Second, func() (string, error) {
		res, err := run("--name", GenerateContainerID(), "-d", "-P", "-p", forward, MongoDBImageName)
		return res, err
	})
	return
}

// ConnectToMongoDB starts a MongoDB image and passes the database url to the connector callback.
// The url will match the ip:port pattern (e.g. 123.123.123.123:4241)
func ConnectToMongoDB(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupMongoContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up MongoDB container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up MongoDB container.")
}
