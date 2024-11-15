// Copyright 2014-2021 Aerospike, Inc.
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
	"encoding/binary"
	"fmt"
	"time"

	"github.com/aerospike/aerospike-client-go/v5/pkg/bcrypt"
	"github.com/aerospike/aerospike-client-go/v5/types"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

const (
	// Commands
	_AUTHENTICATE      byte = 0
	_CREATE_USER       byte = 1
	_DROP_USER         byte = 2
	_SET_PASSWORD      byte = 3
	_CHANGE_PASSWORD   byte = 4
	_GRANT_ROLES       byte = 5
	_REVOKE_ROLES      byte = 6
	_QUERY_USERS       byte = 9
	_CREATE_ROLE       byte = 10
	_DROP_ROLE         byte = 11
	_GRANT_PRIVILEGES  byte = 12
	_REVOKE_PRIVILEGES byte = 13
	_SET_WHITELIST     byte = 14
	_SET_QUOTAS        byte = 15
	_QUERY_ROLES       byte = 16
	_LOGIN             byte = 20

	// Field IDs
	_USER           byte = 0
	_PASSWORD       byte = 1
	_OLD_PASSWORD   byte = 2
	_CREDENTIAL     byte = 3
	_CLEAR_PASSWORD byte = 4
	_SESSION_TOKEN  byte = 5
	_SESSION_TTL    byte = 6
	_ROLES          byte = 10
	_ROLE           byte = 11
	_PRIVILEGES     byte = 12
	_WHITELIST      byte = 13
	_READ_QUOTA     byte = 14
	_WRITE_QUOTA    byte = 15
	_READ_INFO      byte = 16
	_WRITE_INFO     byte = 17
	_CONNECTIONS    byte = 18

	// Misc
	_MSG_VERSION int64 = 2
	_MSG_TYPE    int64 = 2

	_HEADER_SIZE      int = 24
	_HEADER_REMAINING int = 16
	_RESULT_CODE      int = 9
	_QUERY_END        int = 50
)

// AdminCommand allows managing user access to the server
type AdminCommand struct {
	dataBuffer []byte
	dataOffset int
}

// NewAdminCommand returns an AdminCommand object.
func NewAdminCommand(buf []byte) *AdminCommand {
	if buf == nil {
		buf = make([]byte, 10*1024)
	}
	return &AdminCommand{
		dataBuffer: buf,
		dataOffset: 8,
	}
}

