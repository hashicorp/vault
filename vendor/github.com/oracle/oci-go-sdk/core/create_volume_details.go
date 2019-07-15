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

// CreateVolumeDetails The representation of CreateVolumeDetails
type CreateVolumeDetails struct {

	// The availability domain of the volume.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment that contains the volume.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// If provided, specifies the ID of the volume backup policy to assign to the newly
	// created volume. If omitted, no policy will be assigned.
	BackupPolicyId *string `mandatory:"false" json:"backupPolicyId"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The OCID of the KMS key to be used as the master encryption key for the volume.
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`

	// The size of the volume in GBs.
	SizeInGBs *int64 `mandatory:"false" json:"sizeInGBs"`

	// The size of the volume in MBs. The value must be a multiple of 1024.
	// This field is deprecated. Use sizeInGBs instead.
	SizeInMBs *int64 `mandatory:"false" json:"sizeInMBs"`

	// Specifies the volume source details for a new Block volume. The volume source is either another Block volume in the same availability domain or a Block volume backup.
	// This is an optional field. If not specified or set to null, the new Block volume will be empty.
	// When specified, the new Block volume will contain data from the source volume or backup.
	SourceDetails VolumeSourceDetails `mandatory:"false" json:"sourceDetails"`

	// The OCID of the volume backup from which the data should be restored on the newly created volume.
	// This field is deprecated. Use the sourceDetails field instead to specify the
	// backup for the volume.
	VolumeBackupId *string `mandatory:"false" json:"volumeBackupId"`
}

func (m CreateVolumeDetails) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *CreateVolumeDetails) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		BackupPolicyId     *string                           `json:"backupPolicyId"`
		DefinedTags        map[string]map[string]interface{} `json:"definedTags"`
		DisplayName        *string                           `json:"displayName"`
		FreeformTags       map[string]string                 `json:"freeformTags"`
		KmsKeyId           *string                           `json:"kmsKeyId"`
		SizeInGBs          *int64                            `json:"sizeInGBs"`
		SizeInMBs          *int64                            `json:"sizeInMBs"`
		SourceDetails      volumesourcedetails               `json:"sourceDetails"`
		VolumeBackupId     *string                           `json:"volumeBackupId"`
		AvailabilityDomain *string                           `json:"availabilityDomain"`
		CompartmentId      *string                           `json:"compartmentId"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.BackupPolicyId = model.BackupPolicyId
	m.DefinedTags = model.DefinedTags
	m.DisplayName = model.DisplayName
	m.FreeformTags = model.FreeformTags
	m.KmsKeyId = model.KmsKeyId
	m.SizeInGBs = model.SizeInGBs
	m.SizeInMBs = model.SizeInMBs
	nn, e := model.SourceDetails.UnmarshalPolymorphicJSON(model.SourceDetails.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.SourceDetails = nn.(VolumeSourceDetails)
	} else {
		m.SourceDetails = nil
	}
	m.VolumeBackupId = model.VolumeBackupId
	m.AvailabilityDomain = model.AvailabilityDomain
	m.CompartmentId = model.CompartmentId
	return
}
