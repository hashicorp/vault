// +build !go1.14

// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

// Delete after go1.13 is out of maintenance.

type partReadWriter interface {
	part
	numArg() int
	partDecoder
	partEncoder
}
