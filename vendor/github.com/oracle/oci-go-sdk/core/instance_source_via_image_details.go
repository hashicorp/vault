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

// InstanceSourceViaImageDetails The representation of InstanceSourceViaImageDetails
type InstanceSourceViaImageDetails struct {

	// The OCID of the image used to boot the instance.
	ImageId *string `mandatory:"true" json:"imageId"`

	// The size of the boot volume in GBs. Minimum value is 50 GB and maximum value is 16384 GB (16TB).
	BootVolumeSizeInGBs *int64 `mandatory:"false" json:"bootVolumeSizeInGBs"`

	// The OCID of the KMS key to be used as the master encryption key for the boot volume.
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`
}

func (m InstanceSourceViaImageDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m InstanceSourceViaImageDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeInstanceSourceViaImageDetails InstanceSourceViaImageDetails
	s := struct {
		DiscriminatorParam string `json:"sourceType"`
		MarshalTypeInstanceSourceViaImageDetails
	}{
		"image",
		(MarshalTypeInstanceSourceViaImageDetails)(m),
	}

	return json.Marshal(&s)
}
