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

// UpdateBootVolumeKmsKeyDetails The representation of UpdateBootVolumeKmsKeyDetails
type UpdateBootVolumeKmsKeyDetails struct {

	// The OCID of the new KMS key which will be used to protect the specified volume.
	// This key has to be a valid KMS key OCID, and the user must have key delegation policy to allow them to access this key.
	// Even if the new KMS key is the same as the previous KMS key ID, the Block Volume service will use it to regenerate a new volume encryption key.
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`
}

func (m UpdateBootVolumeKmsKeyDetails) String() string {
	return common.PointerString(m)
}
