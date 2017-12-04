package rabbithole

import (
	"encoding/json"
	"net/http"
)

type ClusterName struct {
	Name string `json:"name"`
}

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

func (c *Client) SetClusterName(cn ClusterName) (res *http.Response, err error) {
	body, err := json.Marshal(cn)
	if err != nil {
		return nil, err
	}
	req, err := newRequestWithBody(c, "PUT", "cluster-name", body)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)

	return res, nil
}
