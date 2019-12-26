package cfclient

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type ServiceUsageEvent struct {
	GUID                string `json:"guid"`
	CreatedAt           string `json:"created_at"`
	State               string `json:"state"`
	OrgGUID             string `json:"org_guid"`
	SpaceGUID           string `json:"space_guid"`
	SpaceName           string `json:"space_name"`
	ServiceInstanceGUID string `json:"service_instance_guid"`
	ServiceInstanceName string `json:"service_instance_name"`
	ServiceInstanceType string `json:"service_instance_type"`
	ServicePlanGUID     string `json:"service_plan_guid"`
	ServicePlanName     string `json:"service_plan_name"`
	ServiceGUID         string `json:"service_guid"`
	ServiceLabel        string `json:"service_label"`
	c                   *Client
}

type ServiceUsageEventsResponse struct {
	TotalResults int                         `json:"total_results"`
	Pages        int                         `json:"total_pages"`
	NextURL      string                      `json:"next_url"`
	Resources    []ServiceUsageEventResource `json:"resources"`
}

type ServiceUsageEventResource struct {
	Meta   Meta              `json:"metadata"`
	Entity ServiceUsageEvent `json:"entity"`
}

// ListServiceUsageEventsByQuery lists all events matching the provided query.
func (c *Client) ListServiceUsageEventsByQuery(query url.Values) ([]ServiceUsageEvent, error) {
	var serviceUsageEvents []ServiceUsageEvent
	requestURL := fmt.Sprintf("/v2/service_usage_events?%s", query.Encode())
	for {
		var serviceUsageEventsResponse ServiceUsageEventsResponse
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "error requesting events")
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&serviceUsageEventsResponse); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling events")
		}
		for _, e := range serviceUsageEventsResponse.Resources {
			e.Entity.GUID = e.Meta.Guid
			e.Entity.CreatedAt = e.Meta.CreatedAt
			e.Entity.c = c
			serviceUsageEvents = append(serviceUsageEvents, e.Entity)
		}
		requestURL = serviceUsageEventsResponse.NextURL
		if requestURL == "" {
			break
		}
	}
	return serviceUsageEvents, nil
}

// ListServiceUsageEvents lists all unfiltered events.
func (c *Client) ListServiceUsageEvents() ([]ServiceUsageEvent, error) {
	return c.ListServiceUsageEventsByQuery(nil)
}
