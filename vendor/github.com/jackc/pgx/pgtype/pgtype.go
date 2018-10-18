package pgtype

import (
	"reflect"

	"github.com/pkg/errors"
)

// PostgreSQL oids for common types
const (
	BoolOID             = 16
	ByteaOID            = 17
	CharOID             = 18
	NameOID             = 19
	Int8OID             = 20
	Int2OID             = 21
	Int4OID             = 23
	TextOID             = 25
	OIDOID              = 26
	TIDOID              = 27
	XIDOID              = 28
	CIDOID              = 29
	JSONOID             = 114
	CIDROID             = 650
	CIDRArrayOID        = 651
	Float4OID           = 700
	Float8OID           = 701
	UnknownOID          = 705
	InetOID             = 869
	BoolArrayOID        = 1000
	Int2ArrayOID        = 1005
	Int4ArrayOID        = 1007
	TextArrayOID        = 1009
	ByteaArrayOID       = 1001
	BPCharArrayOID      = 1014
	VarcharArrayOID     = 1015
	Int8ArrayOID        = 1016
	Float4ArrayOID      = 1021
	Float8ArrayOID      = 1022
	ACLItemOID          = 1033
	ACLItemArrayOID     = 1034
	InetArrayOID        = 1041
	BPCharOID           = 1042
	VarcharOID          = 1043
	DateOID             = 1082
	TimestampOID        = 1114
	TimestampArrayOID   = 1115
	DateArrayOID        = 1182
	TimestamptzOID      = 1184
	TimestamptzArrayOID = 1185
	NumericOID          = 1700
	RecordOID           = 2249
	UUIDOID             = 2950
	UUIDArrayOID        = 2951
	JSONBOID            = 3802
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

type InfinityModifier int8

const (
	Infinity         InfinityModifier = 1
	None             InfinityModifier = 0
	NegativeInfinity InfinityModifier = -Infinity
)

func (im InfinityModifier) String() string {
	switch im {
	case None:
		return "none"
	case Infinity:
		return "infinity"
	case NegativeInfinity:
		return "-infinity"
	default:
		return "invalid"
	}
}

type Value interface {
	// Set converts and assigns src to itself.
	Set(src interface{}) error

	// Get returns the simplest representation of Value. If the Value is Null or
	// Undefined that is the return value. If no simpler representation is
	// possible, then Get() returns Value.
	Get() interface{}

	// AssignTo converts and assigns the Value to dst. It MUST make a deep copy of
	// any reference types.
	AssignTo(dst interface{}) error
}

type BinaryDecoder interface {
	// DecodeBinary decodes src into BinaryDecoder. If src is nil then the
	// original SQL value is NULL. BinaryDecoder takes ownership of src. The
	// caller MUST not use it again.
	DecodeBinary(ci *ConnInfo, src []byte) error
}

type TextDecoder interface {
	// DecodeText decodes src into TextDecoder. If src is nil then the original
	// SQL value is NULL. TextDecoder takes ownership of src. The caller MUST not
	// use it again.
	DecodeText(ci *ConnInfo, src []byte) error
}

// BinaryEncoder is implemented by types that can encode themselves into the
// PostgreSQL binary wire format.
type BinaryEncoder interface {
	// EncodeBinary should append the binary format of self to buf. If self is the
	// SQL value NULL then append nothing and return (nil, nil). The caller of
	// EncodeBinary is responsible for writing the correct NULL value or the
	// length of the data written.
	EncodeBinary(ci *ConnInfo, buf []byte) (newBuf []byte, err error)
}

// TextEncoder is implemented by types that can encode themselves into the
// PostgreSQL text wire format.
type TextEncoder interface {
	// EncodeText should append the text format of self to buf. If self is the
	// SQL value NULL then append nothing and return (nil, nil). The caller of
	// EncodeText is responsible for writing the correct NULL value or the
	// length of the data written.
	EncodeText(ci *ConnInfo, buf []byte) (newBuf []byte, err error)
}

var errUndefined = errors.New("cannot encode status undefined")
var errBadStatus = errors.New("invalid status")

type DataType struct {
	Value Value
	Name  string
	OID   OID
}

type ConnInfo struct {
	oidToDataType         map[OID]*DataType
	nameToDataType        map[string]*DataType
	reflectTypeToDataType map[reflect.Type]*DataType
}

func NewConnInfo() *ConnInfo {
	return &ConnInfo{
		oidToDataType:         make(map[OID]*DataType, 256),
		nameToDataType:        make(map[string]*DataType, 256),
		reflectTypeToDataType: make(map[reflect.Type]*DataType, 256),
	}
}

func (ci *ConnInfo) InitializeDataTypes(nameOIDs map[string]OID) {
	for name, oid := range nameOIDs {
		var value Value
		if t, ok := nameValues[name]; ok {
			value = reflect.New(reflect.ValueOf(t).Elem().Type()).Interface().(Value)
		} else {
			value = &GenericText{}
		}
		ci.RegisterDataType(DataType{Value: value, Name: name, OID: oid})
	}
}

func (ci *ConnInfo) RegisterDataType(t DataType) {
	ci.oidToDataType[t.OID] = &t
	ci.nameToDataType[t.Name] = &t
	ci.reflectTypeToDataType[reflect.ValueOf(t.Value).Type()] = &t
}

func (ci *ConnInfo) DataTypeForOID(oid OID) (*DataType, bool) {
	dt, ok := ci.oidToDataType[oid]
	return dt, ok
}

func (ci *ConnInfo) DataTypeForName(name string) (*DataType, bool) {
	dt, ok := ci.nameToDataType[name]
	return dt, ok
}

func (ci *ConnInfo) DataTypeForValue(v Value) (*DataType, bool) {
	dt, ok := ci.reflectTypeToDataType[reflect.ValueOf(v).Type()]
	return dt, ok
}

// DeepCopy makes a deep copy of the ConnInfo.
func (ci *ConnInfo) DeepCopy() *ConnInfo {
	ci2 := &ConnInfo{
		oidToDataType:         make(map[OID]*DataType, len(ci.oidToDataType)),
		nameToDataType:        make(map[string]*DataType, len(ci.nameToDataType)),
		reflectTypeToDataType: make(map[reflect.Type]*DataType, len(ci.reflectTypeToDataType)),
	}

	for _, dt := range ci.oidToDataType {
		ci2.RegisterDataType(DataType{
			Value: reflect.New(reflect.ValueOf(dt.Value).Elem().Type()).Interface().(Value),
			Name:  dt.Name,
			OID:   dt.OID,
		})
	}

	return ci2
}

var nameValues map[string]Value

func init() {
	nameValues = map[string]Value{
		"_aclitem":     &ACLItemArray{},
		"_bool":        &BoolArray{},
		"_bpchar":      &BPCharArray{},
		"_bytea":       &ByteaArray{},
		"_cidr":        &CIDRArray{},
		"_date":        &DateArray{},
		"_float4":      &Float4Array{},
		"_float8":      &Float8Array{},
		"_inet":        &InetArray{},
		"_int2":        &Int2Array{},
		"_int4":        &Int4Array{},
		"_int8":        &Int8Array{},
		"_numeric":     &NumericArray{},
		"_text":        &TextArray{},
		"_timestamp":   &TimestampArray{},
		"_timestamptz": &TimestamptzArray{},
		"_uuid":        &UUIDArray{},
		"_varchar":     &VarcharArray{},
		"aclitem":      &ACLItem{},
		"bit":          &Bit{},
		"bool":         &Bool{},
		"box":          &Box{},
		"bpchar":       &BPChar{},
		"bytea":        &Bytea{},
		"char":         &QChar{},
		"cid":          &CID{},
		"cidr":         &CIDR{},
		"circle":       &Circle{},
		"date":         &Date{},
		"daterange":    &Daterange{},
		"decimal":      &Decimal{},
		"float4":       &Float4{},
		"float8":       &Float8{},
		"hstore":       &Hstore{},
		"inet":         &Inet{},
		"int2":         &Int2{},
		"int4":         &Int4{},
		"int4range":    &Int4range{},
		"int8":         &Int8{},
		"int8range":    &Int8range{},
		"interval":     &Interval{},
		"json":         &JSON{},
		"jsonb":        &JSONB{},
		"line":         &Line{},
		"lseg":         &Lseg{},
		"macaddr":      &Macaddr{},
		"name":         &Name{},
		"numeric":      &Numeric{},
		"numrange":     &Numrange{},
		"oid":          &OIDValue{},
		"path":         &Path{},
		"point":        &Point{},
		"polygon":      &Polygon{},
		"record":       &Record{},
		"text":         &Text{},
		"tid":          &TID{},
		"timestamp":    &Timestamp{},
		"timestamptz":  &Timestamptz{},
		"tsrange":      &Tsrange{},
		"tstzrange":    &Tstzrange{},
		"unknown":      &Unknown{},
		"uuid":         &UUID{},
		"varbit":       &Varbit{},
		"varchar":      &Varchar{},
		"xid":          &XID{},
	}
}
