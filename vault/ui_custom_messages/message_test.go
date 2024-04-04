// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMessageValidateStartAndEndTimes verifies the behaviour of the
// (*Message).ValidateStartAndEndTimes method in each of the conditional
// branches.
func TestMessageValidateStartAndEndTimes(t *testing.T) {
	var (
		time1 = time.Now()
		time2 = time1.Add(time.Minute)
	)

	messageFn := func(t1 time.Time, t2 *time.Time) *Message {
		return &Message{
			StartTime: t1,
			EndTime:   t2,
		}
	}

	for _, testcase := range []struct {
		name      string
		time1     time.Time
		time2     *time.Time
		assertion func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:      "same times",
			time1:     time1,
			time2:     &time1,
			assertion: assert.False,
		},
		{
			name:      "reversed times",
			time1:     time2,
			time2:     &time1,
			assertion: assert.False,
		},
		{
			name:      "proper times",
			time1:     time1,
			time2:     &time2,
			assertion: assert.True,
		},
		{
			name:      "no end time",
			time1:     time1,
			time2:     nil,
			assertion: assert.True,
		},
	} {
		testcase.assertion(t, messageFn(testcase.time1, testcase.time2).HasValidStartAndEndTimes(), testcase.name)
	}
}

// TestMessageValidateMessageType verifies the behaviour of the
// (*Message).ValidateMessageType method in each of its conditional branches.
func TestMessageValidateMessageType(t *testing.T) {
	message := Message{
		Type: BannerMessageType,
	}
	assert.True(t, message.HasValidMessageType())

	message.Type = ModalMessageType
	assert.True(t, message.HasValidMessageType())

	message.Type = "something"
	assert.False(t, message.HasValidMessageType())
}

func TestMessageActive(t *testing.T) {
	var (
		time1 = time.Now().Add(-2 * time.Hour) // now - 2
		time2 = time1.Add(time.Hour)           // now - 1
		time3 = time2.Add(2 * time.Hour)       // now + 1
		time4 = time3.Add(time.Hour)           // now + 2

		trueValue  = true
		falseValue = false
	)

	messageFn := func(t1 time.Time, t2 *time.Time, active *bool) *Message {
		message := &Message{
			StartTime: t1,
			EndTime:   t2,
			active:    active,
		}

		return message
	}

	for _, testcase := range []struct {
		name      string
		startTime time.Time
		endTime   *time.Time
		assertion func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:      "same times in past",
			startTime: time1,
			endTime:   &time1,
			assertion: assert.False,
		},
		{
			name:      "same times in future",
			startTime: time3,
			endTime:   &time3,
			assertion: assert.False,
		},
		{
			name:      "different times in past",
			startTime: time1,
			endTime:   &time2,
			assertion: assert.False,
		},
		{
			name:      "different times in past, reversed",
			startTime: time2,
			endTime:   &time1,
			assertion: assert.False,
		},
		{
			name:      "start in past, end in future",
			startTime: time2,
			endTime:   &time3,
			assertion: assert.True,
		},
		{
			name:      "start in future, end in past",
			startTime: time3,
			endTime:   &time2,
			assertion: assert.False,
		},
		{
			name:      "different times in future",
			startTime: time3,
			endTime:   &time4,
			assertion: assert.False,
		},
		{
			name:      "different times in future, reversed",
			startTime: time4,
			endTime:   &time3,
			assertion: assert.False,
		},
		{
			name:      "no end time, starting in past",
			startTime: time1,
			assertion: assert.True,
		},
		{
			name:      "no end time, starting in future",
			startTime: time3,
			assertion: assert.False,
		},
	} {
		message := messageFn(testcase.startTime, testcase.endTime, nil)
		testcase.assertion(t, message.Active(), testcase.name)
		assert.NotNil(t, message.active, testcase.name)

		message.active = &trueValue
		assert.True(t, message.Active(), testcase.name)

		message.active = &falseValue
		assert.False(t, message.Active(), testcase.name)
	}
}

