// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
// For more information, see Overview of Object Storage (https://docs.cloud.oracle.com/Content/Object/Concepts/objectstorageoverview.htm) and
// Overview of Archive Storage (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm).
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ObjectLifecycleRule To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type ObjectLifecycleRule struct {

	// The name of the lifecycle rule to be applied.
	Name *string `mandatory:"true" json:"name"`

	// The action of the object lifecycle policy rule. Rules using the action 'ARCHIVE' move objects into the
	// Archive Storage tier (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm). Rules using the action
	// 'DELETE' permanently delete objects from buckets. Rules using 'ABORT' abort the uncommitted multipart-uploads
	// and permanently delete their parts from buckets. 'ARCHIVE', 'DELETE' and 'ABORT' are the only three supported
	// actions at this time.
	Action *string `mandatory:"true" json:"action"`

	// Specifies the age of objects to apply the rule to. The timeAmount is interpreted in units defined by the
	// timeUnit parameter, and is calculated in relation to each object's Last-Modified time.
	TimeAmount *int64 `mandatory:"true" json:"timeAmount"`

	// The unit that should be used to interpret timeAmount.  Days are defined as starting and ending at midnight UTC.
	// Years are defined as 365.2425 days long and likewise round up to the next midnight UTC.
	TimeUnit ObjectLifecycleRuleTimeUnitEnum `mandatory:"true" json:"timeUnit"`

	// A Boolean that determines whether this rule is currently enabled.
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	Target *string `mandatory:"false" json:"target"`

	ObjectNameFilter *ObjectNameFilter `mandatory:"false" json:"objectNameFilter"`
}

func (m ObjectLifecycleRule) String() string {
	return common.PointerString(m)
}

// ObjectLifecycleRuleTimeUnitEnum Enum with underlying type: string
type ObjectLifecycleRuleTimeUnitEnum string

// Set of constants representing the allowable values for ObjectLifecycleRuleTimeUnitEnum
const (
	ObjectLifecycleRuleTimeUnitDays  ObjectLifecycleRuleTimeUnitEnum = "DAYS"
	ObjectLifecycleRuleTimeUnitYears ObjectLifecycleRuleTimeUnitEnum = "YEARS"
)

var mappingObjectLifecycleRuleTimeUnit = map[string]ObjectLifecycleRuleTimeUnitEnum{
	"DAYS":  ObjectLifecycleRuleTimeUnitDays,
	"YEARS": ObjectLifecycleRuleTimeUnitYears,
}

// GetObjectLifecycleRuleTimeUnitEnumValues Enumerates the set of values for ObjectLifecycleRuleTimeUnitEnum
func GetObjectLifecycleRuleTimeUnitEnumValues() []ObjectLifecycleRuleTimeUnitEnum {
	values := make([]ObjectLifecycleRuleTimeUnitEnum, 0)
	for _, v := range mappingObjectLifecycleRuleTimeUnit {
		values = append(values, v)
	}
	return values
}
