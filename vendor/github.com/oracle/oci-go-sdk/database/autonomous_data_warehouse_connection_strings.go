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

// AutonomousDataWarehouseConnectionStrings **Deprecated.** For information about connection strings to connect to an Oracle Autonomous Data Warehouse, see AutonomousDatabaseConnectionStrings.
type AutonomousDataWarehouseConnectionStrings struct {

	// The High database service provides the highest level of resources to each SQL statement resulting in the highest performance, but supports the fewest number of concurrent SQL statements.
	High *string `mandatory:"false" json:"high"`

	// The Medium database service provides a lower level of resources to each SQL statement potentially resulting a lower level of performance, but supports more concurrent SQL statements.
	Medium *string `mandatory:"false" json:"medium"`

	// The Low database service provides the least level of resources to each SQL statement, but supports the most number of concurrent SQL statements.
	Low *string `mandatory:"false" json:"low"`

	// Returns all connection strings that can be used to connect to the Autonomous Data Warehouse.
	// For more information, please see Predefined Database Service Names for Autonomous Transaction Processing (https://docs.oracle.com/en/cloud/paas/atp-cloud/atpug/connect-predefined.html#GUID-9747539B-FD46-44F1-8FF8-F5AC650F15BE)
	AllConnectionStrings map[string]string `mandatory:"false" json:"allConnectionStrings"`
}

func (m AutonomousDataWarehouseConnectionStrings) String() string {
	return common.PointerString(m)
}
