// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Patch The representation of Patch
type Patch struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the patch.
	Id *string `mandatory:"true" json:"id"`

	// The text describing this patch package.
	Description *string `mandatory:"true" json:"description"`

	// The date and time that the patch was released.
	TimeReleased *common.SDKTime `mandatory:"true" json:"timeReleased"`

	// The version of this patch package.
	Version *string `mandatory:"true" json:"version"`

	// Action that is currently being performed or was completed last.
	LastAction PatchLastActionEnum `mandatory:"false" json:"lastAction,omitempty"`

	// Actions that can possibly be performed using this patch.
	AvailableActions []PatchAvailableActionsEnum `mandatory:"false" json:"availableActions,omitempty"`

	// A descriptive text associated with the lifecycleState.
	// Typically can contain additional displayable text.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The current state of the patch as a result of lastAction.
	LifecycleState PatchLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`
}

func (m Patch) String() string {
	return common.PointerString(m)
}

// PatchLastActionEnum Enum with underlying type: string
type PatchLastActionEnum string

// Set of constants representing the allowable values for PatchLastActionEnum
const (
	PatchLastActionApply    PatchLastActionEnum = "APPLY"
	PatchLastActionPrecheck PatchLastActionEnum = "PRECHECK"
)

var mappingPatchLastAction = map[string]PatchLastActionEnum{
	"APPLY":    PatchLastActionApply,
	"PRECHECK": PatchLastActionPrecheck,
}

// GetPatchLastActionEnumValues Enumerates the set of values for PatchLastActionEnum
func GetPatchLastActionEnumValues() []PatchLastActionEnum {
	values := make([]PatchLastActionEnum, 0)
	for _, v := range mappingPatchLastAction {
		values = append(values, v)
	}
	return values
}

// PatchAvailableActionsEnum Enum with underlying type: string
type PatchAvailableActionsEnum string

// Set of constants representing the allowable values for PatchAvailableActionsEnum
const (
	PatchAvailableActionsApply    PatchAvailableActionsEnum = "APPLY"
	PatchAvailableActionsPrecheck PatchAvailableActionsEnum = "PRECHECK"
)

var mappingPatchAvailableActions = map[string]PatchAvailableActionsEnum{
	"APPLY":    PatchAvailableActionsApply,
	"PRECHECK": PatchAvailableActionsPrecheck,
}

// GetPatchAvailableActionsEnumValues Enumerates the set of values for PatchAvailableActionsEnum
func GetPatchAvailableActionsEnumValues() []PatchAvailableActionsEnum {
	values := make([]PatchAvailableActionsEnum, 0)
	for _, v := range mappingPatchAvailableActions {
		values = append(values, v)
	}
	return values
}

// PatchLifecycleStateEnum Enum with underlying type: string
type PatchLifecycleStateEnum string

// Set of constants representing the allowable values for PatchLifecycleStateEnum
const (
	PatchLifecycleStateAvailable  PatchLifecycleStateEnum = "AVAILABLE"
	PatchLifecycleStateSuccess    PatchLifecycleStateEnum = "SUCCESS"
	PatchLifecycleStateInProgress PatchLifecycleStateEnum = "IN_PROGRESS"
	PatchLifecycleStateFailed     PatchLifecycleStateEnum = "FAILED"
)

var mappingPatchLifecycleState = map[string]PatchLifecycleStateEnum{
	"AVAILABLE":   PatchLifecycleStateAvailable,
	"SUCCESS":     PatchLifecycleStateSuccess,
	"IN_PROGRESS": PatchLifecycleStateInProgress,
	"FAILED":      PatchLifecycleStateFailed,
}

// GetPatchLifecycleStateEnumValues Enumerates the set of values for PatchLifecycleStateEnum
func GetPatchLifecycleStateEnumValues() []PatchLifecycleStateEnum {
	values := make([]PatchLifecycleStateEnum, 0)
	for _, v := range mappingPatchLifecycleState {
		values = append(values, v)
	}
	return values
}
