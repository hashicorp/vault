// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=transactionFlagType

//transaction flags
type transactionFlagType int8

const (
	tfRolledback                     transactionFlagType = 0
	tfCommited                       transactionFlagType = 1
	tfNewIsolationLevel              transactionFlagType = 2
	tfDDLCommitmodeChanged           transactionFlagType = 3
	tfWriteTransactionStarted        transactionFlagType = 4
	tfNowriteTransactionStarted      transactionFlagType = 5
	tfSessionClosingTransactionError transactionFlagType = 6
)
