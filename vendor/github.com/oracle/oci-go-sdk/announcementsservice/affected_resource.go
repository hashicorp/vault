// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Announcements Service API
//
// Manage Oracle Cloud Infrastructure console announcements.
//

package announcementsservice

import (
	"github.com/oracle/oci-go-sdk/common"
)

// AffectedResource The resource affected by the event described in the announcement.
type AffectedResource struct {

	// The OCID of the affected resource.
	ResourceId *string `mandatory:"true" json:"resourceId"`

	// The friendly name of the resource.
	ResourceName *string `mandatory:"true" json:"resourceName"`

	// The region where the affected resource exists.
	Region *string `mandatory:"true" json:"region"`
}

func (m AffectedResource) String() string {
	return common.PointerString(m)
}
