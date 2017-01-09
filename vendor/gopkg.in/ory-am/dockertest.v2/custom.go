package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupCustomContainer sets up a real an instance of the given image for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupCustomContainer(imageName string, exposedPort int, timeOut time.Duration, extraDockerArgs ...string) (c ContainerID, ip string, localPort int, err error) {
	localPort = RandomPort()
	forward := fmt.Sprintf("%d:%d", localPort, exposedPort)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(imageName, localPort, timeOut, func() (string, error) {
		args := make([]string, 0, len(extraDockerArgs)+7)
		args = append(args, "--name", GenerateContainerID(), "-d", "-P", "-p", forward)
		args = append(args, extraDockerArgs...)
		args = append(args, imageName)
		return run(args...)
	})
	return
}

// ConnectToCustomContainer attempts to connect to a custom container until successful or the maximum number of tries is reached.
func ConnectToCustomContainer(url string, tries int, delay time.Duration, connector func(url string) bool) error {
	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		if connector(url) {
			return nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return errors.New("Could not set up custom container.")
}
