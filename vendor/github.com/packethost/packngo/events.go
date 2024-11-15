package packngo

import (
	"path"
)

const eventBasePath = "/events"

// Event struct
type Event struct {
	ID            string     `json:"id,omitempty"`
	State         string     `json:"state,omitempty"`
	Type          string     `json:"type,omitempty"`
	Body          string     `json:"body,omitempty"`
	Relationships []Href     `json:"relationships,omitempty"`
	Interpolated  string     `json:"interpolated,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	Href          string     `json:"href,omitempty"`
}

type eventsRoot struct {
	Events []Event `json:"events,omitempty"`
	Meta   meta    `json:"meta,omitempty"`
}

// EventService interface defines available event functions
type EventService interface {
	List(*ListOptions) ([]Event, *Response, error)
	Get(string, *GetOptions) (*Event, *Response, error)
}

// EventServiceOp implements EventService
type EventServiceOp struct {
	client *Client
}

// List returns all events
func (s *EventServiceOp) List(listOpt *ListOptions) ([]Event, *Response, error) {
	return listEvents(s.client, eventBasePath, listOpt)
}

// Get returns an event by ID
func (s *EventServiceOp) Get(eventID string, getOpt *GetOptions) (*Event, *Response, error) {
	if validateErr := ValidateUUID(eventID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(eventBasePath, eventID)
	return get(s.client, apiPath, getOpt)
}

// list helper function for all event functions
func listEvents(client requestDoer, endpointPath string, opts *ListOptions) (events []Event, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(eventsRoot)

		resp, err = client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		events = append(events, subset.Events...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}

}

func get(client *Client, endpointPath string, opts *GetOptions) (*Event, *Response, error) {
	event := new(Event)

	apiPathQuery := opts.WithQuery(endpointPath)

	resp, err := client.DoRequest("GET", apiPathQuery, nil, event)
	if err != nil {
		return nil, resp, err
	}

	return event, resp, err
}
