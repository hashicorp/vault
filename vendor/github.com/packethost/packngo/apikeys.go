package packngo

import (
	"fmt"
	"path"
)

const (
	apiKeyBasePath = "/api-keys"
)

// APIKeyService interface defines available device methods
type APIKeyService interface {
	UserList(*ListOptions) ([]APIKey, *Response, error)
	ProjectList(string, *ListOptions) ([]APIKey, *Response, error)
	UserGet(string, *GetOptions) (*APIKey, error)
	ProjectGet(string, string, *GetOptions) (*APIKey, error)
	Create(*APIKeyCreateRequest) (*APIKey, *Response, error)
	Delete(string) (*Response, error)
}

type apiKeyRoot struct {
	APIKeys []APIKey `json:"api_keys"`
}

type APIKey struct {
	// ID is the UUIDv4 representing an API key in API requests and responses.
	ID string `json:"id"`

	// Description is any text description of the key. This can be used to
	// describe the purpose of the key.
	Description string `json:"description"`

	// Token is a sensitive credential that can be used as a `Client.APIKey` to
	// access Equinix Metal resources.
	Token string `json:"token"`

	// ReadOnly keys can not create new resources.
	ReadOnly bool `json:"read_only"`

	// Created is the creation date of the API key.
	Created string `json:"created_at"`

	// Updated is the last-update date of the API key.
	Updated string `json:"updated_at"`

	// User will be non-nil when getting or listing an User API key.
	User *User `json:"user"`

	// Project will be non-nil when getting or listing a Project API key
	Project *Project `json:"project"`
}

// APIKeyCreateRequest type used to create an api key.
type APIKeyCreateRequest struct {
	// Description is any text description of the key. This can be used to
	// describe the purpose of the key.
	Description string `json:"description"`

	// ReadOnly keys can not create new resources.
	ReadOnly bool `json:"read_only"`

	// ProjectID when non-empty will result in the creation of a Project API
	// key.
	ProjectID string `json:"-"`
}

func (s APIKeyCreateRequest) String() string {
	return Stringify(s)
}

// APIKeyServiceOp implements APIKeyService
type APIKeyServiceOp struct {
	client *Client
}

func (s *APIKeyServiceOp) list(url string, opts *ListOptions) ([]APIKey, *Response, error) {
	root := new(apiKeyRoot)
	apiPathQuery := opts.WithQuery(url)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.APIKeys, resp, err
}

// ProjectList lists the API keys associated with a project having `projectID`
// match `Project.ID`.
func (s *APIKeyServiceOp) ProjectList(projectID string, opts *ListOptions) ([]APIKey, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, apiKeyBasePath)
	return s.list(endpointPath, opts)
}

// UserList returns the API keys for the User associated with the
// `Client.APIKey`.
//
// When `Client.APIKey` is a Project API key, this method will return an access
// denied error.
func (s *APIKeyServiceOp) UserList(opts *ListOptions) ([]APIKey, *Response, error) {
	endpointPath := path.Join(userBasePath, apiKeyBasePath)
	return s.list(endpointPath, opts)
}

// ProjectGet returns the Project API key with the given `APIKey.ID`.
//
// In other methods, it is typical for a Response to be returned, which could
// include a StatusCode of `http.StatusNotFound` (404 error) when the resource
// was not found. The Equinix Metal API does not expose a get by ID endpoint for
// APIKeys.  That is why in this method, all API keys are listed and compared
// for a match. Therefor, the Response is not returned and a custom error will
// be returned when the key is not found.
func (s *APIKeyServiceOp) ProjectGet(projectID, apiKeyID string, opts *GetOptions) (*APIKey, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, validateErr
	}
	if validateErr := ValidateUUID(apiKeyID); validateErr != nil {
		return nil, validateErr
	}
	pkeys, _, err := s.ProjectList(projectID, opts)
	if err != nil {
		return nil, err
	}
	for _, k := range pkeys {
		if k.ID == apiKeyID {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("Project (%s) API key %s not found", projectID, apiKeyID)
}

// UserGet returns the User API key with the given `APIKey.ID`.
//
// In other methods, it is typical for a Response to be returned, which could
// include a StatusCode of `http.StatusNotFound` (404 error) when the resource
// was not found. The Equinix Metal API does not expose a get by ID endpoint for
// APIKeys.  That is why in this method, all API keys are listed and compared
// for a match. Therefor, the Response is not returned and a custom error will
// be returned when the key is not found.
func (s *APIKeyServiceOp) UserGet(apiKeyID string, opts *GetOptions) (*APIKey, error) {
	if validateErr := ValidateUUID(apiKeyID); validateErr != nil {
		return nil, validateErr
	}
	ukeys, _, err := s.UserList(opts)
	if err != nil {
		return nil, err
	}
	for _, k := range ukeys {
		if k.ID == apiKeyID {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("User API key %s not found", apiKeyID)
}

// Create creates a new API key.
//
// The API key can be either an User API key or a Project API key, determined by
// the value (or emptiness) of `APIKeyCreateRequest.ProjectID`. Either `User` or
// `Project` will be non-nil in the `APIKey` depending on this factor.
func (s *APIKeyServiceOp) Create(createRequest *APIKeyCreateRequest) (*APIKey, *Response, error) {
	apiPath := path.Join(userBasePath, apiKeyBasePath)
	if createRequest.ProjectID != "" {
		apiPath = path.Join(projectBasePath, createRequest.ProjectID, apiKeyBasePath)
	}
	apiKey := new(APIKey)

	resp, err := s.client.DoRequest("POST", apiPath, createRequest, apiKey)
	if err != nil {
		return nil, resp, err
	}

	return apiKey, resp, err
}

// Delete deletes an API key by `APIKey.ID`
//
// The API key can be either an User API key or a Project API key.
//
// Project API keys can not be used to delete themselves.
func (s *APIKeyServiceOp) Delete(apiKeyID string) (*Response, error) {
	if validateErr := ValidateUUID(apiKeyID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(userBasePath, apiKeyBasePath, apiKeyID)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
