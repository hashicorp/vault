// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

// TODO Sniffer
/*
sniffer:
- complete for go-hdb: especially call with table parameters
- delete caches for statement and result
- don't ignore part read error
  - example: read scramsha256InitialReply got silently stuck because methodname check failed
- test with python client and handle surprises
  - analyze for not ignoring part read errors
*/

import (
	"bufio"
	"io"
	"net"
	"sync"
)

// A Sniffer is a simple proxy for logging hdb protocol requests and responses.
type Sniffer struct {
	conn   net.Conn
	dbConn net.Conn

	//client
	clRd *bufio.Reader
	clWr *bufio.Writer
	//database
	dbRd *bufio.Reader
	dbWr *bufio.Writer

	// reader
	upRd   *sniffUpReader
	downRd *sniffDownReader
}

// NewSniffer creates a new sniffer instance. The conn parameter is the net.Conn connection, where the Sniffer
// is listening for hdb protocol calls. The dbAddr is the hdb host port address in "host:port" format.
func NewSniffer(conn net.Conn, dbConn net.Conn) *Sniffer {

	//TODO - review setting values here
	trace = true
	debug = true

	s := &Sniffer{
		conn:   conn,
		dbConn: dbConn,
		// buffered write to client
		clWr: bufio.NewWriter(conn),
		// buffered write to db
		dbWr: bufio.NewWriter(dbConn),
	}

	//read from client connection and write to db buffer
	s.clRd = bufio.NewReader(io.TeeReader(conn, s.dbWr))
	//read from db and write to client connection buffer
	s.dbRd = bufio.NewReader(io.TeeReader(dbConn, s.clWr))

	s.upRd = newSniffUpReader(s.clRd)
	s.downRd = newSniffDownReader(s.dbRd)

	return s
}

// Do starts the protocol request and response logging.
func (s *Sniffer) Do() error {
	defer s.dbConn.Close()
	defer s.conn.Close()

	if err := s.upRd.pr.readProlog(); err != nil {
		return err
	}
	if err := s.dbWr.Flush(); err != nil {
		return err
	}
	if err := s.downRd.pr.readProlog(); err != nil {
		return err
	}
	if err := s.clWr.Flush(); err != nil {
		return err
	}

	for {
		//up stream
		if err := s.upRd.readMsg(); err != nil {
			return err // err == io.EOF: connection closed by client
		}
		if err := s.dbWr.Flush(); err != nil {
			return err
		}
		//down stream
		if err := s.downRd.readMsg(); err != nil {
			if _, ok := err.(*hdbErrors); !ok { //if hdbErrors continue
				return err
			}
		}
		if err := s.clWr.Flush(); err != nil {
			return err
		}
	}
}

type sniffReader struct {
	pr *protocolReader
}

func newSniffReader(upStream bool, rd *bufio.Reader) *sniffReader {
	return &sniffReader{pr: newProtocolReader(upStream, rd)}
}

type sniffUpReader struct{ *sniffReader }

func newSniffUpReader(rd *bufio.Reader) *sniffUpReader {
	return &sniffUpReader{sniffReader: newSniffReader(true, rd)}
}

type resMetaCache struct {
	mu    sync.RWMutex
	cache map[uint64]*resultMetadata
}

func newResMetaCache() *resMetaCache {
	return &resMetaCache{cache: make(map[uint64]*resultMetadata)}
}

func (c *resMetaCache) put(stmtID uint64, resMeta *resultMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[stmtID] = resMeta
}

type prmMetaCache struct {
	mu    sync.RWMutex
	cache map[uint64]*parameterMetadata
}

func newPrmMetaCache() *prmMetaCache {
	return &prmMetaCache{cache: make(map[uint64]*parameterMetadata)}
}

func (c *prmMetaCache) put(stmtID uint64, prmMeta *parameterMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[stmtID] = prmMeta
}

func (c *prmMetaCache) get(stmtID uint64) *parameterMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache[stmtID]
}

var _resMetaCache = newResMetaCache()
var _prmMetaCache = newPrmMetaCache()

func (r *sniffUpReader) readMsg() error {
	var stmtID uint64

	return r.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkStatementID:
			r.pr.read((*statementID)(&stmtID))
		// case pkResultMetadata:
		// 	r.pr.read(resMeta)
		case pkParameters:
			prmMeta := _prmMetaCache.get(stmtID)
			prms := &inputParameters{inputFields: prmMeta.parameterFields} // TODO only input parameters
			r.pr.read(prms)
		}
	})
}

type sniffDownReader struct {
	*sniffReader
	resMeta *resultMetadata
	prmMeta *parameterMetadata
}

func newSniffDownReader(rd *bufio.Reader) *sniffDownReader {
	return &sniffDownReader{
		sniffReader: newSniffReader(false, rd),
		resMeta:     &resultMetadata{},
		prmMeta:     &parameterMetadata{},
	}
}

func (r *sniffDownReader) readMsg() error {
	var stmtID uint64
	//resMeta := &resultMetadata{}
	//prmMeta := &parameterMetadata{}

	if err := r.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkStatementID:
			r.pr.read((*statementID)(&stmtID))
		case pkResultMetadata:
			r.pr.read(r.resMeta)
		case pkParameterMetadata:
			r.pr.read(r.prmMeta)
		case pkOutputParameters:
			outFields := []*ParameterField{}
			for _, f := range r.prmMeta.parameterFields {
				if f.Out() {
					outFields = append(outFields, f)
				}
			}
			outPrms := &outputParameters{outputFields: outFields}
			r.pr.read(outPrms)
		case pkResultset:
			resSet := &resultset{resultFields: r.resMeta.resultFields}
			r.pr.read(resSet)
		}
	}); err != nil {
		return err
	}
	_resMetaCache.put(stmtID, r.resMeta)
	_prmMetaCache.put(stmtID, r.prmMeta)
	return nil
}
