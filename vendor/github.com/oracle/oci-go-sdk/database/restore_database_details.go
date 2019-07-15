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

// RestoreDatabaseDetails The representation of RestoreDatabaseDetails
type RestoreDatabaseDetails struct {

	// Restores using the backup with the System Change Number (SCN) specified.
	DatabaseSCN *string `mandatory:"false" json:"databaseSCN"`

	// Restores to the timestamp specified.
	Timestamp *common.SDKTime `mandatory:"false" json:"timestamp"`

	// Restores to the last known good state with the least possible data loss.
	Latest *bool `mandatory:"false" json:"latest"`
}

func (m RestoreDatabaseDetails) String() string {
	return common.PointerString(m)
}
