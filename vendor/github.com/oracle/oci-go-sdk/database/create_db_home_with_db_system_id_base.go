// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateDbHomeWithDbSystemIdBase Details for creating a database home.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateDbHomeWithDbSystemIdBase interface {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system.
	GetDbSystemId() *string

	// The user-provided name of the database home.
	GetDisplayName() *string
}

type createdbhomewithdbsystemidbase struct {
	JsonData    []byte
	DbSystemId  *string `mandatory:"false" json:"dbSystemId"`
	DisplayName *string `mandatory:"false" json:"displayName"`
	Source      string  `json:"source"`
}

// UnmarshalJSON unmarshals json
func (m *createdbhomewithdbsystemidbase) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreatedbhomewithdbsystemidbase createdbhomewithdbsystemidbase
	s := struct {
		Model Unmarshalercreatedbhomewithdbsystemidbase
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.DbSystemId = s.Model.DbSystemId
	m.DisplayName = s.Model.DisplayName
	m.Source = s.Model.Source

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createdbhomewithdbsystemidbase) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Source {
	case "DB_BACKUP":
		mm := CreateDbHomeWithDbSystemIdFromBackupDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "NONE":
		mm := CreateDbHomeWithDbSystemIdDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetDbSystemId returns DbSystemId
func (m createdbhomewithdbsystemidbase) GetDbSystemId() *string {
	return m.DbSystemId
}

//GetDisplayName returns DisplayName
func (m createdbhomewithdbsystemidbase) GetDisplayName() *string {
	return m.DisplayName
}

func (m createdbhomewithdbsystemidbase) String() string {
	return common.PointerString(m)
}
