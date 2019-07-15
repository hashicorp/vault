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

// InstanceConfigurationIscsiAttachVolumeDetails The representation of InstanceConfigurationIscsiAttachVolumeDetails
type InstanceConfigurationIscsiAttachVolumeDetails struct {

	// A user-friendly name. Does not have to be unique, and it cannot be changed. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Whether the attachment should be created in read-only mode.
	IsReadOnly *bool `mandatory:"false" json:"isReadOnly"`

	// Whether to use CHAP authentication for the volume attachment. Defaults to false.
	UseChap *bool `mandatory:"false" json:"useChap"`
}

//GetDisplayName returns DisplayName
func (m InstanceConfigurationIscsiAttachVolumeDetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetIsReadOnly returns IsReadOnly
func (m InstanceConfigurationIscsiAttachVolumeDetails) GetIsReadOnly() *bool {
	return m.IsReadOnly
}

func (m InstanceConfigurationIscsiAttachVolumeDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m InstanceConfigurationIscsiAttachVolumeDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeInstanceConfigurationIscsiAttachVolumeDetails InstanceConfigurationIscsiAttachVolumeDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeInstanceConfigurationIscsiAttachVolumeDetails
	}{
		"iscsi",
		(MarshalTypeInstanceConfigurationIscsiAttachVolumeDetails)(m),
	}

	return json.Marshal(&s)
}
