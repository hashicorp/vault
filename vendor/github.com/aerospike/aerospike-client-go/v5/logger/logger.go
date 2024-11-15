// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"log"
	"os"
)

// LogPriority specifies the logging level for the client
type LogPriority int

const (
	// DEBUG log level
	DEBUG LogPriority = iota - 1
	// INFO log level
	INFO
	// WARNING log level
	WARNING
	// ERR log level
	ERR
	// OFF log level
	OFF LogPriority = 999
)

type genericLogger interface {
	Printf(format string, v ...interface{})
}

type logger struct {
	Logger genericLogger

	level LogPriority
}

// Logger is the default logger instance
var Logger = newLogger()

func newLogger() *logger {
	return &logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  OFF,
	}
}

// SetLogger sets the *log.Logger object where log messages should be sent to.
// This method is not goroutine-safe, and is not designed to be accessed
// from multiple goroutines.
func (lgr *logger) SetLogger(l genericLogger) {
	lgr.Logger = l
}

// SetLevel sets logging level. Default is ERR.
// This method is not goroutine-safe, and is not designed to be accessed
// from multiple goroutines.
func (lgr *logger) SetLevel(level LogPriority) {
	lgr.level = level
}

// LogAtLevel will logs a message at the level requested.
func (lgr *logger) LogAtLevel(level LogPriority, format string, v ...interface{}) {
	switch level {
	case DEBUG:
		lgr.Debug(format, v...)
	case INFO:
		lgr.Info(format, v...)
	case WARNING:
		lgr.Warn(format, v...)
	case ERR:
		lgr.Error(format, v...)
	case OFF:
	}
}

// Debug logs a message if log level allows to do so.
func (lgr *logger) Debug(format string, v ...interface{}) {
	if lgr.level <= DEBUG {
		if l, ok := lgr.Logger.(*log.Logger); ok {
			l.Output(2, fmt.Sprintf(format, v...))
		} else {
			lgr.Logger.Printf(format, v...)
		}
	}
}

// Info logs a message if log level allows to do so.
func (lgr *logger) Info(format string, v ...interface{}) {
	if lgr.level <= INFO {
		if l, ok := lgr.Logger.(*log.Logger); ok {
			l.Output(2, fmt.Sprintf(format, v...))
		} else {
			lgr.Logger.Printf(format, v...)
		}
	}
}

// Warn logs a message if log level allows to do so.
func (lgr *logger) Warn(format string, v ...interface{}) {
	if lgr.level <= WARNING {
		if l, ok := lgr.Logger.(*log.Logger); ok {
			l.Output(2, fmt.Sprintf(format, v...))
		} else {
			lgr.Logger.Printf(format, v...)
		}
	}
}

// Error logs a message if log level allows to do so.
func (lgr *logger) Error(format string, v ...interface{}) {
	if lgr.level <= ERR {
		if l, ok := lgr.Logger.(*log.Logger); ok {
			l.Output(2, fmt.Sprintf(format, v...))
		} else {
			lgr.Logger.Printf(format, v...)
		}
	}
}
