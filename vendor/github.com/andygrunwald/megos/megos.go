package megos

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"net/http"
)

// TODO Support new mesos version
// @link http://mesos.apache.org/documentation/latest/upgrades/

// Client manages the communication with the Mesos cluster.
type Client struct {
	sync.Mutex

	// Master is the list of Mesos master nodes in the cluster.
	Master []*url.URL
	// Leader is the PID reference to the Leader of the Cluster (of Master URLs)
	Leader *Pid
	State  *State
	Http *http.Client
}

// Pid is the process if per machine.
type Pid struct {
	// Role of a PID
	Role string
	// Host / IP of the PID
	Host string
	// Port of the PID.
	// If no Port is available the standard port (5050) will be used.
	Port int
}

// NewClient returns a new Megos / Mesos information client.
// addresses has to be the the URL`s of the single nodes of the
// Mesos cluster. It is recommended to apply all nodes in case of failures.
func NewClient(addresses []*url.URL, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := &Client{
		Master: addresses,
		Http: httpClient,
	}

	return client
}

// DetermineLeader will return the leader of several master nodes of
// the Mesos cluster. Only one leader is chosen per time.
// This leader will be returned.
func (c *Client) DetermineLeader() (*Pid, error) {
	state, err := c.GetStateFromCluster()
	if err != nil {
		return nil, err
	}

	pid, err := c.ParsePidInformation(state.Leader)
	if err != nil {
		return nil, err
	}

	c.Lock()
	c.Leader = pid
	c.Unlock()

	return c.Leader, nil
}

// ParsePidInformation will split up a single PID of format node@ip:port
// into a Pid structure to access single parts of the PID on its own.
//
// Example pid: master@10.1.1.12:5050
func (c *Client) ParsePidInformation(pid string) (*Pid, error) {
	firstPart := strings.Split(pid, "@")
	if len(firstPart) != 2 {
		return nil, errors.New("Invalid master pid.")
	}
	secondPart := strings.Split(firstPart[1], ":")

	var port int
	var err error
	if len(secondPart) == 2 {
		port, err = strconv.Atoi(secondPart[1])
	}

	// If we got an error during conversion or no port is available, set the default port
	if err != nil || len(secondPart) < 2 {
		port = 5050
	}

	return &Pid{
		Role: firstPart[0],
		Host: secondPart[0],
		Port: port,
	}, nil
}

// String implements the Stringer interface for PID.
// It can be named as the opposite of Client.ParsePidInformation,
// because it transfers a single Pid structure into its original form
// with format node@ip:port.
func (p *Pid) String() string {
	s := p.Role + "@" + p.Host + ":" + strconv.Itoa(p.Port)
	return s
}
