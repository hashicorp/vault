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

// FaultDomain A Fault Domain is a logical grouping of hardware and infrastructure within an Availability Domain that can become
// unavailable in its entirety either due to hardware failure such as Top-of-rack (TOR) switch failure or due to
// planned software maintenance such as security updates that reboot your instances.
type FaultDomain struct {

	// The name of the Fault Domain.
	Name *string `mandatory:"false" json:"name"`

	// The OCID of the Fault Domain.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the compartment. Currently only tenancy (root) compartment can be provided.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The name of the availabilityDomain where the Fault Domain belongs.
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`
}

func (m FaultDomain) String() string {
	return common.PointerString(m)
}
