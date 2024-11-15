package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	integrationBasePath = "api/atlas/v1.0/groups/%s/integrations"
)

// IntegrationsService is an interface for interfacing with the Third-Party Integrations
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings/
type IntegrationsService interface {
	Create(context.Context, string, string, *ThirdPartyIntegration) (*ThirdPartyIntegrations, *Response, error)
	Replace(context.Context, string, string, *ThirdPartyIntegration) (*ThirdPartyIntegrations, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (*ThirdPartyIntegration, *Response, error)
	List(context.Context, string) (*ThirdPartyIntegrations, *Response, error)
}

// IntegrationsServiceOp handles communication with the third-party integrations related methods of the MongoDB Atlas API.
type IntegrationsServiceOp service

var _ IntegrationsService = &IntegrationsServiceOp{}

// ThirdPartyIntegration contains parameters for different third-party services.
type ThirdPartyIntegration struct {
	Type                     string `json:"type,omitempty"`
	LicenseKey               string `json:"licenseKey,omitempty"`
	AccountID                string `json:"accountId,omitempty"`
	WriteToken               string `json:"writeToken,omitempty"`
	ReadToken                string `json:"readToken,omitempty"`
	APIKey                   string `json:"apiKey,omitempty"`
	Region                   string `json:"region,omitempty"`
	ServiceKey               string `json:"serviceKey,omitempty"`
	APIToken                 string `json:"apiToken,omitempty"`
	TeamName                 string `json:"teamName,omitempty"`
	ChannelName              string `json:"channelName,omitempty"`
	RoutingKey               string `json:"routingKey,omitempty"`
	FlowName                 string `json:"flowName,omitempty"`
	OrgName                  string `json:"orgName,omitempty"`
	URL                      string `json:"url,omitempty"`
	Secret                   string `json:"secret,omitempty"`
	Name                     string `json:"name,omitempty"`
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	UserName                 string `json:"username,omitempty"`
	Password                 string `json:"password,omitempty"`
	ServiceDiscovery         string `json:"serviceDiscovery,omitempty"`
	Scheme                   string `json:"scheme,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
}

// ThirdPartyIntegrations contains the response from the endpoint.
type ThirdPartyIntegrations struct {
	Links      []*Link                  `json:"links"`
	Results    []*ThirdPartyIntegration `json:"results"`
	TotalCount int                      `json:"totalCount"`
}

// Create adds a new third-party integration configuration.
//
// See more: https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-create/index.html
func (s *IntegrationsServiceOp) Create(ctx context.Context, projectID, integrationType string, body *ThirdPartyIntegration) (*ThirdPartyIntegrations, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	if integrationType == "" {
		return nil, nil, NewArgError("integrationType", "must be set")
	}

	basePath := fmt.Sprintf(integrationBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, integrationType)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(ThirdPartyIntegrations)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Replace replaces the third-party integration configuration with a new configuration, or add a new configuration if there is no configuration.
//
// https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-update/
func (s *IntegrationsServiceOp) Replace(ctx context.Context, projectID, integrationType string, body *ThirdPartyIntegration) (*ThirdPartyIntegrations, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	if integrationType == "" {
		return nil, nil, NewArgError("integrationType", "must be set")
	}

	basePath := fmt.Sprintf(integrationBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, integrationType)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(ThirdPartyIntegrations)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Delete removes the third-party integration configuration
//
// https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-delete/
func (s *IntegrationsServiceOp) Delete(ctx context.Context, projectID, integrationType string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	if integrationType == "" {
		return nil, NewArgError("integrationType", "must be set")
	}

	basePath := fmt.Sprintf(integrationBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, integrationType)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	return resp, err
}

// Get retrieves a specific third-party integration configuration
//
// https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-get-one/
func (s *IntegrationsServiceOp) Get(ctx context.Context, projectID, integrationType string) (*ThirdPartyIntegration, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	if integrationType == "" {
		return nil, nil, NewArgError("integrationType", "must be set")
	}

	basePath := fmt.Sprintf(integrationBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, integrationType)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ThirdPartyIntegration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// List retrieves all third-party integration configurations.
//
// See more: https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-get-all/
func (s *IntegrationsServiceOp) List(ctx context.Context, projectID string) (*ThirdPartyIntegrations, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf(integrationBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ThirdPartyIntegrations)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}
