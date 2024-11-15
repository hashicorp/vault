// Copyright (c) 2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"crypto/rand"
	"fmt"
	"strconv"
)

const rfc4122 = 0x40

// UUID is a RFC4122 compliant uuid type
type UUID [16]byte

var nilUUID UUID

// NewUUID creates a new snowflake UUID
func NewUUID() UUID {
	var u UUID
	rand.Read(u[:])
	u[8] = (u[8] | rfc4122) & 0x7F

	var version byte = 4
	u[6] = (u[6] & 0xF) | (version << 4)
	return u
}

func getChar(str string) byte {
	i, _ := strconv.ParseUint(str, 16, 8)
	return byte(i)
}

// ParseUUID parses a string of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx into its UUID form
func ParseUUID(str string) UUID {
	return UUID{
		getChar(str[0:2]), getChar(str[2:4]), getChar(str[4:6]), getChar(str[6:8]),
		getChar(str[9:11]), getChar(str[11:13]),
		getChar(str[14:16]), getChar(str[16:18]),
		getChar(str[19:21]), getChar(str[21:23]),
		getChar(str[24:26]), getChar(str[26:28]), getChar(str[28:30]), getChar(str[30:32]), getChar(str[32:34]), getChar(str[34:36]),
	}
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
