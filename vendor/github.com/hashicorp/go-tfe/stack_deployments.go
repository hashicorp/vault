// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// StackDeployments describes all the stacks deployments-related methods that the
// HCP Terraform API supports.
// NOTE WELL: This is a beta feature and is subject to change until noted otherwise in the
// release notes.
type StackDeployments interface {
	// Read returns a stack deployment by its name.
	Read(ctx context.Context, stackID, deployment string) (*StackDeployment, error)
}

type stackDeployments struct {
	client *Client
}

func (s stackDeployments) Read(ctx context.Context, stackID, deploymentName string) (*StackDeployment, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("stacks/%s/stack-deployments/%s", url.PathEscape(stackID), url.PathEscape(deploymentName)), nil)
	if err != nil {
		return nil, err
	}

	deployment := &StackDeployment{}
	err = req.Do(ctx, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}
