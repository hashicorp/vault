package megos

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// GetStateFromCluster will return the current state of one of the cluster nodes.
func (c *Client) GetStateFromCluster() (*State, error) {
	resp, err := c.GetHTTPResponseFromCluster(c.GetURLForStateFile)
	return c.parseStateResponse(resp, err)
}

// GetStateFromLeader will return the current state of the leader node of the cluster.
func (c *Client) GetStateFromLeader() (*State, error) {
	resp, err := c.GetHTTPResponseFromLeader(c.GetURLForStateFilePid)
	return c.parseStateResponse(resp, err)
}

// GetStateFromPid will return the current state of the process id per machine (PID).
func (c *Client) GetStateFromPid(pid *Pid) (*State, error) {
	u := c.GetURLForStateFilePid(*pid)
	resp, err := c.GetHTTPResponse(&u)
	return c.parseStateResponse(resp, err)
}

// GetStateSummaryFromCluster will return the current state summary of one of the cluster nodes.
func (c *Client) GetStateSummaryFromCluster() (*State, error) {
	resp, err := c.GetHTTPResponseFromCluster(c.GetURLForStateSummaryFile)
	return c.parseStateResponse(resp, err)
}

// GetStateSummaryFromLeader will return the current state summary of the leader node of the cluster.
func (c *Client) GetStateSummaryFromLeader() (*State, error) {
	resp, err := c.GetHTTPResponseFromLeader(c.GetURLForStateSummaryFilePid)
	return c.parseStateResponse(resp, err)
}

// GetStateSummaryFromPid will return the current state summary of the process id per machine (PID).
func (c *Client) GetStateSummaryFromPid(pid *Pid) (*State, error) {
	u := c.GetURLForStateSummaryFilePid(*pid)
	resp, err := c.GetHTTPResponse(&u)
	return c.parseStateResponse(resp, err)
}

// GetSlavesFromCluster will return the current slaves of one of the cluster nodes.
func (c *Client) GetSlavesFromCluster() (*State, error) {
	resp, err := c.GetHTTPResponseFromCluster(c.GetURLForSlavesFile)
	return c.parseStateResponse(resp, err)
}

// GetSlavesFromLeader will return the current slaves of the leader node of the cluster.
func (c *Client) GetSlavesFromLeader() (*State, error) {
	resp, err := c.GetHTTPResponseFromLeader(c.GetURLForSlavesFilePid)
	return c.parseStateResponse(resp, err)
}

// GetSlavesFromPid will return the current slaves of the process id per machine (PID).
func (c *Client) GetSlavesFromPid(pid *Pid) (*State, error) {
	u := c.GetURLForSlavesFilePid(*pid)
	resp, err := c.GetHTTPResponse(&u)
	return c.parseStateResponse(resp, err)
}

// parseStateResponse will transform a http.Response into a State object
func (c *Client) parseStateResponse(resp *http.Response, err error) (*State, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var state State
	err = json.Unmarshal(body, &state)
	if err != nil {
		return nil, err
	}

	c.Lock()
	c.State = &state
	c.Unlock()

	return c.State, nil
}

// GetURLForStateFile will return the URL for the state file of a node
func (c *Client) GetURLForStateFile(instance url.URL) url.URL {
	instance.Path = "master/state"
	return instance
}

// GetURLForStateFilePid will return the URL for the state file of a node
// based on a PID
func (c *Client) GetURLForStateFilePid(pid Pid) url.URL {
	return c.getURLForFilePid(pid, "state")
}

// GetURLForStateSummaryFile will return the URL for the state-summary file of a node
func (c *Client) GetURLForStateSummaryFile(instance url.URL) url.URL {
	instance.Path = "master/state-summary"
	return instance
}

// GetURLForStateSummaryFilePid will return the URL for the state-summary file of a node
// based on a PID
func (c *Client) GetURLForStateSummaryFilePid(pid Pid) url.URL {
	return c.getURLForFilePid(pid, "state-summary")
}

// GetURLForSlavesFile will return the URL for the slaves file of a node
func (c *Client) GetURLForSlavesFile(instance url.URL) url.URL {
	instance.Path = "master/slaves"
	return instance
}

// GetURLForSlavesFilePid will return the URL for the slaves file of a node
// based on a PID
func (c *Client) GetURLForSlavesFilePid(pid Pid) url.URL {
	return c.getURLForFilePid(pid, "slaves")
}

func (c *Client) getURLForFilePid(pid Pid, filename string) url.URL {
	host := pid.Host + ":" + strconv.Itoa(pid.Port)

	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   pid.Role + "/" + filename,
	}

	return u
}
