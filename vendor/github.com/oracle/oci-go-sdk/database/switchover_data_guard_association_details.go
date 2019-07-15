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

// SwitchoverDataGuardAssociationDetails The Data Guard association switchover parameters.
type SwitchoverDataGuardAssociationDetails struct {

	// The DB system administrator password.
	DatabaseAdminPassword *string `mandatory:"true" json:"databaseAdminPassword"`
}

func (m SwitchoverDataGuardAssociationDetails) String() string {
	return common.PointerString(m)
}
