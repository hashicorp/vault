// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Vault Service Key Management API
//
// API for managing and performing operations with keys and vaults. (For the API for managing secrets, see the Vault Service
// Secret Management API. For the API for retrieving secrets, see the Vault Service Secret Retrieval API.)
//

package keymanagement

import (
	"encoding/json"
	"fmt"
	"github.com/oracle/oci-go-sdk/v60/common"
	"strings"
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

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m backuplocation) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// BackupLocationDestinationEnum Enum with underlying type: string
type BackupLocationDestinationEnum string

// Set of constants representing the allowable values for BackupLocationDestinationEnum
const (
	BackupLocationDestinationBucket                     BackupLocationDestinationEnum = "BUCKET"
	BackupLocationDestinationPreAuthenticatedRequestUri BackupLocationDestinationEnum = "PRE_AUTHENTICATED_REQUEST_URI"
)

var mappingBackupLocationDestinationEnum = map[string]BackupLocationDestinationEnum{
	"BUCKET":                        BackupLocationDestinationBucket,
	"PRE_AUTHENTICATED_REQUEST_URI": BackupLocationDestinationPreAuthenticatedRequestUri,
}

var mappingBackupLocationDestinationEnumLowerCase = map[string]BackupLocationDestinationEnum{
	"bucket":                        BackupLocationDestinationBucket,
	"pre_authenticated_request_uri": BackupLocationDestinationPreAuthenticatedRequestUri,
}

// GetBackupLocationDestinationEnumValues Enumerates the set of values for BackupLocationDestinationEnum
func GetBackupLocationDestinationEnumValues() []BackupLocationDestinationEnum {
	values := make([]BackupLocationDestinationEnum, 0)
	for _, v := range mappingBackupLocationDestinationEnum {
		values = append(values, v)
	}
	return values
}

// GetBackupLocationDestinationEnumStringValues Enumerates the set of values in String for BackupLocationDestinationEnum
func GetBackupLocationDestinationEnumStringValues() []string {
	return []string{
		"BUCKET",
		"PRE_AUTHENTICATED_REQUEST_URI",
	}
}

// GetMappingBackupLocationDestinationEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingBackupLocationDestinationEnum(val string) (BackupLocationDestinationEnum, bool) {
	enum, ok := mappingBackupLocationDestinationEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
