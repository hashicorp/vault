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

// ExadataIormConfigUpdateDetails IORM Setting details for this Exadata System to be updated
type ExadataIormConfigUpdateDetails struct {

	// Value for the IORM objective
	// Default is "Auto"
	Objective ExadataIormConfigUpdateDetailsObjectiveEnum `mandatory:"false" json:"objective,omitempty"`

	// Array of IORM Setting for all the database in
	// this Exadata DB System
	DbPlans []DbIormConfigUpdateDetail `mandatory:"false" json:"dbPlans"`
}

func (m ExadataIormConfigUpdateDetails) String() string {
	return common.PointerString(m)
}

// ExadataIormConfigUpdateDetailsObjectiveEnum Enum with underlying type: string
type ExadataIormConfigUpdateDetailsObjectiveEnum string

// Set of constants representing the allowable values for ExadataIormConfigUpdateDetailsObjectiveEnum
const (
	ExadataIormConfigUpdateDetailsObjectiveLowLatency     ExadataIormConfigUpdateDetailsObjectiveEnum = "LOW_LATENCY"
	ExadataIormConfigUpdateDetailsObjectiveHighThroughput ExadataIormConfigUpdateDetailsObjectiveEnum = "HIGH_THROUGHPUT"
	ExadataIormConfigUpdateDetailsObjectiveBalanced       ExadataIormConfigUpdateDetailsObjectiveEnum = "BALANCED"
	ExadataIormConfigUpdateDetailsObjectiveAuto           ExadataIormConfigUpdateDetailsObjectiveEnum = "AUTO"
	ExadataIormConfigUpdateDetailsObjectiveBasic          ExadataIormConfigUpdateDetailsObjectiveEnum = "BASIC"
)

var mappingExadataIormConfigUpdateDetailsObjective = map[string]ExadataIormConfigUpdateDetailsObjectiveEnum{
	"LOW_LATENCY":     ExadataIormConfigUpdateDetailsObjectiveLowLatency,
	"HIGH_THROUGHPUT": ExadataIormConfigUpdateDetailsObjectiveHighThroughput,
	"BALANCED":        ExadataIormConfigUpdateDetailsObjectiveBalanced,
	"AUTO":            ExadataIormConfigUpdateDetailsObjectiveAuto,
	"BASIC":           ExadataIormConfigUpdateDetailsObjectiveBasic,
}

// GetExadataIormConfigUpdateDetailsObjectiveEnumValues Enumerates the set of values for ExadataIormConfigUpdateDetailsObjectiveEnum
func GetExadataIormConfigUpdateDetailsObjectiveEnumValues() []ExadataIormConfigUpdateDetailsObjectiveEnum {
	values := make([]ExadataIormConfigUpdateDetailsObjectiveEnum, 0)
	for _, v := range mappingExadataIormConfigUpdateDetailsObjective {
		values = append(values, v)
	}
	return values
}
