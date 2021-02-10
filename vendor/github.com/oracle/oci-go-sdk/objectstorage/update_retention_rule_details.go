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

// UpdateRetentionRuleDetails The details to update a retention rule.
type UpdateRetentionRuleDetails struct {

	// A user-specified name for the retention rule. Names can be helpful in identifying retention rules.
	DisplayName *string `mandatory:"false" json:"displayName"`

	Duration *Duration `mandatory:"false" json:"duration"`

	// The date and time as per RFC 3339 (https://tools.ietf.org/html/rfc3339) after which this rule is locked
	// and can only be deleted by deleting the bucket. Once a rule is locked, only increases in the duration are
	// allowed and no other properties can be changed. This property cannot be updated for rules that are in a
	// locked state. Specifying it when a duration is not specified is considered an error.
	TimeRuleLocked *common.SDKTime `mandatory:"false" json:"timeRuleLocked"`
}

func (m UpdateRetentionRuleDetails) String() string {
	return common.PointerString(m)
}
