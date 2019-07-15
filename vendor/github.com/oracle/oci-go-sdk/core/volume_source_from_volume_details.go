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

// VolumeSourceFromVolumeDetails Specifies the source volume.
type VolumeSourceFromVolumeDetails struct {

	// The OCID of the volume.
	Id *string `mandatory:"true" json:"id"`
}

func (m VolumeSourceFromVolumeDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m VolumeSourceFromVolumeDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeVolumeSourceFromVolumeDetails VolumeSourceFromVolumeDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeVolumeSourceFromVolumeDetails
	}{
		"volume",
		(MarshalTypeVolumeSourceFromVolumeDetails)(m),
	}

	return json.Marshal(&s)
}
