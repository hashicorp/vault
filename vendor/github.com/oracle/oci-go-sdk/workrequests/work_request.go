// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Work Requests API
//
// A description of the work requests API
//

package workrequests

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequest The status of a work request.
type WorkRequest struct {

	// The asynchronous operation tracked by this work request.
	OperationType *string `mandatory:"true" json:"operationType"`

	// The status of the work request.
	Status WorkRequestStatusEnum `mandatory:"true" json:"status"`

	// OCID identifying this work request.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment containing this work request.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The resources affected by this work request.
	Resources []WorkRequestResource `mandatory:"true" json:"resources"`

	// The amount of work done relative to the total amount of work.
	PercentComplete *float32 `mandatory:"true" json:"percentComplete"`

	// When the work request was created.
	TimeAccepted *common.SDKTime `mandatory:"true" json:"timeAccepted"`

	// When the work request transitioned from ACCEPTED to IN_PROGRESS.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// When the work request reached a terminal state (FAILED or SUCCEEDED).
	TimeFinished *common.SDKTime `mandatory:"false" json:"timeFinished"`
}

func (m WorkRequest) String() string {
	return common.PointerString(m)
}

// WorkRequestStatusEnum Enum with underlying type: string
type WorkRequestStatusEnum string

// Set of constants representing the allowable values for WorkRequestStatusEnum
const (
	WorkRequestStatusAccepted   WorkRequestStatusEnum = "ACCEPTED"
	WorkRequestStatusInProgress WorkRequestStatusEnum = "IN_PROGRESS"
	WorkRequestStatusFailed     WorkRequestStatusEnum = "FAILED"
	WorkRequestStatusSucceeded  WorkRequestStatusEnum = "SUCCEEDED"
	WorkRequestStatusCanceling  WorkRequestStatusEnum = "CANCELING"
	WorkRequestStatusCanceled   WorkRequestStatusEnum = "CANCELED"
)

var mappingWorkRequestStatus = map[string]WorkRequestStatusEnum{
	"ACCEPTED":    WorkRequestStatusAccepted,
	"IN_PROGRESS": WorkRequestStatusInProgress,
	"FAILED":      WorkRequestStatusFailed,
	"SUCCEEDED":   WorkRequestStatusSucceeded,
	"CANCELING":   WorkRequestStatusCanceling,
	"CANCELED":    WorkRequestStatusCanceled,
}

// GetWorkRequestStatusEnumValues Enumerates the set of values for WorkRequestStatusEnum
func GetWorkRequestStatusEnumValues() []WorkRequestStatusEnum {
	values := make([]WorkRequestStatusEnum, 0)
	for _, v := range mappingWorkRequestStatus {
		values = append(values, v)
	}
	return values
}
