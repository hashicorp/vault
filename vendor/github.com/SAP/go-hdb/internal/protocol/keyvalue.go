// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type keyValues map[string]string

func (kv keyValues) size() int {
	size := 0
	for k, v := range kv {
		size += cesu8Type.prmSize(k)
		size += cesu8Type.prmSize(v)
	}
	return size
}

func (kv keyValues) decode(dec *encoding.Decoder, cnt int) {
	for i := 0; i < cnt; i++ {
		k, err := cesu8Type.decodeRes(dec)
		if err != nil {
			plog.Fatalf(err.Error())
		}
		v, err := cesu8Type.decodeRes(dec)
		if err != nil {
			plog.Fatalf(err.Error())
		}
		kv[string(k.([]byte))] = string(v.([]byte)) // set key value
	}
}

func (kv keyValues) encode(enc *encoding.Encoder) {
	for k, v := range kv {
		if err := cesu8Type.encodePrm(enc, k); err != nil {
			plog.Fatalf(err.Error())
		}
		if err := cesu8Type.encodePrm(enc, v); err != nil {
			plog.Fatalf(err.Error())
		}
	}
}
