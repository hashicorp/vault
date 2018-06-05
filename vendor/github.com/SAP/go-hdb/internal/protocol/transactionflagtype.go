/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
