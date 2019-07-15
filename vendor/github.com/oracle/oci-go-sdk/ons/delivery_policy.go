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

// DeliveryPolicy The subscription delivery policy.
type DeliveryPolicy struct {
	BackoffRetryPolicy *BackoffRetryPolicy `mandatory:"false" json:"backoffRetryPolicy"`
}

func (m DeliveryPolicy) String() string {
	return common.PointerString(m)
}
