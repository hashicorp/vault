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
	"bytes"
	"context"
	"fmt"
	"net/http"
)

const x509AuthDBUsersPath = "api/atlas/v1.0/groups/%s/databaseUsers/%s/certs"
const x509CustomerAuthDBUserPath = "api/atlas/v1.0/groups/%s/userSecurity"

// X509AuthDBUsersService is an interface for interfacing with the x509 Authentication Database Users.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration/
type X509AuthDBUsersService interface {
	CreateUserCertificate(context.Context, string, string, int) (*UserCertificate, *Response, error)
	GetUserCertificates(context.Context, string, string, *ListOptions) ([]UserCertificate, *Response, error)
	SaveConfiguration(context.Context, string, *CustomerX509) (*CustomerX509, *Response, error)
	GetCurrentX509Conf(context.Context, string) (*CustomerX509, *Response, error)
	DisableCustomerX509(context.Context, string) (*Response, error)
}

// X509AuthDBUsersServiceOp handles communication with the  X509AuthDBUsers related methods
// of the MongoDB Atlas API.
type X509AuthDBUsersServiceOp service

var _ X509AuthDBUsersService = &X509AuthDBUsersServiceOp{}

// UserCertificate represents an X.509 Certificate for a User.
type UserCertificate struct {
	Username              string `json:"username,omitempty"`              // Username of the database user to create a certificate for.
	MonthsUntilExpiration int    `json:"monthsUntilExpiration,omitempty"` // A number of months that the created certificate is valid for before expiry, up to 24 months.default 3.
	Certificate           string `json:"certificate,omitempty"`

	ID        *int64 `json:"_id,omitempty"`       // Serial number of this certificate.
	CreatedAt string `json:"createdAt,omitempty"` // Timestamp in ISO 8601 date and time format in UTC when Atlas created this X.509 certificate.
	GroupID   string `json:"groupId,omitempty"`   // Unique identifier of the Atlas project to which this certificate belongs.
	NotAfter  string `json:"notAfter,omitempty"`  // Timestamp in ISO 8601 date and time format in UTC when this certificate expires.
	Subject   string `json:"subject,omitempty"`   // Fully distinguished name of the database user to which this certificate belongs. To learn more, see RFC 2253.
}

// UserCertificates is Array of objects where each details one unexpired database user certificate.
type UserCertificates struct {
	Links      []*Link           `json:"links"`      // One or more links to sub-resources and/or related resources.
	Results    []UserCertificate `json:"results"`    // Array of objects where each details one unexpired database user certificate.
	TotalCount int               `json:"totalCount"` // Total number of unexpired certificates returned in this response.
}

// UserSecurity represents the wrapper CustomerX509 struct.
type UserSecurity struct {
	CustomerX509 CustomerX509 `json:"customerX509,omitempty"` // CustomerX509 represents Customer-managed X.509 configuration for an Atlas project.
}

// CustomerX509 represents Customer-managed X.509 configuration for an Atlas project.
type CustomerX509 struct {
	Cas string `json:"cas,omitempty"` // PEM string containing one or more customer CAs for database user authentication.
}

// CreateUserCertificate generates an Atlas-managed X.509 certificate for a MongoDB user that authenticates using X.509 certificates.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration-create-certificate/
func (s *X509AuthDBUsersServiceOp) CreateUserCertificate(ctx context.Context, groupID, username string, monthsUntilExpiration int) (*UserCertificate, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}
	if monthsUntilExpiration == 0 {
		return nil, nil, NewArgError("monthsUntilExpiration", "must be set")
	}

	userCertificate := &UserCertificate{
		MonthsUntilExpiration: monthsUntilExpiration,
	}

	path := fmt.Sprintf(x509AuthDBUsersPath, groupID, username)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, userCertificate)
	if err != nil {
		return nil, nil, err
	}

	var cer bytes.Buffer
	resp, err := s.Client.Do(ctx, req, &cer)
	if err != nil {
		return nil, resp, err
	}

	userCertificate.Username = username
	userCertificate.Certificate = cer.String()

	return userCertificate, resp, err
}

// GetUserCertificates gets a list of all Atlas-managed, unexpired certificates for a user.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-certificates/
func (s *X509AuthDBUsersServiceOp) GetUserCertificates(ctx context.Context, groupID, username string, listOptions *ListOptions) ([]UserCertificate, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}

	path := fmt.Sprintf(x509AuthDBUsersPath, groupID, username)
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(UserCertificates)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, err
}

// SaveConfiguration saves a customer-managed X.509 configuration for an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration-save/
func (s *X509AuthDBUsersServiceOp) SaveConfiguration(ctx context.Context, groupID string, customerX509 *CustomerX509) (*CustomerX509, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if customerX509 == nil {
		return nil, nil, NewArgError("customerX509", "cannot be nil")
	}

	path := fmt.Sprintf(x509CustomerAuthDBUserPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, &UserSecurity{CustomerX509: *customerX509})
	if err != nil {
		return nil, nil, err
	}

	root := new(UserSecurity)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.CustomerX509, resp, err
}

// GetCurrentX509Conf gets the current customer-managed X.509 configuration details for an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration-get-current/
func (s *X509AuthDBUsersServiceOp) GetCurrentX509Conf(ctx context.Context, groupID string) (*CustomerX509, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(x509CustomerAuthDBUserPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(UserSecurity)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.CustomerX509, resp, err
}

// DisableCustomerX509 clears customer-managed X.509 settings on a project. This disables customer-managed X.509.
//
// See more: https://docs.atlas.mongodb.com/reference/api/x509-configuration-disable-advanced/
func (s *X509AuthDBUsersServiceOp) DisableCustomerX509(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(x509CustomerAuthDBUserPath+"/customerX509", groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}
