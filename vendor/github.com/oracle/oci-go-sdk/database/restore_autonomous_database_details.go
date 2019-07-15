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

// RestoreAutonomousDatabaseDetails Details to restore an Oracle Autonomous Database.
type RestoreAutonomousDatabaseDetails struct {

	// The time to restore the database to.
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`

	// Restores using the backup with the System Change Number (SCN) specified.
	DatabaseSCN *string `mandatory:"false" json:"databaseSCN"`

	// Restores to the last known good state with the least possible data loss.
	Latest *bool `mandatory:"false" json:"latest"`
}

func (m RestoreAutonomousDatabaseDetails) String() string {
	return common.PointerString(m)
}
