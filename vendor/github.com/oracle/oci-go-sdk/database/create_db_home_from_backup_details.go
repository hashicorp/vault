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

// CreateDbHomeFromBackupDetails Details for creating a database home if you are creating a database by restoring from a database backup.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateDbHomeFromBackupDetails struct {
	Database *CreateDatabaseFromBackupDetails `mandatory:"true" json:"database"`

	// The user-provided name of the database home.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

func (m CreateDbHomeFromBackupDetails) String() string {
	return common.PointerString(m)
}
