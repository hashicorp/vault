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
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

// cesu8 command
type command []byte

func (c command) kind() partKind {
	return pkCommand
}

func (c command) size() (int, error) {
	return cesu8.Size(c), nil
}

func (c command) numArg() int {
	return 1
}

func (c command) write(wr *bufio.Writer) error {
	wr.WriteCesu8(c)

	if trace {
		outLogger.Printf("command: %s", c)
	}

	return nil
}
