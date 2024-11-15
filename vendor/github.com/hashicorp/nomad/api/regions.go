// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import "sort"

// Regions is used to query the regions in the cluster.
type Regions struct {
	client *Client
}

// Regions returns a handle on the regions endpoints.
func (c *Client) Regions() *Regions {
	return &Regions{client: c}
}

// List returns a list of all of the regions from the server
// that serves the request. It is never forwarded to a leader.
func (r *Regions) List() ([]string, error) {
	var resp []string
	if _, err := r.client.query("/v1/regions", &resp, nil); err != nil {
		return nil, err
	}
	sort.Strings(resp)
	return resp, nil
}
