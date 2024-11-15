package cfclient

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// EventsResponse is a type that wraps a collection of event resources.
type EventsResponse struct {
	TotalResults int             `json:"total_results"`
	Pages        int             `json:"total_pages"`
	NextURL      string          `json:"next_url"`
	Resources    []EventResource `json:"resources"`
}

// EventResource is a type that contains metadata and the entity for an event.
type EventResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Event `json:"entity"`
}

// Event is a type that contains event data.
type Event struct {
	GUID             string                 `json:"guid"`
	Type             string                 `json:"type"`
	CreatedAt        string                 `json:"created_at"`
	Actor            string                 `json:"actor"`
	ActorType        string                 `json:"actor_type"`
	ActorName        string                 `json:"actor_name"`
	ActorUsername    string                 `json:"actor_username"`
	Actee            string                 `json:"actee"`
	ActeeType        string                 `json:"actee_type"`
	ActeeName        string                 `json:"actee_name"`
	OrganizationGUID string                 `json:"organization_guid"`
	SpaceGUID        string                 `json:"space_guid"`
	Metadata         map[string]interface{} `json:"metadata"`
	c                *Client
}

// ListEventsByQuery lists all events matching the provided query.
func (c *Client) ListEventsByQuery(query url.Values) ([]Event, error) {
	var events []Event
	requestURL := fmt.Sprintf("/v2/events?%s", query.Encode())
	for {
		var eventResp EventsResponse
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "error requesting events")
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&eventResp); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling events")
		}
		for _, e := range eventResp.Resources {
			e.Entity.GUID = e.Meta.Guid
			e.Entity.CreatedAt = e.Meta.CreatedAt
			e.Entity.c = c
			events = append(events, e.Entity)
		}
		requestURL = eventResp.NextURL
		if requestURL == "" {
			break
		}
	}
	return events, nil
}

// ListEvents lists all unfiltered events.
func (c *Client) ListEvents() ([]Event, error) {
	return c.ListEventsByQuery(nil)
}

// TotalEventsByQuery returns the number of events matching the provided query.
func (c *Client) TotalEventsByQuery(query url.Values) (int, error) {
	r := c.NewRequest("GET", fmt.Sprintf("/v2/events?%s", query.Encode()))
	resp, err := c.DoRequest(r)
	if err != nil {
		return 0, errors.Wrap(err, "error requesting events")
	}
	defer resp.Body.Close()
	var apiResp EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, errors.Wrap(err, "error unmarshaling events")
	}
	return apiResp.TotalResults, nil
}

// TotalEvents returns the number of unfiltered events.
func (c *Client) TotalEvents() (int, error) {
	return c.TotalEventsByQuery(nil)
}
