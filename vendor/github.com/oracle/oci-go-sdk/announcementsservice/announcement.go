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

// Announcement A message about an impactful operational event.
type Announcement struct {

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

	// A detailed explanation of the event, expressed by using Markdown language. Avoid entering
	// confidential information.
	Description *string `mandatory:"false" json:"description"`

	// Additional information about the event, expressed by using Markdown language and included in the
	// details view of an announcement. Additional information might include remediation steps or
	// answers to frequently asked questions. Avoid entering confidential information.
	AdditionalInformation *string `mandatory:"false" json:"additionalInformation"`

	// The list of resources, if any, affected by the event described in the announcement.
	AffectedResources []AffectedResource `mandatory:"false" json:"affectedResources"`

	// The type of announcement. An announcement's type signals its severity.
	AnnouncementType BaseAnnouncementAnnouncementTypeEnum `mandatory:"true" json:"announcementType"`

	// The current lifecycle state of the announcement.
	LifecycleState BaseAnnouncementLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`
}

//GetId returns Id
func (m Announcement) GetId() *string {
	return m.Id
}

//GetReferenceTicketNumber returns ReferenceTicketNumber
func (m Announcement) GetReferenceTicketNumber() *string {
	return m.ReferenceTicketNumber
}

//GetSummary returns Summary
func (m Announcement) GetSummary() *string {
	return m.Summary
}

//GetTimeOneTitle returns TimeOneTitle
func (m Announcement) GetTimeOneTitle() *string {
	return m.TimeOneTitle
}

//GetTimeOneValue returns TimeOneValue
func (m Announcement) GetTimeOneValue() *common.SDKTime {
	return m.TimeOneValue
}

//GetTimeTwoTitle returns TimeTwoTitle
func (m Announcement) GetTimeTwoTitle() *string {
	return m.TimeTwoTitle
}

//GetTimeTwoValue returns TimeTwoValue
func (m Announcement) GetTimeTwoValue() *common.SDKTime {
	return m.TimeTwoValue
}

//GetServices returns Services
func (m Announcement) GetServices() []string {
	return m.Services
}

//GetAffectedRegions returns AffectedRegions
func (m Announcement) GetAffectedRegions() []string {
	return m.AffectedRegions
}

//GetAnnouncementType returns AnnouncementType
func (m Announcement) GetAnnouncementType() BaseAnnouncementAnnouncementTypeEnum {
	return m.AnnouncementType
}

//GetLifecycleState returns LifecycleState
func (m Announcement) GetLifecycleState() BaseAnnouncementLifecycleStateEnum {
	return m.LifecycleState
}

//GetIsBanner returns IsBanner
func (m Announcement) GetIsBanner() *bool {
	return m.IsBanner
}

//GetTimeCreated returns TimeCreated
func (m Announcement) GetTimeCreated() *common.SDKTime {
	return m.TimeCreated
}

//GetTimeUpdated returns TimeUpdated
func (m Announcement) GetTimeUpdated() *common.SDKTime {
	return m.TimeUpdated
}

func (m Announcement) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m Announcement) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeAnnouncement Announcement
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeAnnouncement
	}{
		"Announcement",
		(MarshalTypeAnnouncement)(m),
	}

	return json.Marshal(&s)
}
