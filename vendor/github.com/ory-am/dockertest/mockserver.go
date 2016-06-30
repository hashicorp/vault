package dockertest

import (
	"fmt"
	"time"
	"log"
	"github.com/go-errors/errors"
)

// SetupMockserverContainer sets up a real Mockserver instance for testing purposes
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupMockserverContainer() (c ContainerID, ip string, mockPort, proxyPort int, err error) {
	mockPort = RandomPort()
	proxyPort = RandomPort()

	mockForward := fmt.Sprintf("%d:%d", mockPort, 1080)
	proxyForward := fmt.Sprintf("%d:%d", proxyPort, 1090)

	if BindDockerToLocalhost != "" {
		mockForward = "127.0.0.1:" + mockForward
		proxyForward = "127.0.0.1:" + proxyForward
	}

	c, ip, err = SetupMultiportContainer(RabbitMQImageName, []int{ mockPort, proxyPort}, 10*time.Second, func() (string, error) {
		res, err := run("--name", GenerateContainerID(), "-d", "-P", "-p", mockForward, "-p", proxyForward, MockserverImageName)
		return res, err
	})
	return
}

// ConnectToMockserver starts a Mockserver image and passes the mock and proxy urls to the connector callback functions.
// The urls will match the http://ip:port pattern (e.g. http://123.123.123.123:4241)
func ConnectToMockserver(tries int, delay time.Duration, mockConnector func(url string) bool, proxyConnector func(url string) bool) (c ContainerID, err error) {
	c, ip, mockPort, proxyPort, err := SetupMockserverContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up Mockserver container: %v", err)
	}

	var mockOk, proxyOk bool

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)

		if !mockOk {
			if mockConnector(fmt.Sprintf("http://%s:%d", ip, mockPort)) {
				mockOk = true
			} else {
				log.Printf("Try %d failed for mock. Retrying.", try)
			}
		}
		if !proxyOk {
			if proxyConnector(fmt.Sprintf("http://%s:%d", ip, proxyPort)) {
				proxyOk = true
			} else {
				log.Printf("Try %d failed for proxy. Retrying.", try)
			}
		}
	}

	if mockOk && proxyOk {
		return c, nil
	} else {
		return c, errors.New("Could not set up Mockserver container.")
	}
}