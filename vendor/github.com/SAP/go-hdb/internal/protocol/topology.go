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

type topologyInformation struct {
	mlo     multiLineOptions
	_numArg int
}

func newTopologyInformation() *topologyInformation {
	return &topologyInformation{
		mlo: multiLineOptions{},
	}
}

func (o *topologyInformation) String() string {
	mlo := make([]map[topologyOption]interface{}, len(o.mlo))
	for i, po := range o.mlo {
		typedPo := make(map[topologyOption]interface{})
		for k, v := range po {
			typedPo[topologyOption(k)] = v
		}
		mlo[i] = typedPo
	}
	return fmt.Sprintf("%s", mlo)
}

func (o *topologyInformation) kind() partKind {
	return pkTopologyInformation
}

func (o *topologyInformation) size() int {
	return o.mlo.size()
}

func (o *topologyInformation) numArg() int {
	return len(o.mlo)
}

func (o *topologyInformation) setNumArg(numArg int) {
	o._numArg = numArg
}

func (o *topologyInformation) read(rd *bufio.Reader) error {
	o.mlo.read(rd, o._numArg)

	if trace {
		outLogger.Printf("topology options: %v", o)
	}

	return rd.GetError()
}

func (o *topologyInformation) write(wr *bufio.Writer) error {
	for _, m := range o.mlo {
		wr.WriteInt16(int16(len(m)))
		o.mlo.write(wr)
	}

	if trace {
		outLogger.Printf("topology options: %v", o)
	}

	return nil
}
