// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	// . "github.com/aerospike/aerospike-client-go/logger"

	"io"
	"strconv"

	. "github.com/aerospike/aerospike-client-go/types"
)

var aeroerr error = NewAerospikeError(PARSE_ERROR, "Error parsing peers list.")

func parsePeers(cluster *Cluster, node *Node) (*peerListParser, error) {
	cmd := cluster.clientPolicy.peersString()

	info, err := node.RequestInfo(&cluster.infoPolicy, cmd)
	if err != nil {
		return nil, err
	}

	peersStr, exists := info[cmd]
	if !exists {
		return nil, NewAerospikeError(PARSE_ERROR, "Info Command response was empty.")
	}

	p := peerListParser{buf: []byte(peersStr)}
	if err := p.Parse(); err != nil {
		return nil, err
	}

	return &p, nil
}

type peerListParser struct {
	buf []byte
	pos int

	defPort *int64
	gen     *int64
	peers   []*peer
}

func (p *peerListParser) generation() int64 {
	if p.gen != nil {
		return *p.gen
	}
	return 0
}

func (p *peerListParser) Expect(ch byte) bool {
	if p.pos == len(p.buf) {
		return false
	}

	if p.buf[p.pos] == ch {
		p.pos++
		return true
	}
	return false
}

func (p *peerListParser) readByte() *byte {
	if p.pos == len(p.buf) {
		return nil
	}

	ch := p.buf[p.pos]
	p.pos++
	return &ch
}

func (p *peerListParser) PeekByte() *byte {
	if p.pos == len(p.buf) {
		return nil
	}

	ch := p.buf[p.pos]
	return &ch
}

func (p *peerListParser) readInt64() (*int64, error) {
	if p.pos == len(p.buf) {
		return nil, io.EOF
	}

	if p.buf[p.pos] == ',' {
		return nil, nil
	}

	begin := p.pos
	for p.pos < len(p.buf) {
		ch := p.buf[p.pos]
		if ch == ',' {
			break
		}
		p.pos++
	}

	num, err := strconv.ParseInt(string(p.buf[begin:p.pos]), 10, 64)
	if err != nil {
		return nil, err
	}
	return &num, nil
}

func (p *peerListParser) readString() (string, error) {
	if p.pos == len(p.buf) {
		return "", io.EOF
	}

	if p.buf[p.pos] == ',' {
		return "", nil
	}

	begin := p.pos
	bracket := p.buf[p.pos] == '['
	for p.pos < len(p.buf) {
		ch := p.buf[p.pos]
		if ch == ',' {
			break
		}

		if ch == ']' {
			if !bracket {
				break
			}
			bracket = false
		}
		p.pos++
	}

	return string(p.buf[begin:p.pos]), nil
}

func (p *peerListParser) ParseHost(host string) (*Host, error) {
	ppos := -1
	bpos := -1
	for i := 0; i < len(host); i++ {
		switch host[i] {
		case ':':
			ppos = i
		case ']':
			ppos = -1
			bpos = i
		}
	}

	port := 0
	if p.defPort != nil {
		port = int(*p.defPort)
	}
	var err error
	if ppos >= 0 {
		portStr := host[ppos+1:]
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	var addr string
	if bpos >= 0 {
		addr = host[1:bpos]
	} else {
		if ppos >= 0 {
			addr = host[:ppos]
		} else {
			addr = host
		}
	}

	return NewHost(addr, port), nil
}

func (p *peerListParser) readHosts(tlsName string) ([]*Host, error) {
	if !p.Expect('[') {
		return nil, aeroerr
	}

	hostList := []*Host{}
	for {
		hostStr, err := p.readString()
		if err != nil {
			return nil, err
		}

		if hostStr == "" {
			break
		}

		host, err := p.ParseHost(hostStr)
		if err != nil {
			return nil, aeroerr
		}

		host.TLSName = tlsName
		hostList = append(hostList, host)

		if !p.Expect(',') {
			break
		}
	}

	if !p.Expect(']') {
		return nil, aeroerr
	}

	return hostList, nil
}

func (p *peerListParser) readPeer() (*peer, error) {
	if !p.Expect('[') {
		return nil, nil
	}

	nodeName, err := p.readString()
	if err != nil {
		return nil, err
	}

	if !p.Expect(',') {
		return nil, aeroerr
	}
	tlsName, err := p.readString()
	if err != nil {
		return nil, err
	}

	if !p.Expect(',') {
		return nil, aeroerr
	}

	hostList, err := p.readHosts(tlsName)
	if err != nil {
		return nil, err
	}

	if !p.Expect(']') {
		return nil, aeroerr
	}

	nodeData := &peer{nodeName: nodeName, tlsName: tlsName, hosts: hostList}
	return nodeData, nil
}

func (p *peerListParser) readNodeList() ([]*peer, error) {
	ch := p.readByte()
	if ch == nil {
		return nil, nil
	}

	if *ch != '[' {
		return nil, aeroerr
	}

	nodeList := []*peer{}
	for {
		node, err := p.readPeer()
		if err != nil {
			return nil, err
		}

		if node == nil {
			break
		}

		nodeList = append(nodeList, node)

		if !p.Expect(',') {
			break
		}
	}

	if !p.Expect(']') {
		return nil, aeroerr
	}

	return nodeList, nil
}

func (p *peerListParser) Parse() error {
	var err error
	p.gen, err = p.readInt64()
	if err != nil {
		return err
	}

	if !p.Expect(',') {
		return aeroerr
	}

	p.defPort, err = p.readInt64()
	if err != nil {
		return err
	}

	if !p.Expect(',') {
		return aeroerr
	}

	p.peers, err = p.readNodeList()
	if err != nil {
		return err
	}

	return nil
}
