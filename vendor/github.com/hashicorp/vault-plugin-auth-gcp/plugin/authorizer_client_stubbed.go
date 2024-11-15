// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import "context"

var _ client = (*stubbedClient)(nil)

// stubbedClient is a simple client to use for testing where the API calls to
// GCP are "stubbed" instead of hitting the actual API.
type stubbedClient struct {
	instanceGroupsByZone          map[string][]string
	instanceGroupsByRegion        map[string][]string
	instanceGroupContainsInstance bool
	saId, saEmail                 string
}

func (c *stubbedClient) InstanceGroups(_ context.Context, _ string, _ []string) (map[string][]string, map[string][]string, error) {
	return c.instanceGroupsByZone, c.instanceGroupsByRegion, nil
}

func (c *stubbedClient) InstanceGroupContainsInstance(_ context.Context, _, _, _, _, _ string) (bool, error) {
	return c.instanceGroupContainsInstance, nil
}

func (c *stubbedClient) ServiceAccount(_ context.Context, _ string) (string, string, error) {
	return c.saId, c.saEmail, nil
}
