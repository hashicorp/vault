// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type clientContext plainOptions

func (c clientContext) String() string {
	m := make(map[clientContextOption]interface{})
	for k, v := range c {
		m[clientContextOption(k)] = v
	}
	return fmt.Sprintf("client context %s", m)
}

func (c clientContext) size() int   { return plainOptions(c).size() }
func (c clientContext) numArg() int { return len(c) }

func (c *clientContext) decode(dec *encoding.Decoder, ph *partHeader) error {
	*c = clientContext{} // no reuse of maps - create new one
	plainOptions(*c).decode(dec, ph.numArg())
	return dec.Error()
}

func (c clientContext) encode(enc *encoding.Encoder) error {
	plainOptions(c).encode(enc)
	return nil
}
