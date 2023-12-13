// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import "time"

const (
	BannerMessageType = "banner"
	ModalMessageType  = "modal"
)

var allowedMessageTypes = [2]string{
	BannerMessageType,
	ModalMessageType,
}

// Message is a struct that contains the fields of a particular custom message.
type Message struct {
	// ID is the unique identifier for this Message automatically generated and
	// assigned by the (*Entry).CreateMessage method.
	ID string `json:"id"`
	// Title contains the title of the message, which the UI uses to identify
	// this Message to the end user.
	Title string `json:"title"`
	// Message contains the body of information of this Message.
	Message string `json:"message"`
	// Type is used to store the presentation type of this message. Refer to
	// allowedMessageTypes above for the valid values.
	Type string `json:"type"`
	// Authenticated is used to indicate if the Message is intended to be
	// presented to the user before authentication (false) or after (true).
	Authenticated bool `json:"authenticated"`
	// The time when the Message begins to be active.
	StartTime time.Time `json:"start_time"`
	// The time when the Message ceases to be active.
	EndTime *time.Time `json:"end_time"`
	// Options holds additional properties for the Message.
	Options map[string]any `json:"options"`
	// Link can hold a MessageLink struct to represent a hyperlink in the
	// Message.
	Link *MessageLink `json:"link"`

	active *bool
}

// Active determines if the active field of the receiver Message has yet to be
// set (if it's nil), and if so calculates it by checking if the current time
// is after the StartTime field value AND in the case where the endTimeSet
// field is true, that the current time is before the EndTime field value. Once
// the calculation is complete, the active field is set accordingly. Finally,
// the value of the active field is returned.
func (m *Message) Active() bool {
	if m.active == nil {
		now := time.Now()

		activeValue := now.After(m.StartTime)
		if activeValue && m.EndTime != nil {
			activeValue = now.Before(*m.EndTime)
		}

		m.active = &activeValue
	}

	return *m.active
}

// Matches determines if the receiver Message struct meets the criteria
// specified by the provided FindFilter struct. A bool value is returned
// indicating if the receiver Message is a match.
func (m *Message) Matches(filters FindFilter) bool {
	if filters.authenticated != nil && *filters.authenticated != m.Authenticated {
		return false
	}

	if len(filters.messageType) != 0 && filters.messageType != m.Type {
		return false
	}

	if filters.active != nil && *filters.active != m.Active() {
		return false
	}

	return true
}

// HasValidMessageType checks if the Type field of the receiver Message
// appears in the array of allowed values and returns a bool value to indicate
// if the value is allowed.
func (m *Message) HasValidMessageType() bool {
	for _, el := range allowedMessageTypes {
		if m.Type == el {
			return true
		}
	}

	return false
}

// HasValidStartAndEndTimes evaluates the StartTime and EndTime fields of the
// receiver Message. This method returns true if the EndTime field does not
// point to a time.Time value, otherwise it returns a bool value to indicate if
// the StartTime is before the value pointed to by EndTime.
func (m *Message) HasValidStartAndEndTimes() bool {
	if m.EndTime == nil {
		return true
	}

	return m.StartTime.Before(*m.EndTime)
}

// MessageLink is a structure that represents a hyperlink included into a
// Message.
type MessageLink struct {
	// Title is the displayed hyperlink text.
	Title string `json:"title"`
	// Href is the hyperlink's URL address.
	Href string `json:"href"`
}
