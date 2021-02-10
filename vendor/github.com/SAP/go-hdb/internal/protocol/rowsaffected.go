// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

//rows affected
const (
	raSuccessNoInfo   = -2
	raExecutionFailed = -3
)

//rows affected
type rowsAffected []int32

func (r rowsAffected) String() string {
	return fmt.Sprintf("%v", []int32(r))
}

func (r *rowsAffected) reset(numArg int) {
	if r == nil || numArg > cap(*r) {
		*r = make(rowsAffected, numArg)
	} else {
		*r = (*r)[:numArg]
	}
}

func (r *rowsAffected) decode(dec *encoding.Decoder, ph *partHeader) error {
	r.reset(ph.numArg())

	for i := 0; i < ph.numArg(); i++ {
		(*r)[i] = dec.Int32()
	}
	return dec.Error()
}

func (r rowsAffected) total() int64 {
	if r == nil {
		return 0
	}

	total := int64(0)
	for _, rows := range r {
		if rows > 0 {
			total += int64(rows)
		}
	}
	return total
}
