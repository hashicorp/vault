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

package protocol

import (
	"github.com/SAP/go-hdb/internal/bufio"
)

//fetch size
type fetchsize int32

func (s fetchsize) kind() partKind {
	return pkFetchSize
}

func (s fetchsize) size() (int, error) {
	return 4, nil
}

func (s fetchsize) numArg() int {
	return 1
}

func (s fetchsize) write(wr *bufio.Writer) error {
	wr.WriteInt32(int32(s))

	if trace {
		outLogger.Printf("fetchsize: %d", s)
	}

	return nil
}
