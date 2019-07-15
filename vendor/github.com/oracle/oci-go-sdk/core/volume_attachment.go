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

// VolumeAttachment A base object for all types of attachments between a storage volume and an instance.
// For specific details about iSCSI attachments, see
// IScsiVolumeAttachment.
// For general information about volume attachments, see
// Overview of Block Volume Storage (https://docs.cloud.oracle.com/Content/Block/Concepts/overview.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type VolumeAttachment interface {

	// The availability domain of an instance.
	// Example: `Uocm:PHX-AD-1`
	GetAvailabilityDomain() *string

	// The OCID of the compartment.
	GetCompartmentId() *string

	// The OCID of the volume attachment.
	GetId() *string

	// The OCID of the instance the volume is attached to.
	GetInstanceId() *string

	// The current state of the volume attachment.
	GetLifecycleState() VolumeAttachmentLifecycleStateEnum

	// The date and time the volume was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	GetTimeCreated() *common.SDKTime

	// The OCID of the volume.
	GetVolumeId() *string

	// The device name.
	GetDevice() *string

	// A user-friendly name. Does not have to be unique, and it cannot be changed.
	// Avoid entering confidential information.
	// Example: `My volume attachment`
	GetDisplayName() *string

	// Whether the attachment was created in read-only mode.
	GetIsReadOnly() *bool

	// Whether in-transit encryption for the data volume's paravirtualized attachment is enabled or not.
	GetIsPvEncryptionInTransitEnabled() *bool
}

type volumeattachment struct {
	JsonData                       []byte
	AvailabilityDomain             *string                            `mandatory:"true" json:"availabilityDomain"`
	CompartmentId                  *string                            `mandatory:"true" json:"compartmentId"`
	Id                             *string                            `mandatory:"true" json:"id"`
	InstanceId                     *string                            `mandatory:"true" json:"instanceId"`
	LifecycleState                 VolumeAttachmentLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`
	TimeCreated                    *common.SDKTime                    `mandatory:"true" json:"timeCreated"`
	VolumeId                       *string                            `mandatory:"true" json:"volumeId"`
	Device                         *string                            `mandatory:"false" json:"device"`
	DisplayName                    *string                            `mandatory:"false" json:"displayName"`
	IsReadOnly                     *bool                              `mandatory:"false" json:"isReadOnly"`
	IsPvEncryptionInTransitEnabled *bool                              `mandatory:"false" json:"isPvEncryptionInTransitEnabled"`
	AttachmentType                 string                             `json:"attachmentType"`
}

// UnmarshalJSON unmarshals json
func (m *volumeattachment) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalervolumeattachment volumeattachment
	s := struct {
		Model Unmarshalervolumeattachment
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.AvailabilityDomain = s.Model.AvailabilityDomain
	m.CompartmentId = s.Model.CompartmentId
	m.Id = s.Model.Id
	m.InstanceId = s.Model.InstanceId
	m.LifecycleState = s.Model.LifecycleState
	m.TimeCreated = s.Model.TimeCreated
	m.VolumeId = s.Model.VolumeId
	m.Device = s.Model.Device
	m.DisplayName = s.Model.DisplayName
	m.IsReadOnly = s.Model.IsReadOnly
	m.IsPvEncryptionInTransitEnabled = s.Model.IsPvEncryptionInTransitEnabled
	m.AttachmentType = s.Model.AttachmentType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *volumeattachment) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.AttachmentType {
	case "iscsi":
		mm := IScsiVolumeAttachment{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "emulated":
		mm := EmulatedVolumeAttachment{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "paravirtualized":
		mm := ParavirtualizedVolumeAttachment{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetAvailabilityDomain returns AvailabilityDomain
func (m volumeattachment) GetAvailabilityDomain() *string {
	return m.AvailabilityDomain
}

//GetCompartmentId returns CompartmentId
func (m volumeattachment) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetId returns Id
func (m volumeattachment) GetId() *string {
	return m.Id
}

//GetInstanceId returns InstanceId
func (m volumeattachment) GetInstanceId() *string {
	return m.InstanceId
}

//GetLifecycleState returns LifecycleState
func (m volumeattachment) GetLifecycleState() VolumeAttachmentLifecycleStateEnum {
	return m.LifecycleState
}

//GetTimeCreated returns TimeCreated
func (m volumeattachment) GetTimeCreated() *common.SDKTime {
	return m.TimeCreated
}

//GetVolumeId returns VolumeId
func (m volumeattachment) GetVolumeId() *string {
	return m.VolumeId
}

//GetDevice returns Device
func (m volumeattachment) GetDevice() *string {
	return m.Device
}

//GetDisplayName returns DisplayName
func (m volumeattachment) GetDisplayName() *string {
	return m.DisplayName
}

//GetIsReadOnly returns IsReadOnly
func (m volumeattachment) GetIsReadOnly() *bool {
	return m.IsReadOnly
}

//GetIsPvEncryptionInTransitEnabled returns IsPvEncryptionInTransitEnabled
func (m volumeattachment) GetIsPvEncryptionInTransitEnabled() *bool {
	return m.IsPvEncryptionInTransitEnabled
}

func (m volumeattachment) String() string {
	return common.PointerString(m)
}

// VolumeAttachmentLifecycleStateEnum Enum with underlying type: string
type VolumeAttachmentLifecycleStateEnum string

// Set of constants representing the allowable values for VolumeAttachmentLifecycleStateEnum
const (
	VolumeAttachmentLifecycleStateAttaching VolumeAttachmentLifecycleStateEnum = "ATTACHING"
	VolumeAttachmentLifecycleStateAttached  VolumeAttachmentLifecycleStateEnum = "ATTACHED"
	VolumeAttachmentLifecycleStateDetaching VolumeAttachmentLifecycleStateEnum = "DETACHING"
	VolumeAttachmentLifecycleStateDetached  VolumeAttachmentLifecycleStateEnum = "DETACHED"
)

var mappingVolumeAttachmentLifecycleState = map[string]VolumeAttachmentLifecycleStateEnum{
	"ATTACHING": VolumeAttachmentLifecycleStateAttaching,
	"ATTACHED":  VolumeAttachmentLifecycleStateAttached,
	"DETACHING": VolumeAttachmentLifecycleStateDetaching,
	"DETACHED":  VolumeAttachmentLifecycleStateDetached,
}

// GetVolumeAttachmentLifecycleStateEnumValues Enumerates the set of values for VolumeAttachmentLifecycleStateEnum
func GetVolumeAttachmentLifecycleStateEnumValues() []VolumeAttachmentLifecycleStateEnum {
	values := make([]VolumeAttachmentLifecycleStateEnum, 0)
	for _, v := range mappingVolumeAttachmentLifecycleState {
		values = append(values, v)
	}
	return values
}
