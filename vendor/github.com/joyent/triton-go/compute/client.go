//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute

import (
	"net/http"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/client"
)

type ComputeClient struct {
	Client *client.Client
}

func newComputeClient(client *client.Client) *ComputeClient {
	return &ComputeClient{
		Client: client,
	}
}

// NewClient returns a new client for working with Compute endpoints and
// resources within CloudAPI
func NewClient(config *triton.ClientConfig) (*ComputeClient, error) {
	// TODO: Utilize config interface within the function itself
	client, err := client.New(
		config.TritonURL,
		config.MantaURL,
		config.AccountName,
		config.Signers...,
	)
	if err != nil {
		return nil, err
	}
	return newComputeClient(client), nil
}

// SetHeaders allows a consumer of the current client to set custom headers for
// the next backend HTTP request sent to CloudAPI
func (c *ComputeClient) SetHeader(header *http.Header) {
	c.Client.RequestHeader = header
}

// Datacenters returns a Compute client used for accessing functions pertaining
// to DataCenter functionality in the Triton API.
func (c *ComputeClient) Datacenters() *DataCentersClient {
	return &DataCentersClient{c.Client}
}

// Images returns a Compute client used for accessing functions pertaining to
// Images functionality in the Triton API.
func (c *ComputeClient) Images() *ImagesClient {
	return &ImagesClient{c.Client}
}

// Machine returns a Compute client used for accessing functions pertaining to
// machine functionality in the Triton API.
func (c *ComputeClient) Instances() *InstancesClient {
	return &InstancesClient{c.Client}
}

// Packages returns a Compute client used for accessing functions pertaining to
// Packages functionality in the Triton API.
func (c *ComputeClient) Packages() *PackagesClient {
	return &PackagesClient{c.Client}
}

// Services returns a Compute client used for accessing functions pertaining to
// Services functionality in the Triton API.
func (c *ComputeClient) Services() *ServicesClient {
	return &ServicesClient{c.Client}
}

// Snapshots returns a Compute client used for accessing functions pertaining to
// Snapshots functionality in the Triton API.
func (c *ComputeClient) Snapshots() *SnapshotsClient {
	return &SnapshotsClient{c.Client}
}

// Snapshots returns a Compute client used for accessing functions pertaining to
// Snapshots functionality in the Triton API.
func (c *ComputeClient) Volumes() *VolumesClient {
	return &VolumesClient{c.Client}
}
