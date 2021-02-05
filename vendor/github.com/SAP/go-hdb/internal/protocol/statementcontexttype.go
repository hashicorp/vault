// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=statementContextType

type statementContextType int8

const (
	scStatementSequenceInfo statementContextType = 1
	scServerExecutionTime   statementContextType = 2
)
