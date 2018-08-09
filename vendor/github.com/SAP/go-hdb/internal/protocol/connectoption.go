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

//go:generate stringer -type=connectOption

type connectOption int8

const (
	coConnectionID                connectOption = 1
	coCompleteArrayExecution      connectOption = 2
	coClientLocale                connectOption = 3
	coSupportsLargeBulkOperations connectOption = 4
	// duplicate in docu: coDataFormatVersion2          connectOption = 5
	// 6-9 reserved: do not use
	coLargeNumberOfParameterSupport connectOption = 10
	coSystemID                      connectOption = 11
	// 12 reserved: do not use
	coAbapVarcharMode             connectOption = 13
	coSelectForUpdateSupported    connectOption = 14
	coClientDistributionMode      connectOption = 15
	coEngineDataFormatVersion     connectOption = 16
	coDistributionProtocolVersion connectOption = 17
	coSplitBatchCommands          connectOption = 18
	coUseTransactionFlagsOnly     connectOption = 19
	//coRowAndColumnOptimizedFormat  connectOption = 20 reserved: do not use
	coIgnoreUnknownParts           connectOption = 21
	coTableOutputParameter         connectOption = 22
	coDataFormatVersion2           connectOption = 23
	coItabParameter                connectOption = 24
	coDescribeTableOutputParameter connectOption = 25
	coColumnarResultset            connectOption = 26
	coScrollablResultSet           connectOption = 27
	coClientInfoNullValueSupported connectOption = 28
	coAssociatedConnectionID       connectOption = 29
	coNoTransactionalPrepare       connectOption = 30
	coFDAEnabled                   connectOption = 31
	coOSUser                       connectOption = 32
	coRowslotImageResult           connectOption = 33
	coEndianess                    connectOption = 34
	// 35, 36 reserved: do not use
	coImplicitLobStreaming connectOption = 37
)
