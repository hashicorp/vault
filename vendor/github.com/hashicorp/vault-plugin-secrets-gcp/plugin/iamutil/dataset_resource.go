// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iamutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
)

// NOTE: BigQuery does not conform to the typical REST for IAM policies
// instead it has an access array with bindings on the dataset
// object. https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets#Dataset
type AccessBinding struct {
	Role         string `json:"role,omitempty"`
	UserByEmail  string `json:"userByEmail,omitempty"`
	GroupByEmail string `json:"groupByEmail,omitempty"`
}

type Dataset struct {
	Access []*AccessBinding `json:"access,omitempty"`
	Etag   string           `json:"etag,omitempty"`
}

// NOTE: DatasetResource implements IamResource.
// This is because bigquery datasets have their own
// ACLs instead of an IAM policy
type DatasetResource struct {
	relativeId *gcputil.RelativeResourceName
	config     *RestResource
}

func (r *DatasetResource) GetConfig() *RestResource {
	return r.config
}

func (r *DatasetResource) GetRelativeId() *gcputil.RelativeResourceName {
	return r.relativeId
}

func (r *DatasetResource) GetIamPolicy(ctx context.Context, h *ApiHandle) (*Policy, error) {
	var dataset Dataset
	if err := h.DoGetRequest(ctx, r, &dataset); err != nil {
		return nil, errwrap.Wrapf("unable to get BigQuery Dataset ACL: {{err}}", err)
	}
	p := datasetAsPolicy(&dataset)
	return p, nil
}

func (r *DatasetResource) SetIamPolicy(ctx context.Context, h *ApiHandle, p *Policy) (*Policy, error) {
	var jsonP []byte
	ds, err := policyAsDataset(p)
	if err != nil {
		return nil, err
	}
	jsonP, err = json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	reqJson := fmt.Sprintf(r.config.SetMethod.RequestFormat, jsonP)
	if !json.Valid([]byte(reqJson)) {
		return nil, fmt.Errorf("request format from generated BigQuery Dataset config invalid JSON: %s", reqJson)
	}

	var dataset Dataset
	if err := h.DoSetRequest(ctx, r, strings.NewReader(reqJson), &dataset); err != nil {
		return nil, errwrap.Wrapf("unable to set BigQuery Dataset ACL: {{err}}", err)
	}
	policy := datasetAsPolicy(&dataset)

	return policy, nil
}

func policyAsDataset(p *Policy) (*Dataset, error) {
	if p == nil {
		return nil, errors.New("Policy cannot be nil")
	}

	ds := &Dataset{Etag: p.Etag}
	for _, binding := range p.Bindings {
		if binding.Condition != nil {
			return nil, errors.New("Bigquery Datasets do not support conditional IAM")
		}
		for _, member := range binding.Members {
			var email, iamType string
			memberSplit := strings.Split(member, ":")
			if len(memberSplit) == 2 {
				iamType = memberSplit[0]
				email = memberSplit[1]
			} else {
				email = member
			}

			if email != "" {
				binding := &AccessBinding{Role: binding.Role}
				if iamType == "group" {
					binding.GroupByEmail = email
				} else {
					binding.UserByEmail = email
				}
				ds.Access = append(ds.Access, binding)
			}
		}
	}
	return ds, nil
}

func datasetAsPolicy(ds *Dataset) *Policy {
	if ds == nil {
		return &Policy{}
	}

	policy := &Policy{Etag: ds.Etag}
	bindingMap := make(map[string]*Binding)
	for _, accessBinding := range ds.Access {
		var iamMember string

		// Role mapping must be applied for datasets in order to properly
		// detect when to change bindings (via RemoveBindings()) after a
		// modification or deletion occurs. This is due to BigQuery
		// access roles accepting both legacy (e.g., OWNER) and current
		// (e.g., roles/bigquery.dataOwner) role references. The API will
		// only return the legacy format, so this mapping allows us to properly
		// diff the current and desired roles to set the access policy.
		//
		// See the access[].role description in the following document for details
		// https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets#Dataset
		role := mapLegacyRoles(accessBinding.Role)

		//NOTE: Can either have GroupByEmail or UserByEmail but not both
		if accessBinding.GroupByEmail != "" {
			iamMember = fmt.Sprintf("group:%s", accessBinding.GroupByEmail)
		} else if strings.HasSuffix(accessBinding.UserByEmail, "gserviceaccount.com") {
			iamMember = fmt.Sprintf("serviceAccount:%s", accessBinding.UserByEmail)
		} else {
			iamMember = fmt.Sprintf("user:%s", accessBinding.UserByEmail)
		}
		if binding, ok := bindingMap[role]; ok {
			binding.Members = append(binding.Members, iamMember)
		} else {
			bindingMap[role] = &Binding{
				Role:    role,
				Members: []string{iamMember},
			}
		}
	}
	for _, v := range bindingMap {
		policy.Bindings = append(policy.Bindings, v)
	}
	return policy
}

// mapLegacyRoles returns a current role name given a legacy role name.
//
// The following role mappings will be applied:
// - OWNER -> roles/bigquery.dataOwner
// - WRITER -> roles/bigquery.dataEditor
// - READER -> roles/bigquery.dataViewer
//
// See the access[].role description in the following document for details
// https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets#Dataset
//
// Returns the given role if no mapping applies.
func mapLegacyRoles(role string) string {
	switch role {
	case "OWNER":
		return "roles/bigquery.dataOwner"
	case "WRITER":
		return "roles/bigquery.dataEditor"
	case "READER":
		return "roles/bigquery.dataViewer"
	default:
		return role
	}
}
