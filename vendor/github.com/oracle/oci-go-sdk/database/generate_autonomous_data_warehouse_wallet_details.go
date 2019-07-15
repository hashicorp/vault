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

// GenerateAutonomousDataWarehouseWalletDetails **Deprecated.** See GenerateAutonomousDatabaseWalletDetails for reference information about creating and downloading a wallet for an Oracle Autonomous Data Warehouse.
type GenerateAutonomousDataWarehouseWalletDetails struct {

	// The password to encrypt the keys inside the wallet. The password must be at least 8 characters long and must include at least 1 letter and either 1 numeric character or 1 special character.
	Password *string `mandatory:"true" json:"password"`
}

func (m GenerateAutonomousDataWarehouseWalletDetails) String() string {
	return common.PointerString(m)
}
