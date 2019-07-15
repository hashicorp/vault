// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// Volume A detachable block volume device that allows you to dynamically expand
// the storage capacity of an instance. For more information, see
// Overview of Cloud Volume Storage (https://docs.cloud.oracle.com/Content/Block/Concepts/overview.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type Volume struct {

	// The availability domain of the volume.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment that contains the volume.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the volume.
	Id *string `mandatory:"true" json:"id"`

	// The current state of a volume.
	LifecycleState VolumeLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The size of the volume in MBs. This field is deprecated. Use sizeInGBs instead.
	SizeInMBs *int64 `mandatory:"true" json:"sizeInMBs"`

	// The date and time the volume was created. Format defined by RFC3339.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Specifies whether the cloned volume's data has finished copying from the source volume or backup.
	IsHydrated *bool `mandatory:"false" json:"isHydrated"`

	// The OCID of the KMS key which is the master encryption key for the volume.
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`

	// The size of the volume in GBs.
	SizeInGBs *int64 `mandatory:"false" json:"sizeInGBs"`

	// The volume source, either an existing volume in the same availability domain or a volume backup.
	// If null, an empty volume is created.
	SourceDetails VolumeSourceDetails `mandatory:"false" json:"sourceDetails"`

	// The OCID of the source volume group.
	VolumeGroupId *string `mandatory:"false" json:"volumeGroupId"`
}

func (m Volume) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *Volume) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		DefinedTags        map[string]map[string]interface{} `json:"definedTags"`
		FreeformTags       map[string]string                 `json:"freeformTags"`
		IsHydrated         *bool                             `json:"isHydrated"`
		KmsKeyId           *string                           `json:"kmsKeyId"`
		SizeInGBs          *int64                            `json:"sizeInGBs"`
		SourceDetails      volumesourcedetails               `json:"sourceDetails"`
		VolumeGroupId      *string                           `json:"volumeGroupId"`
		AvailabilityDomain *string                           `json:"availabilityDomain"`
		CompartmentId      *string                           `json:"compartmentId"`
		DisplayName        *string                           `json:"displayName"`
		Id                 *string                           `json:"id"`
		LifecycleState     VolumeLifecycleStateEnum          `json:"lifecycleState"`
		SizeInMBs          *int64                            `json:"sizeInMBs"`
		TimeCreated        *common.SDKTime                   `json:"timeCreated"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.DefinedTags = model.DefinedTags
	m.FreeformTags = model.FreeformTags
	m.IsHydrated = model.IsHydrated
	m.KmsKeyId = model.KmsKeyId
	m.SizeInGBs = model.SizeInGBs
	nn, e := model.SourceDetails.UnmarshalPolymorphicJSON(model.SourceDetails.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.SourceDetails = nn.(VolumeSourceDetails)
	} else {
		m.SourceDetails = nil
	}
	m.VolumeGroupId = model.VolumeGroupId
	m.AvailabilityDomain = model.AvailabilityDomain
	m.CompartmentId = model.CompartmentId
	m.DisplayName = model.DisplayName
	m.Id = model.Id
	m.LifecycleState = model.LifecycleState
	m.SizeInMBs = model.SizeInMBs
	m.TimeCreated = model.TimeCreated
	return
}

// VolumeLifecycleStateEnum Enum with underlying type: string
type VolumeLifecycleStateEnum string

// Set of constants representing the allowable values for VolumeLifecycleStateEnum
const (
	VolumeLifecycleStateProvisioning VolumeLifecycleStateEnum = "PROVISIONING"
	VolumeLifecycleStateRestoring    VolumeLifecycleStateEnum = "RESTORING"
	VolumeLifecycleStateAvailable    VolumeLifecycleStateEnum = "AVAILABLE"
	VolumeLifecycleStateTerminating  VolumeLifecycleStateEnum = "TERMINATING"
	VolumeLifecycleStateTerminated   VolumeLifecycleStateEnum = "TERMINATED"
	VolumeLifecycleStateFaulty       VolumeLifecycleStateEnum = "FAULTY"
)

var mappingVolumeLifecycleState = map[string]VolumeLifecycleStateEnum{
	"PROVISIONING": VolumeLifecycleStateProvisioning,
	"RESTORING":    VolumeLifecycleStateRestoring,
	"AVAILABLE":    VolumeLifecycleStateAvailable,
	"TERMINATING":  VolumeLifecycleStateTerminating,
	"TERMINATED":   VolumeLifecycleStateTerminated,
	"FAULTY":       VolumeLifecycleStateFaulty,
}

// GetVolumeLifecycleStateEnumValues Enumerates the set of values for VolumeLifecycleStateEnum
func GetVolumeLifecycleStateEnumValues() []VolumeLifecycleStateEnum {
	values := make([]VolumeLifecycleStateEnum, 0)
	for _, v := range mappingVolumeLifecycleState {
		values = append(values, v)
	}
	return values
}
