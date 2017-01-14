package dockertest

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// AwaitReachable tries to make a TCP connection to addr regularly.
// It returns an error if it's unable to make a connection before maxWait.
func AwaitReachable(addr string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}

// ContainerID represents a container and offers methods like Kill or IP.
type ContainerID string

// IP retrieves the container's IP address.
func (c ContainerID) IP() (string, error) {
	return IP(string(c))
}

// Kill runs "docker kill" on the container.
func (c ContainerID) Kill() error {
	return KillContainer(string(c))
}

// Start runs "docker start" on the container.
func (c ContainerID) Start() error {
	return StartContainer(string(c))
}

// Stop runs "docker stop" on the container.
func (c ContainerID) Stop() error {
	return StopContainer(string(c))
}

// Remove runs "docker rm" on the container
func (c ContainerID) Remove() error {
	if Debug || c == "nil" {
		return nil
	}
	return runDockerCommand("docker", "rm", "-v", string(c)).Run()
}

// KillRemove calls Kill on the container, and then Remove if there was
// no error.
func (c ContainerID) KillRemove() error {
	if err := c.Kill(); err != nil {
		return err
	}
	return c.Remove()
}

// lookup retrieves the ip address of the container, and tries to reach
// before timeout the tcp address at this ip and given port.
func (c ContainerID) lookup(ports []int, timeout time.Duration) (ip string, err error) {
	if DockerMachineAvailable {
		var out []byte
		out, err = exec.Command("docker-machine", "ip", DockerMachineName).Output()
		ip = strings.TrimSpace(string(out))
	} else if BindDockerToLocalhost != "" {
		ip = "127.0.0.1"
	} else {
		ip, err = c.IP()
	}
	if err != nil {
		err = fmt.Errorf("error getting IP: %v", err)
		return
	}
	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", ip, port)
		err = AwaitReachable(addr, timeout)
		if err != nil {
			return
		}
	}
	return
}
