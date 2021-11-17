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
	"encoding/binary"
	"fmt"
	"time"

	// . "github.com/aerospike/aerospike-client-go/logger"
	"github.com/aerospike/aerospike-client-go/pkg/bcrypt"
	. "github.com/aerospike/aerospike-client-go/types"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
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

	// Misc
	_MSG_VERSION int64 = 2
	_MSG_TYPE    int64 = 2

	_HEADER_SIZE      int = 24
	_HEADER_REMAINING int = 16
	_RESULT_CODE      int = 9
	_QUERY_END        int = 50
)

type adminCommand struct {
	dataBuffer []byte
	dataOffset int
}

func newAdminCommand(buf []byte) *adminCommand {
	if buf == nil {
		buf = make([]byte, 10*1024)
	}
	return &adminCommand{
		dataBuffer: buf,
		dataOffset: 8,
	}
}

func (acmd *adminCommand) createUser(cluster *Cluster, policy *AdminPolicy, user string, password []byte, roles []string) error {
	acmd.writeHeader(_CREATE_USER, 3)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_PASSWORD, password)
	acmd.writeRoles(roles)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) dropUser(cluster *Cluster, policy *AdminPolicy, user string) error {
	acmd.writeHeader(_DROP_USER, 1)
	acmd.writeFieldStr(_USER, user)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) setPassword(cluster *Cluster, policy *AdminPolicy, user string, password []byte) error {
	acmd.writeHeader(_SET_PASSWORD, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_PASSWORD, password)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) changePassword(cluster *Cluster, policy *AdminPolicy, user string, password []byte) error {
	acmd.writeHeader(_CHANGE_PASSWORD, 3)
	acmd.writeFieldStr(_USER, user)
	acmd.writeFieldBytes(_OLD_PASSWORD, cluster.Password())
	acmd.writeFieldBytes(_PASSWORD, password)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) grantRoles(cluster *Cluster, policy *AdminPolicy, user string, roles []string) error {
	acmd.writeHeader(_GRANT_ROLES, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeRoles(roles)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) revokeRoles(cluster *Cluster, policy *AdminPolicy, user string, roles []string) error {
	acmd.writeHeader(_REVOKE_ROLES, 2)
	acmd.writeFieldStr(_USER, user)
	acmd.writeRoles(roles)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) createRole(cluster *Cluster, policy *AdminPolicy, roleName string, privileges []Privilege, whitelist []string) error {
	fieldcount := 1
	if len(privileges) > 1 {
		fieldcount++
	}
	if len(whitelist) > 1 {
		fieldcount++
	}
	acmd.writeHeader(_CREATE_ROLE, fieldcount)
	acmd.writeFieldStr(_ROLE, roleName)

	if len(privileges) > 0 {
		if err := acmd.writePrivileges(privileges); err != nil {
			return err
		}
	}

	if len(whitelist) > 0 {
		acmd.writeWhitelist(whitelist)
	}

	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) dropRole(cluster *Cluster, policy *AdminPolicy, roleName string) error {
	acmd.writeHeader(_DROP_ROLE, 1)
	acmd.writeFieldStr(_ROLE, roleName)
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) grantPrivileges(cluster *Cluster, policy *AdminPolicy, roleName string, privileges []Privilege) error {
	acmd.writeHeader(_GRANT_PRIVILEGES, 2)
	acmd.writeFieldStr(_ROLE, roleName)
	if err := acmd.writePrivileges(privileges); err != nil {
		return err
	}
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) revokePrivileges(cluster *Cluster, policy *AdminPolicy, roleName string, privileges []Privilege) error {
	acmd.writeHeader(_REVOKE_PRIVILEGES, 2)
	acmd.writeFieldStr(_ROLE, roleName)
	if err := acmd.writePrivileges(privileges); err != nil {
		return err
	}
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) setWhitelist(cluster *Cluster, policy *AdminPolicy, roleName string, whitelist []string) error {
	fieldCount := 1
	if len(whitelist) > 0 {
		fieldCount++
	}
	acmd.writeHeader(_SET_WHITELIST, fieldCount)
	acmd.writeFieldStr(_ROLE, roleName)
	if len(whitelist) > 0 {
		acmd.writeWhitelist(whitelist)
	}
	return acmd.executeCommand(cluster, policy)
}

func (acmd *adminCommand) queryUser(cluster *Cluster, policy *AdminPolicy, user string) (*UserRoles, error) {
	acmd.writeHeader(_QUERY_USERS, 1)
	acmd.writeFieldStr(_USER, user)
	list, err := acmd.readUsers(cluster, policy)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (acmd *adminCommand) queryUsers(cluster *Cluster, policy *AdminPolicy) ([]*UserRoles, error) {
	acmd.writeHeader(_QUERY_USERS, 0)
	list, err := acmd.readUsers(cluster, policy)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (acmd *adminCommand) queryRole(cluster *Cluster, policy *AdminPolicy, roleName string) (*Role, error) {
	acmd.writeHeader(_QUERY_ROLES, 1)
	acmd.writeFieldStr(_ROLE, roleName)
	list, err := acmd.readRoles(cluster, policy)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (acmd *adminCommand) queryRoles(cluster *Cluster, policy *AdminPolicy) ([]*Role, error) {
	acmd.writeHeader(_QUERY_ROLES, 0)
	list, err := acmd.readRoles(cluster, policy)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (acmd *adminCommand) writeRoles(roles []string) {
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

func (acmd *adminCommand) writePrivileges(privileges []Privilege) error {
	offset := acmd.dataOffset + int(_FIELD_HEADER_SIZE)
	acmd.dataBuffer[offset] = byte(len(privileges))
	offset++

	for _, privilege := range privileges {
		code := privilege.code()

		acmd.dataBuffer[offset] = byte(code)
		offset++

		if privilege.canScope() {

			if len(privilege.SetName) > 0 && len(privilege.Namespace) == 0 {
				return NewAerospikeError(INVALID_PRIVILEGE, fmt.Sprintf("Admin privilege '%v' has a set scope with an empty namespace.", privilege))
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
				return NewAerospikeError(INVALID_PRIVILEGE, fmt.Sprintf("Admin global rivilege '%v' can't have a namespace or set.", privilege))
			}
		}
	}

	size := offset - acmd.dataOffset - int(_FIELD_HEADER_SIZE)
	acmd.writeFieldHeader(_PRIVILEGES, size)
	acmd.dataOffset = offset

	return nil
}

func (acmd *adminCommand) writeWhitelist(whitelist []string) {
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

func (acmd *adminCommand) writeSize() {
	// Write total size of message which is the current offset.
	var size = int64(acmd.dataOffset-8) | (_MSG_VERSION << 56) | (_MSG_TYPE << 48)
	binary.BigEndian.PutUint64(acmd.dataBuffer[0:], uint64(size))
}

func (acmd *adminCommand) writeHeader(command byte, fieldCount int) {
	// Authenticate header is almost all zeros
	for i := acmd.dataOffset; i < acmd.dataOffset+16; i++ {
		acmd.dataBuffer[i] = 0
	}
	acmd.dataBuffer[acmd.dataOffset+2] = command
	acmd.dataBuffer[acmd.dataOffset+3] = byte(fieldCount)
	acmd.dataOffset += 16
}

func (acmd *adminCommand) writeFieldStr(id byte, str string) {
	len := copy(acmd.dataBuffer[acmd.dataOffset+int(_FIELD_HEADER_SIZE):], str)
	acmd.writeFieldHeader(id, len)
	acmd.dataOffset += len
}

func (acmd *adminCommand) writeFieldBytes(id byte, bytes []byte) {
	copy(acmd.dataBuffer[acmd.dataOffset+int(_FIELD_HEADER_SIZE):], bytes)
	acmd.writeFieldHeader(id, len(bytes))
	acmd.dataOffset += len(bytes)
}

func (acmd *adminCommand) writeFieldHeader(id byte, size int) {
	// Buffer.Int32ToBytes(int32(size+1), acmd.dataBuffer, acmd.dataOffset)
	binary.BigEndian.PutUint32(acmd.dataBuffer[acmd.dataOffset:], uint32(size+1))

	acmd.dataOffset += 4
	acmd.dataBuffer[acmd.dataOffset] = id
	acmd.dataOffset++
}

func (acmd *adminCommand) executeCommand(cluster *Cluster, policy *AdminPolicy) error {
	acmd.writeSize()
	node, err := cluster.GetRandomNode()
	if err != nil {
		return err
	}
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	node.tendConnLock.Lock()
	defer node.tendConnLock.Unlock()

	if err := node.initTendConn(timeout); err != nil {
		return err
	}

	conn := node.tendConn
	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if _, err := conn.Read(acmd.dataBuffer, _HEADER_SIZE); err != nil {
		return err
	}

	result := acmd.dataBuffer[_RESULT_CODE]
	if result != 0 {
		return NewAerospikeError(ResultCode(result))
	}

	return nil
}

func (acmd *adminCommand) readUsers(cluster *Cluster, policy *AdminPolicy) ([]*UserRoles, error) {
	acmd.writeSize()
	node, err := cluster.GetRandomNode()
	if err != nil {
		return nil, err
	}
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	node.tendConnLock.Lock()
	defer node.tendConnLock.Unlock()

	if err := node.initTendConn(timeout); err != nil {
		return nil, err
	}

	conn := node.tendConn
	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return nil, err
	}

	status, list, err := acmd.readUserBlocks(conn)
	if err != nil {
		return nil, err
	}

	if status > 0 {
		return nil, NewAerospikeError(ResultCode(status))
	}
	return list, nil
}

func (acmd *adminCommand) readUserBlocks(conn *Connection) (status int, rlist []*UserRoles, err error) {

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

func (acmd *adminCommand) parseUsers(receiveSize int) (int, []*UserRoles, error) {
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
			len := int(Buffer.BytesToInt32(acmd.dataBuffer, acmd.dataOffset))
			acmd.dataOffset += 4
			id := acmd.dataBuffer[acmd.dataOffset]
			acmd.dataOffset++
			len--

			if id == _USER {
				userRoles.User = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
				acmd.dataOffset += len
			} else if id == _ROLES {
				acmd.parseRoles(userRoles)
			} else {
				acmd.dataOffset += len
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

func (acmd *adminCommand) parseRoles(userRoles *UserRoles) {
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

func hashPassword(password string) ([]byte, error) {
	// Hashing the password with the cost of 10, with a static salt
	const salt = "$2a$10$7EqJtq98hPqEX7fNZaFWoO"
	hashedPassword, err := bcrypt.Hash(password, salt)
	if err != nil {
		return nil, err
	}
	return []byte(hashedPassword), nil
}

func (acmd *adminCommand) readRoles(cluster *Cluster, policy *AdminPolicy) ([]*Role, error) {
	acmd.writeSize()
	node, err := cluster.GetRandomNode()
	if err != nil {
		return nil, err
	}
	timeout := 1 * time.Second
	if policy != nil && policy.Timeout > 0 {
		timeout = policy.Timeout
	}

	node.tendConnLock.Lock()
	defer node.tendConnLock.Unlock()

	if err := node.initTendConn(timeout); err != nil {
		return nil, err
	}

	conn := node.tendConn
	if _, err := conn.Write(acmd.dataBuffer[:acmd.dataOffset]); err != nil {
		return nil, err
	}

	status, list, err := acmd.readRoleBlocks(conn)
	if err != nil {
		return nil, err
	}

	if status > 0 {
		return nil, NewAerospikeError(ResultCode(status))
	}
	return list, nil
}

func (acmd *adminCommand) readRoleBlocks(conn *Connection) (status int, rlist []*Role, err error) {

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

func (acmd *adminCommand) parseRolesFull(receiveSize int) (int, []*Role, error) {
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

			if id == _ROLE {
				role.Name = string(acmd.dataBuffer[acmd.dataOffset : acmd.dataOffset+len])
				acmd.dataOffset += len
			} else if id == _PRIVILEGES {
				acmd.parsePrivileges(role)
			} else if id == _WHITELIST {
				role.Whitelist = acmd.parseWhitelist(len)
			} else {
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

func (acmd *adminCommand) parsePrivileges(role *Role) {
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

func (acmd *adminCommand) parseWhitelist(length int) []string {
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
