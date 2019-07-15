// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequestError An object returned in the event of a work request error.
type WorkRequestError struct {
	ErrorCode WorkRequestErrorErrorCodeEnum `mandatory:"true" json:"errorCode"`

	// A human-readable error string.
	Message *string `mandatory:"true" json:"message"`
}

func (m WorkRequestError) String() string {
	return common.PointerString(m)
}

// WorkRequestErrorErrorCodeEnum Enum with underlying type: string
type WorkRequestErrorErrorCodeEnum string

// Set of constants representing the allowable values for WorkRequestErrorErrorCodeEnum
const (
	WorkRequestErrorErrorCodeBadInput      WorkRequestErrorErrorCodeEnum = "BAD_INPUT"
	WorkRequestErrorErrorCodeInternalError WorkRequestErrorErrorCodeEnum = "INTERNAL_ERROR"
)

var mappingWorkRequestErrorErrorCode = map[string]WorkRequestErrorErrorCodeEnum{
	"BAD_INPUT":      WorkRequestErrorErrorCodeBadInput,
	"INTERNAL_ERROR": WorkRequestErrorErrorCodeInternalError,
}

// GetWorkRequestErrorErrorCodeEnumValues Enumerates the set of values for WorkRequestErrorErrorCodeEnum
func GetWorkRequestErrorErrorCodeEnumValues() []WorkRequestErrorErrorCodeEnum {
	values := make([]WorkRequestErrorErrorCodeEnum, 0)
	for _, v := range mappingWorkRequestErrorErrorCode {
		values = append(values, v)
	}
	return values
}
