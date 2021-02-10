// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=clientContextOption

type clientContextOption int8

const (
	ccoClientVersion            clientContextOption = 1
	ccoClientType               clientContextOption = 2
	ccoClientApplicationProgram clientContextOption = 3
)
