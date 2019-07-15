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

// WorkRequest Many of the API requests you use to create and configure WAAS policies do not take effect immediately. In these cases, the request spawns an asynchronous work flow to fulfill the request. `WorkRequest` objects provide visibility for in-progress work flows. For more information about work requests, see Viewing the State of a Work Request (https://docs.cloud.oracle.com/Content/Balance/Tasks/viewingworkrequest.htm).
type WorkRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the work request.
	Id *string `mandatory:"true" json:"id"`

	// A description of the operation requested by the work request.
	OperationType WorkRequestOperationTypeEnum `mandatory:"true" json:"operationType"`

	// The current status of the work request.
	Status WorkRequestStatusEnum `mandatory:"true" json:"status"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment that contains the work request.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The date and time the work request was created, in the format defined by RFC3339.
	TimeAccepted *common.SDKTime `mandatory:"true" json:"timeAccepted"`

	// The date and time the work request moved from the `ACCEPTED` state to the `IN_PROGRESS` state, expressed in RFC 3339 timestamp format.
	TimeStarted *common.SDKTime `mandatory:"true" json:"timeStarted"`

	// The date and time the work request was fulfilled or terminated, expressed in RFC 3339 timestamp format.
	TimeFinished *common.SDKTime `mandatory:"true" json:"timeFinished"`

	// The resources being used to complete the work request operation.
	Resources []WorkRequestResource `mandatory:"false" json:"resources"`

	// The percentage of work completed by the work request.
	PercentComplete *int `mandatory:"false" json:"percentComplete"`

	// The list of log entries from the work request workflow.
	Logs []WorkRequestLogEntry `mandatory:"false" json:"logs"`

	// The list of errors that occurred while fulfilling the work request.
	Errors []WorkRequestError `mandatory:"false" json:"errors"`
}

func (m WorkRequest) String() string {
	return common.PointerString(m)
}

// WorkRequestOperationTypeEnum Enum with underlying type: string
type WorkRequestOperationTypeEnum string

// Set of constants representing the allowable values for WorkRequestOperationTypeEnum
const (
	WorkRequestOperationTypeCreateWaasPolicy WorkRequestOperationTypeEnum = "CREATE_WAAS_POLICY"
	WorkRequestOperationTypeUpdateWaasPolicy WorkRequestOperationTypeEnum = "UPDATE_WAAS_POLICY"
	WorkRequestOperationTypeDeleteWaasPolicy WorkRequestOperationTypeEnum = "DELETE_WAAS_POLICY"
	WorkRequestOperationTypePurgeWaasPolicy  WorkRequestOperationTypeEnum = "PURGE_WAAS_POLICY"
)

var mappingWorkRequestOperationType = map[string]WorkRequestOperationTypeEnum{
	"CREATE_WAAS_POLICY": WorkRequestOperationTypeCreateWaasPolicy,
	"UPDATE_WAAS_POLICY": WorkRequestOperationTypeUpdateWaasPolicy,
	"DELETE_WAAS_POLICY": WorkRequestOperationTypeDeleteWaasPolicy,
	"PURGE_WAAS_POLICY":  WorkRequestOperationTypePurgeWaasPolicy,
}

// GetWorkRequestOperationTypeEnumValues Enumerates the set of values for WorkRequestOperationTypeEnum
func GetWorkRequestOperationTypeEnumValues() []WorkRequestOperationTypeEnum {
	values := make([]WorkRequestOperationTypeEnum, 0)
	for _, v := range mappingWorkRequestOperationType {
		values = append(values, v)
	}
	return values
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
