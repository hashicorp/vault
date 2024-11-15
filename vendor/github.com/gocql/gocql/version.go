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

import "runtime/debug"

const (
	defaultDriverName = "github.com/apache/cassandra-gocql-driver"

	// This string MUST have this value since we explicitly test against the
	// current main package returned by runtime/debug below.  Also note the
	// package name used here may change in a future (2.x) release; in that case
	// this constant will be updated as well.
	mainPackage = "github.com/gocql/gocql"
)

var driverName string

var driverVersion string

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, d := range buildInfo.Deps {
			if d.Path == mainPackage {
				driverName = defaultDriverName
				driverVersion = d.Version
				// If there's a replace directive in play for the gocql package
				// then use that information for path and version instead.  This
				// will allow forks or other local packages to clearly identify
				// themselves as distinct from mainPackage above.
				if d.Replace != nil {
					driverName = d.Replace.Path
					driverVersion = d.Replace.Version
				}
				break
			}
		}
	}
}
