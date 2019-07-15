// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateServiceGatewayDetails The representation of CreateServiceGatewayDetails
type CreateServiceGatewayDetails struct {

	// The OCID  (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment to contain the service gateway.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// List of the OCIDs of the Service objects to
	// enable for the service gateway. This list can be empty if you don't want to enable any
	// `Service` objects when you create the gateway. You can enable a `Service`
	// object later by using either AttachServiceId
	// or UpdateServiceGateway.
	// For each enabled `Service`, make sure there's a route rule with the `Service` object's `cidrBlock`
	// as the rule's destination and the service gateway as the rule's target. See
	// RouteTable.
	Services []ServiceIdRequestDetails `mandatory:"true" json:"services"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the VCN.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m CreateServiceGatewayDetails) String() string {
	return common.PointerString(m)
}
