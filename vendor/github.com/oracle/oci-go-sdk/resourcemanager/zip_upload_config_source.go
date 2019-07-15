// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// ZipUploadConfigSource File path to the location of the zip file that contains the Terraform configuration.
type ZipUploadConfigSource struct {

	// File path to the directory from which Terraform runs.
	// If not specified, we use the root directory.
	WorkingDirectory *string `mandatory:"false" json:"workingDirectory"`
}

//GetWorkingDirectory returns WorkingDirectory
func (m ZipUploadConfigSource) GetWorkingDirectory() *string {
	return m.WorkingDirectory
}

func (m ZipUploadConfigSource) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ZipUploadConfigSource) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeZipUploadConfigSource ZipUploadConfigSource
	s := struct {
		DiscriminatorParam string `json:"configSourceType"`
		MarshalTypeZipUploadConfigSource
	}{
		"ZIP_UPLOAD",
		(MarshalTypeZipUploadConfigSource)(m),
	}

	return json.Marshal(&s)
}
