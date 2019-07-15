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

// DatabaseConnectionStrings Connection strings to connect to an Oracle Database.
type DatabaseConnectionStrings struct {

	// Host name based CDB Connection String.
	CdbDefault *string `mandatory:"false" json:"cdbDefault"`

	// IP based CDB Connection String.
	CdbIpDefault *string `mandatory:"false" json:"cdbIpDefault"`

	// All connection strings to use to connect to the Database.
	AllConnectionStrings map[string]string `mandatory:"false" json:"allConnectionStrings"`
}

func (m DatabaseConnectionStrings) String() string {
	return common.PointerString(m)
}
