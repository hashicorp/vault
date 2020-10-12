// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use acmd file except in compliance with the License.
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
	"time"

	. "github.com/aerospike/aerospike-client-go/logger"
	// "github.com/aerospike/aerospike-client-go/pkg/bcrypt"
	. "github.com/aerospike/aerospike-client-go/types"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
)

// Login command authenticates to the server.
// If the authentication is external, Session Information will be returned.
type loginCommand struct {
	adminCommand

	// SessionToken for the current session on the external authentication server.
	SessionToken []byte

	// SessionExpiration for the current session on the external authentication server.
	SessionExpiration time.Time
}

func newLoginCommand(buf []byte) *loginCommand {
	return &loginCommand{
		adminCommand: *newAdminCommand(buf),
	}
}

// Login tries to authenticate to the aerospike server. Depending on the server configuration and ClientPolicy,
// the session information will be returned.
func (lcmd *loginCommand) Login(policy *ClientPolicy, conn *Connection) error {
	hashedPass, err := hashPassword(policy.Password)
	if err != nil {
		return err
	}

	return lcmd.login(policy, conn, hashedPass)
}

// Login tries to authenticate to the aerospike server. Depending on the server configuration and ClientPolicy,
// the session information will be returned.
func (lcmd *loginCommand) login(policy *ClientPolicy, conn *Connection, hashedPass []byte) error {
	switch policy.AuthMode {
	case AuthModeExternal:
		lcmd.writeHeader(_LOGIN, 3)
		lcmd.writeFieldStr(_USER, policy.User)
		lcmd.writeFieldBytes(_CREDENTIAL, hashedPass)
		lcmd.writeFieldStr(_CLEAR_PASSWORD, policy.Password)
	case AuthModeInternal:
		lcmd.writeHeader(_LOGIN, 2)
		lcmd.writeFieldStr(_USER, policy.User)
		lcmd.writeFieldBytes(_CREDENTIAL, hashedPass)
	default:
		return NewAerospikeError(ResultCode(INVALID_COMMAND), "Invalid ClientPolicy.AuthMode.")
	}

	lcmd.writeSize()

	var deadline time.Time
	if policy.LoginTimeout > 0 {
		deadline = time.Now().Add(policy.Timeout)
	}
	conn.SetTimeout(deadline, policy.LoginTimeout)

	if _, err := conn.Write(lcmd.dataBuffer[:lcmd.dataOffset]); err != nil {
		return err
	}

	if _, err := conn.Read(lcmd.dataBuffer, _HEADER_SIZE); err != nil {
		return err
	}

	result := lcmd.dataBuffer[_RESULT_CODE] & 0xFF
	if result != 0 {
		if int(result) == int(INVALID_COMMAND) {
			// New login not supported.  Try old authentication.
			return lcmd.authenticateInternal(conn, policy.User, hashedPass)
		}

		if int(result) == int(SECURITY_NOT_ENABLED) {
			// Server does not require login.
			return nil
		}

		return NewAerospikeError(ResultCode(result))
	}

	// Read session token.
	sz := Buffer.BytesToInt64(lcmd.dataBuffer, 0)
	receiveSize := int((sz & 0xFFFFFFFFFFFF) - int64(_HEADER_REMAINING))
	fieldCount := int(lcmd.dataBuffer[11] & 0xFF)

	if receiveSize <= 0 || receiveSize > len(lcmd.dataBuffer) || fieldCount <= 0 {
		return NewAerospikeError(ResultCode(result), "Node failed to retrieve session token")
	}

	if len(lcmd.dataBuffer) < receiveSize {
		lcmd.dataBuffer = make([]byte, receiveSize)
	}

	_, err := conn.Read(lcmd.dataBuffer, receiveSize)
	if err != nil {
		Logger.Debug("Error reading data from connection for login command: %s", err.Error())
		return err
	}

	lcmd.dataOffset = 0
	for i := 0; i < fieldCount; i++ {
		mlen := int(Buffer.BytesToUint32(lcmd.dataBuffer, lcmd.dataOffset))
		lcmd.dataOffset += 4
		id := lcmd.dataBuffer[lcmd.dataOffset]
		lcmd.dataOffset++
		mlen--

		switch id {
		case _SESSION_TOKEN:
			lcmd.SessionToken = lcmd.dataBuffer[lcmd.dataOffset : lcmd.dataOffset+mlen]
		case _SESSION_TTL:
			// Subtract 60 seconds from TTL so client session expires before server session.
			seconds := int(Buffer.BytesToUint32(lcmd.dataBuffer, lcmd.dataOffset) - 60)

			if seconds > 0 {
				lcmd.SessionExpiration = time.Now().Add(time.Duration(seconds) * time.Second)
			} else {
				Logger.Warn("Invalid session TTL: %d", seconds)
			}
		}

		lcmd.dataOffset += mlen
	}

	if lcmd.SessionToken == nil {
		return NewAerospikeError(ResultCode(result), "Node failed to retrieve session token")
	}
	return nil
}

func (lcmd *loginCommand) authenticateInternal(conn *Connection, user string, passwordHash []byte) error {
	lcmd.dataOffset = 8
	lcmd.writeHeader(_AUTHENTICATE, 2)
	lcmd.writeFieldStr(_USER, user)
	lcmd.writeFieldBytes(_CREDENTIAL, passwordHash)
	lcmd.writeSize()

	if _, err := conn.Write(lcmd.dataBuffer[:lcmd.dataOffset]); err != nil {
		return err
	}

	if _, err := conn.Read(lcmd.dataBuffer, _HEADER_SIZE); err != nil {
		return err
	}

	result := lcmd.dataBuffer[_RESULT_CODE] & 0xFF
	if result != 0 && int(result) != int(SECURITY_NOT_ENABLED) {
		return NewAerospikeError(ResultCode(result), "Authentication failed")
	}

	return nil
}

func (lcmd *loginCommand) authenticateViaToken(policy *ClientPolicy, conn *Connection, sessionToken []byte) error {
	lcmd.setAuthenticate(policy, sessionToken)

	if _, err := conn.Write(lcmd.dataBuffer[:lcmd.dataOffset]); err != nil {
		return err
	}

	if _, err := conn.Read(lcmd.dataBuffer, _HEADER_SIZE); err != nil {
		return err
	}

	result := lcmd.dataBuffer[_RESULT_CODE] & 0xFF
	if result != 0 && int(result) != int(SECURITY_NOT_ENABLED) {
		return NewAerospikeError(ResultCode(result), "Authentication failed")
	}

	return nil
}

func (lcmd *loginCommand) setAuthenticate(policy *ClientPolicy, sessionToken []byte) error {
	lcmd.writeHeader(_AUTHENTICATE, 2)
	lcmd.writeFieldStr(_USER, policy.User)

	if sessionToken != nil {
		// New authentication.
		lcmd.writeFieldBytes(_SESSION_TOKEN, sessionToken)
	} else {
		hashedPass, err := hashPassword(policy.Password)
		if err != nil {
			return err
		}

		// Old authentication.
		lcmd.writeFieldBytes(_CREDENTIAL, hashedPass)
	}

	lcmd.writeSize()

	return nil
}
