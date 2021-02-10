// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type statementID uint64

func (id statementID) String() string { return fmt.Sprintf("%d", id) }
func (id *statementID) decode(dec *encoding.Decoder, ph *partHeader) error {
	*id = statementID(dec.Uint64())
	return dec.Error()
}
func (id statementID) encode(enc *encoding.Encoder) error { enc.Uint64(uint64(id)); return nil }
