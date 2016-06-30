package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupNSQdContainer sets up a real NSQ instance for testing purposes
// using a Docker container and executing `/nsqd`. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupNSQdContainer() (c ContainerID, ip string, tcpPort int, httpPort int, err error) {
	// --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd --broadcast-address=192.168.99.100 --lookupd-tcp-address=192.168.99.100:4160
	tcpPort = RandomPort()
	httpPort = RandomPort()
	tcpForward := fmt.Sprintf("%d:%d", tcpPort, 4150)
	if BindDockerToLocalhost != "" {
		tcpForward = "127.0.0.1:" + tcpForward
	}

	httpForward := fmt.Sprintf("%d:%d", httpPort, 4151)
	if BindDockerToLocalhost != "" {
		httpForward = "127.0.0.1:" + httpForward
	}

	c, ip, err = SetupContainer(NSQImageName, tcpPort, 15*time.Second, func() (string, error) {
		return run("--name", GenerateContainerID(), "-d", "-P", "-p", tcpForward, "-p", httpForward, NSQImageName, "/nsqd", fmt.Sprintf("--broadcast-address=%s", ip), fmt.Sprintf("--lookupd-tcp-address=%s:4160", ip))
	})
	return
}

// SetupNSQLookupdContainer sets up a real NSQ instance for testing purposes
// using a Docker container and executing `/nsqlookupd`. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupNSQLookupdContainer() (c ContainerID, ip string, tcpPort int, httpPort int, err error) {
	// docker run --name lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd
	tcpPort = RandomPort()
	httpPort = RandomPort()
	tcpForward := fmt.Sprintf("%d:%d", tcpPort, 4160)
	if BindDockerToLocalhost != "" {
		tcpForward = "127.0.0.1:" + tcpForward
	}

	httpForward := fmt.Sprintf("%d:%d", httpPort, 4161)
	if BindDockerToLocalhost != "" {
		httpForward = "127.0.0.1:" + httpForward
	}

	c, ip, err = SetupContainer(NSQImageName, tcpPort, 15*time.Second, func() (string, error) {
		return run("--name", GenerateContainerID(), "-d", "-P", "-p", tcpForward, "-p", httpForward, NSQImageName, "/nsqlookupd")
	})
	return
}

// ConnectToNSQLookupd starts a NSQ image with `/nsqlookupd` running and passes the IP, HTTP port, and TCP port to the connector callback function.
// The url will match the ip pattern (e.g. 123.123.123.123).
func ConnectToNSQLookupd(tries int, delay time.Duration, connector func(ip string, httpPort int, tcpPort int) bool) (c ContainerID, err error) {
	c, ip, tcpPort, httpPort, err := SetupNSQLookupdContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up NSQLookupd container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		if connector(ip, httpPort, tcpPort) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up NSQLookupd container.")
}

// ConnectToNSQd starts a NSQ image with `/nsqd` running and passes the IP, HTTP port, and TCP port to the connector callback function.
// The url will match the ip pattern (e.g. 123.123.123.123).
func ConnectToNSQd(tries int, delay time.Duration, connector func(ip string, httpPort int, tcpPort int) bool) (c ContainerID, err error) {
	c, ip, tcpPort, httpPort, err := SetupNSQdContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up NSQd container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		if connector(ip, httpPort, tcpPort) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up NSQd container.")
}
