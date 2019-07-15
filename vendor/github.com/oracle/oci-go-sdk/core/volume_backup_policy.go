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
	"github.com/oracle/oci-go-sdk/common"
)

// VolumeBackupPolicy A policy for automatically creating volume backups according to a
// recurring schedule. Has a set of one or more schedules that control when and
// how backups are created.
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type VolumeBackupPolicy struct {

	// A user-friendly name for the volume backup policy. Does not have to be unique and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the volume backup policy.
	Id *string `mandatory:"true" json:"id"`

	// The collection of schedules that this policy will apply.
	Schedules []VolumeBackupSchedule `mandatory:"true" json:"schedules"`

	// The date and time the volume backup policy was created. Format defined by RFC3339.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`
}

func (m VolumeBackupPolicy) String() string {
	return common.PointerString(m)
}
