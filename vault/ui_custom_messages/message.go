// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import "time"

var allowedMessageTypes = [2]string{
	"banner",
	"modal",
}

// Message is a struct that contains the fields of a particular custom message.
type Message struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Message       string         `json:"message"`
	Type          string         `json:"type"`
	Authenticated bool           `json:"authenticated"`
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time"`
	Options       map[string]any `json:"options"`
	Link          map[string]any `json:"link"`

	active *bool
}

// Active determines if the active field of the receiver Message has yet to be
// set (if it's nil), and if so calculates it by checking if the current time
// is after the StartTime field value AND before the EndTime field value and
// sets the active field accordingly. Finally, the value of the active field
// is returned.
func (m *Message) Active() bool {
	if m.active == nil {

		now := time.Now()

		activeValue := now.After(m.StartTime) && now.Before(m.EndTime)
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

// ValidateMessageType checks if the Type field of the receiver Message appears
// in the array of allowed values and returns a bool value to indicate if the
// value is allowed.
func (m *Message) ValidateMessageType() bool {
	for _, el := range allowedMessageTypes {
		if m.Type == el {
			return true
		}
	}

	return false
}

// ValidateStartAndEndTimes evaluates the StartTime and EndTime fields of the
// receiver Message and returns a bool value to indicate if the StartTime is
// before the EndTime.
func (m *Message) ValidateStartAndEndTimes() bool {
	return m.StartTime.Before(m.EndTime)
}
