// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	eventsPathProjects     = "api/atlas/v1.0/groups/%s/events"
	eventsPathOrganization = "api/atlas/v1.0/orgs/%s/events"
)

// EventsService is an interface for interfacing with the Events
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/events/
type EventsService interface {
	ListOrganizationEvents(context.Context, string, *EventListOptions) (*EventResponse, *Response, error)
	GetOrganizationEvent(context.Context, string, string) (*Event, *Response, error)
	ListProjectEvents(context.Context, string, *EventListOptions) (*EventResponse, *Response, error)
	GetProjectEvent(context.Context, string, string) (*Event, *Response, error)
}

// EventsServiceOp handles communication with the Event related methods
// of the MongoDB Atlas API.
type EventsServiceOp service

var _ EventsService = &EventsServiceOp{}

// Event represents an event of the MongoDB Atlas API.
type Event struct {
	AlertID         string        `json:"alertId"`
	AlertConfigID   string        `json:"alertConfigId"`
	APIKeyID        string        `json:"apiKeyId,omitempty"`
	Collection      string        `json:"collection,omitempty"`
	Created         string        `json:"created"`
	CurrentValue    *CurrentValue `json:"currentValue,omitempty"`
	Database        string        `json:"database,omitempty"`
	EventTypeName   string        `json:"eventTypeName"`
	GroupID         string        `json:"groupId,omitempty"`
	Hostname        string        `json:"hostname"`
	ID              string        `json:"id"`
	InvoiceID       string        `json:"invoiceId,omitempty"`
	IsGlobalAdmin   bool          `json:"isGlobalAdmin,omitempty"`
	Links           []*Link       `json:"links"`
	MetricName      string        `json:"metricName,omitempty"`
	OpType          string        `json:"opType,omitempty"`
	OrgID           string        `json:"orgId,omitempty"`
	PaymentID       string        `json:"paymentId,omitempty"`
	Port            int           `json:"Port,omitempty"`
	PublicKey       string        `json:"publicKey,omitempty"`
	RemoteAddress   string        `json:"remoteAddress,omitempty"`
	ReplicaSetName  string        `json:"replicaSetName,omitempty"`
	ShardName       string        `json:"shardName,omitempty"`
	TargetPublicKey string        `json:"targetPublicKey,omitempty"`
	TargetUsername  string        `json:"targetUsername,omitempty"`
	TeamID          string        `json:"teamId,omitempty"`
	UserID          string        `json:"userId,omitempty"`
	Username        string        `json:"username,omitempty"`
	WhitelistEntry  string        `json:"whitelistEntry,omitempty"`
}

// EventResponse is the response from the EventsService.List.
type EventResponse struct {
	Links      []*Link  `json:"links,omitempty"`
	Results    []*Event `json:"results,omitempty"`
	TotalCount int      `json:"totalCount,omitempty"`
}

// EventListOptions specifies the optional parameters to the Event List methods.
type EventListOptions struct {
	ListOptions
	EventType []string `url:"eventType,omitempty"`
	MinDate   string   `url:"minDate,omitempty"`
	MaxDate   string   `url:"maxDate,omitempty"`
}

// ListOrganizationEvents lists all events in the organization associated to {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/events-orgs-get-all/
func (s *EventsServiceOp) ListOrganizationEvents(ctx context.Context, orgID string, listOptions *EventListOptions) (*EventResponse, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	path := fmt.Sprintf(eventsPathOrganization, orgID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(EventResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// GetOrganizationEvent gets the alert specified to {EVENT-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-one-event-for-org/
func (s *EventsServiceOp) GetOrganizationEvent(ctx context.Context, orgID, eventID string) (*Event, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if eventID == "" {
		return nil, nil, NewArgError("eventID", "must be set")
	}
	basePath := fmt.Sprintf(eventsPathOrganization, orgID)
	path := fmt.Sprintf("%s/%s", basePath, eventID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Event)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// ListProjectEvents lists all events in the project associated to {PROJECT-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-all-events-for-project/
func (s *EventsServiceOp) ListProjectEvents(ctx context.Context, groupID string, listOptions *EventListOptions) (*EventResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(eventsPathProjects, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(EventResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// GetProjectEvent gets the alert specified to {EVENT-ID} from the project associated to {PROJECT-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-one-event-for-project/
func (s *EventsServiceOp) GetProjectEvent(ctx context.Context, groupID, eventID string) (*Event, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if eventID == "" {
		return nil, nil, NewArgError("eventID", "must be set")
	}
	basePath := fmt.Sprintf(eventsPathProjects, groupID)
	path := fmt.Sprintf("%s/%s", basePath, eventID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Event)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
