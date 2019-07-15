// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Announcements Service API
//
// Manage Oracle Cloud Infrastructure console announcements.
//

package announcementsservice

import (
	"github.com/oracle/oci-go-sdk/common"
)

// AnnouncementsCollection A list of announcements that match filter criteria, if any. Results contain both the announcements and the user-specific status of the announcements.
type AnnouncementsCollection struct {

	// A collection of announcements.
	Items []AnnouncementSummary `mandatory:"false" json:"items"`

	// The user-specific status for found announcements.
	UserStatuses []AnnouncementUserStatusDetails `mandatory:"false" json:"userStatuses"`
}

func (m AnnouncementsCollection) String() string {
	return common.PointerString(m)
}
