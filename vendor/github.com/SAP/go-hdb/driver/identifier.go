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

package driver

import (
	"crypto/rand"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var reSimple = regexp.MustCompile("^[_A-Z][_#$A-Z0-9]*$")

// Identifier in hdb SQL statements like schema or table name.
type Identifier string

// RandomIdentifier returns a random Identifier prefixed by the prefix parameter.
// This function is used to generate database objects with random names for test and example code.
func RandomIdentifier(prefix string) Identifier {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return Identifier(fmt.Sprintf("%s%x", prefix, b))
}

// String implements Stringer interface.
func (i Identifier) String() string {
	s := string(i)
	if reSimple.MatchString(s) {
		return s
	}
	return strconv.Quote(s)
}
