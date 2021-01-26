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
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	. "github.com/aerospike/aerospike-client-go/logger"
	. "github.com/aerospike/aerospike-client-go/types"
)

const (
	_DEFAULT_TIMEOUT = 2 * time.Second
	_NO_TIMEOUT      = 365 * 24 * time.Hour
)

// Access server's info monitoring protocol.
type info struct {
	msg *Message
}

// Send multiple commands to server and store results.
// Timeout should already be set on the connection.
func newInfo(conn *Connection, commands ...string) (*info, error) {
	commandStr := strings.Trim(strings.Join(commands, "\n"), " ")
	if strings.Trim(commandStr, " ") != "" {
		commandStr += "\n"
	}
	newInfo := &info{
		msg: NewMessage(MSG_INFO, []byte(commandStr)),
	}

	if err := newInfo.sendCommand(conn); err != nil {
		return nil, err
	}
	return newInfo, nil
}

// RequestInfo gets info values by name from the specified connection.
// Timeout should already be set on the connection.
func RequestInfo(conn *Connection, names ...string) (map[string]string, error) {
	info, err := newInfo(conn, names...)
	if err != nil {
		return nil, err
	}
	return info.parseMultiResponse()
}

// Issue request and set results buffer. This method is used internally.
// The static request methods should be used instead.
func (nfo *info) sendCommand(conn *Connection) error {
	// Write.
	if _, err := conn.Write(nfo.msg.Serialize()); err != nil {
		Logger.Debug("Failed to send command.")
		return err
	}

	// Read - reuse input buffer.
	header := bytes.NewBuffer(make([]byte, MSG_HEADER_SIZE))
	if _, err := conn.Read(header.Bytes(), MSG_HEADER_SIZE); err != nil {
		return err
	}
	if err := binary.Read(header, binary.BigEndian, &nfo.msg.MessageHeader); err != nil {
		Logger.Debug("Failed to read command response.")
		return err
	}

	// Logger.Debug("Header Response: %v %v %v %v", t.Type, t.Version, t.Length(), t.DataLen)
	if err := nfo.msg.Resize(nfo.msg.Length()); err != nil {
		return err
	}
	_, err := conn.Read(nfo.msg.Data, len(nfo.msg.Data))
	return err
}

func (nfo *info) parseMultiResponse() (map[string]string, error) {
	responses := make(map[string]string)
	data := strings.Trim(string(nfo.msg.Data), "\n")

	keyValuesArr := strings.Split(data, "\n")
	for _, keyValueStr := range keyValuesArr {
		KeyValArr := strings.Split(keyValueStr, "\t")

		switch len(KeyValArr) {
		case 1:
			responses[KeyValArr[0]] = ""
		case 2:
			responses[KeyValArr[0]] = KeyValArr[1]
		default:
			Logger.Error("Requested info buffer does not adhere to the protocol: %s", data)
		}
	}

	return responses, nil
}
