// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// validateFields verifies that no bad arguments were given to the request.
func validateFields(req *logical.Request, data *framework.FieldData) error {
	var unknownFields []string
	for k := range req.Data {
		if _, ok := data.Schema[k]; !ok {
			unknownFields = append(unknownFields, k)
		}
	}

	if len(unknownFields) > 0 {
		// Sort since this is a human error
		sort.Strings(unknownFields)

		return fmt.Errorf("unknown fields: %q", unknownFields)
	}

	return nil
}

// zoneToRegion converts a zone name to its corresponding region. From
// https://cloud.google.com/compute/docs/regions-zones/, the FQDN of a zone is
// always <region>-<zone>. Instead of doing an API call, this function uses
// string parsing as an opimization.
//
// If the zone is a self-link, it is converted into a human name first. If the
// zone cannot be converted to a region, an error is returned.
func zoneToRegion(input string) (string, error) {
	zone, _, err := zoneOrRegionFromSelfLink(input)
	if err != nil {
		return "", err
	}

	if i := strings.LastIndex(zone, "-"); i > -1 {
		return zone[0:i], nil
	}
	return "", fmt.Errorf("failed to extract region from zone name %q", input)
}

// zoneOrRegionFromSelfLink converts a zone or region self-link into the human
// zone or region names.
func zoneOrRegionFromSelfLink(selfLink string) (string, string, error) {
	zPrefix := "zones/"
	rPrefix := "regions/"
	var zone, region string

	if selfLink == "" {
		return "", "", fmt.Errorf("failed to extract zone or region from self-link %q", selfLink)
	}

	if strings.Contains(selfLink, "/") {
		if i := strings.LastIndex(selfLink, zPrefix); i > -1 {
			zone = selfLink[i+len(zPrefix) : len(selfLink)]
		} else if i := strings.LastIndex(selfLink, rPrefix); i > -1 {
			region = selfLink[i+len(rPrefix) : len(selfLink)]
		} else {
			return "", "", fmt.Errorf("failed to extract zone or region from self-link %q", selfLink)
		}
	} else {
		return selfLink, "", nil
	}

	return zone, region, nil
}
