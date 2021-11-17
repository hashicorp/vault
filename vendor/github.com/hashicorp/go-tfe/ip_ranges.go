package tfe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-retryablehttp"
)

// Compile-time proof of interface implementation.
var _ IPRanges = (*ipRanges)(nil)

// IP Ranges provides a list of Terraform Cloud and Enterprise's IP ranges.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/ip-ranges.html
type IPRanges interface {
	// Retrieve TFC IP ranges. If `modifiedSince` is not an empty string
	// then it will only return the IP ranges changes since that date.
	// The format for `modifiedSince` can be found here:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since
	Read(ctx context.Context, modifiedSince string) (*IPRange, error)
}

type ipRanges struct {
	client *Client
}

type IPRange struct {
	// List of IP ranges in CIDR notation used for connections from user site to Terraform Cloud APIs
	API []string `json:"api"`
	// List of IP ranges in CIDR notation used for notifications
	Notifications []string `json:"notifications"`
	// List of IP ranges in CIDR notation used for outbound requests from Sentinel policies
	Sentinel []string `json:"sentinel"`
	// List of IP ranges in CIDR notation used for connecting to VCS providers
	VCS []string `json:"vcs"`
}

func (i *ipRanges) Read(ctx context.Context, modifiedSince string) (*IPRange, error) {
	req, err := i.client.newRequest("GET", "/api/meta/ip-ranges", nil)
	if err != nil {
		return nil, err
	}

	if modifiedSince != "" {
		req.Header.Add("If-Modified-Since", modifiedSince)
	}

	ir := &IPRange{}
	err = i.customDo(ctx, req, ir)
	if err != nil {
		return nil, err
	}

	return ir, nil
}

// The IP ranges API is not returning jsonapi like every other endpoint
// which means we need to handle it differently.
func (i *ipRanges) customDo(ctx context.Context, req *retryablehttp.Request, ir *IPRange) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if err := i.client.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add the context to the request.
	req = req.WithContext(ctx)

	// Execute the request and check the response.
	resp, err := i.client.http.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 400 {
		return fmt.Errorf("error HTTP response while retrieving IP ranges: %d", resp.StatusCode)
	} else if resp.StatusCode == 304 {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(ir)
	if err != nil {
		return err
	}
	return nil
}
