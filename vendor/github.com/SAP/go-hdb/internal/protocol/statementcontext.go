// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type statementContext plainOptions

func (c statementContext) String() string {
	typedSc := make(map[statementContextType]interface{})
	for k, v := range c {
		typedSc[statementContextType(k)] = v
	}
	return fmt.Sprintf("options %s", typedSc)
}

func (c *statementContext) decode(dec *encoding.Decoder, ph *partHeader) error {
	*c = statementContext{} // no reuse of maps - create new one
	plainOptions(*c).decode(dec, ph.numArg())
	return dec.Error()
}
