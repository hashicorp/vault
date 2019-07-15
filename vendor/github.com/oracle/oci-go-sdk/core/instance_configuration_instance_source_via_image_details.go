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

// InstanceConfigurationInstanceSourceViaImageDetails The representation of InstanceConfigurationInstanceSourceViaImageDetails
type InstanceConfigurationInstanceSourceViaImageDetails struct {

	// The size of the boot volume in GBs. The minimum value is 50 GB and the maximum value is 16384 GB (16TB).
	BootVolumeSizeInGBs *int64 `mandatory:"false" json:"bootVolumeSizeInGBs"`

	// The OCID of the image used to boot the instance.
	ImageId *string `mandatory:"false" json:"imageId"`
}

func (m InstanceConfigurationInstanceSourceViaImageDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m InstanceConfigurationInstanceSourceViaImageDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeInstanceConfigurationInstanceSourceViaImageDetails InstanceConfigurationInstanceSourceViaImageDetails
	s := struct {
		DiscriminatorParam string `json:"sourceType"`
		MarshalTypeInstanceConfigurationInstanceSourceViaImageDetails
	}{
		"image",
		(MarshalTypeInstanceConfigurationInstanceSourceViaImageDetails)(m),
	}

	return json.Marshal(&s)
}
