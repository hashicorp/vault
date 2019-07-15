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

// AttachServiceDeterminedVolumeDetails The representation of AttachServiceDeterminedVolumeDetails
type AttachServiceDeterminedVolumeDetails struct {

	// The OCID of the instance.
	InstanceId *string `mandatory:"true" json:"instanceId"`

	// The OCID of the volume.
	VolumeId *string `mandatory:"true" json:"volumeId"`

	// The device name.
	Device *string `mandatory:"false" json:"device"`

	// A user-friendly name. Does not have to be unique, and it cannot be changed. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Whether the attachment was created in read-only mode.
	IsReadOnly *bool `mandatory:"false" json:"isReadOnly"`
}

//GetDevice returns Device
func (m AttachServiceDeterminedVolumeDetails) GetDevice() *string {
	return m.Device
}

//GetDisplayName returns DisplayName
func (m AttachServiceDeterminedVolumeDetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetInstanceId returns InstanceId
func (m AttachServiceDeterminedVolumeDetails) GetInstanceId() *string {
	return m.InstanceId
}

//GetIsReadOnly returns IsReadOnly
func (m AttachServiceDeterminedVolumeDetails) GetIsReadOnly() *bool {
	return m.IsReadOnly
}

//GetVolumeId returns VolumeId
func (m AttachServiceDeterminedVolumeDetails) GetVolumeId() *string {
	return m.VolumeId
}

func (m AttachServiceDeterminedVolumeDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m AttachServiceDeterminedVolumeDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeAttachServiceDeterminedVolumeDetails AttachServiceDeterminedVolumeDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeAttachServiceDeterminedVolumeDetails
	}{
		"service_determined",
		(MarshalTypeAttachServiceDeterminedVolumeDetails)(m),
	}

	return json.Marshal(&s)
}
