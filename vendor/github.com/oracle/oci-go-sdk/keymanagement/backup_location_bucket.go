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

// BackupLocationBucket Object storage bucket details to upload or download the backup
type BackupLocationBucket struct {
	Namespace *string `mandatory:"true" json:"namespace"`

	BucketName *string `mandatory:"true" json:"bucketName"`

	ObjectName *string `mandatory:"true" json:"objectName"`
}

func (m BackupLocationBucket) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m BackupLocationBucket) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeBackupLocationBucket BackupLocationBucket
	s := struct {
		DiscriminatorParam string `json:"destination"`
		MarshalTypeBackupLocationBucket
	}{
		"BUCKET",
		(MarshalTypeBackupLocationBucket)(m),
	}

	return json.Marshal(&s)
}
