// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateDbHomeWithDbSystemIdFromBackupDetails The representation of CreateDbHomeWithDbSystemIdFromBackupDetails
type CreateDbHomeWithDbSystemIdFromBackupDetails struct {
	Database *CreateDatabaseFromBackupDetails `mandatory:"true" json:"database"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system.
	DbSystemId *string `mandatory:"false" json:"dbSystemId"`

	// The user-provided name of the database home.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

//GetDbSystemId returns DbSystemId
func (m CreateDbHomeWithDbSystemIdFromBackupDetails) GetDbSystemId() *string {
	return m.DbSystemId
}

//GetDisplayName returns DisplayName
func (m CreateDbHomeWithDbSystemIdFromBackupDetails) GetDisplayName() *string {
	return m.DisplayName
}

func (m CreateDbHomeWithDbSystemIdFromBackupDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m CreateDbHomeWithDbSystemIdFromBackupDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeCreateDbHomeWithDbSystemIdFromBackupDetails CreateDbHomeWithDbSystemIdFromBackupDetails
	s := struct {
		DiscriminatorParam string `json:"source"`
		MarshalTypeCreateDbHomeWithDbSystemIdFromBackupDetails
	}{
		"DB_BACKUP",
		(MarshalTypeCreateDbHomeWithDbSystemIdFromBackupDetails)(m),
	}

	return json.Marshal(&s)
}
