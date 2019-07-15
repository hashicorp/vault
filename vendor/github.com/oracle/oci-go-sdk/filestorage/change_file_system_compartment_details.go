// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// File Storage Service API
//
// The API for the File Storage Service.
//

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ChangeFileSystemCompartmentDetails Details for changing the compartment.
type ChangeFileSystemCompartmentDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment to move the file system to.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeFileSystemCompartmentDetails) String() string {
	return common.PointerString(m)
}
