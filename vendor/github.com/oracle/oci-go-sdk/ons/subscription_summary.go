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

// SubscriptionSummary The subscription's configuration summary.
type SubscriptionSummary struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subscription.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the associated topic.
	TopicId *string `mandatory:"true" json:"topicId"`

	// The protocol used for the subscription. Valid values: EMAIL, HTTPS.
	Protocol *string `mandatory:"true" json:"protocol"`

	// The endpoint of the subscription. Valid values depend on the protocol.
	// For EMAIL, only an email address is valid. For HTTPS, only a PagerDuty URL is valid. A URL cannot exceed 512 characters.
	// Avoid entering confidential information.
	Endpoint *string `mandatory:"true" json:"endpoint"`

	// The lifecycle state of the subscription. Default value for a newly created subscription: PENDING.
	LifecycleState SubscriptionSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The time when this suscription was created.
	CreatedTime *int64 `mandatory:"false" json:"createdTime"`

	DeliveryPolicy *DeliveryPolicy `mandatory:"false" json:"deliveryPolicy"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `mandatory:"false" json:"etag"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m SubscriptionSummary) String() string {
	return common.PointerString(m)
}

// SubscriptionSummaryLifecycleStateEnum Enum with underlying type: string
type SubscriptionSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for SubscriptionSummaryLifecycleStateEnum
const (
	SubscriptionSummaryLifecycleStatePending SubscriptionSummaryLifecycleStateEnum = "PENDING"
	SubscriptionSummaryLifecycleStateActive  SubscriptionSummaryLifecycleStateEnum = "ACTIVE"
	SubscriptionSummaryLifecycleStateDeleted SubscriptionSummaryLifecycleStateEnum = "DELETED"
)

var mappingSubscriptionSummaryLifecycleState = map[string]SubscriptionSummaryLifecycleStateEnum{
	"PENDING": SubscriptionSummaryLifecycleStatePending,
	"ACTIVE":  SubscriptionSummaryLifecycleStateActive,
	"DELETED": SubscriptionSummaryLifecycleStateDeleted,
}

// GetSubscriptionSummaryLifecycleStateEnumValues Enumerates the set of values for SubscriptionSummaryLifecycleStateEnum
func GetSubscriptionSummaryLifecycleStateEnumValues() []SubscriptionSummaryLifecycleStateEnum {
	values := make([]SubscriptionSummaryLifecycleStateEnum, 0)
	for _, v := range mappingSubscriptionSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
