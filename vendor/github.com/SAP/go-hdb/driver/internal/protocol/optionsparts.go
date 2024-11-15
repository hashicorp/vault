package protocol

import (
	"fmt"
	"slices"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

// ClientContextOption represents a client context option.
type clientContextOption int8

func (k clientContextOption) valueString(v any) string {
	return fmt.Sprintf("%s: %v", k, v)
}

// ClientContextOption constants.
const (
	ccoVersion            clientContextOption = 1
	ccoType               clientContextOption = 2
	ccoApplicationProgram clientContextOption = 3
)

// ClientContext represents a client context part.
type ClientContext struct {
	options[clientContextOption]
}

// SetVersion sets the client version option.
func (cc *ClientContext) SetVersion(v string) { cc.options.set(ccoVersion, v) }

// SetType sets the client type option.
func (cc *ClientContext) SetType(v string) { cc.options.set(ccoType, v) }

// SetApplicationProgram sets the client application program option.
func (cc *ClientContext) SetApplicationProgram(v string) { cc.options.set(ccoApplicationProgram, v) }

// Cdm represents a ConnectOption ClientDistributionMode.
type Cdm byte

// ConnectOption ClientDistributionMode constants.
const (
	CdmOff                 Cdm = 0
	CdmConnection          Cdm = 1
	CdmStatement           Cdm = 2
	CdmConnectionStatement Cdm = 3
)

// dpv represents a ConnectOption DistributionProtocolVersion.
type dpv byte

// distribution protocol version

// ConnectOption DistributionProtocolVersion constants.
const (
	dpvBaseline                       dpv = 0
	dpvClientHandlesStatementSequence dpv = 1
)

// ConnectOption represents a connect option.
type connectOption int8

func (k connectOption) valueString(v any) string {
	// TODO: sub options
	return fmt.Sprintf("%s: %v", k, v)
}

// ConnectOption constants.
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

// ConnectOptions represents a connect options part.
type ConnectOptions struct {
	options[connectOption]
}

// DataFormatVersion2OrZero returns the data format version2 option if available, the zero value otherwise.
func (co *ConnectOptions) DataFormatVersion2OrZero() int {
	var v int32
	co.options.get(coDataFormatVersion2, &v)
	return int(v)
}

// SetDataFormatVersion2 sets the data format version 2 option.
func (co *ConnectOptions) SetDataFormatVersion2(v int) {
	co.options.set(coDataFormatVersion2, int32(v)) //nolint: gosec
}

// SetClientDistributionMode sets the client distribution mode option.
func (co *ConnectOptions) SetClientDistributionMode(v Cdm) {
	co.options.set(coClientDistributionMode, int32(v))
}

// SetSelectForUpdateSupported sets the select for update supported option.
func (co *ConnectOptions) SetSelectForUpdateSupported(v bool) {
	co.options.set(coSelectForUpdateSupported, v)
}

// DatabaseNameOrZero returns the database name option if available, the zero value otherwise.
func (co *ConnectOptions) DatabaseNameOrZero() string {
	var v string
	co.options.get(coDatabaseName, &v)
	return v
}

// FullVersionOrZero returns the full version option if available, the zero value otherwise.
func (co *ConnectOptions) FullVersionOrZero() string {
	var v string
	co.options.get(coFullVersionString, &v)
	return v
}

// SetClientLocale sets the client locale option.
func (co *ConnectOptions) SetClientLocale(v string) { co.options.set(coClientLocale, v) }

// DBConnectInfoType represents a database connect info type.
type dbConnectInfoType int8

func (k dbConnectInfoType) valueString(v any) string {
	return fmt.Sprintf("%s: %v", k, v)
}

// DBConnectInfoType constants.
const (
	ciDatabaseName dbConnectInfoType = 1 // string
	ciHost         dbConnectInfoType = 2 // string
	ciPort         dbConnectInfoType = 3 // int4
	ciIsConnected  dbConnectInfoType = 4 // bool
)

// DBConnectInfo represents a database connect info part.
type DBConnectInfo struct {
	options[dbConnectInfoType]
}

// SetDatabaseName sets the database name option.
func (ci *DBConnectInfo) SetDatabaseName(v string) { ci.options.set(ciDatabaseName, v) }

// HostOrZero returns the host option, the zero value otherwise.
func (ci *DBConnectInfo) HostOrZero() string { var v string; ci.options.get(ciHost, &v); return v }

// PortOrZero returns the port option, the zero value otherwise.
func (ci *DBConnectInfo) PortOrZero() int { var v int32; ci.options.get(ciPort, &v); return int(v) }

// IsConnectedOrZero returns this IsConnected option, the zero value otherwise.
func (ci *DBConnectInfo) IsConnectedOrZero() bool {
	var v bool
	ci.options.get(ciIsConnected, &v)
	return v
}

type statementContextType int8

func (k statementContextType) valueString(v any) string {
	return fmt.Sprintf("%s: %v", k, v)
}

const (
	scStatementSequenceInfo         statementContextType = 1
	scServerProcessingTime          statementContextType = 2
	scSchemaName                    statementContextType = 3
	scFlagSet                       statementContextType = 4
	scQueryTimeout                  statementContextType = 5
	scClientReconnectionWaitTimeout statementContextType = 6
	scServerCPUTime                 statementContextType = 7
	scServerMemoryUsage             statementContextType = 8
)

type statementContext struct {
	options[statementContextType]
}

// transaction flags.
type transactionFlagType int8

func (k transactionFlagType) valueString(v any) string {
	return fmt.Sprintf("%s: %v", k, v)
}

const (
	tfRolledback                      transactionFlagType = 0
	tfCommited                        transactionFlagType = 1
	tfNewIsolationLevel               transactionFlagType = 2
	tfDDLCommitmodeChanged            transactionFlagType = 3
	tfWriteTransactionStarted         transactionFlagType = 4
	tfNowriteTransactionStarted       transactionFlagType = 5
	tfSessionClosingTransactionError  transactionFlagType = 6
	tfSessionClosingTransactionErrror transactionFlagType = 7
	tfReadOnlyMode                    transactionFlagType = 8
)

type transactionFlags struct {
	options[transactionFlagType]
}

type topologyOption int8

func (k topologyOption) valueString(v any) string {
	switch k {
	case toServiceType:
		v := v.(int32)
		return fmt.Sprintf("%s: %v", k, ServiceType(v))
	default:
		return fmt.Sprintf("%s: %v", k, v)
	}
}

const (
	toHostName         topologyOption = 1
	toHostPortnumber   topologyOption = 2
	toTenantName       topologyOption = 3
	toLoadfactor       topologyOption = 4
	toVolumeID         topologyOption = 5
	toIsPrimary        topologyOption = 6
	toIsCurrentSession topologyOption = 7
	toServiceType      topologyOption = 8
	toNetworkDomain    topologyOption = 9 // deprecated
	toIsStandby        topologyOption = 10
	toAllIPAddresses   topologyOption = 11 // deprecated
	toAllHostNames     topologyOption = 12 // deprecated
	toSiteType         topologyOption = 13
)

// ServiceType represents a service type.
type ServiceType int32

// Service type constants.
const (
	StOther            ServiceType = 0
	StNameServer       ServiceType = 1
	StPreprocessor     ServiceType = 2
	StIndexServer      ServiceType = 3
	StStatisticsServer ServiceType = 4
	StXSEngine         ServiceType = 5
	StReserved6        ServiceType = 6
	StCompileServer    ServiceType = 7
	StDPServer         ServiceType = 8
	StDIServer         ServiceType = 9
	StComputeServer    ServiceType = 10
	StScriptServer     ServiceType = 11
)

// TopologyInformation represents a topology information part.
type TopologyInformation struct {
	hosts []*options[topologyOption]
}

func (ti TopologyInformation) String() string { return fmt.Sprintf("%v", ti.hosts) }

func (ti *TopologyInformation) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	ti.hosts = resizeSlice(ti.hosts, numArg)
	for i := range numArg {
		host := &options[topologyOption]{}
		ti.hosts[i] = host
		hostNumArg := int(dec.Int16())
		if err := host.decodeNumArg(dec, hostNumArg); err != nil {
			return err
		}
	}
	return dec.Error()
}