func TestMessageMatches(t *testing.T) {
	var (
		time1 = time.Now().Add(-2 * time.Hour) // now - 2
		time2 = time1.Add(time.Hour)           // now - 1
		time3 = time2.Add(2 * time.Hour)       // now + 1
		time4 = time3.Add(time.Hour)           // now + 2

		trueValue  = true
		falseValue = false
	)

	filterFn := func(authenticated, active *bool, messageType string) FindFilter {
		return FindFilter{
			authenticated: authenticated,
			messageType:   messageType,
			active:        active,
		}
	}

	messageFn := func(t1 time.Time, t2 *time.Time, authenticated bool, messageType string) *Message {
		return &Message{
			StartTime:     t1,
			EndTime:       t2,
			Authenticated: authenticated,
			Type:          messageType,
		}
	}

	for _, testcase := range []struct {
		name                 string
		starttime            time.Time
		endtime              *time.Time
		messageAuthenticated bool
		messageType          string
		filter               FindFilter
		assertion            func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:      "empty filter",
			filter:    filterFn(nil, nil, ""),
			assertion: assert.True,
		},
		{
			name:      "active: filter active-true",
			starttime: time2,
			endtime:   &time3,
			filter:    filterFn(nil, &trueValue, ""),
			assertion: assert.True,
		},
		{
			name:      "active (no end time): filter active-true",
			starttime: time2,
			filter:    filterFn(nil, &trueValue, ""),
			assertion: assert.True,
		},
		{
			name:      "active: filter active-false",
			starttime: time2,
			endtime:   &time3,
			filter:    filterFn(nil, &falseValue, ""),
			assertion: assert.False,
		},
		{
			name:      "active (no end time): filter active-false",
			starttime: time2,
			filter:    filterFn(nil, &falseValue, ""),
			assertion: assert.False,
		},
		{
			name:      "inactive: filter active-true",
			starttime: time3,
			endtime:   &time4,
			filter:    filterFn(nil, &trueValue, ""),
			assertion: assert.False,
		},
		{
			name:      "inactive: filter active-false",
			starttime: time3,
			endtime:   &time4,
			filter:    filterFn(nil, &falseValue, ""),
			assertion: assert.True,
		},
		{
			name:      "pre-login: filter authenticated-false",
			filter:    filterFn(&falseValue, nil, ""),
			assertion: assert.True,
		},
		{
			name:      "pre-login: filter authenticated-true",
			filter:    filterFn(&trueValue, nil, ""),
			assertion: assert.False,
		},
		{
			name:                 "post-login: filter authenticated-false",
			messageAuthenticated: true,
			filter:               filterFn(&falseValue, nil, ""),
			assertion:            assert.False,
		},
		{
			name:                 "post-login: filter authenticated-true",
			messageAuthenticated: true,
			filter:               filterFn(&trueValue, nil, ""),
			assertion:            assert.True,
		},
		{
			name:        "banner: filter type-banner",
			messageType: BannerMessageType,
			filter:      filterFn(nil, nil, BannerMessageType),
			assertion:   assert.True,
		},
		{
			name:        "banner: filter type-modal",
			messageType: BannerMessageType,
			filter:      filterFn(nil, nil, ModalMessageType),
			assertion:   assert.False,
		},
		{
			name:        "modal: filter type-banner",
			messageType: ModalMessageType,
			filter:      filterFn(nil, nil, BannerMessageType),
			assertion:   assert.False,
		},
		{
			name:        "modal: filter type-modal",
			messageType: ModalMessageType,
			filter:      filterFn(nil, nil, ModalMessageType),
			assertion:   assert.True,
		},
	} {
		message := messageFn(testcase.starttime, testcase.endtime, testcase.messageAuthenticated, testcase.messageType)

		testcase.assertion(t, message.Matches(testcase.filter), testcase.name)
	}
}
