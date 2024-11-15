// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package logger

import "strings"

// DiffToInfo is the number of levels in the Go Driver that come before the
// "Info" level. This should ensure that "Info" is the 0th level passed to the
// sink.
const DiffToInfo = 1

// Level is an enumeration representing the log severity levels supported by
// the driver. The order of the logging levels is important. The driver expects
// that a user will likely use the "logr" package to create a LogSink, which
// defaults InfoLevel as 0. Any additions to the Level enumeration before the
// InfoLevel will need to also update the "diffToInfo" constant.
type Level int

const (
	// LevelOff suppresses logging.
	LevelOff Level = iota

	// LevelInfo enables logging of informational messages. These logs are
	// high-level information about normal driver behavior.
	LevelInfo

	// LevelDebug enables logging of debug messages. These logs can be
	// voluminous and are intended for detailed information that may be
	// helpful when debugging an application.
	LevelDebug
)

const (
	levelLiteralOff       = "off"
	levelLiteralEmergency = "emergency"
	levelLiteralAlert     = "alert"
	levelLiteralCritical  = "critical"
	levelLiteralError     = "error"
	levelLiteralWarning   = "warning"
	levelLiteralNotice    = "notice"
	levelLiteralInfo      = "info"
	levelLiteralDebug     = "debug"
	levelLiteralTrace     = "trace"
)

var LevelLiteralMap = map[string]Level{
	levelLiteralOff:       LevelOff,
	levelLiteralEmergency: LevelInfo,
	levelLiteralAlert:     LevelInfo,
	levelLiteralCritical:  LevelInfo,
	levelLiteralError:     LevelInfo,
	levelLiteralWarning:   LevelInfo,
	levelLiteralNotice:    LevelInfo,
	levelLiteralInfo:      LevelInfo,
	levelLiteralDebug:     LevelDebug,
	levelLiteralTrace:     LevelDebug,
}

// ParseLevel will check if the given string is a valid environment variable
// for a logging severity level. If it is, then it will return the associated
// driver's Level. The default Level is “LevelOff”.
func ParseLevel(str string) Level {
	for literal, level := range LevelLiteralMap {
		if strings.EqualFold(literal, str) {
			return level
		}
	}

	return LevelOff
}
