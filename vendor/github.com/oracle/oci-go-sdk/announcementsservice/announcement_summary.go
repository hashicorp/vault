// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Announcements Service API
//
// Manage Oracle Cloud Infrastructure console announcements.
//

package announcementsservice

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// AnnouncementSummary Summary representation of an announcement.
type AnnouncementSummary struct {

	// The OCID of the announcement.
	Id *string `mandatory:"true" json:"id"`

	// The reference Jira ticket number.
	ReferenceTicketNumber *string `mandatory:"true" json:"referenceTicketNumber"`

	// A summary of the issue. A summary might appear in the console banner view of the announcement or in
	// an email subject line. Avoid entering confidential information.
	Summary *string `mandatory:"true" json:"summary"`

	// Impacted Oracle Cloud Infrastructure services.
	Services []string `mandatory:"true" json:"services"`

	// Impacted regions.
	AffectedRegions []string `mandatory:"true" json:"affectedRegions"`

	// Whether the announcement is displayed as a banner in the console.
	IsBanner *bool `mandatory:"true" json:"isBanner"`

	// The label associated with an initial time value.
	// Example: `Time Started`
	TimeOneTitle *string `mandatory:"false" json:"timeOneTitle"`

	// The actual value of the first time value for the event. Typically, this is the time an event started, but the meaning
	// can vary, depending on the announcement type.
	TimeOneValue *common.SDKTime `mandatory:"false" json:"timeOneValue"`

	// The label associated with a second time value.
	// Example: `Time Ended`
	TimeTwoTitle *string `mandatory:"false" json:"timeTwoTitle"`

	// The actual value of the second time value. Typically, this is the time an event ended, but the meaning
	// can vary, depending on the announcement type.
	TimeTwoValue *common.SDKTime `mandatory:"false" json:"timeTwoValue"`

	// The date and time the announcement was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-01-01T17:43:01.389+0000`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The date and time the announcement was last updated, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-01-01T17:43:01.389+0000`
	TimeUpdated *common.SDKTime `mandatory:"false" json:"timeUpdated"`

	// The type of announcement. An announcement's type signals its severity.
	AnnouncementType BaseAnnouncementAnnouncementTypeEnum `mandatory:"true" json:"announcementType"`

	// The current lifecycle state of the announcement.
	LifecycleState BaseAnnouncementLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`
}

//GetId returns Id
func (m AnnouncementSummary) GetId() *string {
	return m.Id
}

//GetReferenceTicketNumber returns ReferenceTicketNumber
func (m AnnouncementSummary) GetReferenceTicketNumber() *string {
	return m.ReferenceTicketNumber
}

//GetSummary returns Summary
func (m AnnouncementSummary) GetSummary() *string {
	return m.Summary
}

//GetTimeOneTitle returns TimeOneTitle
func (m AnnouncementSummary) GetTimeOneTitle() *string {
	return m.TimeOneTitle
}

//GetTimeOneValue returns TimeOneValue
func (m AnnouncementSummary) GetTimeOneValue() *common.SDKTime {
	return m.TimeOneValue
}

//GetTimeTwoTitle returns TimeTwoTitle
func (m AnnouncementSummary) GetTimeTwoTitle() *string {
	return m.TimeTwoTitle
}

//GetTimeTwoValue returns TimeTwoValue
func (m AnnouncementSummary) GetTimeTwoValue() *common.SDKTime {
	return m.TimeTwoValue
}

//GetServices returns Services
func (m AnnouncementSummary) GetServices() []string {
	return m.Services
}

//GetAffectedRegions returns AffectedRegions
func (m AnnouncementSummary) GetAffectedRegions() []string {
	return m.AffectedRegions
}

//GetAnnouncementType returns AnnouncementType
func (m AnnouncementSummary) GetAnnouncementType() BaseAnnouncementAnnouncementTypeEnum {
	return m.AnnouncementType
}

//GetLifecycleState returns LifecycleState
func (m AnnouncementSummary) GetLifecycleState() BaseAnnouncementLifecycleStateEnum {
	return m.LifecycleState
}

//GetIsBanner returns IsBanner
func (m AnnouncementSummary) GetIsBanner() *bool {
	return m.IsBanner
}

//GetTimeCreated returns TimeCreated
func (m AnnouncementSummary) GetTimeCreated() *common.SDKTime {
	return m.TimeCreated
}

//GetTimeUpdated returns TimeUpdated
func (m AnnouncementSummary) GetTimeUpdated() *common.SDKTime {
	return m.TimeUpdated
}

func (m AnnouncementSummary) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m AnnouncementSummary) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeAnnouncementSummary AnnouncementSummary
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeAnnouncementSummary
	}{
		"AnnouncementSummary",
		(MarshalTypeAnnouncementSummary)(m),
	}

	return json.Marshal(&s)
}
