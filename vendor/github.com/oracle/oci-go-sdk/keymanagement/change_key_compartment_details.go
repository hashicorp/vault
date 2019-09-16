// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ChangeKeyCompartmentDetails The representation of ChangeKeyCompartmentDetails
type ChangeKeyCompartmentDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment into which the key should be moved.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeKeyCompartmentDetails) String() string {
	return common.PointerString(m)
}
