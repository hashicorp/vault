// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iamutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
)

// IamResource implements Resource.
type IamResource struct {
	relativeId *gcputil.RelativeResourceName
	config     *RestResource
}

func (r *IamResource) GetConfig() *RestResource {
	return r.config
}

func (r *IamResource) GetRelativeId() *gcputil.RelativeResourceName {
	return r.relativeId
}

func (r *IamResource) GetIamPolicy(ctx context.Context, h *ApiHandle) (*Policy, error) {
	var p Policy
	if err := h.DoGetRequest(ctx, r, &p); err != nil {
		return nil, errwrap.Wrapf("unable to get policy: {{err}}", err)
	}
	return &p, nil
}

func (r *IamResource) SetIamPolicy(ctx context.Context, h *ApiHandle, p *Policy) (*Policy, error) {
	jsonP, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	reqJson := fmt.Sprintf(r.config.SetMethod.RequestFormat, jsonP)
	if !json.Valid([]byte(reqJson)) {
		return nil, fmt.Errorf("request format from generated IAM config invalid JSON: %s", reqJson)
	}

	var policy Policy
	if err := h.DoSetRequest(ctx, r, strings.NewReader(reqJson), &policy); err != nil {
		return nil, errwrap.Wrapf("unable to set policy: {{err}}", err)
	}
	return &policy, nil
}
