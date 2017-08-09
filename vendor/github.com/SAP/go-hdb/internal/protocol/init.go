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
	var err error

	if err := rd.Skip(initRequestFillerSize); err != nil { //filler
		return err
	}

	if r.product.major, err = rd.ReadInt8(); err != nil {
		return err
	}
	if r.product.minor, err = rd.ReadInt16(); err != nil {
		return err
	}
	if r.protocol.major, err = rd.ReadInt8(); err != nil {
		return err
	}
	if r.protocol.minor, err = rd.ReadInt16(); err != nil {
		return err
	}
	if err := rd.Skip(1); err != nil { //reserved filler
		return err
	}
	if r.numOptions, err = rd.ReadInt8(); err != nil {
		return err
	}

	switch r.numOptions {
	default:
		outLogger.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		if err := rd.Skip(2); err != nil {
			return err
		}

	case 1:
		if cnt, err := rd.ReadInt8(); err == nil {
			if cnt != 1 {
				return fmt.Errorf("endianess %d - 1 expected", cnt)
			}
		} else {
			return err
		}
		_endianess, err := rd.ReadInt8()
		if err != nil {
			return err
		}
		r.endianess = endianess(_endianess)
	}

	if trace {
		outLogger.Printf("read %s", r)
	}

	return nil
}

func (r *initRequest) write(wr *bufio.Writer) error {

	if err := wr.WriteUint32(initRequestFiller); err != nil {
		return err
	}
	if err := wr.WriteInt8(r.product.major); err != nil {
		return err
	}
	if err := wr.WriteInt16(r.product.minor); err != nil {
		return err
	}
	if err := wr.WriteInt8(r.protocol.major); err != nil {
		return err
	}
	if err := wr.WriteInt16(r.protocol.minor); err != nil {
		return err
	}

	switch r.numOptions {
	default:
		outLogger.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		if err := wr.WriteZeroes(4); err != nil {
			return err
		}

	case 1:
		// reserved
		if err := wr.WriteZeroes(1); err != nil {
			return err
		}
		if err := wr.WriteInt8(r.numOptions); err != nil {
			return err
		}
		if err := wr.WriteInt8(int8(okEndianess)); err != nil {
			return err
		}
		if err := wr.WriteInt8(int8(r.endianess)); err != nil {
			return err
		}
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
	var err error

	if r.product.major, err = rd.ReadInt8(); err != nil {
		return err
	}
	if r.product.minor, err = rd.ReadInt16(); err != nil {
		return err
	}
	if r.protocol.major, err = rd.ReadInt8(); err != nil {
		return err
	}
	if r.protocol.minor, err = rd.ReadInt16(); err != nil {
		return err
	}

	if err := rd.Skip(2); err != nil { //commitInitReplySize
		return err
	}

	if trace {
		outLogger.Printf("read %s", r)
	}

	return nil
}

func (r *initReply) write(wr *bufio.Writer) error {
	if err := wr.WriteInt8(r.product.major); err != nil {
		return err
	}
	if err := wr.WriteInt16(r.product.minor); err != nil {
		return err
	}
	if err := wr.WriteInt8(r.product.major); err != nil {
		return err
	}
	if err := wr.WriteInt16(r.protocol.minor); err != nil {
		return err
	}

	if err := wr.WriteZeroes(2); err != nil { // commitInitReplySize
		return err
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
