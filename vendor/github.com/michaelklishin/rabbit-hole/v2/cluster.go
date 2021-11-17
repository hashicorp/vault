package rabbithole

import (
	"encoding/json"
	"net/http"
)

// ClusterName represents a RabbitMQ cluster name (identifier).
type ClusterName struct {
	Name string `json:"name"`
}

// GetClusterName returns current cluster name.
func (c *Client) GetClusterName() (rec *ClusterName, err error) {
	req, err := newGETRequest(c, "cluster-name/")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// SetClusterName sets cluster name.
func (c *Client) SetClusterName(cn ClusterName) (res *http.Response, err error) {
	body, err := json.Marshal(cn)
	if err != nil {
		return nil, err
	}
	req, err := newRequestWithBody(c, "PUT", "cluster-name", body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
