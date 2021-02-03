// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type topologyInformation multiLineOptions

func (o topologyInformation) String() string {
	mlo := make([]map[topologyOption]interface{}, len(o))
	for i, po := range o {
		typedPo := make(map[topologyOption]interface{})
		for k, v := range po {
			typedPo[topologyOption(k)] = v
		}
		mlo[i] = typedPo
	}
	return fmt.Sprintf("options %s", mlo)
}

func (o *topologyInformation) decode(dec *encoding.Decoder, ph *partHeader) error {
	(*multiLineOptions)(o).decode(dec, ph.numArg())
	return dec.Error()
}
