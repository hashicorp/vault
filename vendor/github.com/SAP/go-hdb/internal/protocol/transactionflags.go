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

type transactionFlags struct {
	options plainOptions
	_numArg int
}

func newTransactionFlags() *transactionFlags {
	return &transactionFlags{
		options: plainOptions{},
	}
}

func (f *transactionFlags) String() string {
	typedSc := make(map[transactionFlagType]interface{})
	for k, v := range f.options {
		typedSc[transactionFlagType(k)] = v
	}
	return fmt.Sprintf("%s", typedSc)
}

func (f *transactionFlags) kind() partKind {
	return pkTransactionFlags
}

func (f *transactionFlags) setNumArg(numArg int) {
	f._numArg = numArg
}

func (f *transactionFlags) read(rd *bufio.Reader) error {
	f.options.read(rd, f._numArg)

	if trace {
		outLogger.Printf("transaction flags: %v", f)
	}

	return rd.GetError()
}
