package megos

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GetHTTPResponseFromCluster will return a http.Response from one of the Mesos master nodes.
// In a cluster the master nodes can be online or offline.
// With GetHTTPResponseFromCluster you will receive a response from one of the nodes.
func (c *Client) GetHTTPResponseFromCluster(f func(url.URL) url.URL) (*http.Response, error) {
	for _, instance := range c.Master {
		u := f(*instance)
		resp, err := c.GetHTTPResponse(&u)

		// If there is no error, we hit an instance / master that is online
		if err == nil {
			return resp, nil
		}
	}

	return nil, errors.New("No master online.")
}

// GetHTTPResponseFromLeader will return a http.Response from the determined leader
// of the master nodes.
func (c *Client) GetHTTPResponseFromLeader(f func(Pid) url.URL) (*http.Response, error) {
	if c.Leader == nil {
		return nil, errors.New("No leader set.")
	}
	u := f(*c.Leader)
	return c.GetHTTPResponse(&u)
}

// GetHTTPResponse will return a http.Response from a URL
func (c *Client) GetHTTPResponse(u *url.URL) (*http.Response, error) {
	resp, err := c.Http.Get(u.String())

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBodyOfHTTPResponse will return the request body of the requested url u.
func (c *Client) GetBodyOfHTTPResponse(u *url.URL) ([]byte, error) {
	resp, err := c.GetHTTPResponse(u)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
