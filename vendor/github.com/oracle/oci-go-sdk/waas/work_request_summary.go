// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequestSummary The summarized details of a work request.
type WorkRequestSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the work request.
	Id *string `mandatory:"true" json:"id"`

	// A description of the operation requested by the work request.
	OperationType WorkRequestSummaryOperationTypeEnum `mandatory:"true" json:"operationType"`

	// The current status of the work request.
	Status WorkRequestSummaryStatusEnum `mandatory:"true" json:"status"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment that contains the work request.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The date and time the work request was created, expressed in RFC 3339 timestamp format.
	TimeAccepted *common.SDKTime `mandatory:"true" json:"timeAccepted"`

	// The date and time the work request moved from the `ACCEPTED` state to the `IN_PROGRESS` state, expressed in RFC 3339 timestamp format.
	TimeStarted *common.SDKTime `mandatory:"true" json:"timeStarted"`

	// The date and time the work request was fulfilled or terminated, in the format defined by RFC3339.
	TimeFinished *common.SDKTime `mandatory:"true" json:"timeFinished"`

	// The resources being used to complete the work request operation.
	Resources []WorkRequestResource `mandatory:"false" json:"resources"`

	// The percentage of work completed by the work request.
	PercentComplete *int `mandatory:"false" json:"percentComplete"`
}

func (m WorkRequestSummary) String() string {
	return common.PointerString(m)
}

// WorkRequestSummaryOperationTypeEnum Enum with underlying type: string
type WorkRequestSummaryOperationTypeEnum string

// Set of constants representing the allowable values for WorkRequestSummaryOperationTypeEnum
const (
	WorkRequestSummaryOperationTypeCreateWaasPolicy WorkRequestSummaryOperationTypeEnum = "CREATE_WAAS_POLICY"
	WorkRequestSummaryOperationTypeUpdateWaasPolicy WorkRequestSummaryOperationTypeEnum = "UPDATE_WAAS_POLICY"
	WorkRequestSummaryOperationTypeDeleteWaasPolicy WorkRequestSummaryOperationTypeEnum = "DELETE_WAAS_POLICY"
	WorkRequestSummaryOperationTypePurgeWaasPolicy  WorkRequestSummaryOperationTypeEnum = "PURGE_WAAS_POLICY"
)

var mappingWorkRequestSummaryOperationType = map[string]WorkRequestSummaryOperationTypeEnum{
	"CREATE_WAAS_POLICY": WorkRequestSummaryOperationTypeCreateWaasPolicy,
	"UPDATE_WAAS_POLICY": WorkRequestSummaryOperationTypeUpdateWaasPolicy,
	"DELETE_WAAS_POLICY": WorkRequestSummaryOperationTypeDeleteWaasPolicy,
	"PURGE_WAAS_POLICY":  WorkRequestSummaryOperationTypePurgeWaasPolicy,
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
	WorkRequestSummaryStatusSucceeded  WorkRequestSummaryStatusEnum = "SUCCEEDED"
	WorkRequestSummaryStatusCanceling  WorkRequestSummaryStatusEnum = "CANCELING"
	WorkRequestSummaryStatusCanceled   WorkRequestSummaryStatusEnum = "CANCELED"
)

var mappingWorkRequestSummaryStatus = map[string]WorkRequestSummaryStatusEnum{
	"ACCEPTED":    WorkRequestSummaryStatusAccepted,
	"IN_PROGRESS": WorkRequestSummaryStatusInProgress,
	"FAILED":      WorkRequestSummaryStatusFailed,
	"SUCCEEDED":   WorkRequestSummaryStatusSucceeded,
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
