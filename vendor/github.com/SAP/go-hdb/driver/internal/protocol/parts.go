package protocol

import (
	"reflect"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

// Part represents a protocol part.
type Part interface {
	String() string // should support Stringer interface
	kind() PartKind
}

type partDecoder interface {
	Part
	decode(dec *encoding.Decoder) error
}
type numArgPartDecoder interface {
	Part
	decodeNumArg(dec *encoding.Decoder, numArg int) error
}
type bufLenPartDecoder interface {
	Part
	decodeBufLen(dec *encoding.Decoder, bufLen int) error
}
type resultPartDecoder interface {
	Part
	decodeResult(dec *encoding.Decoder, numArg int, lobReader LobReader, lobChunkSize int) error
}

// PartEncoder represents a protocol part the driver is able to encode.
type PartEncoder interface {
	Part
	numArg() int
	size() int
	encode(enc *encoding.Encoder) error
}

func (*HdbErrors) kind() PartKind           { return pkError }
func (*AuthInitRequest) kind() PartKind     { return PkAuthentication }
func (*AuthInitReply) kind() PartKind       { return PkAuthentication }
func (*AuthFinalRequest) kind() PartKind    { return PkAuthentication }
func (*AuthFinalReply) kind() PartKind      { return PkAuthentication }
func (ClientID) kind() PartKind             { return PkClientID }
func (clientInfo) kind() PartKind           { return PkClientInfo }
func (*TopologyInformation) kind() PartKind { return PkTopologyInformation }
func (Command) kind() PartKind              { return PkCommand }
func (*rowsAffected) kind() PartKind        { return pkRowsAffected }
func (StatementID) kind() PartKind          { return PkStatementID }
func (*ParameterMetadata) kind() PartKind   { return PkParameterMetadata }
func (*InputParameters) kind() PartKind     { return PkParameters }
func (*OutputParameters) kind() PartKind    { return PkOutputParameters }
func (*ResultMetadata) kind() PartKind      { return PkResultMetadata }
func (ResultsetID) kind() PartKind          { return PkResultsetID }
func (*Resultset) kind() PartKind           { return PkResultset }
func (Fetchsize) kind() PartKind            { return PkFetchSize }
func (*ReadLobRequest) kind() PartKind      { return PkReadLobRequest }
func (*ReadLobReply) kind() PartKind        { return PkReadLobReply }
func (*WriteLobRequest) kind() PartKind     { return PkWriteLobRequest }
func (*WriteLobReply) kind() PartKind       { return PkWriteLobReply }
func (*ClientContext) kind() PartKind       { return PkClientContext }
func (*ConnectOptions) kind() PartKind      { return PkConnectOptions }
func (*DBConnectInfo) kind() PartKind       { return PkDBConnectInfo }
func (*statementContext) kind() PartKind    { return PkStatementContext }
func (*transactionFlags) kind() PartKind    { return PkTransactionFlags }

// numArg methods (result == 1).
func (*AuthInitRequest) numArg() int  { return 1 }
func (*AuthFinalRequest) numArg() int { return 1 }
func (ClientID) numArg() int          { return 1 }
func (Command) numArg() int           { return 1 }
func (StatementID) numArg() int       { return 1 }
func (ResultsetID) numArg() int       { return 1 }
func (Fetchsize) numArg() int         { return 1 }
func (*ReadLobRequest) numArg() int   { return 1 }

// size methods (fixed size).
const (
	statementIDSize    = 8
	resultsetIDSize    = 8
	fetchsizeSize      = 4
	readLobRequestSize = 24
)

func (StatementID) size() int    { return statementIDSize }
func (ResultsetID) size() int    { return resultsetIDSize }
func (Fetchsize) size() int      { return fetchsizeSize }
func (ReadLobRequest) size() int { return readLobRequestSize }

// func (lobFlags) size() int       { return tinyintFieldSize }

// check if part types implement the part encoder interface.
var (
	_ PartEncoder = (*AuthInitRequest)(nil)
	_ PartEncoder = (*AuthFinalRequest)(nil)
	_ PartEncoder = (*ClientID)(nil)
	_ PartEncoder = (*clientInfo)(nil)
	_ PartEncoder = (*Command)(nil)
	_ PartEncoder = (*StatementID)(nil)
	_ PartEncoder = (*InputParameters)(nil)
	_ PartEncoder = (*ResultsetID)(nil)
	_ PartEncoder = (*Fetchsize)(nil)
	_ PartEncoder = (*ReadLobRequest)(nil)
	_ PartEncoder = (*WriteLobRequest)(nil)
	_ PartEncoder = (*ClientContext)(nil)
	_ PartEncoder = (*ConnectOptions)(nil)
	_ PartEncoder = (*DBConnectInfo)(nil)
)

// check if part types implement the right part decoder interface.
var (
	_ numArgPartDecoder = (*HdbErrors)(nil)
	_ partDecoder       = (*AuthInitRequest)(nil)
	_ partDecoder       = (*AuthInitReply)(nil)
	_ partDecoder       = (*AuthFinalRequest)(nil)
	_ partDecoder       = (*AuthFinalReply)(nil)
	_ bufLenPartDecoder = (*ClientID)(nil)
	_ numArgPartDecoder = (*clientInfo)(nil)
	_ numArgPartDecoder = (*TopologyInformation)(nil)
	_ bufLenPartDecoder = (*Command)(nil)
	_ numArgPartDecoder = (*rowsAffected)(nil)
	_ partDecoder       = (*StatementID)(nil)
	_ numArgPartDecoder = (*ParameterMetadata)(nil)
	_ numArgPartDecoder = (*InputParameters)(nil)
	_ resultPartDecoder = (*OutputParameters)(nil)
	_ numArgPartDecoder = (*ResultMetadata)(nil)
	_ partDecoder       = (*ResultsetID)(nil)
	_ resultPartDecoder = (*Resultset)(nil)
	_ partDecoder       = (*Fetchsize)(nil)
	_ partDecoder       = (*ReadLobRequest)(nil)
	_ numArgPartDecoder = (*WriteLobRequest)(nil)
	_ numArgPartDecoder = (*ReadLobReply)(nil)
	_ numArgPartDecoder = (*WriteLobReply)(nil)
	_ numArgPartDecoder = (*ClientContext)(nil)
	_ numArgPartDecoder = (*ConnectOptions)(nil)
	_ numArgPartDecoder = (*DBConnectInfo)(nil)
	_ numArgPartDecoder = (*statementContext)(nil)
	_ numArgPartDecoder = (*transactionFlags)(nil)
)

var genPartTypeMap = map[PartKind]reflect.Type{
	pkError:               reflect.TypeFor[HdbErrors](),
	PkClientID:            reflect.TypeFor[ClientID](),
	PkClientInfo:          reflect.TypeFor[clientInfo](),
	PkTopologyInformation: reflect.TypeFor[TopologyInformation](),
	PkCommand:             reflect.TypeFor[Command](),
	pkRowsAffected:        reflect.TypeFor[rowsAffected](),
	PkStatementID:         reflect.TypeFor[StatementID](),
	PkResultsetID:         reflect.TypeFor[ResultsetID](),
	PkFetchSize:           reflect.TypeFor[Fetchsize](),
	PkReadLobRequest:      reflect.TypeFor[ReadLobRequest](),
	PkReadLobReply:        reflect.TypeFor[ReadLobReply](),
	PkWriteLobReply:       reflect.TypeFor[WriteLobReply](),
	PkWriteLobRequest:     reflect.TypeFor[WriteLobRequest](),
	PkClientContext:       reflect.TypeFor[ClientContext](),
	PkConnectOptions:      reflect.TypeFor[ConnectOptions](),
	PkTransactionFlags:    reflect.TypeFor[transactionFlags](),
	PkStatementContext:    reflect.TypeFor[statementContext](),
	PkDBConnectInfo:       reflect.TypeFor[DBConnectInfo](),
	/*
	   parts that cannot be used generically as additional parameters are needed

	   PkParameterMetadata
	   PkParameters
	   PkOutputParameters
	   PkResultMetadata
	   PkResultset
	*/
}

// to be implemented by parts needing initialization
// in case the part is instatiated generically.
type initer interface {
	init()
}

// newGenPartReader returns a generic part reader.
func newGenPartReader(kind PartKind) Part {
	if kind == PkAuthentication {
		return nil // cannot instantiate generically
	}
	pt, ok := genPartTypeMap[kind]
	if !ok {
		// whether part cannot be instantiated generically or
		// part is not (yet) known to the driver
		return nil
	}
	// create instance
	part, ok := reflect.New(pt).Interface().(Part)
	if !ok {
		panic("part kind does not implement part reader interface") // should never happen
	}
	if part, ok := part.(initer); ok {
		part.init()
	}
	return part
}
