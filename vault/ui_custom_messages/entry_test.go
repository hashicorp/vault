// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEntryFindMessages verifies that the (*Entry).FindMessages method behaves
// correctly in the different edge cases that could occur.
func TestEntryFindMessages(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)
	earlier := now.Add(-1 * time.Hour)

	testMessagesMap := map[string]Message{
		"post-login-message": {
			Authenticated: true,
		},
		"modal-message": {
			Type: ModalMessageType,
		},
		"active-message": {
			StartTime: earlier,
			EndTime:   &later,
		},
	}

	for _, testcase := range []struct {
		name             string
		entry            Entry
		filter           FindFilter
		expectResultsLen int
	}{
		{
			name:             "messages map is nil: filter empty",
			entry:            Entry{},
			filter:           FindFilter{},
			expectResultsLen: 0,
		},
		{
			name: "filter matches",
			entry: Entry{
				Messages: testMessagesMap,
			},
			filter: FindFilter{
				messageType: ModalMessageType,
			},
			expectResultsLen: 1,
		},
		{
			name: "filter does not match",
			entry: Entry{
				Messages: testMessagesMap,
			},
			filter: FindFilter{
				messageType: BannerMessageType,
			},
			expectResultsLen: 0,
		},
	} {
		result := testcase.entry.findMessages(testcase.filter)

		assert.NotNil(t, result, testcase.name)
		assert.Equal(t, testcase.expectResultsLen, len(result))
	}
}

// TestEntryCreateMessage verifies that the (*Entry).CreateMessage method
// behaves correctly in the different edge cases that could occur.
func TestEntryCreateMessage(t *testing.T) {
	var (
		testEntry        = Entry{}
		time1            = time.Now()
		time2            = time1.Add(time.Hour)
		testValidMessage = Message{
			StartTime: time1,
			EndTime:   &time2,
			Type:      BannerMessageType,
		}
		testInvalidTimesMessage = Message{
			StartTime: time2,
			EndTime:   &time1,
			Type:      BannerMessageType,
		}
		testInvalidTypeMessage = Message{
			StartTime: time1,
			EndTime:   &time2,
			Type:      "watermark",
		}
	)

	for _, testcase := range []struct {
		name                 string
		messagesMap          map[string]Message
		message              Message
		expectedError        bool
		expectedErrorKeyword string
	}{
		{
			name:    "uninitialized messages map",
			message: testValidMessage,
		},
		{
			name:        "empty messages map, valid message",
			messagesMap: make(map[string]Message),
			message:     testValidMessage,
		},
		{
			name:                 "empty messages map, invalid times message",
			messagesMap:          make(map[string]Message),
			message:              testInvalidTimesMessage,
			expectedError:        true,
			expectedErrorKeyword: "must occur before",
		},
		{
			name:                 "empty messages map, invalid type message",
			messagesMap:          make(map[string]Message),
			message:              testInvalidTypeMessage,
			expectedError:        true,
			expectedErrorKeyword: "unrecognized",
		},
		{
			name:                 "full messages map, valid message",
			messagesMap:          buildMessagesMap(testValidMessage, MaximumMessageCountPerNamespace),
			message:              testValidMessage,
			expectedError:        true,
			expectedErrorKeyword: "maximum number",
		},
		{
			name:                 "full messages map, invalid times message",
			messagesMap:          buildMessagesMap(testValidMessage, MaximumMessageCountPerNamespace),
			message:              testInvalidTimesMessage,
			expectedError:        true,
			expectedErrorKeyword: "must occur before",
		},
		{
			name:                 "full messages map, invalid type message",
			messagesMap:          buildMessagesMap(testValidMessage, MaximumMessageCountPerNamespace),
			message:              testInvalidTypeMessage,
			expectedError:        true,
			expectedErrorKeyword: "unrecognized",
		},
		{
			name:        "nearly full messages map, valid message",
			messagesMap: buildMessagesMap(testValidMessage, MaximumMessageCountPerNamespace-1),
			message:     testValidMessage,
		},
	} {
		// Set the Messages field to the testcase's messagesMap field.
		testEntry.Messages = testcase.messagesMap

		// Count the number of messages to compare after the CreateMessage call.
		previousMessageCount := len(testcase.messagesMap)

		err := testEntry.addMessage(&testcase.message)
		if testcase.expectedError {
			assert.Error(t, err, testcase.name)
			assert.Contains(t, err.Error(), testcase.expectedErrorKeyword)
			assert.Equal(t, previousMessageCount, len(testEntry.Messages), testcase.name)
		} else {
			assert.NoError(t, err, testcase.name)
			assert.Equal(t, previousMessageCount+1, len(testEntry.Messages), testcase.name)
			assert.NotEmpty(t, testcase.message.ID, testcase.name)
		}
	}
}

