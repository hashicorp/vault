// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

// Entry is a struct that contains a map of custom message ID to Message struct
// that can be marshalled and unmarshalled using the JSON encoding to/from a
// slice of byte to translate to a logical.StorageEntry struct. An Entry
// consists of data for a particular namespace.
type Entry struct {
	// Messages is the map of custom message ID to Message struct.
	Messages map[string]Message `json:"messages"`
}

// findMessages searches through all of the custom messages in the receiver
// Entry struct and only returns those that match the criteria set in the
// provided FindFilter struct.
func (e *Entry) findMessages(filter FindFilter) []Message {
	result := []Message{}

	for _, message := range e.Messages {
		if message.Matches(filter) {
			result = append(result, message)
		}
	}

	return result
}

// addMessage adds a custom message in the receiver Entry struct using the
// provided Message struct to populate its properties. If the either the
// start/end times are invalid or the maximum number of messages already exists,
// then the message is not added.
func (e *Entry) addMessage(message *Message) error {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	if !message.HasValidStartAndEndTimes() {
		return errors.New("message start time must occur before end time")
	}

	if !message.HasValidMessageType() {
		return errors.New("unrecognized message type")
	}

	// This condition should be evaluated last, because if anything else was to
	// prevent the creation of the message, there's no use bringing up the
	// limit.
	if len(e.Messages) >= MaximumMessageCountPerNamespace {
		return errors.New("maximum number of messages already exist")
	}

	message.ID = uuid

	if e.Messages == nil {
		e.Messages = make(map[string]Message)
	}
	e.Messages[uuid] = *message

	return nil
}

// UpdateMessage updates the Message struct stored in the receiver's Messages
// map with the provided Message struct. If the start/end times are invalid or
// if the type is invalid, the message is not updated. The Messages map is not
// changed if it does not contain the key message.ID and an error is returned.
func (e *Entry) updateMessage(message *Message) error {
	if e.Messages == nil {
		e.Messages = make(map[string]Message)
	}

	if _, ok := e.Messages[message.ID]; !ok {
		return fmt.Errorf("custom message %w", logical.ErrNotFound)
	}

	if !message.HasValidStartAndEndTimes() {
		return errors.New("message start time must occur before end time")
	}

	if !message.HasValidMessageType() {
		return errors.New("unrecognized message type")
	}

	e.Messages[message.ID] = *message

	return nil
}
