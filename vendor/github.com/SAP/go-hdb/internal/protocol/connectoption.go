// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=connectOption

type connectOption int8

const (
	coConnectionID                        connectOption = 1
	coCompleteArrayExecution              connectOption = 2  //!< @deprecated Array execution semantics, always true.
	coClientLocale                        connectOption = 3  //!< Client locale information.
	coSupportsLargeBulkOperations         connectOption = 4  //!< Bulk operations >32K are supported.
	coDistributionEnabled                 connectOption = 5  //!< @deprecated Distribution (topology & call routing) enabled
	coPrimaryConnectionID                 connectOption = 6  //!< @deprecated Id of primary connection (unused).
	coPrimaryConnectionHost               connectOption = 7  //!< @deprecated Primary connection host name (unused).
	coPrimaryConnectionPort               connectOption = 8  //!< @deprecated Primary connection port (unused).
	coCompleteDatatypeSupport             connectOption = 9  //!< @deprecated All data types supported (always on).
	coLargeNumberOfParametersSupport      connectOption = 10 //!< Number of parameters >32K is supported.
	coSystemID                            connectOption = 11 //!< SID of SAP HANA Database system (output only).
	coDataFormatVersion                   connectOption = 12 //!< Version of data format used in communication (@see DataFormatVersionEnum).
	coAbapVarcharMode                     connectOption = 13 //!< ABAP varchar mode is enabled (trailing blanks in string constants are trimmed off).
	coSelectForUpdateSupported            connectOption = 14 //!< SELECT FOR UPDATE function code understood by client
	coClientDistributionMode              connectOption = 15 //!< client distribution mode
	coEngineDataFormatVersion             connectOption = 16 //!< Engine version of data format used in communication (@see DataFormatVersionEnum).
	coDistributionProtocolVersion         connectOption = 17 //!< version of distribution protocol handling (@see DistributionProtocolVersionEnum)
	coSplitBatchCommands                  connectOption = 18 //!< permit splitting of batch commands
	coUseTransactionFlagsOnly             connectOption = 19 //!< use transaction flags only for controlling transaction
	coRowSlotImageParameter               connectOption = 20 //!< row-slot image parameter passing
	coIgnoreUnknownParts                  connectOption = 21 //!< server does not abort on unknown parts
	coTableOutputParameterMetadataSupport connectOption = 22 //!< support table type output parameter metadata.
	coDataFormatVersion2                  connectOption = 23 //!< Version of data format used in communication (as DataFormatVersion used wrongly in old servers)
	coItabParameter                       connectOption = 24 //!< bool option to signal abap itab parameter support
	coDescribeTableOutputParameter        connectOption = 25 //!< override "omit table output parameter" setting in this session
	coColumnarResultSet                   connectOption = 26 //!< column wise result passing
	coScrollableResultSet                 connectOption = 27 //!< scrollable result set
	coClientInfoNullValueSupported        connectOption = 28 //!< can handle null values in client info
	coAssociatedConnectionID              connectOption = 29 //!< associated connection id
	coNonTransactionalPrepare             connectOption = 30 //!< can handle and uses non-transactional prepare
	coFdaEnabled                          connectOption = 31 //!< Fast Data Access at all enabled
	coOSUser                              connectOption = 32 //!< client OS user name
	coRowSlotImageResultSet               connectOption = 33 //!< row-slot image result passing
	coEndianness                          connectOption = 34 //!< endianness (@see EndiannessEnumType)
	coUpdateTopologyAnwhere               connectOption = 35 //!< Allow update of topology from any reply
	coEnableArrayType                     connectOption = 36 //!< Enable supporting Array data type
	coImplicitLobStreaming                connectOption = 37 //!< implicit lob streaming
	coCachedViewProperty                  connectOption = 38 //!< provide cached view timestamps to the client
	coXOpenXAProtocolSupported            connectOption = 39 //!< JTA(X/Open XA) Protocol
	coPrimaryCommitRedirectionSupported   connectOption = 40 //!< S2PC routing control
	coActiveActiveProtocolVersion         connectOption = 41 //!< Version of Active/Active protocol
	coActiveActiveConnectionOriginSite    connectOption = 42 //!< Tell where is the anchor connection located. This is unidirectional property from client to server.
	coQueryTimeoutSupported               connectOption = 43 //!< support query timeout (e.g., Statement.setQueryTimeout)
	coFullVersionString                   connectOption = 44 //!< Full version string of the client or server (the sender) (added to hana2sp0)
	coDatabaseName                        connectOption = 45 //!< Database name (string) that we connected to (sent by server) (added to hana2sp0)
	coBuildPlatform                       connectOption = 46 //!< Build platform of the client or server (the sender) (added to hana2sp0)
	coImplicitXASessionSupported          connectOption = 47 //!< S2PC routing control - implicit XA join support after prepare and before execute in MessageType_Prepare, MessageType_Execute and MessageType_PrepareAndExecute
	coClientSideColumnEncryptionVersion   connectOption = 48 //!< Version of client-side column encryption
	coCompressionLevelAndFlags            connectOption = 49 //!< Network compression level and flags (added to hana2sp02)
	coClientSideReExecutionSupported      connectOption = 50 //!< support client-side re-execution for client-side encryption (added to hana2sp03)
	coClientReconnectWaitTimeout          connectOption = 51 //!< client reconnection wait timeout for transparent session recovery
	coOriginalAnchorConnectionID          connectOption = 52 //!< original anchor connectionID to notify client's RECONNECT
	coFlagSet1                            connectOption = 53 //!< flags for aggregating several options
	coTopologyNetworkGroup                connectOption = 54 //!< NetworkGroup name sent by client to choose topology mapping (added to hana2sp04)
	coIPAddress                           connectOption = 55 //!< IP Address of the sender (added to hana2sp04)
	coLRRPingTime                         connectOption = 56 //!< Long running request ping time
)
