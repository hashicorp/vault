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

// CreateSubscriptionDetails The configuration details for creating the subscription.
type CreateSubscriptionDetails struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the topic for the subscription.
	TopicId *string `mandatory:"true" json:"topicId"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment for the subscription.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The protocol to use for delivering messages. Valid values: EMAIL, HTTPS.
	Protocol *string `mandatory:"true" json:"protocol"`

	// The endpoint of the subscription. Valid values depend on the protocol.
	// For EMAIL, only an email address is valid. For HTTPS, only a PagerDuty URL is valid. A URL cannot exceed 512 characters.
	// Avoid entering confidential information.
	Endpoint *string `mandatory:"true" json:"endpoint"`

	// Metadata for the subscription. Avoid entering confidential information.
	Metadata *string `mandatory:"false" json:"metadata"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateSubscriptionDetails) String() string {
	return common.PointerString(m)
}
