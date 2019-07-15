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

// DbIormConfig IORM Config setting response for this database
type DbIormConfig struct {

	// Database Name. For default DbPlan, the dbName will always be `default`
	DbName *string `mandatory:"false" json:"dbName"`

	// Relative priority of a database
	Share *int `mandatory:"false" json:"share"`

	// Flash Cache limit, internally configured based on shares
	FlashCacheLimit *string `mandatory:"false" json:"flashCacheLimit"`
}

func (m DbIormConfig) String() string {
	return common.PointerString(m)
}
