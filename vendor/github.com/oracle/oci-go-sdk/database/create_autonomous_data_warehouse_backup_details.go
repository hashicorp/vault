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

// CreateAutonomousDataWarehouseBackupDetails **Deprecated.** See CreateAutonomousDatabaseBackupDetails for reference information about creating Autonomous Data Warehouse backups.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateAutonomousDataWarehouseBackupDetails struct {

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Data Warehouse backup.
	AutonomousDataWarehouseId *string `mandatory:"true" json:"autonomousDataWarehouseId"`
}

func (m CreateAutonomousDataWarehouseBackupDetails) String() string {
	return common.PointerString(m)
}
