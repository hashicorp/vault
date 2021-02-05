// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// BackupLocation Backup upload location
type BackupLocation interface {
}

type backuplocation struct {
	JsonData    []byte
	Destination string `json:"destination"`
}

// UnmarshalJSON unmarshals json
func (m *backuplocation) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerbackuplocation backuplocation
	s := struct {
		Model Unmarshalerbackuplocation
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.Destination = s.Model.Destination

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *backuplocation) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Destination {
	case "BUCKET":
		mm := BackupLocationBucket{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "PRE_AUTHENTICATED_REQUEST_URI":
		mm := BackupLocationUri{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

func (m backuplocation) String() string {
	return common.PointerString(m)
}

// BackupLocationDestinationEnum Enum with underlying type: string
type BackupLocationDestinationEnum string

// Set of constants representing the allowable values for BackupLocationDestinationEnum
const (
	BackupLocationDestinationBucket                     BackupLocationDestinationEnum = "BUCKET"
	BackupLocationDestinationPreAuthenticatedRequestUri BackupLocationDestinationEnum = "PRE_AUTHENTICATED_REQUEST_URI"
)

var mappingBackupLocationDestination = map[string]BackupLocationDestinationEnum{
	"BUCKET":                        BackupLocationDestinationBucket,
	"PRE_AUTHENTICATED_REQUEST_URI": BackupLocationDestinationPreAuthenticatedRequestUri,
}

// GetBackupLocationDestinationEnumValues Enumerates the set of values for BackupLocationDestinationEnum
func GetBackupLocationDestinationEnumValues() []BackupLocationDestinationEnum {
	values := make([]BackupLocationDestinationEnum, 0)
	for _, v := range mappingBackupLocationDestination {
		values = append(values, v)
	}
	return values
}
