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

// UpdateAutonomousDataWarehouseDetails **Deprecated.** See UpdateAutonomousDatabaseDetails for reference information about updating an Autonomous Data Warehouse.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type UpdateAutonomousDataWarehouseDetails struct {

	// The number of CPU cores to be made available to the database.
	CpuCoreCount *int `mandatory:"false" json:"cpuCoreCount"`

	// Size, in terabytes, of the data volume that will be attached to the database.
	DataStorageSizeInTBs *int `mandatory:"false" json:"dataStorageSizeInTBs"`

	// The user-friendly name for the Autonomous Data Warehouse. The name does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The password must be between 12 and 30 characters long, and must contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing. It must be different from the last four passwords and it must not be a password used within the last 24 hours.
	AdminPassword *string `mandatory:"false" json:"adminPassword"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m UpdateAutonomousDataWarehouseDetails) String() string {
	return common.PointerString(m)
}
