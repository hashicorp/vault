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

// MessageDetails The content of the message to be published.
type MessageDetails struct {

	// The body of the message to be published.
	// For `messageType` of JSON, a default key-value pair is required. Example: `{"default": "Alarm breached", "Email": "Alarm breached: <url>"}.`
	// Avoid entering confidential information.
	Body *string `mandatory:"true" json:"body"`

	// The title of the message to be published. Avoid entering confidential information.
	Title *string `mandatory:"false" json:"title"`
}

func (m MessageDetails) String() string {
	return common.PointerString(m)
}
