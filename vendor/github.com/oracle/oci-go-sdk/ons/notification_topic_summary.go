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

// NotificationTopicSummary A summary of the properties that define a topic.
type NotificationTopicSummary struct {

	// The name of the topic. Avoid entering confidential information.
	Name *string `mandatory:"true" json:"name"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the topic.
	TopicId *string `mandatory:"true" json:"topicId"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment for the topic.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The lifecycle state of the topic.
	LifecycleState NotificationTopicSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The time the topic was created.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The endpoint for managing topic subscriptions or publishing messages to the topic.
	ApiEndpoint *string `mandatory:"true" json:"apiEndpoint"`

	// The description of the topic. Avoid entering confidential information.
	Description *string `mandatory:"false" json:"description"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `mandatory:"false" json:"etag"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m NotificationTopicSummary) String() string {
	return common.PointerString(m)
}

// NotificationTopicSummaryLifecycleStateEnum Enum with underlying type: string
type NotificationTopicSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for NotificationTopicSummaryLifecycleStateEnum
const (
	NotificationTopicSummaryLifecycleStateActive   NotificationTopicSummaryLifecycleStateEnum = "ACTIVE"
	NotificationTopicSummaryLifecycleStateDeleting NotificationTopicSummaryLifecycleStateEnum = "DELETING"
	NotificationTopicSummaryLifecycleStateCreating NotificationTopicSummaryLifecycleStateEnum = "CREATING"
)

var mappingNotificationTopicSummaryLifecycleState = map[string]NotificationTopicSummaryLifecycleStateEnum{
	"ACTIVE":   NotificationTopicSummaryLifecycleStateActive,
	"DELETING": NotificationTopicSummaryLifecycleStateDeleting,
	"CREATING": NotificationTopicSummaryLifecycleStateCreating,
}

// GetNotificationTopicSummaryLifecycleStateEnumValues Enumerates the set of values for NotificationTopicSummaryLifecycleStateEnum
func GetNotificationTopicSummaryLifecycleStateEnumValues() []NotificationTopicSummaryLifecycleStateEnum {
	values := make([]NotificationTopicSummaryLifecycleStateEnum, 0)
	for _, v := range mappingNotificationTopicSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
