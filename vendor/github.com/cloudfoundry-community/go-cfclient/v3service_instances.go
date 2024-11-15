package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type V3ServiceInstance struct {
	Guid          string                         `json:"guid"`
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
	Name          string                         `json:"name"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Metadata      Metadata                       `json:"metadata"`
	Links         map[string]Link                `json:"links"`
}

type listV3ServiceInstancesResponse struct {
	Pagination Pagination          `json:"pagination,omitempty"`
	Resources  []V3ServiceInstance `json:"resources,omitempty"`
}

func (c *Client) ListV3ServiceInstances() ([]V3ServiceInstance, error) {
	return c.ListV3ServiceInstancesByQuery(nil)
}

func (c *Client) ListV3ServiceInstancesByQuery(query url.Values) ([]V3ServiceInstance, error) {
	var svcInstances []V3ServiceInstance
	requestURL := "/v3/service_instances"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 service instances")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing v3 service instances, response code: %d", resp.StatusCode)
		}

		var data listV3ServiceInstancesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 service instances")
		}

		svcInstances = append(svcInstances, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 service instances")
		}
	}

	return svcInstances, nil
}
