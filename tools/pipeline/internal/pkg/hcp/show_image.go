// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

// ShowImageReq is a request to wait for an image to be available in the
// image service.
type ShowImageReq struct {
	Req                 *GetLatestProductVersionReq `json:"req,omitempty"`
	WriteToGithubOutput bool                        `json:"write_to_github_output,omitempty"`
}

// ShowImageRes is a response to a ShowImageReq.
type ShowImageRes struct {
	Res *GetLatestProductVersionRes `json:"res,omitempty"`
}

// Run runs the wait for image request.
func (r *ShowImageReq) Run(ctx context.Context, client *Client) (*ShowImageRes, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	slog.Default().DebugContext(ctx, "showing HCP image")

	res := &ShowImageRes{}
	var err error
	res.Res, err = r.Req.Run(ctx, client)

	return res, err
}

// ToGithubOutput marshals just the artifact response to JSON.
func (r *ShowImageRes) ToGithubOutput() ([]byte, error) {
	if r == nil || r.Res == nil {
		return nil, fmt.Errorf("unable to marshal unitialized response to GITHUB_OUTPUT")
	}

	b, err := json.Marshal(r.Res.Image)
	if err != nil {
		return nil, fmt.Errorf("marshaling show-image to GITHUB_OUTPUT JSON: %w", err)
	}

	return b, nil
}
