// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"os"
	"strconv"
	"strings"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type clientID []byte

func newClientID() clientID {
	if h, err := os.Hostname(); err == nil {
		return clientID(strings.Join([]string{strconv.Itoa(os.Getpid()), h}, "@"))
	}
	return clientID(strconv.Itoa(os.Getpid()))
}

func (id clientID) String() string { return string(id) }
func (id *clientID) resize(size int) {
	if id == nil || size > cap(*id) {
		*id = make([]byte, size)
	} else {
		*id = (*id)[:size]
	}
}
func (id clientID) size() int { return len(id) }
func (id *clientID) decode(dec *encoding.Decoder, ph *partHeader) error {
	id.resize(int(ph.bufferLength))
	dec.Bytes(*id)
	return dec.Error()
}
func (id clientID) encode(enc *encoding.Encoder) error { enc.Bytes(id); return nil }
