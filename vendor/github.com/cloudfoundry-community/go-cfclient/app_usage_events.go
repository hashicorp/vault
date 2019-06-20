package cfclient

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type AppUsageEvent struct {
	GUID                          string `json:"guid"`
	CreatedAt                     string `json:"created_at"`
	State                         string `json:"state"`
	PreviousState                 string `json:"previous_state"`
	MemoryInMbPerInstance         int    `json:"memory_in_mb_per_instance"`
	PreviousMemoryInMbPerInstance int    `json:"previous_memory_in_mb_per_instance"`
	InstanceCount                 int    `json:"instance_count"`
	PreviousInstanceCount         int    `json:"previous_instance_count"`
	AppGUID                       string `json:"app_guid"`
	SpaceGUID                     string `json:"space_guid"`
	SpaceName                     string `json:"space_name"`
	OrgGUID                       string `json:"org_guid"`
	BuildpackGUID                 string `json:"buildpack_guid"`
	BuildpackName                 string `json:"buildpack_name"`
	PackageState                  string `json:"package_state"`
	PreviousPackageState          string `json:"previous_package_state"`
	ParentAppGUID                 string `json:"parent_app_guid"`
	ParentAppName                 string `json:"parent_app_name"`
	ProcessType                   string `json:"process_type"`
	TaskName                      string `json:"task_name"`
	TaskGUID                      string `json:"task_guid"`
	c                             *Client
}

type AppUsageEventsResponse struct {
	TotalResults int                     `json:"total_results"`
	Pages        int                     `json:"total_pages"`
	NextURL      string                  `json:"next_url"`
	Resources    []AppUsageEventResource `json:"resources"`
}

type AppUsageEventResource struct {
	Meta   Meta          `json:"metadata"`
	Entity AppUsageEvent `json:"entity"`
}

// ListAppUsageEventsByQuery lists all events matching the provided query.
func (c *Client) ListAppUsageEventsByQuery(query url.Values) ([]AppUsageEvent, error) {
	var appUsageEvents []AppUsageEvent
	requestURL := fmt.Sprintf("/v2/app_usage_events?%s", query.Encode())
	for {
		var appUsageEventsResponse AppUsageEventsResponse
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "error requesting events")
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&appUsageEventsResponse); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling events")
		}
		for _, e := range appUsageEventsResponse.Resources {
			e.Entity.GUID = e.Meta.Guid
			e.Entity.CreatedAt = e.Meta.CreatedAt
			e.Entity.c = c
			appUsageEvents = append(appUsageEvents, e.Entity)
		}
		requestURL = appUsageEventsResponse.NextURL
		if requestURL == "" {
			break
		}
	}
	return appUsageEvents, nil
}

// ListAppUsageEvents lists all unfiltered events.
func (c *Client) ListAppUsageEvents() ([]AppUsageEvent, error) {
	return c.ListAppUsageEventsByQuery(nil)
}
