// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// Scope models a GitHub authorization scope.
//
// GitHub API docs:https://developer.github.com/v3/oauth/#scopes
type Scope string

// This is the set of scopes for GitHub API V3
const (
	ScopeNone           Scope = "(no scope)" // REVISIT: is this actually returned, or just a documentation artifact?
	ScopeUser           Scope = "user"
	ScopeUserEmail      Scope = "user:email"
	ScopeUserFollow     Scope = "user:follow"
	ScopePublicRepo     Scope = "public_repo"
	ScopeRepo           Scope = "repo"
	ScopeRepoDeployment Scope = "repo_deployment"
	ScopeRepoStatus     Scope = "repo:status"
	ScopeDeleteRepo     Scope = "delete_repo"
	ScopeNotifications  Scope = "notifications"
	ScopeGist           Scope = "gist"
	ScopeReadRepoHook   Scope = "read:repo_hook"
	ScopeWriteRepoHook  Scope = "write:repo_hook"
	ScopeAdminRepoHook  Scope = "admin:repo_hook"
	ScopeAdminOrgHook   Scope = "admin:org_hook"
	ScopeReadOrg        Scope = "read:org"
	ScopeWriteOrg       Scope = "write:org"
	ScopeAdminOrg       Scope = "admin:org"
	ScopeReadPublicKey  Scope = "read:public_key"
	ScopeWritePublicKey Scope = "write:public_key"
	ScopeAdminPublicKey Scope = "admin:public_key"
)

// AuthorizationsService handles communication with the authorization related
// methods of the GitHub API.
//
// This service requires HTTP Basic Authentication; it cannot be accessed using
// an OAuth token.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/
type AuthorizationsService struct {
	client *Client
}

// Authorization represents an individual GitHub authorization.
type Authorization struct {
	ID             *int              `json:"id,omitempty"`
	URL            *string           `json:"url,omitempty"`
	Scopes         []Scope           `json:"scopes,omitempty"`
	Token          *string           `json:"token,omitempty"`
	TokenLastEight *string           `json:"token_last_eight,omitempty"`
	HashedToken    *string           `json:"hashed_token,omitempty"`
	App            *AuthorizationApp `json:"app,omitempty"`
	Note           *string           `json:"note,omitempty"`
	NoteURL        *string           `json:"note_url,omitempty"`
	UpdateAt       *Timestamp        `json:"updated_at,omitempty"`
	CreatedAt      *Timestamp        `json:"created_at,omitempty"`
	Fingerprint    *string           `json:"fingerprint,omitempty"`

	// User is only populated by the Check and Reset methods.
	User *User `json:"user,omitempty"`
}

func (a Authorization) String() string {
	return Stringify(a)
}

// AuthorizationApp represents an individual GitHub app (in the context of authorization).
type AuthorizationApp struct {
	URL      *string `json:"url,omitempty"`
	Name     *string `json:"name,omitempty"`
	ClientID *string `json:"client_id,omitempty"`
}

func (a AuthorizationApp) String() string {
	return Stringify(a)
}

// AuthorizationRequest represents a request to create an authorization.
type AuthorizationRequest struct {
	Scopes       []Scope `json:"scopes,omitempty"`
	Note         *string `json:"note,omitempty"`
	NoteURL      *string `json:"note_url,omitempty"`
	ClientID     *string `json:"client_id,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
	Fingerprint  *string `json:"fingerprint,omitempty"`
}

func (a AuthorizationRequest) String() string {
	return Stringify(a)
}

// AuthorizationUpdateRequest represents a request to update an authorization.
//
// Note that for any one update, you must only provide one of the "scopes"
// fields.  That is, you may provide only one of "Scopes", or "AddScopes", or
// "RemoveScopes".
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
type AuthorizationUpdateRequest struct {
	Scopes       []string `json:"scopes,omitempty"`
	AddScopes    []string `json:"add_scopes,omitempty"`
	RemoveScopes []string `json:"remove_scopes,omitempty"`
	Note         *string  `json:"note,omitempty"`
	NoteURL      *string  `json:"note_url,omitempty"`
	Fingerprint  *string  `json:"fingerprint,omitempty"`
}

func (a AuthorizationUpdateRequest) String() string {
	return Stringify(a)
}

// List the authorizations for the authenticated user.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#list-your-authorizations
func (s *AuthorizationsService) List(opt *ListOptions) ([]Authorization, *Response, error) {
	u := "authorizations"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	auths := new([]Authorization)
	resp, err := s.client.Do(req, auths)
	if err != nil {
		return nil, resp, err
	}
	return *auths, resp, err
}

// Get a single authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
func (s *AuthorizationsService) Get(id int) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Create a new authorization for the specified OAuth application.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
func (s *AuthorizationsService) Create(auth *AuthorizationRequest) (*Authorization, *Response, error) {
	u := "authorizations"

	req, err := s.client.NewRequest("POST", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// GetOrCreateForApp creates a new authorization for the specified OAuth
// application, only if an authorization for that application doesn’t already
// exist for the user.
//
// If a new token is created, the HTTP status code will be "201 Created", and
// the returned Authorization.Token field will be populated. If an existing
// token is returned, the status code will be "200 OK" and the
// Authorization.Token field will be empty.
//
// clientID is the OAuth Client ID with which to create the token.
//
// GitHub API docs:
// - https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app
// - https://developer.github.com/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app-and-fingerprint
func (s *AuthorizationsService) GetOrCreateForApp(clientID string, auth *AuthorizationRequest) (*Authorization, *Response, error) {
	var u string
	if auth.Fingerprint == nil || *auth.Fingerprint == "" {
		u = fmt.Sprintf("authorizations/clients/%v", clientID)
	} else {
		u = fmt.Sprintf("authorizations/clients/%v/%v", clientID, *auth.Fingerprint)
	}

	req, err := s.client.NewRequest("PUT", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Edit a single authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#update-an-existing-authorization
func (s *AuthorizationsService) Edit(id int, auth *AuthorizationUpdateRequest) (*Authorization, *Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("PATCH", u, auth)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Delete a single authorization.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#delete-an-authorization
func (s *AuthorizationsService) Delete(id int) (*Response, error) {
	u := fmt.Sprintf("authorizations/%d", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Check if an OAuth token is valid for a specific app.
//
// Note that this operation requires the use of BasicAuth, but where the
// username is the OAuth application clientID, and the password is its
// clientSecret. Invalid tokens will return a 404 Not Found.
//
// The returned Authorization.User field will be populated.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#check-an-authorization
func (s *AuthorizationsService) Check(clientID string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Reset is used to reset a valid OAuth token without end user involvement.
// Applications must save the "token" property in the response, because changes
// take effect immediately.
//
// Note that this operation requires the use of BasicAuth, but where the
// username is the OAuth application clientID, and the password is its
// clientSecret. Invalid tokens will return a 404 Not Found.
//
// The returned Authorization.User field will be populated.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#reset-an-authorization
func (s *AuthorizationsService) Reset(clientID string, token string) (*Authorization, *Response, error) {
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Authorization)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

// Revoke an authorization for an application.
//
// Note that this operation requires the use of BasicAuth, but where the
// username is the OAuth application clientID, and the password is its
// clientSecret. Invalid tokens will return a 404 Not Found.
//
// GitHub API docs: https://developer.github.com/v3/oauth_authorizations/#revoke-an-authorization-for-an-application
func (s *AuthorizationsService) Revoke(clientID string, token string) (*Response, error) {
	u := fmt.Sprintf("applications/%v/tokens/%v", clientID, token)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
