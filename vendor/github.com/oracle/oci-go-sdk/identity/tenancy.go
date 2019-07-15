// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Tenancy The root compartment that contains all of your organization's compartments and other
// Oracle Cloud Infrastructure cloud resources. When you sign up for Oracle Cloud Infrastructure,
// Oracle creates a tenancy for your company, which is a secure and isolated partition
// where you can create, organize, and administer your cloud resources.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access,
// see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type Tenancy struct {

	// The OCID of the tenancy.
	Id *string `mandatory:"false" json:"id"`

	// The name of the tenancy.
	Name *string `mandatory:"false" json:"name"`

	// The description of the tenancy.
	Description *string `mandatory:"false" json:"description"`

	// The region key for the tenancy's home region. For more information about regions, see
	// Regions and Availability Domains (https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm).
	// Allowed values are:
	// - `IAD`
	// - `PHX`
	// - `FRA`
	// - `LHR`
	HomeRegionKey *string `mandatory:"false" json:"homeRegionKey"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Tenancy) String() string {
	return common.PointerString(m)
}
