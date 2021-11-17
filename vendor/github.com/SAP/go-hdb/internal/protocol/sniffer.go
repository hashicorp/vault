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
	"log"
	"net"

	"github.com/SAP/go-hdb/internal/bufio"
)

type dir bool

const (
	maxBinarySize = 128
)

type fragment interface {
	read(rd *bufio.Reader) error
	write(wr *bufio.Writer) error
}

func (d dir) String() string {
	if d {
		return "->"
	}
	return "<-"
}

// A Sniffer is a simple proxy for logging hdb protocol requests and responses.
type Sniffer struct {
	conn   net.Conn
	dbAddr string
	dbConn net.Conn

	//client
	clRd *bufio.Reader
	clWr *bufio.Writer
	//database
	dbRd *bufio.Reader
	dbWr *bufio.Writer

	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	buf []byte
}

// NewSniffer creates a new sniffer instance. The conn parameter is the net.Conn connection, where the Sniffer
// is listening for hdb protocol calls. The dbAddr is the hdb host port address in "host:port" format.
func NewSniffer(conn net.Conn, dbAddr string) (*Sniffer, error) {
	s := &Sniffer{
		conn:   conn,
		dbAddr: dbAddr,
		clRd:   bufio.NewReader(conn),
		clWr:   bufio.NewWriter(conn),
		mh:     &messageHeader{},
		sh:     &segmentHeader{},
		ph:     &partHeader{},
		buf:    make([]byte, 0),
	}

	dbConn, err := net.Dial("tcp", s.dbAddr)
	if err != nil {
		return nil, err
	}

	s.dbRd = bufio.NewReader(dbConn)
	s.dbWr = bufio.NewWriter(dbConn)
	s.dbConn = dbConn
	return s, nil
}

func (s *Sniffer) getBuffer(size int) []byte {
	if cap(s.buf) < size {
		s.buf = make([]byte, size)
	}
	return s.buf[:size]
}

// Go starts the protocol request and response logging.
func (s *Sniffer) Go() {
	defer s.dbConn.Close()
	defer s.conn.Close()

	req := newInitRequest()
	if err := s.streamFragment(dir(true), s.clRd, s.dbWr, req); err != nil {
		return
	}

	rep := newInitReply()
	if err := s.streamFragment(dir(false), s.dbRd, s.clWr, rep); err != nil {
		return
	}

	for {
		//up stream
		if err := s.stream(dir(true), s.clRd, s.dbWr); err != nil {
			return
		}
		//down stream
		if err := s.stream(dir(false), s.dbRd, s.clWr); err != nil {
			return
		}
	}
}

func (s *Sniffer) stream(d dir, from *bufio.Reader, to *bufio.Writer) error {

	if err := s.streamFragment(d, from, to, s.mh); err != nil {
		return err
	}

	size := int(s.mh.varPartLength)

	for i := 0; i < int(s.mh.noOfSegm); i++ {

		if err := s.streamFragment(d, from, to, s.sh); err != nil {
			return err
		}

		size -= int(s.sh.segmentLength)

		for j := 0; j < int(s.sh.noOfParts); j++ {

			if err := s.streamFragment(d, from, to, s.ph); err != nil {
				return err
			}

			// protocol error workaraound
			padding := (size == 0) || (j != (int(s.sh.noOfParts) - 1))

			if err := s.streamPart(d, from, to, s.ph, padding); err != nil {
				return err
			}
		}
	}
	return to.Flush()
}

func (s *Sniffer) streamPart(d dir, from *bufio.Reader, to *bufio.Writer, ph *partHeader, padding bool) error {

	switch ph.partKind {

	default:
		return s.streamBinary(d, from, to, int(ph.bufferLength), padding)
	}
}

func (s *Sniffer) streamBinary(d dir, from *bufio.Reader, to *bufio.Writer, size int, padding bool) error {
	var b []byte

	//protocol error workaraound
	if padding {
		pad := padBytes(size)
		b = s.getBuffer(size + pad)
	} else {
		b = s.getBuffer(size)
	}

	from.ReadFull(b)
	err := from.GetError()
	if err != nil {
		log.Print(err)
		return err
	}

	if size > maxBinarySize {
		log.Printf("%s %v", d, b[:maxBinarySize])
	} else {
		log.Printf("%s %v", d, b[:size])
	}
	to.Write(b)
	return nil
}

func (s *Sniffer) streamFragment(d dir, from *bufio.Reader, to *bufio.Writer, f fragment) error {
	if err := f.read(from); err != nil {
		log.Print(err)
		return err
	}
	log.Printf("%s %s", d, f)
	if err := f.write(to); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
