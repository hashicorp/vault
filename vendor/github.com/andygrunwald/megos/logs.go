package megos

import (
	"net/url"
	"strconv"
)

// GetStdOutOfTask will return Stdout of a task.
//
// pid is a single mesos slave node.
// directory is the directory of a single executor,
func (c *Client) GetStdOutOfTask(pid *Pid, directory string) ([]byte, error) {
	u := c.getBaseURLForLogs("stdout", pid, directory)
	body, err := c.GetBodyOfHTTPResponse(&u)
	return body, err
}

// GetStdErrOfTask will return Stdout of a task.
//
// pid is a single mesos slave node.
// directory is the directory of a single executor,
func (c *Client) GetStdErrOfTask(pid *Pid, directory string) ([]byte, error) {
	u := c.getBaseURLForLogs("stderr", pid, directory)
	body, err := c.GetBodyOfHTTPResponse(&u)
	return body, err
}

// getBaseURLForLogs will build the URL to get Stdout or Stderr.
func (c *Client) getBaseURLForLogs(mode string, pid *Pid, directory string) url.URL {
	u := url.URL{
		// TODO How to support https cluster?
		Scheme:   "http",
		Host:     pid.Host + ":" + strconv.Itoa(pid.Port),
		Path:     "files/download.json",
		RawQuery: "path=" + directory + "/" + mode,
	}

	return u
}
