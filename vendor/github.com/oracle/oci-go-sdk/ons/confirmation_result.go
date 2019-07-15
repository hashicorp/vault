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

// ConfirmationResult The confirmation result.
type ConfirmationResult struct {

	// The name of the subscribed topic.
	TopicName *string `mandatory:"true" json:"topicName"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the topic to delete.
	TopicId *string `mandatory:"true" json:"topicId"`

	// The endpoint of the subscription. Valid values depend on the protocol.
	// For EMAIL, only an email address is valid. For HTTPS, only a PagerDuty URL is valid. A URL cannot exceed 512 characters.
	// Avoid entering confidential information.
	Endpoint *string `mandatory:"true" json:"endpoint"`

	// The URL user can use to unsubscribe the topic.
	UnsubscribeUrl *string `mandatory:"true" json:"unsubscribeUrl"`

	// Human readable text which tells the user if the confirmation succeeds.
	Message *string `mandatory:"true" json:"message"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subscription.
	SubscriptionId *string `mandatory:"false" json:"subscriptionId"`
}

func (m ConfirmationResult) String() string {
	return common.PointerString(m)
}
