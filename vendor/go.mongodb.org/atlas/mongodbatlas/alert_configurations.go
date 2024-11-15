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
	alertConfigurationPath         = "api/atlas/v1.0/groups/%s/alertConfigs"
	alertConfigurationMatchersPath = "api/atlas/v1.0/alertConfigs/matchers/fieldNames"
)

// AlertConfigurationsService provides access to the alert configuration related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations
type AlertConfigurationsService interface {
	Create(context.Context, string, *AlertConfiguration) (*AlertConfiguration, *Response, error)
	EnableAnAlertConfig(context.Context, string, string, *bool) (*AlertConfiguration, *Response, error)
	GetAnAlertConfig(context.Context, string, string) (*AlertConfiguration, *Response, error)
	GetOpenAlertsConfig(context.Context, string, string) ([]AlertConfiguration, *Response, error)
	List(context.Context, string, *ListOptions) ([]AlertConfiguration, *Response, error)
	ListMatcherFields(ctx context.Context) ([]string, *Response, error)
	Update(context.Context, string, string, *AlertConfiguration) (*AlertConfiguration, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// AlertConfigurationsServiceOp handles communication with the AlertConfiguration related methods
// of the MongoDB Atlas API.
type AlertConfigurationsServiceOp service

var _ AlertConfigurationsService = &AlertConfigurationsServiceOp{}

// AlertConfiguration represents MongoDB Alert Configuration.
type AlertConfiguration struct {
	ID                     string           `json:"id,omitempty"`                     // Unique identifier.
	GroupID                string           `json:"groupId,omitempty"`                // Unique identifier of the project that owns this alert configuration.
	AlertConfigID          string           `json:"alertConfigId,omitempty"`          // ID of the alert configuration that triggered this alert.
	EventTypeName          string           `json:"eventTypeName,omitempty"`          // The type of event that will trigger an alert.
	Created                string           `json:"created,omitempty"`                // Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
	Status                 string           `json:"status,omitempty"`                 // The current state of the alert. Possible values are: TRACKING, OPEN, CLOSED, CANCELLED
	AcknowledgedUntil      string           `json:"acknowledgedUntil,omitempty"`      // The date through which the alert has been acknowledged. Will not be present if the alert has never been acknowledged.
	AcknowledgementComment string           `json:"acknowledgementComment,omitempty"` // The comment left by the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	AcknowledgingUsername  string           `json:"acknowledgingUsername,omitempty"`  // The username of the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	Updated                string           `json:"updated,omitempty"`                // Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.
	Resolved               string           `json:"resolved,omitempty"`               // When the alert was closed. Only present if the status is CLOSED.
	LastNotified           string           `json:"lastNotified,omitempty"`           // When the last notification was sent for this alert. Only present if notifications have been sent.
	HostnameAndPort        string           `json:"hostnameAndPort,omitempty"`        // The hostname and port of each host to which the alert applies. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.
	HostID                 string           `json:"hostId,omitempty"`                 // ID of the host to which the metric pertains. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.
	ReplicaSetName         string           `json:"replicaSetName,omitempty"`         // Name of the replica set. Only present for alerts of type HOST, HOST_METRIC, BACKUP, and REPLICA_SET.
	MetricName             string           `json:"metricName,omitempty"`             // The name of the measurement whose value went outside the threshold. Only present if eventTypeName is set to OUTSIDE_METRIC_THRESHOLD.
	Enabled                *bool            `json:"enabled,omitempty"`                // If omitted, the configuration is disabled.
	ClusterID              string           `json:"clusterId,omitempty"`              // The ID of the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.
	ClusterName            string           `json:"clusterName,omitempty"`            // The name the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.
	SourceTypeName         string           `json:"sourceTypeName,omitempty"`         // For alerts of the type BACKUP, the type of server being backed up.
	CurrentValue           *CurrentValue    `json:"currentValue,omitempty"`           // CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.
	Matchers               []Matcher        `json:"matchers,omitempty"`               // You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.
	MetricThreshold        *MetricThreshold `json:"metricThreshold,omitempty"`        // MetricThreshold  causes an alert to be triggered.
	Threshold              *Threshold       `json:"threshold,omitempty"`              // Threshold  causes an alert to be triggered.
	Notifications          []Notification   `json:"notifications,omitempty"`          // Notifications are sending when an alert condition is detected.
}

// Matcher represents the Rules to apply when matching an object against this alert configuration.
// Only entities that match all these rules are checked for an alert condition.
type Matcher struct {
	FieldName string `json:"fieldName,omitempty"` // Name of the field in the target object to match on.
	Operator  string `json:"operator,omitempty"`  // The operator to test the field’s value.
	Value     string `json:"value,omitempty"`     // Value to test with the specified operator.
}

// MetricThreshold  causes an alert to be triggered. Required if "eventTypeName" : "OUTSIDE_METRIC_THRESHOLD".
type MetricThreshold struct {
	MetricName string  `json:"metricName,omitempty"` // Name of the metric to check.
	Operator   string  `json:"operator,omitempty"`   // Operator to apply when checking the current metric value against the threshold value.
	Threshold  float64 `json:"threshold"`            // Threshold value outside of which an alert will be triggered.
	Units      string  `json:"units,omitempty"`      // The units for the threshold value.
	Mode       string  `json:"mode,omitempty"`       // This must be set to AVERAGE. Atlas computes the current metric value as an average.
}

// Threshold that triggers an alert. Don’t include if "eventTypeName" : "OUTSIDE_METRIC_THRESHOLD".
type Threshold struct {
	Operator  string  `json:"operator,omitempty"`  // Operator to apply when checking the current metric value against the threshold value. it accepts the following values: GREATER_THAN, LESS_THAN
	Units     string  `json:"units,omitempty"`     // The units for the threshold value.
	Threshold float64 `json:"threshold,omitempty"` // Threshold value outside of which an alert will be triggered.
}

// Notification sends when an alert condition is detected.
type Notification struct {
	APIToken                 string   `json:"apiToken,omitempty"`                 // Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	ChannelName              string   `json:"channelName,omitempty"`              // Slack channel name. Populated for the SLACK notifications type.
	DatadogAPIKey            string   `json:"datadogApiKey,omitempty"`            // Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.
	DatadogRegion            string   `json:"datadogRegion,omitempty"`            // Region that indicates which API URL to use
	DelayMin                 *int     `json:"delayMin,omitempty"`                 // Number of minutes to wait after an alert condition is detected before sending out the first notification.
	EmailAddress             string   `json:"emailAddress,omitempty"`             // Email address to which alert notifications are sent. Populated for the EMAIL notifications type.
	EmailEnabled             *bool    `json:"emailEnabled,omitempty"`             // Flag indicating if email notifications should be sent. Populated for ORG, GROUP, and USER notifications types.
	FlowdockAPIToken         string   `json:"flowdockApiToken,omitempty"`         // The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	FlowName                 string   `json:"flowName,omitempty"`                 // Flowdock flow namse in lower-case letters.
	IntervalMin              int      `json:"intervalMin,omitempty"`              // Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.
	MobileNumber             string   `json:"mobileNumber,omitempty"`             // Mobile number to which alert notifications are sent. Populated for the SMS notifications type.
	OpsGenieAPIKey           string   `json:"opsGenieApiKey,omitempty"`           // Opsgenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	OpsGenieRegion           string   `json:"opsGenieRegion,omitempty"`           // Region that indicates which API URL to use.
	OrgName                  string   `json:"orgName,omitempty"`                  // Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.
	ServiceKey               string   `json:"serviceKey,omitempty"`               // PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	SMSEnabled               *bool    `json:"smsEnabled,omitempty"`               // Flag indicating if text message notifications should be sent. Populated for ORG, GROUP, and USER notifications types.
	TeamID                   string   `json:"teamId,omitempty"`                   // Unique identifier of a team.
	TeamName                 string   `json:"teamName,omitempty"`                 // Label for the team that receives this notification.
	NotifierID               string   `json:"notifierId,omitempty"`               // The notifierId is a system-generated unique identifier assigned to each notification method.
	TypeName                 string   `json:"typeName,omitempty"`                 // Type of alert notification.
	Username                 string   `json:"username,omitempty"`                 // Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Populated for the USER notifications type.
	VictorOpsAPIKey          string   `json:"victorOpsApiKey,omitempty"`          // VictorOps API key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	VictorOpsRoutingKey      string   `json:"victorOpsRoutingKey,omitempty"`      // VictorOps routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	Roles                    []string `json:"roles,omitempty"`                    // The following roles grant privileges within a project.
	MicrosoftTeamsWebhookURL string   `json:"microsoftTeamsWebhookUrl,omitempty"` // Microsoft Teams Wewbhook URL
	WebhookSecret            string   `json:"webhookSecret,omitempty"`            // Webhook Secret
	WebhookURL               string   `json:"webhookUrl,omitempty"`               // Webhook URL
}

// AlertConfigurationsResponse is the response from the AlertConfigurationsService.List.
type AlertConfigurationsResponse struct {
	Links      []*Link              `json:"links"`
	Results    []AlertConfiguration `json:"results"`
	TotalCount int                  `json:"totalCount"`
}

// CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.
type CurrentValue struct {
	Number *float64 `json:"number,omitempty"` // The value of the metric.
	Units  string   `json:"units,omitempty"`  // The units for the value. Depends on the type of metric.
}

// Create creates an alert configuration for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-create-config/
func (s *AlertConfigurationsServiceOp) Create(ctx context.Context, groupID string, createReq *AlertConfiguration) (*AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}

	path := fmt.Sprintf(alertConfigurationPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// EnableAnAlertConfig Enables/disables the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-enable-disable-config/
func (s *AlertConfigurationsServiceOp) EnableAnAlertConfig(ctx context.Context, groupID, alertConfigID string, enabled *bool) (*AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, AlertConfiguration{Enabled: enabled})
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetAnAlertConfig gets the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-config/
func (s *AlertConfigurationsServiceOp) GetAnAlertConfig(ctx context.Context, groupID, alertConfigID string) (*AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOpenAlertsConfig gets all open alerts for the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-open-alerts/
func (s *AlertConfigurationsServiceOp) GetOpenAlertsConfig(ctx context.Context, groupID, alertConfigID string) ([]AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s/alerts", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfigurationsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}
	return root.Results, resp, err
}

// List gets all alert configurations for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-all-configs/
func (s *AlertConfigurationsServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(alertConfigurationPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfigurationsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Update the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-update-config/
func (s *AlertConfigurationsServiceOp) Update(ctx context.Context, groupID, alertConfigID string, updateReq *AlertConfiguration) (*AlertConfiguration, *Response, error) {
	if updateReq == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, err
	}

	root := new(AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-delete-config/
func (s *AlertConfigurationsServiceOp) Delete(ctx context.Context, groupID, alertConfigID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ListMatcherFields gets all field names that the matchers.fieldName parameter accepts when you create or update an Alert Configuration.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-matchers-field-names/
func (s *AlertConfigurationsServiceOp) ListMatcherFields(ctx context.Context) ([]string, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, alertConfigurationMatchersPath, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []string
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
