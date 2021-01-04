package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const customAWSDNSPath = "groups/%s/awsCustomDNS"

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
