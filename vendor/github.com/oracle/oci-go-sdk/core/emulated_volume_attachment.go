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

// EmulatedVolumeAttachment An Emulated volume attachment.
type EmulatedVolumeAttachment struct {

	// The availability domain of an instance.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the volume attachment.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the instance the volume is attached to.
	InstanceId *string `mandatory:"true" json:"instanceId"`

	// The date and time the volume was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the volume.
	VolumeId *string `mandatory:"true" json:"volumeId"`

	// The device name.
	Device *string `mandatory:"false" json:"device"`

	// A user-friendly name. Does not have to be unique, and it cannot be changed.
	// Avoid entering confidential information.
	// Example: `My volume attachment`
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Whether the attachment was created in read-only mode.
	IsReadOnly *bool `mandatory:"false" json:"isReadOnly"`

	// Whether in-transit encryption for the data volume's paravirtualized attachment is enabled or not.
	IsPvEncryptionInTransitEnabled *bool `mandatory:"false" json:"isPvEncryptionInTransitEnabled"`

	// The current state of the volume attachment.
	LifecycleState VolumeAttachmentLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`
}

//GetAvailabilityDomain returns AvailabilityDomain
func (m EmulatedVolumeAttachment) GetAvailabilityDomain() *string {
	return m.AvailabilityDomain
}

//GetCompartmentId returns CompartmentId
func (m EmulatedVolumeAttachment) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetDevice returns Device
func (m EmulatedVolumeAttachment) GetDevice() *string {
	return m.Device
}

//GetDisplayName returns DisplayName
func (m EmulatedVolumeAttachment) GetDisplayName() *string {
	return m.DisplayName
}

//GetId returns Id
func (m EmulatedVolumeAttachment) GetId() *string {
	return m.Id
}

//GetInstanceId returns InstanceId
func (m EmulatedVolumeAttachment) GetInstanceId() *string {
	return m.InstanceId
}

//GetIsReadOnly returns IsReadOnly
func (m EmulatedVolumeAttachment) GetIsReadOnly() *bool {
	return m.IsReadOnly
}

//GetLifecycleState returns LifecycleState
func (m EmulatedVolumeAttachment) GetLifecycleState() VolumeAttachmentLifecycleStateEnum {
	return m.LifecycleState
}

//GetTimeCreated returns TimeCreated
func (m EmulatedVolumeAttachment) GetTimeCreated() *common.SDKTime {
	return m.TimeCreated
}

//GetVolumeId returns VolumeId
func (m EmulatedVolumeAttachment) GetVolumeId() *string {
	return m.VolumeId
}

//GetIsPvEncryptionInTransitEnabled returns IsPvEncryptionInTransitEnabled
func (m EmulatedVolumeAttachment) GetIsPvEncryptionInTransitEnabled() *bool {
	return m.IsPvEncryptionInTransitEnabled
}

func (m EmulatedVolumeAttachment) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m EmulatedVolumeAttachment) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeEmulatedVolumeAttachment EmulatedVolumeAttachment
	s := struct {
		DiscriminatorParam string `json:"attachmentType"`
		MarshalTypeEmulatedVolumeAttachment
	}{
		"emulated",
		(MarshalTypeEmulatedVolumeAttachment)(m),
	}

	return json.Marshal(&s)
}
