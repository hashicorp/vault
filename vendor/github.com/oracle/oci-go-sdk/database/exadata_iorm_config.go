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

// ExadataIormConfig Response details which has IORM Settings for this Exadata System
type ExadataIormConfig struct {

	// The current config state of IORM settings for this Exadata System.
	LifecycleState ExadataIormConfigLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Additional information about the current lifecycleState.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// Value for the IORM objective
	// Default is "Auto"
	Objective ExadataIormConfigObjectiveEnum `mandatory:"false" json:"objective,omitempty"`

	// Array of IORM Setting for all the database in
	// this Exadata DB System
	DbPlans []DbIormConfig `mandatory:"false" json:"dbPlans"`
}

func (m ExadataIormConfig) String() string {
	return common.PointerString(m)
}

// ExadataIormConfigLifecycleStateEnum Enum with underlying type: string
type ExadataIormConfigLifecycleStateEnum string

// Set of constants representing the allowable values for ExadataIormConfigLifecycleStateEnum
const (
	ExadataIormConfigLifecycleStateBootstrapping ExadataIormConfigLifecycleStateEnum = "BOOTSTRAPPING"
	ExadataIormConfigLifecycleStateEnabled       ExadataIormConfigLifecycleStateEnum = "ENABLED"
	ExadataIormConfigLifecycleStateDisabled      ExadataIormConfigLifecycleStateEnum = "DISABLED"
	ExadataIormConfigLifecycleStateUpdating      ExadataIormConfigLifecycleStateEnum = "UPDATING"
	ExadataIormConfigLifecycleStateFailed        ExadataIormConfigLifecycleStateEnum = "FAILED"
)

var mappingExadataIormConfigLifecycleState = map[string]ExadataIormConfigLifecycleStateEnum{
	"BOOTSTRAPPING": ExadataIormConfigLifecycleStateBootstrapping,
	"ENABLED":       ExadataIormConfigLifecycleStateEnabled,
	"DISABLED":      ExadataIormConfigLifecycleStateDisabled,
	"UPDATING":      ExadataIormConfigLifecycleStateUpdating,
	"FAILED":        ExadataIormConfigLifecycleStateFailed,
}

// GetExadataIormConfigLifecycleStateEnumValues Enumerates the set of values for ExadataIormConfigLifecycleStateEnum
func GetExadataIormConfigLifecycleStateEnumValues() []ExadataIormConfigLifecycleStateEnum {
	values := make([]ExadataIormConfigLifecycleStateEnum, 0)
	for _, v := range mappingExadataIormConfigLifecycleState {
		values = append(values, v)
	}
	return values
}

// ExadataIormConfigObjectiveEnum Enum with underlying type: string
type ExadataIormConfigObjectiveEnum string

// Set of constants representing the allowable values for ExadataIormConfigObjectiveEnum
const (
	ExadataIormConfigObjectiveLowLatency     ExadataIormConfigObjectiveEnum = "LOW_LATENCY"
	ExadataIormConfigObjectiveHighThroughput ExadataIormConfigObjectiveEnum = "HIGH_THROUGHPUT"
	ExadataIormConfigObjectiveBalanced       ExadataIormConfigObjectiveEnum = "BALANCED"
	ExadataIormConfigObjectiveAuto           ExadataIormConfigObjectiveEnum = "AUTO"
	ExadataIormConfigObjectiveBasic          ExadataIormConfigObjectiveEnum = "BASIC"
)

var mappingExadataIormConfigObjective = map[string]ExadataIormConfigObjectiveEnum{
	"LOW_LATENCY":     ExadataIormConfigObjectiveLowLatency,
	"HIGH_THROUGHPUT": ExadataIormConfigObjectiveHighThroughput,
	"BALANCED":        ExadataIormConfigObjectiveBalanced,
	"AUTO":            ExadataIormConfigObjectiveAuto,
	"BASIC":           ExadataIormConfigObjectiveBasic,
}

// GetExadataIormConfigObjectiveEnumValues Enumerates the set of values for ExadataIormConfigObjectiveEnum
func GetExadataIormConfigObjectiveEnumValues() []ExadataIormConfigObjectiveEnum {
	values := make([]ExadataIormConfigObjectiveEnum, 0)
	for _, v := range mappingExadataIormConfigObjective {
		values = append(values, v)
	}
	return values
}
