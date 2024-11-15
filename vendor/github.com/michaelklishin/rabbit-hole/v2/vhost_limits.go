package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// VhostLimitsValues are properties used to modify virtual host limits (max-connections, max-queues)
type VhostLimitsValues map[string]int

// VhostLimits are properties used to delete virtual host limits (max-connections, max-queues)
type VhostLimits []string

// VhostLimitsInfo holds information about the current virtual host limits
type VhostLimitsInfo struct {
	Vhost string            `json:"vhost"`
	Value VhostLimitsValues `json:"value"`
}

// GetAllVhostLimits gets all virtual hosts limits.
func (c *Client) GetAllVhostLimits() (rec []VhostLimitsInfo, err error) {
	req, err := newGETRequest(c, "vhost-limits")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// GetVhostLimits gets a virtual host limits.
func (c *Client) GetVhostLimits(vhostname string) (rec []VhostLimitsInfo, err error) {
	req, err := newGETRequest(c, "vhost-limits/"+url.PathEscape(vhostname))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// PutVhostLimits puts limits of a virtual host.
func (c *Client) PutVhostLimits(vhostname string, limits VhostLimitsValues) (res *http.Response, err error) {
	for limitName, limitValue := range limits {
		body, err := json.Marshal(struct {
			Value int `json:"value"`
		}{Value: limitValue})
		if err != nil {
			return nil, err
		}

		req, err := newRequestWithBody(c, "PUT", "vhost-limits/"+url.PathEscape(vhostname)+"/"+limitName, body)
		if err != nil {
			return nil, err
		}

		if res, err = executeRequest(c, req); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// DeleteVhostLimits deletes limits of a virtual host.
func (c *Client) DeleteVhostLimits(vhostname string, limits VhostLimits) (res *http.Response, err error) {
	for _, limit := range limits {
		req, err := newRequestWithBody(c, "DELETE", "vhost-limits/"+url.PathEscape(vhostname)+"/"+limit, nil)
		if err != nil {
			return nil, err
		}

		if res, err = executeRequest(c, req); err != nil {
			return nil, err
		}
	}

	return res, nil
}