// TestEntryUpdateMessage verifies that the (*Entry).UpdateMessage method
// behaves correctly in different edge cases that could occur.
func TestEntryUpdateMessage(t *testing.T) {
	var (
		testEntry = Entry{}

		now     = time.Now()
		later   = now.Add(time.Hour)
		earlier = now.Add(-1 * time.Hour)

		testValidMessage = Message{
			StartTime: now,
			EndTime:   &later,
			Type:      BannerMessageType,
		}
		testInvalidTimesMessage = Message{
			StartTime: now,
			EndTime:   &earlier,
			Type:      BannerMessageType,
		}
		testInvalidTypeMessage = Message{
			StartTime: now,
			EndTime:   &later,
			Type:      "watermark",
		}
	)

	for _, testcase := range []struct {
		name             string
		messagesMap      map[string]Message
		message          Message
		errorAssertion   func(assert.TestingT, error, ...any) bool
		expectedUpdated  bool
		compareAssertion func(assert.TestingT, any, any, ...any) bool
	}{
		{
			name:           "uninitialized messages map",
			messagesMap:    nil,
			message:        testValidMessage,
			errorAssertion: assert.Error,
		},
		{
			name:           "empty messages map",
			messagesMap:    make(map[string]Message),
			message:        testValidMessage,
			errorAssertion: assert.Error,
		},
		{
			name:           "updating existing with invalid times",
			messagesMap:    buildMessagesMap(testValidMessage, 1),
			message:        updateMessageID(testInvalidTimesMessage, "0"),
			errorAssertion: assert.Error,
		},
		{
			name:           "updating existing with invalid type",
			messagesMap:    buildMessagesMap(testValidMessage, 1),
			message:        updateMessageID(testInvalidTypeMessage, "0"),
			errorAssertion: assert.Error,
		},
		{
			name:             "updating existing with valid times and no changes",
			messagesMap:      buildMessagesMap(testValidMessage, 1),
			message:          updateMessageID(testValidMessage, "0"),
			errorAssertion:   assert.NoError,
			expectedUpdated:  true,
			compareAssertion: assert.Equal,
		},
		{
			name:             "updating existing with valid times and changes",
			messagesMap:      buildMessagesMap(testInvalidTimesMessage, 1),
			message:          updateMessageID(testValidMessage, "0"),
			errorAssertion:   assert.NoError,
			expectedUpdated:  true,
			compareAssertion: assert.NotEqual,
		},
	} {
		testEntry.Messages = testcase.messagesMap

		var previousMessage Message

		if testEntry.Messages != nil {
			previousMessage = testEntry.Messages[testcase.message.ID]
		}

		err := testEntry.updateMessage(&testcase.message)
		testcase.errorAssertion(t, err, testcase.name)

		currentMessage := testEntry.Messages[testcase.message.ID]

		if testcase.expectedUpdated {
			testcase.compareAssertion(t, previousMessage, currentMessage, testcase.name)
		} else {
			assert.Equal(t, previousMessage, currentMessage, testcase.name)
		}
	}
}

// buildMessagesMap is a helper that builds a map[string]Message and loads it
// with n elements, where the keys are the string representation of the ordinals
// 0 to n-1 and the values are copies of m.
func buildMessagesMap(m Message, n int) map[string]Message {
	messageMap := make(map[string]Message)

	for i := 0; i < n; i++ {
		m.ID = fmt.Sprintf("%d", i)
		messageMap[m.ID] = m
	}

	return messageMap
}

// updateMessageID is a helper that takes a Message struct and returns a copy of
// it with the ID set to the specified id value.
func updateMessageID(m Message, id string) Message {
	m.ID = id

	return m
}
