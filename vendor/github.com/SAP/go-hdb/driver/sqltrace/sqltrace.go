/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqltrace

import (
	"flag"
	"log"
	"os"
	"sync"
)

type sqlTrace struct {
	mu sync.RWMutex //protects field on
	on bool
	*log.Logger
}

func newSQLTrace() *sqlTrace {
	return &sqlTrace{
		Logger: log.New(os.Stdout, "hdb ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

var tracer = newSQLTrace()

func init() {
	flag.BoolVar(&tracer.on, "hdb.sqlTrace", false, "enabling hdb sql trace")
}

// On returns if tracing methods output is active.
func On() bool {
	tracer.mu.RLock()
	on := tracer.on
	tracer.mu.RUnlock()
	return on
}

// SetOn sets tracing methods output active or inactive.
func SetOn(on bool) {
	tracer.mu.Lock()
	tracer.on = on
	tracer.mu.Unlock()
}

// Trace calls trace logger Print method to print to the trace logger.
func Trace(v ...interface{}) {
	if On() {
		tracer.Print(v...)
	}
}

// Tracef calls trace logger Printf method to print to the trace logger.
func Tracef(format string, v ...interface{}) {
	if On() {
		tracer.Printf(format, v...)
	}
}

// Traceln calls trace logger Println method to print to the trace logger.
func Traceln(v ...interface{}) {
	if On() {
		tracer.Println(v...)
	}
}
