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

const alertPath = "api/atlas/v1.0/groups/%s/alerts"

// AlertsService is an interface for interfacing with the Alerts
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts/
type AlertsService interface {
	List(context.Context, string, *AlertsListOptions) (*AlertsResponse, *Response, error)
	Get(context.Context, string, string) (*Alert, *Response, error)
	Acknowledge(context.Context, string, string, *AcknowledgeRequest) (*Alert, *Response, error)
}

// AlertsServiceOp provides an implementation of AlertsService.
type AlertsServiceOp service

var _ AlertsService = &AlertsServiceOp{}

// Alert represents MongoDB Alert.
type Alert struct {
	ID                     string           `json:"id,omitempty"`                     // Unique identifier.
	GroupID                string           `json:"groupId,omitempty"`                // Unique identifier of the project that owns this alert configuration.
	AlertConfigID          string           `json:"alertConfigId,omitempty"`          // ID of the alert configuration that triggered this alert.
	EventTypeName          string           `json:"eventTypeName,omitempty"`          // The type of event that will trigger an alert.
	Created                string           `json:"created,omitempty"`                // Timestamp in ISO 8601 date and time format in UTC when this alert was opened.
	Updated                string           `json:"updated,omitempty"`                // Timestamp in ISO 8601 date and time format in UTC when this alert was last updated.
	Enabled                *bool            `json:"enabled,omitempty"`                // If omitted, the configuration is disabled.
	Resolved               string           `json:"resolved,omitempty"`               // When the alert was closed. Only present if the status is CLOSED.
	Status                 string           `json:"status,omitempty"`                 // The current state of the alert. Possible values are: TRACKING, OPEN, CLOSED, CANCELLED
	LastNotified           string           `json:"lastNotified,omitempty"`           // When the last notification was sent for this alert. Only present if notifications have been sent.
	AcknowledgedUntil      string           `json:"acknowledgedUntil,omitempty"`      // The date through which the alert has been acknowledged. Will not be present if the alert has never been acknowledged.
	AcknowledgementComment string           `json:"acknowledgementComment,omitempty"` // The comment left by the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	AcknowledgingUsername  string           `json:"acknowledgingUsername,omitempty"`  // The username of the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	HostnameAndPort        string           `json:"hostnameAndPort,omitempty"`        // The hostname and port of each host to which the alert applies. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.
	MetricName             string           `json:"metricName,omitempty"`             // The name of the measurement whose value went outside the threshold. Only present if eventTypeName is set to OUTSIDE_METRIC_THRESHOLD.
	CurrentValue           *CurrentValue    `json:"currentValue,omitempty"`           // CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.
	ReplicaSetName         string           `json:"replicaSetName,omitempty"`         // Name of the replica set. Only present for alerts of type HOST, HOST_METRIC, BACKUP, and REPLICA_SET.
	ClusterName            string           `json:"clusterName,omitempty"`            // The name the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.
	Matchers               []Matcher        `json:"matchers,omitempty"`               // You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.
	MetricThreshold        *MetricThreshold `json:"metricThreshold,omitempty"`        // MetricThreshold  causes an alert to be triggered.
	Notifications          []Notification   `json:"notifications,omitempty"`          // Notifications are sending when an alert condition is detected.
}

// AcknowledgeRequest contains the request Body Parameters.
type AcknowledgeRequest struct {
	AcknowledgedUntil      *string `json:"acknowledgedUntil,omitempty"`      // The date through which the alert has been acknowledged. Will not be present if the alert has never been acknowledged.
	AcknowledgementComment string  `json:"acknowledgementComment,omitempty"` // The comment left by the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
}

// AlertsListOptions contains the list of options for Alerts.
type AlertsListOptions struct {
	Status string `url:"status,omitempty"`
	ListOptions
}

// AlertsResponse is the response from the AlertService.List.
type AlertsResponse struct {
	Links      []*Link `json:"links"`
	Results    []Alert `json:"results"`
	TotalCount int     `json:"totalCount"`
}

// Get gets the alert specified to {ALERT-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-get-alert/
func (s *AlertsServiceOp) Get(ctx context.Context, groupID, alertID string) (*Alert, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if alertID == "" {
		return nil, nil, NewArgError("alertID", "must be set")
	}

	basePath := fmt.Sprintf(alertPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Alert)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// List gets all alert for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-get-all-alerts/
func (s *AlertsServiceOp) List(ctx context.Context, groupID string, listOptions *AlertsListOptions) (*AlertsResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(alertPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Acknowledge allows to acknowledge an alert.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-acknowledge-alert/
func (s *AlertsServiceOp) Acknowledge(ctx context.Context, groupID, alertID string, params *AcknowledgeRequest) (*Alert, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if alertID == "" {
		return nil, nil, NewArgError("alertID", "must be set")
	}

	if params == nil {
		return nil, nil, NewArgError("params", "must be set")
	}

	basePath := fmt.Sprintf(alertPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, params)

	if err != nil {
		return nil, nil, err
	}

	root := new(Alert)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
