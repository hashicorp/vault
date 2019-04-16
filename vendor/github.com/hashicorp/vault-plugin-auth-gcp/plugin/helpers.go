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
func zoneToRegion(zone string) (string, error) {
	zone, err := zoneFromSelfLink(zone)
	if err != nil {
		return "", err
	}

	if i := strings.LastIndex(zone, "-"); i > -1 {
		return zone[0:i], nil
	}
	return "", fmt.Errorf("failed to extract region from zone name %q", zone)
}

// zoneFromSelfLink converts a zone self-link into the human zone name.
func zoneFromSelfLink(zone string) (string, error) {
	prefix := "zones/"

	if zone == "" {
		return "", fmt.Errorf("failed to extract zone from self-link %q", zone)
	}

	if strings.Contains(zone, "/") {
		if i := strings.LastIndex(zone, prefix); i > -1 {
			zone = zone[i+len(prefix) : len(zone)]
		} else {
			return "", fmt.Errorf("failed to extract zone from self-link %q", zone)
		}
	}

	return zone, nil
}
