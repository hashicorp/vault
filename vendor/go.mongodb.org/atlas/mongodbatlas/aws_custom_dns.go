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

const customAWSDNSPath = "api/atlas/v1.0/groups/%s/awsCustomDNS"

// AWSCustomDNSService provides access to the custom AWS DNS related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/aws-custom-dns/
type AWSCustomDNSService interface {
	Get(context.Context, string) (*AWSCustomDNSSetting, *Response, error)
	Update(context.Context, string, *AWSCustomDNSSetting) (*AWSCustomDNSSetting, *Response, error)
}

// AWSCustomDNSServiceOp provides an implementation of the CustomAWSDNS interface.
type AWSCustomDNSServiceOp service

var _ AWSCustomDNSService = &AWSCustomDNSServiceOp{}

// AWSCustomDNSSetting represents the dns settings.
type AWSCustomDNSSetting struct {
	Enabled bool `json:"enabled"`
}

// Get retrieves the custom DNS configuration of an Atlas project’s clusters deployed to AWS.
//
// See more: https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-get/
func (s *AWSCustomDNSServiceOp) Get(ctx context.Context, groupID string) (*AWSCustomDNSSetting, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(customAWSDNSPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *AWSCustomDNSSetting
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// Update updates the custom DNS configuration of an Atlas project’s clusters deployed to AWS.
//
// See more: https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/
func (s *AWSCustomDNSServiceOp) Update(ctx context.Context, groupID string, r *AWSCustomDNSSetting) (*AWSCustomDNSSetting, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(customAWSDNSPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, r)
	if err != nil {
		return nil, nil, err
	}

	var root *AWSCustomDNSSetting
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}
