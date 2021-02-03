// +build go1.14

// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Delete and re-itegrate into parts.go after go1.13 is out of maintenance.

package protocol

type partReadWriter interface {
	partReader
	partWriter
}
