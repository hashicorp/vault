// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

// client distribution mode
const (
	cdmOff                 optIntType = 0
	cdmConnection          optIntType = 1
	cdmStatement           optIntType = 2
	cdmConnectionStatement optIntType = 3
)

// distribution protocol version
const (
	dpvBaseline                       = 0
	dpvClientHandlesStatementSequence = 1
)

type connectOptions plainOptions

func (o connectOptions) String() string {
	m := make(map[connectOption]interface{})
	for k, v := range o {
		m[connectOption(k)] = v
	}
	return fmt.Sprintf("options %s", m)
}

func (o connectOptions) size() int   { return plainOptions(o).size() }
func (o connectOptions) numArg() int { return len(o) }

func (o connectOptions) fullVersionString() (version string) {
	v, ok := o[int8(coFullVersionString)]
	if !ok {
		return
	}
	if s, ok := v.(optStringType); ok {
		return string(s)
	}
	return
}

func (o *connectOptions) decode(dec *encoding.Decoder, ph *partHeader) error {
	*o = connectOptions{} // no reuse of maps - create new one
	plainOptions(*o).decode(dec, ph.numArg())
	return dec.Error()
}

func (o connectOptions) encode(enc *encoding.Encoder) error {
	plainOptions(o).encode(enc)
	return nil
}
