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

	messageFn := func(t1, t2 time.Time) *Message {
		return &Message{
			StartTime: t1,
			EndTime:   t2,
		}
	}

	for _, testcase := range []struct {
		name      string
		time1     time.Time
		time2     time.Time
		assertion func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:      "same times",
			time1:     time1,
			time2:     time1,
			assertion: assert.False,
		},
		{
			name:      "reversed times",
			time1:     time2,
			time2:     time1,
			assertion: assert.False,
		},
		{
			name:      "proper times",
			time1:     time1,
			time2:     time2,
			assertion: assert.True,
		},
	} {
		testcase.assertion(t, messageFn(testcase.time1, testcase.time2).ValidateStartAndEndTimes(), testcase.name)
	}
}

// TestMessageValidateMessageType verifies the behaviour of the
// (*Message).ValidateMessageType method in each of its conditional branches.
func TestMessageValidateMessageType(t *testing.T) {
	message := Message{
		Type: "banner",
	}
	assert.True(t, message.ValidateMessageType())

	message.Type = "modal"
	assert.True(t, message.ValidateMessageType())

	message.Type = "something"
	assert.False(t, message.ValidateMessageType())
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

	messageFn := func(t1, t2 time.Time, active *bool) *Message {
		return &Message{
			StartTime: t1,
			EndTime:   t2,
			active:    active,
		}
	}

	for _, testcase := range []struct {
		name      string
		time1     time.Time
		time2     time.Time
		assertion func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:      "same times in past",
			time1:     time1,
			time2:     time1,
			assertion: assert.False,
		},
		{
			name:      "same times in future",
			time1:     time3,
			time2:     time3,
			assertion: assert.False,
		},
		{
			name:      "different times in past",
			time1:     time1,
			time2:     time2,
			assertion: assert.False,
		},
		{
			name:      "different times in past, reversed",
			time1:     time2,
			time2:     time1,
			assertion: assert.False,
		},
		{
			name:      "start in past, end in future",
			time1:     time2,
			time2:     time3,
			assertion: assert.True,
		},
		{
			name:      "start in future, end in past",
			time1:     time3,
			time2:     time2,
			assertion: assert.False,
		},
		{
			name:      "different times in future",
			time1:     time3,
			time2:     time4,
			assertion: assert.False,
		},
		{
			name:      "different times in future, reversed",
			time1:     time4,
			time2:     time3,
			assertion: assert.False,
		},
	} {
		message := messageFn(testcase.time1, testcase.time2, nil)
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

	messageFn := func(t1, t2 time.Time, authenticated bool, messageType string) *Message {
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
		endtime              time.Time
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
			endtime:   time3,
			filter:    filterFn(nil, &trueValue, ""),
			assertion: assert.True,
		},
		{
			name:      "active: filter active-false",
			starttime: time2,
			endtime:   time3,
			filter:    filterFn(nil, &falseValue, ""),
			assertion: assert.False,
		},
		{
			name:      "inactive: filter active-true",
			starttime: time3,
			endtime:   time4,
			filter:    filterFn(nil, &trueValue, ""),
			assertion: assert.False,
		},
		{
			name:      "inactive: filter active-false",
			starttime: time3,
			endtime:   time4,
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
			messageType: "banner",
			filter:      filterFn(nil, nil, "banner"),
			assertion:   assert.True,
		},
		{
			name:        "banner: filter type-modal",
			messageType: "banner",
			filter:      filterFn(nil, nil, "modal"),
			assertion:   assert.False,
		},
		{
			name:        "modal: filter type-banner",
			messageType: "modal",
			filter:      filterFn(nil, nil, "banner"),
			assertion:   assert.False,
		},
		{
			name:        "modal: filter type-modal",
			messageType: "modal",
			filter:      filterFn(nil, nil, "modal"),
			assertion:   assert.True,
		},
	} {
		message := messageFn(testcase.starttime, testcase.endtime, testcase.messageAuthenticated, testcase.messageType)

		testcase.assertion(t, message.Matches(testcase.filter), testcase.name)
	}
}
