/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"bytes"
	"fmt"
	"log"
)

type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type nopLogger struct{}

func (n nopLogger) Print(_ ...interface{}) {}

func (n nopLogger) Printf(_ string, _ ...interface{}) {}

func (n nopLogger) Println(_ ...interface{}) {}

type testLogger struct {
	capture bytes.Buffer
}

func (l *testLogger) Print(v ...interface{})                 { fmt.Fprint(&l.capture, v...) }
func (l *testLogger) Printf(format string, v ...interface{}) { fmt.Fprintf(&l.capture, format, v...) }
func (l *testLogger) Println(v ...interface{})               { fmt.Fprintln(&l.capture, v...) }
func (l *testLogger) String() string                         { return l.capture.String() }

type defaultLogger struct{}

func (l *defaultLogger) Print(v ...interface{})                 { log.Print(v...) }
func (l *defaultLogger) Printf(format string, v ...interface{}) { log.Printf(format, v...) }
func (l *defaultLogger) Println(v ...interface{})               { log.Println(v...) }

// Logger for logging messages.
// Deprecated: Use ClusterConfig.Logger instead.
var Logger StdLogger = &defaultLogger{}
