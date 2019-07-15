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

// DayOfWeek Day of the week.
type DayOfWeek struct {

	// Name of the day of the week.
	Name DayOfWeekNameEnum `mandatory:"true" json:"name"`
}

func (m DayOfWeek) String() string {
	return common.PointerString(m)
}

// DayOfWeekNameEnum Enum with underlying type: string
type DayOfWeekNameEnum string

// Set of constants representing the allowable values for DayOfWeekNameEnum
const (
	DayOfWeekNameMonday    DayOfWeekNameEnum = "MONDAY"
	DayOfWeekNameTuesday   DayOfWeekNameEnum = "TUESDAY"
	DayOfWeekNameWednesday DayOfWeekNameEnum = "WEDNESDAY"
	DayOfWeekNameThursday  DayOfWeekNameEnum = "THURSDAY"
	DayOfWeekNameFriday    DayOfWeekNameEnum = "FRIDAY"
	DayOfWeekNameSaturday  DayOfWeekNameEnum = "SATURDAY"
	DayOfWeekNameSunday    DayOfWeekNameEnum = "SUNDAY"
)

var mappingDayOfWeekName = map[string]DayOfWeekNameEnum{
	"MONDAY":    DayOfWeekNameMonday,
	"TUESDAY":   DayOfWeekNameTuesday,
	"WEDNESDAY": DayOfWeekNameWednesday,
	"THURSDAY":  DayOfWeekNameThursday,
	"FRIDAY":    DayOfWeekNameFriday,
	"SATURDAY":  DayOfWeekNameSaturday,
	"SUNDAY":    DayOfWeekNameSunday,
}

// GetDayOfWeekNameEnumValues Enumerates the set of values for DayOfWeekNameEnum
func GetDayOfWeekNameEnumValues() []DayOfWeekNameEnum {
	values := make([]DayOfWeekNameEnum, 0)
	for _, v := range mappingDayOfWeekName {
		values = append(values, v)
	}
	return values
}
