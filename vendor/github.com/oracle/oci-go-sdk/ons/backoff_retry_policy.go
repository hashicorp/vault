// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Notification API
//
// Use the Notification API to broadcast messages to distributed components by topic, using a publish-subscribe pattern.
// For information about managing topics, subscriptions, and messages, see Notification Overview (https://docs.cloud.oracle.com/iaas/Content/Notification/Concepts/notificationoverview.htm).
//

package ons

import (
	"github.com/oracle/oci-go-sdk/common"
)

// BackoffRetryPolicy The backoff retry portion of the subscription delivery policy.
type BackoffRetryPolicy struct {

	// The maximum retry duration in milliseconds.
	MaxRetryDuration *int `mandatory:"true" json:"maxRetryDuration"`

	// The type of delivery policy. Default value: EXPONENTIAL.
	PolicyType BackoffRetryPolicyPolicyTypeEnum `mandatory:"true" json:"policyType"`
}

func (m BackoffRetryPolicy) String() string {
	return common.PointerString(m)
}

// BackoffRetryPolicyPolicyTypeEnum Enum with underlying type: string
type BackoffRetryPolicyPolicyTypeEnum string

// Set of constants representing the allowable values for BackoffRetryPolicyPolicyTypeEnum
const (
	BackoffRetryPolicyPolicyTypeExponential BackoffRetryPolicyPolicyTypeEnum = "EXPONENTIAL"
)

var mappingBackoffRetryPolicyPolicyType = map[string]BackoffRetryPolicyPolicyTypeEnum{
	"EXPONENTIAL": BackoffRetryPolicyPolicyTypeExponential,
}

// GetBackoffRetryPolicyPolicyTypeEnumValues Enumerates the set of values for BackoffRetryPolicyPolicyTypeEnum
func GetBackoffRetryPolicyPolicyTypeEnumValues() []BackoffRetryPolicyPolicyTypeEnum {
	values := make([]BackoffRetryPolicyPolicyTypeEnum, 0)
	for _, v := range mappingBackoffRetryPolicyPolicyType {
		values = append(values, v)
	}
	return values
}
