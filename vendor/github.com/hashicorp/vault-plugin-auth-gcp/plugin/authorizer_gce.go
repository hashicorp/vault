// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
)

type client interface {
	InstanceGroups(context.Context, string, []string) (map[string][]string, map[string][]string, error)
	InstanceGroupContainsInstance(context.Context, string, string, string, string, string) (bool, error)
	ServiceAccount(context.Context, string) (string, string, error)
}

type AuthorizeGCEInput struct {
	client client

	project        string
	serviceAccount string

	instanceLabels   map[string]string
	instanceSelfLink string
	instanceZone     string

	boundLabels  map[string]string
	boundRegions []string
	boundZones   []string

	boundInstanceGroups  []string
	boundServiceAccounts []string
}

func AuthorizeGCE(ctx context.Context, i *AuthorizeGCEInput) error {
	// Verify instance has role labels if labels were set on role.
	for k, v := range i.boundLabels {
		if act, ok := i.instanceLabels[k]; !ok || act != v {
			return fmt.Errorf("instance missing bound label \"%s:%s\"", k, v)
		}
	}

	// Parse the zone name from the self-link URI if given; compute
	// instances are always zonal.
	zone, _, err := zoneOrRegionFromSelfLink(i.instanceZone)
	if err != nil {
		return err
	}

	// Convert the zone to a region name.
	region, err := zoneToRegion(zone)
	if err != nil {
		return err
	}

	// Verify the instance is in the zone/region
	switch {
	case len(i.boundZones) > 0:
		if !strutil.StrListContains(i.boundZones, zone) {
			return fmt.Errorf("instance not in bound zones %q", i.boundZones)
		}
	case len(i.boundRegions) > 0:
		if !strutil.StrListContains(i.boundRegions, region) {
			return fmt.Errorf("instance not in bound regions %q", i.boundRegions)
		}
	}

	// For each bound instance group, verify the group exists and that the
	// instance is a member of that group.
	if len(i.boundInstanceGroups) > 0 {
		igz, igr, err := i.client.InstanceGroups(ctx, i.project, i.boundInstanceGroups)
		if err != nil {
			return fmt.Errorf("failed to list instance groups for project %q: %s", i.project, err)
		}

		// Keep track of whether we've successfully found an instance group of
		// which this instance is a member, which meets the zonal/regional criteria.
		found := false

		for _, g := range i.boundInstanceGroups {
			if found {
				break
			}

			var group, zone, region string

			switch {
			case len(i.boundZones) > 0:
				for _, z := range i.boundZones {
					if groups, ok := igz[z]; ok && len(groups) > 0 {
						for _, grp := range groups {
							if grp == g {
								group = g
								zone = z
							}
						}
					}
				}
				if group == "" {
					return fmt.Errorf("instance group %q does not exist in zones %q for project %q",
						g, i.boundZones, i.project)
				}
			case len(i.boundRegions) > 0:
				for _, r := range i.boundRegions {
					for z, groups := range igz {
						if strings.HasPrefix(z, r) { // zone is prefixed with region
							for _, grp := range groups {
								if grp == g {
									group = g
									zone = z
								}
							}
						}
					}
					for r, groups := range igr {
						for _, grp := range groups {
							if grp == g {
								group = g
								region = r
							}
						}
					}
				}
				if group == "" {
					return fmt.Errorf("instance group %q does not exist in regions %q for project %q",
						g, i.boundRegions, i.project)
				}
			default:
				return fmt.Errorf("instance group %q is not bound to any zones or regions", g)
			}

			ok, err := i.client.InstanceGroupContainsInstance(ctx, i.project, zone, region, group, i.instanceSelfLink)
			if err != nil {
				return fmt.Errorf("failed to list instances in instance group %q for project %q: %s",
					group, i.project, err)
			}

			if ok {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("instance is not part of instance groups %q",
				i.boundInstanceGroups)
		}
	}

	// Verify instance is running under one of the allowed service accounts.
	if len(i.boundServiceAccounts) > 0 {
		// ServiceAccount wraps a call to the GCP IAM API to get a service account.
		name := fmt.Sprintf("projects/-/serviceAccounts/%s", i.serviceAccount)

		saId, saEmail, err := i.client.ServiceAccount(ctx, name)
		if err != nil {
			return fmt.Errorf("could not find service account %q: %w", i.serviceAccount, err)
		}

		if !(strutil.StrListContains(i.boundServiceAccounts, saEmail) ||
			strutil.StrListContains(i.boundServiceAccounts, saId)) {
			return fmt.Errorf("service account %q (%q) is not in bound service accounts %q",
				saId, saEmail, i.boundServiceAccounts)
		}
	}

	return nil
}
