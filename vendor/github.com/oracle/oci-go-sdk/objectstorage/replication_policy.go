// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
// For more information, see Overview of Object Storage (https://docs.cloud.oracle.com/Content/Object/Concepts/objectstorageoverview.htm) and
// Overview of Archive Storage (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm).
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ReplicationPolicy The details of a replication policy.
type ReplicationPolicy struct {

	// The id of the replication policy.
	Id *string `mandatory:"true" json:"id"`

	// The name of the policy.
	Name *string `mandatory:"true" json:"name"`

	// The destination region to replicate to, for example "us-ashburn-1".
	DestinationRegionName *string `mandatory:"true" json:"destinationRegionName"`

	// The bucket to replicate to in the destination region. Replication policy creation does not automatically
	// create a destination bucket. Create the destination bucket before creating the policy.
	DestinationBucketName *string `mandatory:"true" json:"destinationBucketName"`

	// The date when the replication policy was created as per RFC 3339 (https://tools.ietf.org/html/rfc3339).
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Changes made to the source bucket before this time has been replicated.
	TimeLastSync *common.SDKTime `mandatory:"true" json:"timeLastSync"`

	// The replication status of the policy. If the status is CLIENT_ERROR, once the user fixes the issue
	// described in the status message, the status will become ACTIVE.
	Status ReplicationPolicyStatusEnum `mandatory:"true" json:"status"`

	// A human-readable description of the status.
	StatusMessage *string `mandatory:"true" json:"statusMessage"`
}

func (m ReplicationPolicy) String() string {
	return common.PointerString(m)
}

// ReplicationPolicyStatusEnum Enum with underlying type: string
type ReplicationPolicyStatusEnum string

// Set of constants representing the allowable values for ReplicationPolicyStatusEnum
const (
	ReplicationPolicyStatusActive      ReplicationPolicyStatusEnum = "ACTIVE"
	ReplicationPolicyStatusClientError ReplicationPolicyStatusEnum = "CLIENT_ERROR"
)

var mappingReplicationPolicyStatus = map[string]ReplicationPolicyStatusEnum{
	"ACTIVE":       ReplicationPolicyStatusActive,
	"CLIENT_ERROR": ReplicationPolicyStatusClientError,
}

// GetReplicationPolicyStatusEnumValues Enumerates the set of values for ReplicationPolicyStatusEnum
func GetReplicationPolicyStatusEnumValues() []ReplicationPolicyStatusEnum {
	values := make([]ReplicationPolicyStatusEnum, 0)
	for _, v := range mappingReplicationPolicyStatus {
		values = append(values, v)
	}
	return values
}
