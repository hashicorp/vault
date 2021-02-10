// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type clientInfo keyValues

func (c clientInfo) String() string {
	return fmt.Sprintf("client info %s", keyValues(c))
}

func (c clientInfo) size() int   { return keyValues(c).size() }
func (c clientInfo) numArg() int { return len(c) }

func (c *clientInfo) setMap(m map[string]string) {
	*c = clientInfo(m)
}

func (c clientInfo) set(k, v string) {
	c[k] = v
}

func (c clientInfo) get(k string) (string, bool) {
	v, ok := c[k]
	return v, ok
}

func (c *clientInfo) decode(dec *encoding.Decoder, ph *partHeader) error {
	*c = clientInfo{} // no reuse of maps - create new one
	keyValues(*c).decode(dec, ph.numArg())
	return dec.Error()
}

func (c clientInfo) encode(enc *encoding.Encoder) error {
	keyValues(c).encode(enc)
	return nil
}