type optionsType interface {
	~int8
	valueString(v any) string
}

// options represents a generic option part.
type options[K optionsType] map[K]any

func (ops options[K]) String() string {
	s := []string{}
	for k, v := range ops {
		s = append(s, k.valueString(v))
	}
	slices.Sort(s)
	return fmt.Sprintf("%v", s)
}

func (ops *options[K]) get(k K, v any) bool {
	if *ops == nil {
		return false
	}
	mv, ok := (*ops)[k]
	if !ok {
		return false
	}
	switch v := v.(type) {
	case *string:
		*v = mv.(string)
	case *bool:
		*v = mv.(bool)
	case *int32:
		*v = mv.(int32)
	default:
		panic("invalid option type")
	}
	return true
}

func (ops *options[K]) set(k K, v any) {
	if *ops == nil {
		*ops = options[K]{}
	}
	(*ops)[k] = v
}

func (ops options[K]) size() int {
	size := 2 * len(ops) // option + type
	for _, v := range ops {
		ot := optTypeViaType(v)
		size += ot.size(v)
	}
	return size
}

func (ops options[K]) numArg() int { return len(ops) }

func (ops *options[K]) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	*ops = options[K]{} // no reuse of maps - create new one
	for range numArg {
		k := K(dec.Int8())
		tc := typeCode(dec.Byte())
		ot := optTypeViaTypeCode(tc)
		(*ops)[k] = ot.decode(dec)
	}
	return dec.Error()
}

func (ops options[K]) encode(enc *encoding.Encoder) error {
	for k, v := range ops {
		enc.Int8(int8(k))
		ot := optTypeViaType(v)
		enc.Int8(int8(ot.typeCode()))
		ot.encode(enc, v)
	}
	return nil
}