func (acmd *AdminCommand) createUser(conn *Connection, policy *AdminPolicy, user string, password []byte, roles []string) Error {
	acmd.writeHeader(_CREATE_USER, 3)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_PASSWORD, password)
	acmd.writeRoles(roles)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) dropUser(conn *Connection, policy *AdminPolicy, user string) Error {
	acmd.writeHeader(_DROP_USER, 1)
	acmd.writeFieldStr(_USER, user)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) setPassword(conn *Connection, policy *AdminPolicy, user string, password []byte) Error {
	acmd.writeHeader(_SET_PASSWORD, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_PASSWORD, password)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) changePassword(conn *Connection, policy *AdminPolicy, user string, oldPass, newPass []byte) Error {
	acmd.writeHeader(_CHANGE_PASSWORD, 3)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_OLD_PASSWORD, oldPass)
	acmd.writeFieldBytes(_PASSWORD, newPass)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) grantRoles(conn *Connection, policy *AdminPolicy, user string, roles []string) Error {
	acmd.writeHeader(_GRANT_ROLES, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeRoles(roles)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) revokeRoles(conn *Connection, policy *AdminPolicy, user string, roles []string) Error {
	acmd.writeHeader(_REVOKE_ROLES, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeRoles(roles)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) createRole(conn *Connection, policy *AdminPolicy, roleName string, privileges []Privilege, whitelist []string, readQuota, writeQuota uint32) Error {
	fieldCount := 1
	if len(privileges) > 1 {
		fieldCount++
	}

	if len(whitelist) > 1 {
		fieldCount++
	}

	if readQuota > 0 {
		fieldCount++
	}

	if writeQuota > 0 {
		fieldCount++
	}

	acmd.writeHeader(_CREATE_ROLE, fieldCount)
	acmd.writeFieldStr(_ROLE, roleName)

	if len(privileges) > 0 {
		if err := acmd.writePrivileges(privileges); err != nil {
			return err
		}
	}

	if len(whitelist) > 0 {
		acmd.writeWhitelist(whitelist)
	}

	if readQuota > 0 {
		acmd.writeFieldUint32(_READ_QUOTA, readQuota)
	}

	if writeQuota > 0 {
		acmd.writeFieldUint32(_WRITE_QUOTA, writeQuota)
	}

	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) dropRole(conn *Connection, policy *AdminPolicy, roleName string) Error {
	acmd.writeHeader(_DROP_ROLE, 1)
	acmd.writeFieldStr(_ROLE, roleName)
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) grantPrivileges(conn *Connection, policy *AdminPolicy, roleName string, privileges []Privilege) Error {
	acmd.writeHeader(_GRANT_PRIVILEGES, 2)
	acmd.writeFieldStr(_ROLE, roleName)
	if err := acmd.writePrivileges(privileges); err != nil {
		return err
	}
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) revokePrivileges(conn *Connection, policy *AdminPolicy, roleName string, privileges []Privilege) Error {
	acmd.writeHeader(_REVOKE_PRIVILEGES, 2)
	acmd.writeFieldStr(_ROLE, roleName)
	if err := acmd.writePrivileges(privileges); err != nil {
		return err
	}
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) setWhitelist(conn *Connection, policy *AdminPolicy, roleName string, whitelist []string) Error {
	fieldCount := 1
	if len(whitelist) > 0 {
		fieldCount++
	}
	acmd.writeHeader(_SET_WHITELIST, fieldCount)
	acmd.writeFieldStr(_ROLE, roleName)
	if len(whitelist) > 0 {
		acmd.writeWhitelist(whitelist)
	}
	return acmd.executeCommand(conn, policy)
}

func (acmd *AdminCommand) setQuotas(conn *Connection, policy *AdminPolicy, roleName string, readQuota, writeQuota uint32) Error {
	acmd.writeHeader(_SET_QUOTAS, 3)
	acmd.writeFieldStr(_ROLE, roleName)
	acmd.writeFieldUint32(_READ_QUOTA, readQuota)
	acmd.writeFieldUint32(_WRITE_QUOTA, writeQuota)
	return acmd.executeCommand(conn, policy)
}

// QueryUser returns user information.
func (acmd *AdminCommand) QueryUser(conn *Connection, policy *AdminPolicy, user string) (*UserRoles, Error) {
	acmd.writeHeader(_QUERY_USERS, 1)
	acmd.writeFieldStr(_USER, user)
	list, err := acmd.readUsers(conn, policy)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// QueryUsers returns user information for all users.
func (acmd *AdminCommand) QueryUsers(conn *Connection, policy *AdminPolicy) ([]*UserRoles, Error) {
	acmd.writeHeader(_QUERY_USERS, 0)
	list, err := acmd.readUsers(conn, policy)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// QueryRole returns role information.
func (acmd *AdminCommand) QueryRole(conn *Connection, policy *AdminPolicy, roleName string) (*Role, Error) {
	acmd.writeHeader(_QUERY_ROLES, 1)
	acmd.writeFieldStr(_ROLE, roleName)
	list, err := acmd.readRoles(conn, policy)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// QueryRoles returns role information for all roles.
func (acmd *AdminCommand) QueryRoles(conn *Connection, policy *AdminPolicy) ([]*Role, Error) {
	acmd.writeHeader(_QUERY_ROLES, 0)
	list, err := acmd.readRoles(conn, policy)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (acmd *AdminCommand) writeRoles(roles []string) {
	offset := acmd.dataOffset + int(_FIELD_HEADER_SIZE)
	acmd.dataBuffer[offset] = byte(len(roles))
	offset++

	for _, role := range roles {
		len := copy(acmd.dataBuffer[offset+1:], role)
		acmd.dataBuffer[offset] = byte(len)
		offset += len + 1
	}

	size := offset - acmd.dataOffset - int(_FIELD_HEADER_SIZE)
	acmd.writeFieldHeader(_ROLES, size)
	acmd.dataOffset = offset
}

func (acmd *AdminCommand) writePrivileges(privileges []Privilege) Error {
	offset := acmd.dataOffset + int(_FIELD_HEADER_SIZE)
	acmd.dataBuffer[offset] = byte(len(privileges))
	offset++

	for _, privilege := range privileges {
		code := privilege.code()

		acmd.dataBuffer[offset] = byte(code)
		offset++

		if privilege.canScope() {

			if len(privilege.SetName) > 0 && len(privilege.Namespace) == 0 {
				return newError(types.INVALID_PRIVILEGE, fmt.Sprintf("Admin privilege '%v' has a set scope with an empty namespace.", privilege))
			}

			acmd.dataBuffer[offset] = byte(len(privilege.Namespace))
			offset++
			copy(acmd.dataBuffer[offset:], privilege.Namespace)
			offset += len(privilege.Namespace)

			acmd.dataBuffer[offset] = byte(len(privilege.SetName))
			offset++
			copy(acmd.dataBuffer[offset:], privilege.SetName)
			offset += len(privilege.SetName)
		} else {
			if len(privilege.Namespace) > 0 || len(privilege.SetName) > 0 {
				return newError(types.INVALID_PRIVILEGE, fmt.Sprintf("Admin global rivilege '%v' can't have a namespace or set.", privilege))
			}
		}
	}

	size := offset - acmd.dataOffset - int(_FIELD_HEADER_SIZE)
	acmd.writeFieldHeader(_PRIVILEGES, size)
	acmd.dataOffset = offset

	return nil
}

func (acmd *AdminCommand) writeWhitelist(whitelist []string) {
	offset := acmd.dataOffset + int(_FIELD_HEADER_SIZE)

	comma := false
	for _, address := range whitelist {
		if comma {
			acmd.dataBuffer[offset] = ','
			offset++
		} else {
			comma = true
		}

		offset += copy(acmd.dataBuffer[offset:], address)
	}

	size := offset - acmd.dataOffset - int(_FIELD_HEADER_SIZE)
	acmd.writeFieldHeader(_WHITELIST, size)
	acmd.dataOffset = offset
}

func (acmd *AdminCommand) writeSize() {
	// Write total size of message which is the current offset.
	var size = int64(acmd.dataOffset-8) | (_MSG_VERSION << 56) | (_MSG_TYPE << 48)
	binary.BigEndian.PutUint64(acmd.dataBuffer[0:], uint64(size))
}

func (acmd *AdminCommand) writeHeader(command byte, fieldCount int) {
	// Authenticate header is almost all zeros
	for i := acmd.dataOffset; i < acmd.dataOffset+16; i++ {
		acmd.dataBuffer[i] = 0
	}
	acmd.dataBuffer[acmd.dataOffset+2] = command
	acmd.dataBuffer[acmd.dataOffset+3] = byte(fieldCount)
	acmd.dataOffset += 16
}

func (acmd *AdminCommand) writeFieldStr(id byte, str string) {
	len := copy(acmd.dataBuffer[acmd.dataOffset+int(_FIELD_HEADER_SIZE):], str)
	acmd.writeFieldHeader(id, len)
	acmd.dataOffset += len
}

func (acmd *AdminCommand) writeFieldBytes(id byte, bytes []byte) {
	copy(acmd.dataBuffer[acmd.dataOffset+int(_FIELD_HEADER_SIZE):], bytes)
	acmd.writeFieldHeader(id, len(bytes))
	acmd.dataOffset += len(bytes)
}

func (acmd *AdminCommand) writeFieldUint32(id byte, size uint32) {
	acmd.writeFieldHeader(id, 4)
	binary.BigEndian.PutUint32(acmd.dataBuffer[acmd.dataOffset:], size)
	acmd.dataOffset += 4
}

func (acmd *AdminCommand) writeFieldHeader(id byte, size int) {
	// Buffer.Int32ToBytes(int32(size+1), acmd.dataBuffer, acmd.dataOffset)
	binary.BigEndian.PutUint32(acmd.dataBuffer[acmd.dataOffset:], uint32(size+1))

	acmd.dataOffset += 4
	acmd.dataBuffer[acmd.dataOffset] = id
	acmd.dataOffset++
}

func (acmd *AdminCommand) executeCommand(conn *Connection, policy *AdminPolicy) Error {
	acmd.writeSize()
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	if err := conn.SetTimeout(time.Now().Add(timeout), timeout); err != nil {
		return err
	}

	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return err
	}

	if _, err := conn.Read(acmd.dataBuffer, _HEADER_SIZE); err != nil {
		return err
	}

	result := acmd.dataBuffer[_RESULT_CODE]
	if result != 0 {
		if conn.node != nil {
			return newCustomNodeError(conn.node, types.ResultCode(result))
		}
		return newError(types.ResultCode(result))
	}

	return nil
}

func (acmd *AdminCommand) readUsers(conn *Connection, policy *AdminPolicy) ([]*UserRoles, Error) {
	acmd.writeSize()
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	if err := conn.SetTimeout(time.Now().Add(timeout), timeout); err != nil {
		return nil, err
	}

	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return nil, err
	}

	status, list, err := acmd.readUserBlocks(conn)
	if err != nil {
		return nil, err
	}

	if status > 0 {
		if conn.node != nil {
			return nil, newCustomNodeError(conn.node, types.ResultCode(status))
		}
		return nil, newError(types.ResultCode(status))
	}
	return list, nil
}

func (acmd *AdminCommand) readUserBlocks(conn *Connection) (status int, rlist []*UserRoles, err Error) {

	var list []*UserRoles

	for status == 0 {
		if _, err = conn.Read(acmd.dataBuffer, 8); err != nil {
			return -1, nil, err
		}

		size := Buffer.BytesToInt64(acmd.dataBuffer, 0)
		receiveSize := (size & 0xFFFFFFFFFFFF)

		if receiveSize > 0 {
			if receiveSize > int64(len(acmd.dataBuffer)) {
				acmd.dataBuffer = make([]byte, receiveSize)
			}
			if _, err = conn.Read(acmd.dataBuffer, int(receiveSize)); err != nil {
				return -1, nil, err
			}
			status, list, err = acmd.parseUsers(int(receiveSize))
			if err != nil {
				return -1, nil, err
			}
			rlist = append(rlist, list...)
		} else {
			break
		}
	}
	return status, rlist, nil
}

func (acmd *AdminCommand) parseUsers(receiveSize int) (int, []*UserRoles, Error) {
	acmd.dataOffset = 0
	list := make([]*UserRoles, 0, 100)

	for acmd.dataOffset < receiveSize {
		resultCode := int(acmd.dataBuffer[acmd.dataOffset+1])

		if resultCode != 0 {
			if resultCode == _QUERY_END {
				return -1, nil, nil
			}
			return resultCode, nil, nil
		}

		userRoles := &UserRoles{}
		fieldCount := int(acmd.dataBuffer[acmd.dataOffset+3])
		acmd.dataOffset += _HEADER_REMAINING

		for i := 0; i < fieldCount; i++ {
			flen := int(Buffer.BytesToInt32(acmd.dataBuffer, acmd.dataOffset))
			acmd.dataOffset += 4
			id := acmd.dataBuffer[acmd.dataOffset]
			acmd.dataOffset++
			flen--

			switch id {
			case _USER:
				userRoles.User = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+flen])
				acmd.dataOffset += flen
			case _ROLES:
				acmd.parseRoles(userRoles)
			case _READ_INFO:
				userRoles.ReadInfo = acmd.parseInfo()
			case _WRITE_INFO:
				userRoles.WriteInfo = acmd.parseInfo()
			case _CONNECTIONS:
				userRoles.ConnsInUse = int(Buffer.BytesToInt32(acmd.dataBuffer, acmd.dataOffset))
				acmd.dataOffset += flen
			default:
				acmd.dataOffset += flen
			}
		}

		if userRoles.User == "" && userRoles.Roles == nil {
			continue
		}

		if userRoles.Roles == nil {
			userRoles.Roles = make([]string, 0)
		}
		list = append(list, userRoles)
	}

	return 0, list, nil
}

func (acmd *AdminCommand) parseRoles(userRoles *UserRoles) {
	size := int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
	acmd.dataOffset++
	userRoles.Roles = make([]string, 0, size)

	for i := 0; i < size; i++ {
		len := int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
		acmd.dataOffset++
		role := string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
		acmd.dataOffset += len
		userRoles.Roles = append(userRoles.Roles, role)
	}
}

func (acmd *AdminCommand) parseInfo() []int {
	size := int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
	acmd.dataOffset++
	list := make([]int, 0, size)

	for i := 0; i < size; i++ {
		val := int(Buffer.BytesToUint32(acmd.dataBuffer, acmd.dataOffset))
		acmd.dataOffset += 4
		list = append(list, val)
	}
	return list
}

func hashPassword(password string) ([]byte, Error) {
	// Hashing the password with the cost of 10, with a static salt
	const salt = "$2a$10$7EqJtq98hPqEX7fNZaFWoO"
	hashedPassword, err := bcrypt.Hash(password, salt)
	if err != nil {
		return nil, newCommonError(err)
	}
	return []byte(hashedPassword), nil
}

func (acmd *AdminCommand) readRoles(conn *Connection, policy *AdminPolicy) ([]*Role, Error) {
	acmd.writeSize()
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	if err := conn.SetTimeout(time.Now().Add(timeout), timeout); err != nil {
		return nil, err
	}

	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return nil, err
	}

	status, list, err := acmd.readRoleBlocks(conn)
	if err != nil {
		return nil, err
	}

	if status > 0 {
		if conn.node != nil {
			return nil, newCustomNodeError(conn.node, types.ResultCode(status))
		}
		return nil, newError(types.ResultCode(status))
	}
	return list, nil
}

func (acmd *AdminCommand) readRoleBlocks(conn *Connection) (status int, rlist []*Role, err Error) {

	var list []*Role

	for status == 0 {
		if _, err = conn.Read(acmd.dataBuffer, 8); err != nil {
			return -1, nil, err
		}

		size := Buffer.BytesToInt64(acmd.dataBuffer, 0)
		receiveSize := int(size & 0xFFFFFFFFFFFF)

		if receiveSize > 0 {
			if receiveSize > len(acmd.dataBuffer) {
				acmd.dataBuffer = make([]byte, receiveSize)
			}
			if _, err = conn.Read(acmd.dataBuffer, receiveSize); err != nil {
				return -1, nil, err
			}
			status, list, err = acmd.parseRolesFull(receiveSize)
			if err != nil {
				return -1, nil, err
			}
			rlist = append(rlist, list...)
		} else {
			break
		}
	}
	return status, rlist, nil
}

func (acmd *AdminCommand) parseRolesFull(receiveSize int) (int, []*Role, Error) {
	acmd.dataOffset = 0

	var list []*Role
	for acmd.dataOffset < receiveSize {
		resultCode := int(acmd.dataBuffer[acmd.dataOffset+1])

		if resultCode != 0 {
			if resultCode == _QUERY_END {
				return -1, nil, nil
			}
			return resultCode, nil, nil
		}

		role := &Role{}
		fieldCount := int(acmd.dataBuffer[acmd.dataOffset+3])
		acmd.dataOffset += _HEADER_REMAINING

		for i := 0; i < fieldCount; i++ {
			len := int(Buffer.BytesToInt32(acmd.dataBuffer, acmd.dataOffset))
			acmd.dataOffset += 4
			id := acmd.dataBuffer[acmd.dataOffset]
			acmd.dataOffset++
			len--

			switch id {
			case _ROLE:
				role.Name = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
				acmd.dataOffset += len
			case _PRIVILEGES:
				acmd.parsePrivileges(role)
			case _WHITELIST:
				role.Whitelist = acmd.parseWhitelist(len)
			case _READ_QUOTA:
				role.ReadQuota = Buffer.BytesToUint32(acmd.dataBuffer, acmd.dataOffset)
				acmd.dataOffset += len
			case _WRITE_QUOTA:
				role.WriteQuota = Buffer.BytesToUint32(acmd.dataBuffer, acmd.dataOffset)
				acmd.dataOffset += len
			default:
				acmd.dataOffset += len
			}
		}

		if len(role.Name) == 0 && len(role.Privileges) == 0 {
			continue
		}

		if role.Privileges == nil {
			role.Privileges = []Privilege{}
		}
		list = append(list, role)
	}
	return 0, list, nil
}

func (acmd *AdminCommand) parsePrivileges(role *Role) {
	size := int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
	acmd.dataOffset++
	role.Privileges = make([]Privilege, 0, size)

	for i := 0; i < size; i++ {
		priv := Privilege{}
		priv.Code = privilegeFrom(acmd.dataBuffer[acmd.dataOffset])
		acmd.dataOffset++

		if priv.canScope() {
			len := int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
			acmd.dataOffset++
			priv.Namespace = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
			acmd.dataOffset += len

			len = int(acmd.dataBuffer[acmd.dataOffset] & 0xFF)
			acmd.dataOffset++
			priv.SetName = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
			acmd.dataOffset += len
		}
		role.Privileges = append(role.Privileges, priv)
	}
}

func (acmd *AdminCommand) parseWhitelist(length int) []string {
	list := []string{}
	begin := acmd.dataOffset
	max := begin + length

	for acmd.dataOffset < max {
		if acmd.dataBuffer[acmd.dataOffset] == ',' {
			l := acmd.dataOffset - begin
			if l > 0 {
				s := string(acmd.dataBuffer[begin : begin+l])
				list = append(list, s)
			}
			acmd.dataOffset++
			begin = acmd.dataOffset
		} else {
			acmd.dataOffset++
		}
	}

	l := acmd.dataOffset - begin
	if l > 0 {
		s := string(acmd.dataBuffer[begin : begin+l])
		list = append(list, s)
	}

	return list
}
