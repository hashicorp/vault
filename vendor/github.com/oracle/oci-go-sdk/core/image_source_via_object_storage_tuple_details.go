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

// ImageSourceViaObjectStorageTupleDetails The representation of ImageSourceViaObjectStorageTupleDetails
type ImageSourceViaObjectStorageTupleDetails struct {

	// The Object Storage bucket for the image.
	BucketName *string `mandatory:"true" json:"bucketName"`

	// The Object Storage namespace for the image.
	NamespaceName *string `mandatory:"true" json:"namespaceName"`

	// The Object Storage name for the image.
	ObjectName *string `mandatory:"true" json:"objectName"`

	OperatingSystem *string `mandatory:"false" json:"operatingSystem"`

	OperatingSystemVersion *string `mandatory:"false" json:"operatingSystemVersion"`

	// The format of the image to be imported.  Only monolithic
	// images are supported. This attribute is not used for exported Oracle images with the OCI image format.
	SourceImageType ImageSourceDetailsSourceImageTypeEnum `mandatory:"false" json:"sourceImageType,omitempty"`
}

//GetOperatingSystem returns OperatingSystem
func (m ImageSourceViaObjectStorageTupleDetails) GetOperatingSystem() *string {
	return m.OperatingSystem
}

//GetOperatingSystemVersion returns OperatingSystemVersion
func (m ImageSourceViaObjectStorageTupleDetails) GetOperatingSystemVersion() *string {
	return m.OperatingSystemVersion
}

//GetSourceImageType returns SourceImageType
func (m ImageSourceViaObjectStorageTupleDetails) GetSourceImageType() ImageSourceDetailsSourceImageTypeEnum {
	return m.SourceImageType
}

func (m ImageSourceViaObjectStorageTupleDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ImageSourceViaObjectStorageTupleDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeImageSourceViaObjectStorageTupleDetails ImageSourceViaObjectStorageTupleDetails
	s := struct {
		DiscriminatorParam string `json:"sourceType"`
		MarshalTypeImageSourceViaObjectStorageTupleDetails
	}{
		"objectStorageTuple",
		(MarshalTypeImageSourceViaObjectStorageTupleDetails)(m),
	}

	return json.Marshal(&s)
}
