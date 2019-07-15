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

// CreateCrossConnectDetails The representation of CreateCrossConnectDetails
type CreateCrossConnectDetails struct {

	// The OCID of the compartment to contain the cross-connect.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name of the FastConnect location where this cross-connect will be installed.
	// To get a list of the available locations, see
	// ListCrossConnectLocations.
	// Example: `CyrusOne, Chandler, AZ`
	LocationName *string `mandatory:"true" json:"locationName"`

	// The port speed for this cross-connect. To get a list of the available port speeds, see
	// ListCrossconnectPortSpeedShapes.
	// Example: `10 Gbps`
	PortSpeedShapeName *string `mandatory:"true" json:"portSpeedShapeName"`

	// The OCID of the cross-connect group to put this cross-connect in.
	CrossConnectGroupId *string `mandatory:"false" json:"crossConnectGroupId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// If you already have an existing cross-connect or cross-connect group at this FastConnect
	// location, and you want this new cross-connect to be on a different router (for the
	// purposes of redundancy), provide the OCID of that existing cross-connect or
	// cross-connect group.
	FarCrossConnectOrCrossConnectGroupId *string `mandatory:"false" json:"farCrossConnectOrCrossConnectGroupId"`

	// If you already have an existing cross-connect or cross-connect group at this FastConnect
	// location, and you want this new cross-connect to be on the same router, provide the
	// OCID of that existing cross-connect or cross-connect group.
	NearCrossConnectOrCrossConnectGroupId *string `mandatory:"false" json:"nearCrossConnectOrCrossConnectGroupId"`

	// A reference name or identifier for the physical fiber connection that this cross-connect
	// uses.
	CustomerReferenceName *string `mandatory:"false" json:"customerReferenceName"`
}

func (m CreateCrossConnectDetails) String() string {
	return common.PointerString(m)
}
