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
	"fmt"

	"github.com/SAP/go-hdb/internal/bufio"
)

// data format version
const (
	dfvBaseline intType = 1
	dfvDoNotUse intType = 3
	dfvSPS06    intType = 4 //see docu
	dfvBINTEXT  intType = 6
)

// client distribution mode
const (
	cdmOff                 intType = 0
	cdmConnection                  = 1
	cdmStatement                   = 2
	cdmConnectionStatement         = 3
)

// distribution protocol version
const (
	dpvBaseline                       = 0
	dpvClientHandlesStatementSequence = 1
)

type connectOptions struct {
	po      plainOptions
	_numArg int
}

func newConnectOptions() *connectOptions {
	return &connectOptions{
		po: plainOptions{},
	}
}

func (o *connectOptions) String() string {
	m := make(map[connectOption]interface{})
	for k, v := range o.po {
		m[connectOption(k)] = v
	}
	return fmt.Sprintf("%s", m)
}

func (o *connectOptions) kind() partKind {
	return pkConnectOptions
}

func (o *connectOptions) size() (int, error) {
	return o.po.size(), nil
}

func (o *connectOptions) numArg() int {
	return len(o.po)
}

func (o *connectOptions) setNumArg(numArg int) {
	o._numArg = numArg
}

func (o *connectOptions) set(k connectOption, v interface{}) {
	o.po[int8(k)] = v
}

func (o *connectOptions) get(k connectOption) (interface{}, bool) {
	v, ok := o.po[int8(k)]
	return v, ok
}

func (o *connectOptions) read(rd *bufio.Reader) error {
	o.po.read(rd, o._numArg)

	if trace {
		outLogger.Printf("connect options: %v", o)
	}

	return rd.GetError()
}

func (o *connectOptions) write(wr *bufio.Writer) error {
	o.po.write(wr)

	if trace {
		outLogger.Printf("connect options: %v", o)
	}

	return nil
}
