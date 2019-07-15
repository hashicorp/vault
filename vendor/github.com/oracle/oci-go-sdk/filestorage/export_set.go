// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// File Storage Service API
//
// The API for the File Storage Service.
//

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ExportSet A set of file systems to export through one or more mount
// targets. Composed of zero or more export resources.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type ExportSet struct {

	// The OCID of the compartment that contains the export set.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	// Example: `My export set`
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the export set.
	Id *string `mandatory:"true" json:"id"`

	// The current state of the export set.
	LifecycleState ExportSetLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the export set was created, expressed
	// in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the virtual cloud network (VCN) the export set is in.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// The availability domain the export set is in. May be unset
	// as a blank or NULL value.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// Controls the maximum `tbytes`, `fbytes`, and `abytes`,
	// values reported by `NFS FSSTAT` calls through any associated
	// mount targets. This is an advanced feature. For most
	// applications, use the default value. The
	// `tbytes` value reported by `FSSTAT` will be
	// `maxFsStatBytes`. The value of `fbytes` and `abytes` will be
	// `maxFsStatBytes` minus the metered size of the file
	// system. If the metered size is larger than `maxFsStatBytes`,
	// then `fbytes` and `abytes` will both be '0'.
	MaxFsStatBytes *int64 `mandatory:"false" json:"maxFsStatBytes"`

	// Controls the maximum `tfiles`, `ffiles`, and `afiles`
	// values reported by `NFS FSSTAT` calls through any associated
	// mount targets. This is an advanced feature. For most
	// applications, use the default value. The
	// `tfiles` value reported by `FSSTAT` will be
	// `maxFsStatFiles`. The value of `ffiles` and `afiles` will be
	// `maxFsStatFiles` minus the metered size of the file
	// system. If the metered size is larger than `maxFsStatFiles`,
	// then `ffiles` and `afiles` will both be '0'.
	MaxFsStatFiles *int64 `mandatory:"false" json:"maxFsStatFiles"`
}

func (m ExportSet) String() string {
	return common.PointerString(m)
}

// ExportSetLifecycleStateEnum Enum with underlying type: string
type ExportSetLifecycleStateEnum string

// Set of constants representing the allowable values for ExportSetLifecycleStateEnum
const (
	ExportSetLifecycleStateCreating ExportSetLifecycleStateEnum = "CREATING"
	ExportSetLifecycleStateActive   ExportSetLifecycleStateEnum = "ACTIVE"
	ExportSetLifecycleStateDeleting ExportSetLifecycleStateEnum = "DELETING"
	ExportSetLifecycleStateDeleted  ExportSetLifecycleStateEnum = "DELETED"
)

var mappingExportSetLifecycleState = map[string]ExportSetLifecycleStateEnum{
	"CREATING": ExportSetLifecycleStateCreating,
	"ACTIVE":   ExportSetLifecycleStateActive,
	"DELETING": ExportSetLifecycleStateDeleting,
	"DELETED":  ExportSetLifecycleStateDeleted,
}

// GetExportSetLifecycleStateEnumValues Enumerates the set of values for ExportSetLifecycleStateEnum
func GetExportSetLifecycleStateEnumValues() []ExportSetLifecycleStateEnum {
	values := make([]ExportSetLifecycleStateEnum, 0)
	for _, v := range mappingExportSetLifecycleState {
		values = append(values, v)
	}
	return values
}
