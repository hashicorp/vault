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

// FailoverDataGuardAssociationDetails The Data Guard association failover parameters.
type FailoverDataGuardAssociationDetails struct {

	// The DB system administrator password.
	DatabaseAdminPassword *string `mandatory:"true" json:"databaseAdminPassword"`
}

func (m FailoverDataGuardAssociationDetails) String() string {
	return common.PointerString(m)
}
