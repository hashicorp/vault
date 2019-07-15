// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Container Engine for Kubernetes API
//
// API for the Container Engine for Kubernetes service. Use this API to build, deploy,
// and manage cloud-native applications. For more information, see
// Overview of Container Engine for Kubernetes (https://docs.cloud.oracle.com/iaas/Content/ContEng/Concepts/contengoverview.htm).
//

package containerengine

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequest An asynchronous work request.
type WorkRequest struct {

	// The OCID of the work request.
	Id *string `mandatory:"false" json:"id"`

	// The type of work the work request is doing.
	OperationType WorkRequestOperationTypeEnum `mandatory:"false" json:"operationType,omitempty"`

	// The current status of the work request.
	Status WorkRequestStatusEnum `mandatory:"false" json:"status,omitempty"`

	// The OCID of the compartment in which the work request exists.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The resources this work request affects.
	Resources []WorkRequestResource `mandatory:"false" json:"resources"`

	// The time the work request was accepted.
	TimeAccepted *common.SDKTime `mandatory:"false" json:"timeAccepted"`

	// The time the work request was started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The time the work request was finished.
	TimeFinished *common.SDKTime `mandatory:"false" json:"timeFinished"`
}

func (m WorkRequest) String() string {
	return common.PointerString(m)
}

// WorkRequestOperationTypeEnum Enum with underlying type: string
type WorkRequestOperationTypeEnum string

// Set of constants representing the allowable values for WorkRequestOperationTypeEnum
const (
	WorkRequestOperationTypeClusterCreate     WorkRequestOperationTypeEnum = "CLUSTER_CREATE"
	WorkRequestOperationTypeClusterUpdate     WorkRequestOperationTypeEnum = "CLUSTER_UPDATE"
	WorkRequestOperationTypeClusterDelete     WorkRequestOperationTypeEnum = "CLUSTER_DELETE"
	WorkRequestOperationTypeNodepoolCreate    WorkRequestOperationTypeEnum = "NODEPOOL_CREATE"
	WorkRequestOperationTypeNodepoolUpdate    WorkRequestOperationTypeEnum = "NODEPOOL_UPDATE"
	WorkRequestOperationTypeNodepoolDelete    WorkRequestOperationTypeEnum = "NODEPOOL_DELETE"
	WorkRequestOperationTypeWorkrequestCancel WorkRequestOperationTypeEnum = "WORKREQUEST_CANCEL"
)

var mappingWorkRequestOperationType = map[string]WorkRequestOperationTypeEnum{
	"CLUSTER_CREATE":     WorkRequestOperationTypeClusterCreate,
	"CLUSTER_UPDATE":     WorkRequestOperationTypeClusterUpdate,
	"CLUSTER_DELETE":     WorkRequestOperationTypeClusterDelete,
	"NODEPOOL_CREATE":    WorkRequestOperationTypeNodepoolCreate,
	"NODEPOOL_UPDATE":    WorkRequestOperationTypeNodepoolUpdate,
	"NODEPOOL_DELETE":    WorkRequestOperationTypeNodepoolDelete,
	"WORKREQUEST_CANCEL": WorkRequestOperationTypeWorkrequestCancel,
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
