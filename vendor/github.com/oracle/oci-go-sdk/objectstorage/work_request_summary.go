// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequestSummary A summary of the status of a work request.
type WorkRequestSummary struct {

	// The type of work request.
	OperationType WorkRequestSummaryOperationTypeEnum `mandatory:"false" json:"operationType,omitempty"`

	// The status of a specified work request.
	Status WorkRequestSummaryStatusEnum `mandatory:"false" json:"status,omitempty"`

	// The id of the work request.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the compartment that contains the work request. Work requests should be scoped to
	// the same compartment as the resource the work request affects. If the work request affects multiple resources,
	// and those resources are not in the same compartment, it is up to the service team to pick the primary
	// resource whose compartment should be used.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	Resources []WorkRequestResource `mandatory:"false" json:"resources"`

	// Percentage of the work request completed.
	PercentComplete *float32 `mandatory:"false" json:"percentComplete"`

	// The date and time the work request was created, as described in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	TimeAccepted *common.SDKTime `mandatory:"false" json:"timeAccepted"`

	// The date and time the work request was started, as described in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the work request was finished, as described in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	TimeFinished *common.SDKTime `mandatory:"false" json:"timeFinished"`
}

func (m WorkRequestSummary) String() string {
	return common.PointerString(m)
}

// WorkRequestSummaryOperationTypeEnum Enum with underlying type: string
type WorkRequestSummaryOperationTypeEnum string

// Set of constants representing the allowable values for WorkRequestSummaryOperationTypeEnum
const (
	WorkRequestSummaryOperationTypeObject WorkRequestSummaryOperationTypeEnum = "COPY_OBJECT"
)

var mappingWorkRequestSummaryOperationType = map[string]WorkRequestSummaryOperationTypeEnum{
	"COPY_OBJECT": WorkRequestSummaryOperationTypeObject,
}

// GetWorkRequestSummaryOperationTypeEnumValues Enumerates the set of values for WorkRequestSummaryOperationTypeEnum
func GetWorkRequestSummaryOperationTypeEnumValues() []WorkRequestSummaryOperationTypeEnum {
	values := make([]WorkRequestSummaryOperationTypeEnum, 0)
	for _, v := range mappingWorkRequestSummaryOperationType {
		values = append(values, v)
	}
	return values
}

// WorkRequestSummaryStatusEnum Enum with underlying type: string
type WorkRequestSummaryStatusEnum string

// Set of constants representing the allowable values for WorkRequestSummaryStatusEnum
const (
	WorkRequestSummaryStatusAccepted   WorkRequestSummaryStatusEnum = "ACCEPTED"
	WorkRequestSummaryStatusInProgress WorkRequestSummaryStatusEnum = "IN_PROGRESS"
	WorkRequestSummaryStatusFailed     WorkRequestSummaryStatusEnum = "FAILED"
	WorkRequestSummaryStatusCompleted  WorkRequestSummaryStatusEnum = "COMPLETED"
	WorkRequestSummaryStatusCanceling  WorkRequestSummaryStatusEnum = "CANCELING"
	WorkRequestSummaryStatusCanceled   WorkRequestSummaryStatusEnum = "CANCELED"
)

var mappingWorkRequestSummaryStatus = map[string]WorkRequestSummaryStatusEnum{
	"ACCEPTED":    WorkRequestSummaryStatusAccepted,
	"IN_PROGRESS": WorkRequestSummaryStatusInProgress,
	"FAILED":      WorkRequestSummaryStatusFailed,
	"COMPLETED":   WorkRequestSummaryStatusCompleted,
	"CANCELING":   WorkRequestSummaryStatusCanceling,
	"CANCELED":    WorkRequestSummaryStatusCanceled,
}

// GetWorkRequestSummaryStatusEnumValues Enumerates the set of values for WorkRequestSummaryStatusEnum
func GetWorkRequestSummaryStatusEnumValues() []WorkRequestSummaryStatusEnum {
	values := make([]WorkRequestSummaryStatusEnum, 0)
	for _, v := range mappingWorkRequestSummaryStatus {
		values = append(values, v)
	}
	return values
}
