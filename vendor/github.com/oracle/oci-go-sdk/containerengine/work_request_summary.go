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

// WorkRequestSummary The properties that define a work request summary.
type WorkRequestSummary struct {

	// The OCID of the work request.
	Id *string `mandatory:"false" json:"id"`

	// The type of work the work request is doing.
	OperationType WorkRequestSummaryOperationTypeEnum `mandatory:"false" json:"operationType,omitempty"`

	// The current status of the work request.
	Status WorkRequestSummaryStatusEnum `mandatory:"false" json:"status,omitempty"`

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

func (m WorkRequestSummary) String() string {
	return common.PointerString(m)
}

// WorkRequestSummaryOperationTypeEnum Enum with underlying type: string
type WorkRequestSummaryOperationTypeEnum string

// Set of constants representing the allowable values for WorkRequestSummaryOperationTypeEnum
const (
	WorkRequestSummaryOperationTypeClusterCreate     WorkRequestSummaryOperationTypeEnum = "CLUSTER_CREATE"
	WorkRequestSummaryOperationTypeClusterUpdate     WorkRequestSummaryOperationTypeEnum = "CLUSTER_UPDATE"
	WorkRequestSummaryOperationTypeClusterDelete     WorkRequestSummaryOperationTypeEnum = "CLUSTER_DELETE"
	WorkRequestSummaryOperationTypeNodepoolCreate    WorkRequestSummaryOperationTypeEnum = "NODEPOOL_CREATE"
	WorkRequestSummaryOperationTypeNodepoolUpdate    WorkRequestSummaryOperationTypeEnum = "NODEPOOL_UPDATE"
	WorkRequestSummaryOperationTypeNodepoolDelete    WorkRequestSummaryOperationTypeEnum = "NODEPOOL_DELETE"
	WorkRequestSummaryOperationTypeWorkrequestCancel WorkRequestSummaryOperationTypeEnum = "WORKREQUEST_CANCEL"
)

var mappingWorkRequestSummaryOperationType = map[string]WorkRequestSummaryOperationTypeEnum{
	"CLUSTER_CREATE":     WorkRequestSummaryOperationTypeClusterCreate,
	"CLUSTER_UPDATE":     WorkRequestSummaryOperationTypeClusterUpdate,
	"CLUSTER_DELETE":     WorkRequestSummaryOperationTypeClusterDelete,
	"NODEPOOL_CREATE":    WorkRequestSummaryOperationTypeNodepoolCreate,
	"NODEPOOL_UPDATE":    WorkRequestSummaryOperationTypeNodepoolUpdate,
	"NODEPOOL_DELETE":    WorkRequestSummaryOperationTypeNodepoolDelete,
	"WORKREQUEST_CANCEL": WorkRequestSummaryOperationTypeWorkrequestCancel,
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
