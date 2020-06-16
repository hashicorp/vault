package openstack

import (
	"github.com/huaweicloud/golangsdk"
	tokens2 "github.com/huaweicloud/golangsdk/openstack/identity/v2/tokens"
	tokens3 "github.com/huaweicloud/golangsdk/openstack/identity/v3/tokens"
)

/*
V2EndpointURL discovers the endpoint URL for a specific service from a
ServiceCatalog acquired during the v2 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V2EndpointURL(catalog *tokens2.ServiceCatalog, opts golangsdk.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	var endpoints = make([]tokens2.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if endpoint.Region == opts.Region || (opts.Region == "" && endpoint.Region == "*") {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		err := &ErrMultipleMatchingEndpointsV2{}
		err.Endpoints = endpoints
		return "", err
	}

	// Extract the appropriate URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		switch opts.Availability {
		case golangsdk.AvailabilityPublic:
			return golangsdk.NormalizeURL(endpoint.PublicURL), nil
		case golangsdk.AvailabilityInternal:
			return golangsdk.NormalizeURL(endpoint.InternalURL), nil
		case golangsdk.AvailabilityAdmin:
			return golangsdk.NormalizeURL(endpoint.AdminURL), nil
		default:
			err := &ErrInvalidAvailabilityProvided{}
			err.Argument = "Availability"
			err.Value = opts.Availability
			return "", err
		}
	}

	// Report an error if there were no matching endpoints.
	err := &golangsdk.ErrEndpointNotFound{}
	return "", err
}

/*
V3EndpointURL discovers the endpoint URL for a specific service from a Catalog
acquired during the v3 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V3EndpointURL(catalog *tokens3.ServiceCatalog, opts golangsdk.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Interface,
	// Name if provided, and Region if provided.
	var endpoints = make([]tokens3.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Availability != golangsdk.AvailabilityAdmin &&
					opts.Availability != golangsdk.AvailabilityPublic &&
					opts.Availability != golangsdk.AvailabilityInternal {
					err := &ErrInvalidAvailabilityProvided{}
					err.Argument = "Availability"
					err.Value = opts.Availability
					return "", err
				}
				if opts.Availability == golangsdk.Availability(endpoint.Interface) &&
					(endpoint.Region == opts.Region ||
						(opts.Region == "" && endpoint.Region == "*")) {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		return "", ErrMultipleMatchingEndpointsV3{Endpoints: endpoints}
	}

	// Extract the URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		return golangsdk.NormalizeURL(endpoint.URL), nil
	}

	// Report an error if there were no matching endpoints.
	err := &golangsdk.ErrEndpointNotFound{}
	return "", err
}
