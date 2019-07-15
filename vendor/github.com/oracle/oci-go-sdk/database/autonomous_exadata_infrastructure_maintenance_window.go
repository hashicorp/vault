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

// AutonomousExadataInfrastructureMaintenanceWindow Autonomous Exadata Infrastructure maintenance window details for quarterly patching.
type AutonomousExadataInfrastructureMaintenanceWindow struct {

	// Day of the week that the patch should be applied to the Autonomous Exadata Infrastructure. Patches are applied during the first week of the quarter.
	DayOfWeek AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum `mandatory:"true" json:"dayOfWeek"`

	// Hour of the day that the patch should be applied.
	HourOfDay *int `mandatory:"false" json:"hourOfDay"`
}

func (m AutonomousExadataInfrastructureMaintenanceWindow) String() string {
	return common.PointerString(m)
}

// AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum Enum with underlying type: string
type AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum string

// Set of constants representing the allowable values for AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum
const (
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekAny       AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "ANY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekSunday    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "SUNDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekMonday    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "MONDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekTuesday   AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "TUESDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekWednesday AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "WEDNESDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekThursday  AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "THURSDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekFriday    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "FRIDAY"
	AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekSaturday  AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum = "SATURDAY"
)

var mappingAutonomousExadataInfrastructureMaintenanceWindowDayOfWeek = map[string]AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum{
	"ANY":       AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekAny,
	"SUNDAY":    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekSunday,
	"MONDAY":    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekMonday,
	"TUESDAY":   AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekTuesday,
	"WEDNESDAY": AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekWednesday,
	"THURSDAY":  AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekThursday,
	"FRIDAY":    AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekFriday,
	"SATURDAY":  AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekSaturday,
}

// GetAutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnumValues Enumerates the set of values for AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum
func GetAutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnumValues() []AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum {
	values := make([]AutonomousExadataInfrastructureMaintenanceWindowDayOfWeekEnum, 0)
	for _, v := range mappingAutonomousExadataInfrastructureMaintenanceWindowDayOfWeek {
		values = append(values, v)
	}
	return values
}
