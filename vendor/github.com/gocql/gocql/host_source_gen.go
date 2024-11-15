//go:build genhostinfo
// +build genhostinfo

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

package main

import (
	"fmt"
	"reflect"
	"sync"

	gocql "github.com/gocql/gocql"
)

func gen(clause, field string) {
	fmt.Printf("if h.%s == %s {\n", field, clause)
	fmt.Printf("\th.%s = from.%s\n", field, field)
	fmt.Println("}")
}

func main() {
	t := reflect.ValueOf(&gocql.HostInfo{}).Elem().Type()
	mu := reflect.TypeOf(sync.RWMutex{})

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type == mu {
			continue
		}

		switch f.Type.Kind() {
		case reflect.Slice:
			gen("nil", f.Name)
		case reflect.String:
			gen(`""`, f.Name)
		case reflect.Int:
			gen("0", f.Name)
		case reflect.Struct:
			gen("("+f.Type.Name()+"{})", f.Name)
		case reflect.Bool, reflect.Int32:
			continue
		default:
			panic(fmt.Sprintf("unknown field: %s", f))
		}
	}

}
