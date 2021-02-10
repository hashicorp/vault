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

// BackupKeyDetails The representation of BackupKeyDetails
type BackupKeyDetails struct {
	BackupLocation BackupLocation `mandatory:"false" json:"backupLocation"`
}

func (m BackupKeyDetails) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *BackupKeyDetails) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		BackupLocation backuplocation `json:"backupLocation"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	var nn interface{}
	nn, e = model.BackupLocation.UnmarshalPolymorphicJSON(model.BackupLocation.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.BackupLocation = nn.(BackupLocation)
	} else {
		m.BackupLocation = nil
	}

	return
}
