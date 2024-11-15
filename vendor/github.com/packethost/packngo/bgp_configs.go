package packngo

import (
	"path"
)

var (
	bgpConfigPostBasePath = "/bgp-configs"
	bgpConfigGetBasePath  = "/bgp-config"
)

// BGPConfigService interface defines available BGP config methods
type BGPConfigService interface {
	Get(projectID string, getOpt *GetOptions) (*BGPConfig, *Response, error)
	Create(projectID string, request CreateBGPConfigRequest) (*Response, error)
	// Delete(configID string) (resp *Response, err error) TODO: Not in Equinix Metal API
}

// BGPConfigServiceOp implements BgpConfigService
type BGPConfigServiceOp struct {
	client *Client
}

// CreateBGPConfigRequest struct
type CreateBGPConfigRequest struct {
	DeploymentType string `json:"deployment_type,omitempty"`
	Asn            int    `json:"asn,omitempty"`
	Md5            string `json:"md5,omitempty"`
	UseCase        string `json:"use_case,omitempty"`
}

// BGPConfig represents an Equinix Metal BGP Config
type BGPConfig struct {
	ID             string       `json:"id,omitempty"`
	Status         string       `json:"status,omitempty"`
	DeploymentType string       `json:"deployment_type,omitempty"`
	Asn            int          `json:"asn,omitempty"`
	RouteObject    string       `json:"route_object,omitempty"`
	Md5            string       `json:"md5,omitempty"`
	MaxPrefix      int          `json:"max_prefix,omitempty"`
	Project        Project      `json:"project,omitempty"`
	CreatedAt      Timestamp    `json:"created_at,omitempty"`
	RequestedAt    Timestamp    `json:"requested_at,omitempty"`
	Sessions       []BGPSession `json:"sessions,omitempty"`
	Href           string       `json:"href,omitempty"`
}

// Create function
func (s *BGPConfigServiceOp) Create(projectID string, request CreateBGPConfigRequest) (*Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, bgpConfigPostBasePath)

	resp, err := s.client.DoRequest("POST", apiPath, request, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Get function
func (s *BGPConfigServiceOp) Get(projectID string, opts *GetOptions) (bgpConfig *BGPConfig, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, bgpConfigGetBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	subset := new(BGPConfig)

	resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	return subset, resp, err
}

// Delete function TODO: this is not implemented in the Equinix Metal API
// func (s *BGPConfigServiceOp) Delete(configID string) (resp *Response, err error) {
// if validateErr := ValidateUUID(configID); validateErr != nil {
//  return nil, validateErr
// }
// 	apiPath := fmt.Sprintf("%ss/%s", bgpConfigBasePath, configID)

// 	resp, err = s.client.DoRequest("DELETE", apiPath, nil, nil)
// 	if err != nil {
// 		return resp, err
// 	}

// 	return resp, err
// }
