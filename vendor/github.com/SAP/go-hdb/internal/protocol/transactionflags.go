// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type transactionFlags plainOptions

func (f transactionFlags) String() string {
	typedSc := make(map[transactionFlagType]interface{})
	for k, v := range f {
		typedSc[transactionFlagType(k)] = v
	}
	return fmt.Sprintf("flags %s", typedSc)
}

func (f *transactionFlags) decode(dec *encoding.Decoder, ph *partHeader) error {
	*f = transactionFlags{} // no reuse of maps - create new one
	plainOptions(*f).decode(dec, ph.numArg())
	return dec.Error()
}
