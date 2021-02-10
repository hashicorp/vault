// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

//fetch size
type fetchsize int32

func (s fetchsize) String() string { return fmt.Sprintf("fetchsize %d", s) }
func (s *fetchsize) decode(dec *encoding.Decoder, ph *partHeader) error {
	*s = fetchsize(dec.Int32())
	return dec.Error()
}
func (s fetchsize) encode(enc *encoding.Encoder) error { enc.Int32(int32(s)); return nil }
