// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/avast/retry-go/v4"
	slogctx "github.com/veqryn/slog-context"
)

// WaitForImageReq is a request to wait for an image to be available in the
// image service.
type WaitForImageReq struct {
	Req                 *GetLatestProductVersionReq `json:"req,omitempty"`
	Delay               time.Duration               `json:"delay,omitempty"`
	WriteToGithubOutput bool                        `json:"write_to_github_output,omitempty"`
}

// WaitForImageRes is a response to a WaitForImageReq.
type WaitForImageRes struct {
	Res *GetLatestProductVersionRes `json:"res,omitempty"`
}

// Run runs the wait for image request.
func (r *WaitForImageReq) Run(ctx context.Context, client *Client) (*WaitForImageRes, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	slog.Default().DebugContext(ctx, "waiting for HCP image to be available")

	res := &WaitForImageRes{
		Res: &GetLatestProductVersionRes{},
	}
	attempt := 0
	err := retry.Do(
		func() error {
			attempt++
			// Limit each request to the image service to 5 seconds max
			reqCtx, reqCancel := context.WithTimeout(ctx, 5*time.Second)
			defer reqCancel()
			reqRes, err := r.Req.Run(reqCtx, client)
			if reqRes != nil {
				if reqRes.Response != nil {
					res.Res.Response = reqRes.Response
				}
				if reqRes.Image != nil {
					res.Res.Image = reqRes.Image
				}
			}

			if err != nil {
				slog.Default().DebugContext(
					slogctx.Append(ctx,
						slog.Int("attempt", attempt),
						slog.String("error", err.Error())),
					"attempt to get HCP image details failed",
				)
			}

			return err
		},
		retry.UntilSucceeded(),
		retry.Context(ctx),
		retry.WrapContextErrorWithLastError(true),
		retry.Delay(r.Delay),
		retry.DelayType(retry.FixedDelay),
	)

	return res, err
}

// ToGithubOutput marshals just the artifact response to JSON.
func (r *WaitForImageRes) ToGithubOutput() ([]byte, error) {
	if r == nil || r.Res == nil {
		return nil, fmt.Errorf("unable to marshal unitialized response to GITHUB_OUTPUT")
	}

	b, err := json.Marshal(r.Res.Image)
	if err != nil {
		return nil, fmt.Errorf("marshaling wait-image to GITHUB_OUTPUT JSON: %w", err)
	}

	return b, nil
}
