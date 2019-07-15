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

// ImageSourceViaObjectStorageUriDetails The representation of ImageSourceViaObjectStorageUriDetails
type ImageSourceViaObjectStorageUriDetails struct {

	// The Object Storage URL for the image.
	SourceUri *string `mandatory:"true" json:"sourceUri"`

	OperatingSystem *string `mandatory:"false" json:"operatingSystem"`

	OperatingSystemVersion *string `mandatory:"false" json:"operatingSystemVersion"`

	// The format of the image to be imported.  Only monolithic
	// images are supported. This attribute is not used for exported Oracle images with the OCI image format.
	SourceImageType ImageSourceDetailsSourceImageTypeEnum `mandatory:"false" json:"sourceImageType,omitempty"`
}

//GetOperatingSystem returns OperatingSystem
func (m ImageSourceViaObjectStorageUriDetails) GetOperatingSystem() *string {
	return m.OperatingSystem
}

//GetOperatingSystemVersion returns OperatingSystemVersion
func (m ImageSourceViaObjectStorageUriDetails) GetOperatingSystemVersion() *string {
	return m.OperatingSystemVersion
}

//GetSourceImageType returns SourceImageType
func (m ImageSourceViaObjectStorageUriDetails) GetSourceImageType() ImageSourceDetailsSourceImageTypeEnum {
	return m.SourceImageType
}

func (m ImageSourceViaObjectStorageUriDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ImageSourceViaObjectStorageUriDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeImageSourceViaObjectStorageUriDetails ImageSourceViaObjectStorageUriDetails
	s := struct {
		DiscriminatorParam string `json:"sourceType"`
		MarshalTypeImageSourceViaObjectStorageUriDetails
	}{
		"objectStorageUri",
		(MarshalTypeImageSourceViaObjectStorageUriDetails)(m),
	}

	return json.Marshal(&s)
}
