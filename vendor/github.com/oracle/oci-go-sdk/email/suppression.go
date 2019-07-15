// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Email Delivery API
//
// API for the Email Delivery service. Use this API to send high-volume, application-generated
// emails. For more information, see Overview of the Email Delivery Service (https://docs.cloud.oracle.com/iaas/Content/Email/Concepts/overview.htm).
//
// **Note:** Write actions (POST, UPDATE, DELETE) may take several minutes to propagate and be reflected by the API. If a subsequent read request fails to reflect your changes, wait a few minutes and try again.
//

package email

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Suppression The full information representing an email suppression.
type Suppression struct {

	// The OCID of the compartment to contain the suppression. Since
	// suppressions are at the customer level, this must be the tenancy
	// OCID.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// Email address of the suppression.
	EmailAddress *string `mandatory:"false" json:"emailAddress"`

	// The unique ID of the suppression.
	Id *string `mandatory:"false" json:"id"`

	// The reason that the email address was suppressed. For more information on the types of bounces, see Suppression List (https://docs.cloud.oracle.com/Content/Email/Concepts/overview.htm#components).
	Reason SuppressionReasonEnum `mandatory:"false" json:"reason,omitempty"`

	// The date and time the suppression was added in "YYYY-MM-ddThh:mmZ"
	// format with a Z offset, as defined by RFC 3339.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m Suppression) String() string {
	return common.PointerString(m)
}

// SuppressionReasonEnum Enum with underlying type: string
type SuppressionReasonEnum string

// Set of constants representing the allowable values for SuppressionReasonEnum
const (
	SuppressionReasonUnknown     SuppressionReasonEnum = "UNKNOWN"
	SuppressionReasonHardbounce  SuppressionReasonEnum = "HARDBOUNCE"
	SuppressionReasonComplaint   SuppressionReasonEnum = "COMPLAINT"
	SuppressionReasonManual      SuppressionReasonEnum = "MANUAL"
	SuppressionReasonSoftbounce  SuppressionReasonEnum = "SOFTBOUNCE"
	SuppressionReasonUnsubscribe SuppressionReasonEnum = "UNSUBSCRIBE"
)

var mappingSuppressionReason = map[string]SuppressionReasonEnum{
	"UNKNOWN":     SuppressionReasonUnknown,
	"HARDBOUNCE":  SuppressionReasonHardbounce,
	"COMPLAINT":   SuppressionReasonComplaint,
	"MANUAL":      SuppressionReasonManual,
	"SOFTBOUNCE":  SuppressionReasonSoftbounce,
	"UNSUBSCRIBE": SuppressionReasonUnsubscribe,
}

// GetSuppressionReasonEnumValues Enumerates the set of values for SuppressionReasonEnum
func GetSuppressionReasonEnumValues() []SuppressionReasonEnum {
	values := make([]SuppressionReasonEnum, 0)
	for _, v := range mappingSuppressionReason {
		values = append(values, v)
	}
	return values
}
