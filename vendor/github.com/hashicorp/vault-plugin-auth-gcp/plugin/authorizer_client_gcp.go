// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

var _ client = (*gcpClient)(nil)

// gcpClient implements client and communicates with the GCP API. It is
// abstracted as an interface for stubbing during testing. See stubbedClient for
// more details.
type gcpClient struct {
	logger     log.Logger
	computeSvc *compute.Service
	iamSvc     *iam.Service
}

func (c *gcpClient) InstanceGroups(ctx context.Context, project string, boundInstanceGroups []string) (map[string][]string, map[string][]string, error) {
	// maps of zone/region names to a slice of instance group names in that
	// location.
	igz := make(map[string][]string)
	igr := make(map[string][]string)

	// AggregatedList, unlike all the other InstanceGroupsService methods,
	// returns both zonal and regional instance groups.
	if err := c.computeSvc.InstanceGroups.
		AggregatedList(project).
		Fields("items/*/instanceGroups/name").
		Pages(ctx, func(l *compute.InstanceGroupAggregatedList) error {
			for k, v := range l.Items {
				zone, region, err := zoneOrRegionFromSelfLink(k)
				if err != nil {
					return err
				}

				for _, g := range v.InstanceGroups {
					if strutil.StrListContains(boundInstanceGroups, g.Name) {
						if zone != "" {
							igz[zone] = append(igz[zone], g.Name)
						}
						if region != "" {
							igr[region] = append(igr[region], g.Name)
						}
					}
				}
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}

	return igz, igr, nil
}

func (c *gcpClient) InstanceGroupContainsInstance(ctx context.Context, project, zone, region, group, instanceSelfLink string) (bool, error) {
	if zone != "" {
		return c.zoneInstanceGroupContainsInstance(ctx, project, zone, group, instanceSelfLink)
	} else {
		return c.regionInstanceGroupContainsInstance(ctx, project, region, group, instanceSelfLink)
	}
}

func (c *gcpClient) zoneInstanceGroupContainsInstance(ctx context.Context, project, zone, group, instanceSelfLink string) (bool, error) {
	var req compute.InstanceGroupsListInstancesRequest
	resp, err := c.computeSvc.InstanceGroups.
		ListInstances(project, zone, group, &req).
		Filter(fmt.Sprintf("instance eq %s", instanceSelfLink)).
		Context(ctx).
		Do()
	if err != nil {
		return false, err
	}

	if resp != nil && len(resp.Items) > 0 {
		return true, nil
	}
	return false, nil
}

func (c *gcpClient) regionInstanceGroupContainsInstance(ctx context.Context, project, region, group, instanceSelfLink string) (bool, error) {
	var req compute.RegionInstanceGroupsListInstancesRequest
	resp, err := c.computeSvc.RegionInstanceGroups.
		ListInstances(project, region, group, &req).
		Filter(fmt.Sprintf("instance eq %s", instanceSelfLink)).
		Context(ctx).
		Do()
	if err != nil {
		return false, err
	}

	if resp != nil && len(resp.Items) > 0 {
		return true, nil
	}
	return false, nil
}

func (c *gcpClient) ServiceAccount(ctx context.Context, name string) (string, string, error) {
	account, err := c.iamSvc.Projects.ServiceAccounts.
		Get(name).
		Fields("uniqueId", "email").
		Context(ctx).
		Do()
	if err != nil {
		return "", "", err
	}

	return account.UniqueId, account.Email, nil
}
