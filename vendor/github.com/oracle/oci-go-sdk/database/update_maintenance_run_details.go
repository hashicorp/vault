// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateMaintenanceRunDetails Describes the modification parameters for the Maintenance Run.
type UpdateMaintenanceRunDetails struct {

	// If set to false, skips the Maintenance Run.
	IsEnabled *bool `mandatory:"false" json:"isEnabled"`
}

func (m UpdateMaintenanceRunDetails) String() string {
	return common.PointerString(m)
}
