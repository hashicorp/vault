package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const whitelistAPIKeysPath = "orgs/%s/apiKeys/%s/whitelist"

// WhitelistAPIKeysService is an interface for interfacing with the Whitelist API Keys
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/#organization-api-key-endpoints
type WhitelistAPIKeysService interface {
	List(context.Context, string, string, *ListOptions) (*WhitelistAPIKeys, *Response, error)
	Get(context.Context, string, string, string) (*WhitelistAPIKey, *Response, error)
	Create(context.Context, string, string, []*WhitelistAPIKeysReq) (*WhitelistAPIKeys, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
}

// WhitelistAPIKeysServiceOp handles communication with the Whitelist API keys related methods of the
// MongoDB Atlas API
type WhitelistAPIKeysServiceOp service

var _ WhitelistAPIKeysService = &WhitelistAPIKeysServiceOp{}

// WhitelistAPIKey represents a Whitelist API key.
type WhitelistAPIKey struct {
	CidrBlock       string  `json:"cidrBlock,omitempty"`       // CIDR-notated range of whitelisted IP addresses.
	Count           int     `json:"count,omitempty"`           // Total number of requests that have originated from this IP address.
	Created         string  `json:"created,omitempty"`         // Date this IP address was added to the whitelist.
	IPAddress       string  `json:"ipAddress,omitempty"`       // Whitelisted IP address.
	LastUsed        string  `json:"lastUsed,omitempty"`        // Date of the most recent request that originated from this IP address. This field only appears if at least one request has originated from this IP address, and is only updated when a whitelisted resource is accessed.
	LastUsedAddress string  `json:"lastUsedAddress,omitempty"` // IP address from which the last call to the API was issued. This field only appears if at least one request has originated from this IP address.
	Links           []*Link `json:"links,omitempty"`           // An array of documents, representing a link to one or more sub-resources and/or related resources such as list pagination. See Linking for more information.}
}

// WhitelistAPIKeys represents all Whitelist API keys.
type WhitelistAPIKeys struct {
	Results    []*WhitelistAPIKey `json:"results,omitempty"`    // Includes one WhitelistAPIKey object for each item detailed in the results array section.
	Links      []*Link            `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// WhitelistAPIKeysReq represents the request to the mehtod create
type WhitelistAPIKeysReq struct {
	IPAddress string `json:"ipAddress,omitempty"` // IP address to be added to the whitelist for the API key.
	CidrBlock string `json:"cidrBlock,omitempty"` // Whitelist entry in CIDR notation to be added for the API key.
}

// List gets all Whitelist API keys.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-org-whitelist-get-all/
func (s *WhitelistAPIKeysServiceOp) List(ctx context.Context, orgID, apiKeyID string, listOptions *ListOptions) (*WhitelistAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}

	path := fmt.Sprintf(whitelistAPIKeysPath, orgID, apiKeyID)
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(WhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets the Whitelist API keys.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-get-one/
func (s *WhitelistAPIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID, ipAddress string) (*WhitelistAPIKey, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, nil, NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(whitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(WhitelistAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create a submit a POST request containing ipAddress or cidrBlock values which are not already present in the whitelist, Atlas adds those entries to the list of existing entries in the whitelist.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-org-whitelist-create/
func (s *WhitelistAPIKeysServiceOp) Create(ctx context.Context, orgID, apiKeyID string, createRequest []*WhitelistAPIKeysReq) (*WhitelistAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(whitelistAPIKeysPath, orgID, apiKeyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(WhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the Whitelist API keys.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-delete-one/
func (s *WhitelistAPIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID, ipAddress string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if apiKeyID == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if ipAddress == "" {
		return nil, NewArgError("snapshotId", "must be set")
	}

	path := fmt.Sprintf(whitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
