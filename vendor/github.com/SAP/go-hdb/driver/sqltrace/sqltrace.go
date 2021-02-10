// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package sqltrace

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync/atomic"
)

func boolToInt64(f bool) int64 {
	if f {
		return 1
	}
	return 0
}

type sqlTrace struct {
	// 64-bit alignment
	on int64 // atomic access (0: false, 1:true)

	*log.Logger
}

func newSQLTrace() *sqlTrace {
	return &sqlTrace{
		Logger: log.New(os.Stdout, "hdb ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
func (t *sqlTrace) On() bool      { return atomic.LoadInt64(&t.on) != 0 }
func (t *sqlTrace) SetOn(on bool) { atomic.StoreInt64(&t.on, boolToInt64(on)) }

type flagValue struct {
	t *sqlTrace
}

func (v flagValue) IsBoolFlag() bool { return true }

func (v flagValue) String() string {
	if v.t == nil {
		return ""
	}
	return strconv.FormatBool(v.t.On())
}

func (v flagValue) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.t.SetOn(f)
	return nil
}

var tracer = newSQLTrace()

func init() {
	flag.Var(&flagValue{t: tracer}, "hdb.sqlTrace", "enabling hdb sql trace")
}

// On returns if tracing methods output is active.
func On() bool { return tracer.On() }

// SetOn sets tracing methods output active or inactive.
func SetOn(on bool) { tracer.SetOn(on) }

// Trace calls trace logger Print method to print to the trace logger.
func Trace(v ...interface{}) { tracer.Print(v...) }

// Tracef calls trace logger Printf method to print to the trace logger.
func Tracef(format string, v ...interface{}) { tracer.Printf(format, v...) }

// Traceln calls trace logger Println method to print to the trace logger.
func Traceln(v ...interface{}) { tracer.Println(v...) }
