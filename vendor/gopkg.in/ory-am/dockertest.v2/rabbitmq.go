package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupRabbitMQContainer sets up a real RabbitMQ instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupRabbitMQContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 5672)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(RabbitMQImageName, port, 10*time.Second, func() (string, error) {
		res, err := run("--name", GenerateContainerID(), "-d", "-P", "-p", forward, RabbitMQImageName)
		return res, err
	})
	return
}

// ConnectToRabbitMQ starts a RabbitMQ image and passes the amqp url to the connector callback.
// The url will match the ip:port pattern (e.g. 123.123.123.123:4241)
func ConnectToRabbitMQ(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupRabbitMQContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up RabbitMQ container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up RabbitMQ container.")
}
