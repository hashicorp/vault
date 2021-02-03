// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=endianess

type endianess int8

const (
	bigEndian    endianess = 0
	littleEndian endianess = 1
)
