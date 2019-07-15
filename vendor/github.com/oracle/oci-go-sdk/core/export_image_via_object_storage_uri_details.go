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

// ExportImageViaObjectStorageUriDetails The representation of ExportImageViaObjectStorageUriDetails
type ExportImageViaObjectStorageUriDetails struct {

	// The Object Storage URL to export the image to. See Object Storage URLs (https://docs.cloud.oracle.com/Content/Compute/Tasks/imageimportexport.htm#URLs)
	// and Using Pre-Authenticated Requests (https://docs.cloud.oracle.com/Content/Object/Tasks/usingpreauthenticatedrequests.htm) for constructing URLs for image import/export.
	DestinationUri *string `mandatory:"true" json:"destinationUri"`
}

func (m ExportImageViaObjectStorageUriDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ExportImageViaObjectStorageUriDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeExportImageViaObjectStorageUriDetails ExportImageViaObjectStorageUriDetails
	s := struct {
		DiscriminatorParam string `json:"destinationType"`
		MarshalTypeExportImageViaObjectStorageUriDetails
	}{
		"objectStorageUri",
		(MarshalTypeExportImageViaObjectStorageUriDetails)(m),
	}

	return json.Marshal(&s)
}
