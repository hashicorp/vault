// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateDbHomeDetails Describes the modification parameters for the database home.
type UpdateDbHomeDetails struct {
	DbVersion *PatchDetails `mandatory:"false" json:"dbVersion"`
}

func (m UpdateDbHomeDetails) String() string {
	return common.PointerString(m)
}
