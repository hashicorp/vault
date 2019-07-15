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

// DbHome The representation of DbHome
type DbHome struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the database home.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-provided name for the database home. The name does not need to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The current state of the database home.
	LifecycleState DbHomeLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The Oracle Database version.
	DbVersion *string `mandatory:"true" json:"dbVersion"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the last patch history. This value is updated as soon as a patch operation is started.
	LastPatchHistoryEntryId *string `mandatory:"false" json:"lastPatchHistoryEntryId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system.
	DbSystemId *string `mandatory:"false" json:"dbSystemId"`

	// The date and time the database home was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m DbHome) String() string {
	return common.PointerString(m)
}

// DbHomeLifecycleStateEnum Enum with underlying type: string
type DbHomeLifecycleStateEnum string

// Set of constants representing the allowable values for DbHomeLifecycleStateEnum
const (
	DbHomeLifecycleStateProvisioning DbHomeLifecycleStateEnum = "PROVISIONING"
	DbHomeLifecycleStateAvailable    DbHomeLifecycleStateEnum = "AVAILABLE"
	DbHomeLifecycleStateUpdating     DbHomeLifecycleStateEnum = "UPDATING"
	DbHomeLifecycleStateTerminating  DbHomeLifecycleStateEnum = "TERMINATING"
	DbHomeLifecycleStateTerminated   DbHomeLifecycleStateEnum = "TERMINATED"
	DbHomeLifecycleStateFailed       DbHomeLifecycleStateEnum = "FAILED"
)

var mappingDbHomeLifecycleState = map[string]DbHomeLifecycleStateEnum{
	"PROVISIONING": DbHomeLifecycleStateProvisioning,
	"AVAILABLE":    DbHomeLifecycleStateAvailable,
	"UPDATING":     DbHomeLifecycleStateUpdating,
	"TERMINATING":  DbHomeLifecycleStateTerminating,
	"TERMINATED":   DbHomeLifecycleStateTerminated,
	"FAILED":       DbHomeLifecycleStateFailed,
}

// GetDbHomeLifecycleStateEnumValues Enumerates the set of values for DbHomeLifecycleStateEnum
func GetDbHomeLifecycleStateEnumValues() []DbHomeLifecycleStateEnum {
	values := make([]DbHomeLifecycleStateEnum, 0)
	for _, v := range mappingDbHomeLifecycleState {
		values = append(values, v)
	}
	return values
}
