// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	timezones           map[int]*time.Location
	updateTimezoneMutex *sync.Mutex
)

// Location returns an offset (minutes) based Location object for Snowflake database.
func Location(offset int) *time.Location {
	updateTimezoneMutex.Lock()
	defer updateTimezoneMutex.Unlock()
	loc := timezones[offset]
	if loc != nil {
		return loc
	}
	loc = genTimezone(offset)
	timezones[offset] = loc
	return loc
}

// LocationWithOffsetString returns an offset based Location object. The offset string must consist of sHHMI where one sign
// character '+'/'-' followed by zero filled hours and minutes.
func LocationWithOffsetString(offsets string) (loc *time.Location, err error) {
	if len(offsets) != 5 {
		return nil, &SnowflakeError{
			Number:      ErrInvalidOffsetStr,
			SQLState:    SQLStateInvalidDataTimeFormat,
			Message:     errMsgInvalidOffsetStr,
			MessageArgs: []interface{}{offsets},
		}
	}
	if offsets[0] != '-' && offsets[0] != '+' {
		return nil, &SnowflakeError{
			Number:      ErrInvalidOffsetStr,
			SQLState:    SQLStateInvalidDataTimeFormat,
			Message:     errMsgInvalidOffsetStr,
			MessageArgs: []interface{}{offsets},
		}
	}
	s := 1
	if offsets[0] == '-' {
		s = -1
	}
	var h, m int64
	h, err = strconv.ParseInt(offsets[1:3], 10, 64)
	if err != nil {
		return
	}
	m, err = strconv.ParseInt(offsets[3:], 10, 64)
	if err != nil {
		return
	}
	offset := s * (int(h)*60 + int(m))
	loc = Location(offset)
	return
}

func genTimezone(offset int) *time.Location {
	var offsetSign string
	var toffset int
	if offset < 0 {
		offsetSign = "-"
		toffset = -offset
	} else {
		offsetSign = "+"
		toffset = offset
	}
	logger.Debugf("offset: %v", offset)
	return time.FixedZone(
		fmt.Sprintf("%v%02d%02d",
			offsetSign, toffset/60, toffset%60), int(offset)*60)
}

func init() {
	updateTimezoneMutex = &sync.Mutex{}
	timezones = make(map[int]*time.Location, 48)
	// pre-generate all common timezones
	for i := -720; i <= 720; i += 30 {
		logger.Debugf("offset: %v", i)
		timezones[i] = genTimezone(i)
	}
}

// retrieve current location based on connection
func getCurrentLocation(params map[string]*string) *time.Location {
	loc := time.Now().Location()
	var err error
	paramsMutex.Lock()
	if tz, ok := params["timezone"]; ok && tz != nil {
		loc, err = time.LoadLocation(*tz)
		if err != nil {
			loc = time.Now().Location()
		}
	}
	paramsMutex.Unlock()
	return loc
}
