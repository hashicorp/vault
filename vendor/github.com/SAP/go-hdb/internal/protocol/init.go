/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/bufio"
)

const (
	okEndianess int8 = 1
)

const (
	initRequestFillerSize = 4
)

var initRequestFiller uint32 = 0xffffffff

type productVersion struct {
	major int8
	minor int16
}

func (v *productVersion) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

type protocolVersion struct {
	major int8
	minor int16
}

func (v *protocolVersion) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

type version struct {
	major int8
	minor int16
}

func (v *version) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

type initRequest struct {
	product    *version
	protocol   *version
	numOptions int8
	endianess  endianess
}

func newInitRequest() *initRequest {
	return &initRequest{
		product:  new(version),
		protocol: new(version),
	}
}

func (r *initRequest) String() string {
	switch r.numOptions {
	default:
		return fmt.Sprintf("init request: product version %s protocol version %s", r.product, r.protocol)
	case 1:
		return fmt.Sprintf("init request: product version %s protocol version %s endianess %s", r.product, r.protocol, r.endianess)
	}
}

func (r *initRequest) read(rd *bufio.Reader) error {
	rd.Skip(initRequestFillerSize) //filler
	r.product.major = rd.ReadInt8()
	r.product.minor = rd.ReadInt16()
	r.protocol.major = rd.ReadInt8()
	r.protocol.minor = rd.ReadInt16()
	rd.Skip(1) //reserved filler
	r.numOptions = rd.ReadInt8()

	switch r.numOptions {
	default:
		outLogger.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		rd.Skip(2)

	case 1:
		cnt := rd.ReadInt8()
		if cnt != 1 {
			outLogger.Fatalf("endianess %d - 1 expected", cnt)
		}
		r.endianess = endianess(rd.ReadInt8())
	}

	if trace {
		outLogger.Printf("read %s", r)
	}

	return rd.GetError()
}

func (r *initRequest) write(wr *bufio.Writer) error {
	wr.WriteUint32(initRequestFiller)
	wr.WriteInt8(r.product.major)
	wr.WriteInt16(r.product.minor)
	wr.WriteInt8(r.protocol.major)
	wr.WriteInt16(r.protocol.minor)

	switch r.numOptions {
	default:
		outLogger.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		wr.WriteZeroes(4)

	case 1:
		// reserved
		wr.WriteZeroes(1)
		wr.WriteInt8(r.numOptions)
		wr.WriteInt8(int8(okEndianess))
		wr.WriteInt8(int8(r.endianess))

	}

	// flush
	if err := wr.Flush(); err != nil {
		return err
	}

	if trace {
		outLogger.Printf("write %s", r)
	}

	return nil
}

type initReply struct {
	product  *version
	protocol *version
}

func newInitReply() *initReply {
	return &initReply{
		product:  new(version),
		protocol: new(version),
	}
}

func (r *initReply) String() string {
	return fmt.Sprintf("init reply: product version %s protocol version %s", r.product, r.protocol)
}

func (r *initReply) read(rd *bufio.Reader) error {
	r.product.major = rd.ReadInt8()
	r.product.minor = rd.ReadInt16()
	r.protocol.major = rd.ReadInt8()
	r.protocol.minor = rd.ReadInt16()
	rd.Skip(2) //commitInitReplySize

	if trace {
		outLogger.Printf("read %s", r)
	}

	return rd.GetError()
}

func (r *initReply) write(wr *bufio.Writer) error {
	wr.WriteInt8(r.product.major)
	wr.WriteInt16(r.product.minor)
	wr.WriteInt8(r.product.major)
	wr.WriteInt16(r.protocol.minor)
	wr.WriteZeroes(2) // commitInitReplySize

	// flush
	if err := wr.Flush(); err != nil {
		return err
	}

	if trace {
		outLogger.Printf("write %s", r)
	}

	return nil
}
