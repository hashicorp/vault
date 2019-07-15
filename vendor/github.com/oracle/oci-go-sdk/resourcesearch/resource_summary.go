// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Search Service API
//
// Search for resources in your cloud network.
//

package resourcesearch

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ResourceSummary A resource that exists in the user's cloud network.
type ResourceSummary struct {

	// The resource type name.
	ResourceType *string `mandatory:"true" json:"resourceType"`

	// The unique identifier for this particular resource, usually an OCID.
	Identifier *string `mandatory:"true" json:"identifier"`

	// The OCID of the compartment that contains this resource.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The time this resource was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The display name (or name) of this resource, if one exists.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The availability domain this resource is located in, if applicable.
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// The lifecycle state of this resource, if applicable.
	LifecycleState *string `mandatory:"false" json:"lifecycleState"`

	// The freeform tags associated with this resource, if any.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The defined tags associated with this resource, if any.
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Contains search context, such as highlighting, for found resources.
	SearchContext *SearchContext `mandatory:"false" json:"searchContext"`
}

func (m ResourceSummary) String() string {
	return common.PointerString(m)
}
