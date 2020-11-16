// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jaeger

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/apache/thrift/lib/go/thrift"

	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
)

// udpPacketMaxLength is the max size of UDP packet we want to send, synced with jaeger-agent
const udpPacketMaxLength = 65000

// agentClientUDP is a UDP client to Jaeger agent that implements gen.Agent interface.
type agentClientUDP struct {
	gen.Agent
	io.Closer

	connUDP       udpConn
	client        *gen.AgentClient
	maxPacketSize int                   // max size of datagram in bytes
	thriftBuffer  *thrift.TMemoryBuffer // buffer used to calculate byte size of a span
}

type udpConn interface {
	Write([]byte) (int, error)
	SetWriteBuffer(int) error
	Close() error
}

type agentClientUDPParams struct {
	HostPort                 string
	MaxPacketSize            int
	Logger                   *log.Logger
	AttemptReconnecting      bool
	AttemptReconnectInterval time.Duration
}

// newAgentClientUDP creates a client that sends spans to Jaeger Agent over UDP.
func newAgentClientUDP(params agentClientUDPParams) (*agentClientUDP, error) {
	// validate hostport
	if _, _, err := net.SplitHostPort(params.HostPort); err != nil {
		return nil, err
	}

	if params.MaxPacketSize <= 0 {
		params.MaxPacketSize = udpPacketMaxLength
	}

	if params.AttemptReconnecting && params.AttemptReconnectInterval <= 0 {
		params.AttemptReconnectInterval = time.Second * 30
	}

	thriftBuffer := thrift.NewTMemoryBufferLen(params.MaxPacketSize)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := gen.NewAgentClientFactory(thriftBuffer, protocolFactory)

	var connUDP udpConn
	var err error

	if params.AttemptReconnecting {
		// host is hostname, setup resolver loop in case host record changes during operation
		connUDP, err = newReconnectingUDPConn(params.HostPort, params.MaxPacketSize, params.AttemptReconnectInterval, net.ResolveUDPAddr, net.DialUDP, params.Logger)
		if err != nil {
			return nil, err
		}
	} else {
		destAddr, err := net.ResolveUDPAddr("udp", params.HostPort)
		if err != nil {
			return nil, err
		}

		connUDP, err = net.DialUDP(destAddr.Network(), nil, destAddr)
		if err != nil {
			return nil, err
		}
	}

	if err := connUDP.SetWriteBuffer(params.MaxPacketSize); err != nil {
		return nil, err
	}

	return &agentClientUDP{
		connUDP:       connUDP,
		client:        client,
		maxPacketSize: params.MaxPacketSize,
		thriftBuffer:  thriftBuffer,
	}, nil
}

// EmitBatch implements EmitBatch() of Agent interface
func (a *agentClientUDP) EmitBatch(batch *gen.Batch) error {
	a.thriftBuffer.Reset()
	a.client.SeqId = 0 // we have no need for distinct SeqIds for our one-way UDP messages
	if err := a.client.EmitBatch(batch); err != nil {
		return err
	}
	if a.thriftBuffer.Len() > a.maxPacketSize {
		return fmt.Errorf("data does not fit within one UDP packet; size %d, max %d, spans %d",
			a.thriftBuffer.Len(), a.maxPacketSize, len(batch.Spans))
	}
	_, err := a.connUDP.Write(a.thriftBuffer.Bytes())
	return err
}

// Close implements Close() of io.Closer and closes the underlying UDP connection.
func (a *agentClientUDP) Close() error {
	return a.connUDP.Close()
}
