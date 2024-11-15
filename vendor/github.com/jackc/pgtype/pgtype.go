package pgtype

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"time"
)

// PostgreSQL oids for common types
const (
	BoolOID             = 16
	ByteaOID            = 17
	QCharOID            = 18
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
	JSONArrayOID        = 199
	PointOID            = 600
	LsegOID             = 601
	PathOID             = 602
	BoxOID              = 603
	PolygonOID          = 604
	LineOID             = 628
	CIDROID             = 650
	CIDRArrayOID        = 651
	Float4OID           = 700
	Float8OID           = 701
	CircleOID           = 718
	UnknownOID          = 705
	MacaddrOID          = 829
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
	TimeOID             = 1083
	TimestampOID        = 1114
	TimestampArrayOID   = 1115
	DateArrayOID        = 1182
	TimestamptzOID      = 1184
	TimestamptzArrayOID = 1185
	IntervalOID         = 1186
	NumericArrayOID     = 1231
	BitOID              = 1560
	VarbitOID           = 1562
	NumericOID          = 1700
	RecordOID           = 2249
	UUIDOID             = 2950
	UUIDArrayOID        = 2951
	JSONBOID            = 3802
	JSONBArrayOID       = 3807
	DaterangeOID        = 3912
	Int4rangeOID        = 3904
	Int4multirangeOID   = 4451
	NumrangeOID         = 3906
	NummultirangeOID    = 4532
	TsrangeOID          = 3908
	TsrangeArrayOID     = 3909
	TstzrangeOID        = 3910
	TstzrangeArrayOID   = 3911
	Int8rangeOID        = 3926
	Int8multirangeOID   = 4536
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

// PostgreSQL format codes
const (
	TextFormatCode   = 0
	BinaryFormatCode = 1
)

// Value translates values to and from an internal canonical representation for the type. To actually be usable a type
// that implements Value should also implement some combination of BinaryDecoder, BinaryEncoder, TextDecoder,
// and TextEncoder.
//
// Operations that update a Value (e.g. Set, DecodeText, DecodeBinary) should entirely replace the value. e.g. Internal
// slices should be replaced not resized and reused. This allows Get and AssignTo to return a slice directly rather
// than incur a usually unnecessary copy.
type Value interface {
	// Set converts and assigns src to itself. Value takes ownership of src.
	Set(src interface{}) error

	// Get returns the simplest representation of Value. Get may return a pointer to an internal value but it must never
	// mutate that value. e.g. If Get returns a []byte Value must never change the contents of the []byte.
	Get() interface{}

	// AssignTo converts and assigns the Value to dst. AssignTo may a pointer to an internal value but it must never
	// mutate that value. e.g. If Get returns a []byte Value must never change the contents of the []byte.
	AssignTo(dst interface{}) error
}

// TypeValue is a Value where instances can represent different PostgreSQL types. This can be useful for
// representing types such as enums, composites, and arrays.
//
// In general, instances of TypeValue should not be used to directly represent a value. It should only be used as an
// encoder and decoder internal to ConnInfo.
type TypeValue interface {
	Value

	// NewTypeValue creates a TypeValue including references to internal type information. e.g. the list of members
	// in an EnumType.
	NewTypeValue() Value

	// TypeName returns the PostgreSQL name of this type.
	TypeName() string
}

// ValueTranscoder is a value that implements the text and binary encoding and decoding interfaces.
type ValueTranscoder interface {
	Value
	TextEncoder
	BinaryEncoder
	TextDecoder
	BinaryDecoder
}

// ResultFormatPreferrer allows a type to specify its preferred result format instead of it being inferred from
// whether it is also a BinaryDecoder.
type ResultFormatPreferrer interface {
	PreferredResultFormat() int16
}

// ParamFormatPreferrer allows a type to specify its preferred param format instead of it being inferred from
// whether it is also a BinaryEncoder.
type ParamFormatPreferrer interface {
	PreferredParamFormat() int16
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

type nullAssignmentError struct {
	dst interface{}
}

func (e *nullAssignmentError) Error() string {
	return fmt.Sprintf("cannot assign NULL to %T", e.dst)
}

type DataType struct {
	Value Value

	textDecoder   TextDecoder
	binaryDecoder BinaryDecoder

	Name string
	OID  uint32
}

type ConnInfo struct {
	oidToDataType         map[uint32]*DataType
	nameToDataType        map[string]*DataType
	reflectTypeToName     map[reflect.Type]string
	oidToParamFormatCode  map[uint32]int16
	oidToResultFormatCode map[uint32]int16

	reflectTypeToDataType map[reflect.Type]*DataType
}

func newConnInfo() *ConnInfo {
	return &ConnInfo{
		oidToDataType:         make(map[uint32]*DataType),
		nameToDataType:        make(map[string]*DataType),
		reflectTypeToName:     make(map[reflect.Type]string),
		oidToParamFormatCode:  make(map[uint32]int16),
		oidToResultFormatCode: make(map[uint32]int16),
	}
}

func NewConnInfo() *ConnInfo {
	ci := newConnInfo()

	ci.RegisterDataType(DataType{Value: &ACLItemArray{}, Name: "_aclitem", OID: ACLItemArrayOID})
	ci.RegisterDataType(DataType{Value: &BoolArray{}, Name: "_bool", OID: BoolArrayOID})
	ci.RegisterDataType(DataType{Value: &BPCharArray{}, Name: "_bpchar", OID: BPCharArrayOID})
	ci.RegisterDataType(DataType{Value: &ByteaArray{}, Name: "_bytea", OID: ByteaArrayOID})
	ci.RegisterDataType(DataType{Value: &CIDRArray{}, Name: "_cidr", OID: CIDRArrayOID})
	ci.RegisterDataType(DataType{Value: &DateArray{}, Name: "_date", OID: DateArrayOID})
	ci.RegisterDataType(DataType{Value: &Float4Array{}, Name: "_float4", OID: Float4ArrayOID})
	ci.RegisterDataType(DataType{Value: &Float8Array{}, Name: "_float8", OID: Float8ArrayOID})
	ci.RegisterDataType(DataType{Value: &InetArray{}, Name: "_inet", OID: InetArrayOID})
	ci.RegisterDataType(DataType{Value: &Int2Array{}, Name: "_int2", OID: Int2ArrayOID})
	ci.RegisterDataType(DataType{Value: &Int4Array{}, Name: "_int4", OID: Int4ArrayOID})
	ci.RegisterDataType(DataType{Value: &Int8Array{}, Name: "_int8", OID: Int8ArrayOID})
	ci.RegisterDataType(DataType{Value: &NumericArray{}, Name: "_numeric", OID: NumericArrayOID})
	ci.RegisterDataType(DataType{Value: &TextArray{}, Name: "_text", OID: TextArrayOID})
	ci.RegisterDataType(DataType{Value: &TimestampArray{}, Name: "_timestamp", OID: TimestampArrayOID})
	ci.RegisterDataType(DataType{Value: &TimestamptzArray{}, Name: "_timestamptz", OID: TimestamptzArrayOID})
	ci.RegisterDataType(DataType{Value: &UUIDArray{}, Name: "_uuid", OID: UUIDArrayOID})
	ci.RegisterDataType(DataType{Value: &VarcharArray{}, Name: "_varchar", OID: VarcharArrayOID})
	ci.RegisterDataType(DataType{Value: &ACLItem{}, Name: "aclitem", OID: ACLItemOID})
	ci.RegisterDataType(DataType{Value: &Bit{}, Name: "bit", OID: BitOID})
	ci.RegisterDataType(DataType{Value: &Bool{}, Name: "bool", OID: BoolOID})
	ci.RegisterDataType(DataType{Value: &Box{}, Name: "box", OID: BoxOID})
	ci.RegisterDataType(DataType{Value: &BPChar{}, Name: "bpchar", OID: BPCharOID})
	ci.RegisterDataType(DataType{Value: &Bytea{}, Name: "bytea", OID: ByteaOID})
	ci.RegisterDataType(DataType{Value: &QChar{}, Name: "char", OID: QCharOID})
	ci.RegisterDataType(DataType{Value: &CID{}, Name: "cid", OID: CIDOID})
	ci.RegisterDataType(DataType{Value: &CIDR{}, Name: "cidr", OID: CIDROID})
	ci.RegisterDataType(DataType{Value: &Circle{}, Name: "circle", OID: CircleOID})
	ci.RegisterDataType(DataType{Value: &Date{}, Name: "date", OID: DateOID})
	ci.RegisterDataType(DataType{Value: &Daterange{}, Name: "daterange", OID: DaterangeOID})
	ci.RegisterDataType(DataType{Value: &Float4{}, Name: "float4", OID: Float4OID})
	ci.RegisterDataType(DataType{Value: &Float8{}, Name: "float8", OID: Float8OID})
	ci.RegisterDataType(DataType{Value: &Inet{}, Name: "inet", OID: InetOID})
	ci.RegisterDataType(DataType{Value: &Int2{}, Name: "int2", OID: Int2OID})
	ci.RegisterDataType(DataType{Value: &Int4{}, Name: "int4", OID: Int4OID})
	ci.RegisterDataType(DataType{Value: &Int4range{}, Name: "int4range", OID: Int4rangeOID})
	ci.RegisterDataType(DataType{Value: &Int4multirange{}, Name: "int4multirange", OID: Int4multirangeOID})
	ci.RegisterDataType(DataType{Value: &Int8{}, Name: "int8", OID: Int8OID})
	ci.RegisterDataType(DataType{Value: &Int8range{}, Name: "int8range", OID: Int8rangeOID})
	ci.RegisterDataType(DataType{Value: &Int8multirange{}, Name: "int8multirange", OID: Int8multirangeOID})
	ci.RegisterDataType(DataType{Value: &Interval{}, Name: "interval", OID: IntervalOID})
	ci.RegisterDataType(DataType{Value: &JSON{}, Name: "json", OID: JSONOID})
	ci.RegisterDataType(DataType{Value: &JSONArray{}, Name: "_json", OID: JSONArrayOID})
	ci.RegisterDataType(DataType{Value: &JSONB{}, Name: "jsonb", OID: JSONBOID})
	ci.RegisterDataType(DataType{Value: &JSONBArray{}, Name: "_jsonb", OID: JSONBArrayOID})
	ci.RegisterDataType(DataType{Value: &Line{}, Name: "line", OID: LineOID})
	ci.RegisterDataType(DataType{Value: &Lseg{}, Name: "lseg", OID: LsegOID})
	ci.RegisterDataType(DataType{Value: &Macaddr{}, Name: "macaddr", OID: MacaddrOID})
	ci.RegisterDataType(DataType{Value: &Name{}, Name: "name", OID: NameOID})
	ci.RegisterDataType(DataType{Value: &Numeric{}, Name: "numeric", OID: NumericOID})
	ci.RegisterDataType(DataType{Value: &Numrange{}, Name: "numrange", OID: NumrangeOID})
	ci.RegisterDataType(DataType{Value: &Nummultirange{}, Name: "nummultirange", OID: NummultirangeOID})
	ci.RegisterDataType(DataType{Value: &OIDValue{}, Name: "oid", OID: OIDOID})
	ci.RegisterDataType(DataType{Value: &Path{}, Name: "path", OID: PathOID})
	ci.RegisterDataType(DataType{Value: &Point{}, Name: "point", OID: PointOID})
	ci.RegisterDataType(DataType{Value: &Polygon{}, Name: "polygon", OID: PolygonOID})
	ci.RegisterDataType(DataType{Value: &Record{}, Name: "record", OID: RecordOID})
	ci.RegisterDataType(DataType{Value: &Text{}, Name: "text", OID: TextOID})
	ci.RegisterDataType(DataType{Value: &TID{}, Name: "tid", OID: TIDOID})
	ci.RegisterDataType(DataType{Value: &Time{}, Name: "time", OID: TimeOID})
	ci.RegisterDataType(DataType{Value: &Timestamp{}, Name: "timestamp", OID: TimestampOID})
	ci.RegisterDataType(DataType{Value: &Timestamptz{}, Name: "timestamptz", OID: TimestamptzOID})
	ci.RegisterDataType(DataType{Value: &Tsrange{}, Name: "tsrange", OID: TsrangeOID})
	ci.RegisterDataType(DataType{Value: &TsrangeArray{}, Name: "_tsrange", OID: TsrangeArrayOID})
	ci.RegisterDataType(DataType{Value: &Tstzrange{}, Name: "tstzrange", OID: TstzrangeOID})
	ci.RegisterDataType(DataType{Value: &TstzrangeArray{}, Name: "_tstzrange", OID: TstzrangeArrayOID})
	ci.RegisterDataType(DataType{Value: &Unknown{}, Name: "unknown", OID: UnknownOID})
	ci.RegisterDataType(DataType{Value: &UUID{}, Name: "uuid", OID: UUIDOID})
	ci.RegisterDataType(DataType{Value: &Varbit{}, Name: "varbit", OID: VarbitOID})
	ci.RegisterDataType(DataType{Value: &Varchar{}, Name: "varchar", OID: VarcharOID})
	ci.RegisterDataType(DataType{Value: &XID{}, Name: "xid", OID: XIDOID})

	registerDefaultPgTypeVariants := func(name, arrayName string, value interface{}) {
		ci.RegisterDefaultPgType(value, name)
		valueType := reflect.TypeOf(value)

		ci.RegisterDefaultPgType(reflect.New(valueType).Interface(), name)

		sliceType := reflect.SliceOf(valueType)
		ci.RegisterDefaultPgType(reflect.MakeSlice(sliceType, 0, 0).Interface(), arrayName)

		ci.RegisterDefaultPgType(reflect.New(sliceType).Interface(), arrayName)
	}

	// Integer types that directly map to a PostgreSQL type
	registerDefaultPgTypeVariants("int2", "_int2", int16(0))
	registerDefaultPgTypeVariants("int4", "_int4", int32(0))
	registerDefaultPgTypeVariants("int8", "_int8", int64(0))

	// Integer types that do not have a direct match to a PostgreSQL type
	registerDefaultPgTypeVariants("int8", "_int8", uint16(0))
	registerDefaultPgTypeVariants("int8", "_int8", uint32(0))
	registerDefaultPgTypeVariants("int8", "_int8", uint64(0))
	registerDefaultPgTypeVariants("int8", "_int8", int(0))
	registerDefaultPgTypeVariants("int8", "_int8", uint(0))

	registerDefaultPgTypeVariants("float4", "_float4", float32(0))
	registerDefaultPgTypeVariants("float8", "_float8", float64(0))

	registerDefaultPgTypeVariants("bool", "_bool", false)
	registerDefaultPgTypeVariants("timestamptz", "_timestamptz", time.Time{})
	registerDefaultPgTypeVariants("text", "_text", "")
	registerDefaultPgTypeVariants("bytea", "_bytea", []byte(nil))

	registerDefaultPgTypeVariants("inet", "_inet", net.IP{})
	ci.RegisterDefaultPgType((*net.IPNet)(nil), "cidr")
	ci.RegisterDefaultPgType([]*net.IPNet(nil), "_cidr")

	return ci
}

func (ci *ConnInfo) InitializeDataTypes(nameOIDs map[string]uint32) {
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
	t.Value = NewValue(t.Value)

	ci.oidToDataType[t.OID] = &t
	ci.nameToDataType[t.Name] = &t

	{
		var formatCode int16
		if pfp, ok := t.Value.(ParamFormatPreferrer); ok {
			formatCode = pfp.PreferredParamFormat()
		} else if _, ok := t.Value.(BinaryEncoder); ok {
			formatCode = BinaryFormatCode
		}
		ci.oidToParamFormatCode[t.OID] = formatCode
	}

	{
		var formatCode int16
		if rfp, ok := t.Value.(ResultFormatPreferrer); ok {
			formatCode = rfp.PreferredResultFormat()
		} else if _, ok := t.Value.(BinaryDecoder); ok {
			formatCode = BinaryFormatCode
		}
		ci.oidToResultFormatCode[t.OID] = formatCode
	}

	if d, ok := t.Value.(TextDecoder); ok {
		t.textDecoder = d
	}

	if d, ok := t.Value.(BinaryDecoder); ok {
		t.binaryDecoder = d
	}

	ci.reflectTypeToDataType = nil // Invalidated by type registration
}

// RegisterDefaultPgType registers a mapping of a Go type to a PostgreSQL type name. Typically the data type to be
// encoded or decoded is determined by the PostgreSQL OID. But if the OID of a value to be encoded or decoded is
// unknown, this additional mapping will be used by DataTypeForValue to determine a suitable data type.
func (ci *ConnInfo) RegisterDefaultPgType(value interface{}, name string) {
	ci.reflectTypeToName[reflect.TypeOf(value)] = name
	ci.reflectTypeToDataType = nil // Invalidated by registering a default type
}

func (ci *ConnInfo) DataTypeForOID(oid uint32) (*DataType, bool) {
	dt, ok := ci.oidToDataType[oid]
	return dt, ok
}

func (ci *ConnInfo) DataTypeForName(name string) (*DataType, bool) {
	dt, ok := ci.nameToDataType[name]
	return dt, ok
}

func (ci *ConnInfo) buildReflectTypeToDataType() {
	ci.reflectTypeToDataType = make(map[reflect.Type]*DataType)

	for _, dt := range ci.oidToDataType {
		if _, is := dt.Value.(TypeValue); !is {
			ci.reflectTypeToDataType[reflect.ValueOf(dt.Value).Type()] = dt
		}
	}

	for reflectType, name := range ci.reflectTypeToName {
		if dt, ok := ci.nameToDataType[name]; ok {
			ci.reflectTypeToDataType[reflectType] = dt
		}
	}
}

// DataTypeForValue finds a data type suitable for v. Use RegisterDataType to register types that can encode and decode
// themselves. Use RegisterDefaultPgType to register that can be handled by a registered data type.
func (ci *ConnInfo) DataTypeForValue(v interface{}) (*DataType, bool) {
	if ci.reflectTypeToDataType == nil {
		ci.buildReflectTypeToDataType()
	}

	if tv, ok := v.(TypeValue); ok {
		dt, ok := ci.nameToDataType[tv.TypeName()]
		return dt, ok
	}

	dt, ok := ci.reflectTypeToDataType[reflect.TypeOf(v)]
	return dt, ok
}

func (ci *ConnInfo) ParamFormatCodeForOID(oid uint32) int16 {
	fc, ok := ci.oidToParamFormatCode[oid]
	if ok {
		return fc
	}
	return TextFormatCode
}

func (ci *ConnInfo) ResultFormatCodeForOID(oid uint32) int16 {
	fc, ok := ci.oidToResultFormatCode[oid]
	if ok {
		return fc
	}
	return TextFormatCode
}

// DeepCopy makes a deep copy of the ConnInfo.
func (ci *ConnInfo) DeepCopy() *ConnInfo {
	ci2 := newConnInfo()

	for _, dt := range ci.oidToDataType {
		ci2.RegisterDataType(DataType{
			Value: NewValue(dt.Value),
			Name:  dt.Name,
			OID:   dt.OID,
		})
	}

	for t, n := range ci.reflectTypeToName {
		ci2.reflectTypeToName[t] = n
	}

	return ci2
}

// ScanPlan is a precompiled plan to scan into a type of destination.
type ScanPlan interface {
	// Scan scans src into dst. If the dst type has changed in an incompatible way a ScanPlan should automatically
	// replan and scan.
	Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error
}

type scanPlanDstBinaryDecoder struct{}

func (scanPlanDstBinaryDecoder) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if d, ok := (dst).(BinaryDecoder); ok {
		return d.DecodeBinary(ci, src)
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanDstTextDecoder struct{}

func (plan scanPlanDstTextDecoder) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if d, ok := (dst).(TextDecoder); ok {
		return d.DecodeText(ci, src)
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanDataTypeSQLScanner DataType

func (plan *scanPlanDataTypeSQLScanner) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	scanner, ok := dst.(sql.Scanner)
	if !ok {
		dv := reflect.ValueOf(dst)
		if dv.Kind() != reflect.Ptr || !dv.Type().Elem().Implements(scannerType) {
			newPlan := ci.PlanScan(oid, formatCode, dst)
			return newPlan.Scan(ci, oid, formatCode, src, dst)
		}
		if src == nil {
			// Ensure the pointer points to a zero version of the value
			dv.Elem().Set(reflect.Zero(dv.Type().Elem()))
			return nil
		}
		dv = dv.Elem()
		// If the pointer is to a nil pointer then set that before scanning
		if dv.Kind() == reflect.Ptr && dv.IsNil() {
			dv.Set(reflect.New(dv.Type().Elem()))
		}
		scanner = dv.Interface().(sql.Scanner)
	}

	dt := (*DataType)(plan)
	var err error
	switch formatCode {
	case BinaryFormatCode:
		err = dt.binaryDecoder.DecodeBinary(ci, src)
	case TextFormatCode:
		err = dt.textDecoder.DecodeText(ci, src)
	}
	if err != nil {
		return err
	}

	sqlSrc, err := DatabaseSQLValue(ci, dt.Value)
	if err != nil {
		return err
	}
	return scanner.Scan(sqlSrc)
}

type scanPlanDataTypeAssignTo DataType

func (plan *scanPlanDataTypeAssignTo) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	dt := (*DataType)(plan)
	var err error
	switch formatCode {
	case BinaryFormatCode:
		err = dt.binaryDecoder.DecodeBinary(ci, src)
	case TextFormatCode:
		err = dt.textDecoder.DecodeText(ci, src)
	}
	if err != nil {
		return err
	}

	assignToErr := dt.Value.AssignTo(dst)
	if assignToErr == nil {
		return nil
	}

	if dstPtr, ok := dst.(*interface{}); ok {
		*dstPtr = dt.Value.Get()
		return nil
	}

	// assignToErr might have failed because the type of destination has changed
	newPlan := ci.PlanScan(oid, formatCode, dst)
	if newPlan, sameType := newPlan.(*scanPlanDataTypeAssignTo); !sameType {
		return newPlan.Scan(ci, oid, formatCode, src, dst)
	}

	return assignToErr
}

type scanPlanSQLScanner struct{}

func (scanPlanSQLScanner) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	scanner, ok := dst.(sql.Scanner)
	if !ok {
		dv := reflect.ValueOf(dst)
		if dv.Kind() != reflect.Ptr || !dv.Type().Elem().Implements(scannerType) {
			newPlan := ci.PlanScan(oid, formatCode, dst)
			return newPlan.Scan(ci, oid, formatCode, src, dst)
		}
		if src == nil {
			// Ensure the pointer points to a zero version of the value
			dv.Elem().Set(reflect.Zero(dv.Elem().Type()))
			return nil
		}
		dv = dv.Elem()
		// If the pointer is to a nil pointer then set that before scanning
		if dv.Kind() == reflect.Ptr && dv.IsNil() {
			dv.Set(reflect.New(dv.Type().Elem()))
		}
		scanner = dv.Interface().(sql.Scanner)
	}
	if src == nil {
		// This is necessary because interface value []byte:nil does not equal nil:nil for the binary format path and the
		// text format path would be converted to empty string.
		return scanner.Scan(nil)
	} else if formatCode == BinaryFormatCode {
		return scanner.Scan(src)
	} else {
		return scanner.Scan(string(src))
	}
}

type scanPlanReflection struct{}

func (scanPlanReflection) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	// We might be given a pointer to something that implements the decoder interface(s),
	// even though the pointer itself doesn't.
	refVal := reflect.ValueOf(dst)
	if refVal.Kind() == reflect.Ptr && refVal.Type().Elem().Kind() == reflect.Ptr {
		// If the database returned NULL, then we set dest as nil to indicate that.
		if src == nil {
			nilPtr := reflect.Zero(refVal.Type().Elem())
			refVal.Elem().Set(nilPtr)
			return nil
		}

		// We need to allocate an element, and set the destination to it
		// Then we can retry as that element.
		elemPtr := reflect.New(refVal.Type().Elem().Elem())
		refVal.Elem().Set(elemPtr)

		plan := ci.PlanScan(oid, formatCode, elemPtr.Interface())
		return plan.Scan(ci, oid, formatCode, src, elemPtr.Interface())
	}

	return scanUnknownType(oid, formatCode, src, dst)
}

type scanPlanBinaryInt16 struct{}

func (scanPlanBinaryInt16) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if len(src) != 2 {
		return fmt.Errorf("invalid length for int2: %v", len(src))
	}

	if p, ok := (dst).(*int16); ok {
		*p = int16(binary.BigEndian.Uint16(src))
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanBinaryInt32 struct{}

func (scanPlanBinaryInt32) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for int4: %v", len(src))
	}

	if p, ok := (dst).(*int32); ok {
		*p = int32(binary.BigEndian.Uint32(src))
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanBinaryInt64 struct{}

func (scanPlanBinaryInt64) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for int8: %v", len(src))
	}

	if p, ok := (dst).(*int64); ok {
		*p = int64(binary.BigEndian.Uint64(src))
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanBinaryFloat32 struct{}

func (scanPlanBinaryFloat32) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for int4: %v", len(src))
	}

	if p, ok := (dst).(*float32); ok {
		n := int32(binary.BigEndian.Uint32(src))
		*p = float32(math.Float32frombits(uint32(n)))
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanBinaryFloat64 struct{}

func (scanPlanBinaryFloat64) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for int8: %v", len(src))
	}

	if p, ok := (dst).(*float64); ok {
		n := int64(binary.BigEndian.Uint64(src))
		*p = float64(math.Float64frombits(uint64(n)))
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanBinaryBytes struct{}

func (scanPlanBinaryBytes) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if p, ok := (dst).(*[]byte); ok {
		*p = src
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

type scanPlanString struct{}

func (scanPlanString) Scan(ci *ConnInfo, oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("cannot scan null into %T", dst)
	}

	if p, ok := (dst).(*string); ok {
		*p = string(src)
		return nil
	}

	newPlan := ci.PlanScan(oid, formatCode, dst)
	return newPlan.Scan(ci, oid, formatCode, src, dst)
}

var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func isScanner(dst interface{}) bool {
	if _, ok := dst.(sql.Scanner); ok {
		return true
	}
	if t := reflect.TypeOf(dst); t != nil && t.Kind() == reflect.Ptr && t.Elem().Implements(scannerType) {
		return true
	}
	return false
}

// PlanScan prepares a plan to scan a value into dst.
func (ci *ConnInfo) PlanScan(oid uint32, formatCode int16, dst interface{}) ScanPlan {
	switch formatCode {
	case BinaryFormatCode:
		switch dst.(type) {
		case *string:
			switch oid {
			case TextOID, VarcharOID:
				return scanPlanString{}
			}
		case *int16:
			if oid == Int2OID {
				return scanPlanBinaryInt16{}
			}
		case *int32:
			if oid == Int4OID {
				return scanPlanBinaryInt32{}
			}
		case *int64:
			if oid == Int8OID {
				return scanPlanBinaryInt64{}
			}
		case *float32:
			if oid == Float4OID {
				return scanPlanBinaryFloat32{}
			}
		case *float64:
			if oid == Float8OID {
				return scanPlanBinaryFloat64{}
			}
		case *[]byte:
			switch oid {
			case ByteaOID, TextOID, VarcharOID, JSONOID:
				return scanPlanBinaryBytes{}
			}
		case BinaryDecoder:
			return scanPlanDstBinaryDecoder{}
		}
	case TextFormatCode:
		switch dst.(type) {
		case *string:
			return scanPlanString{}
		case *[]byte:
			if oid != ByteaOID {
				return scanPlanBinaryBytes{}
			}
		case TextDecoder:
			return scanPlanDstTextDecoder{}
		}
	}

	var dt *DataType

	if oid == 0 {
		if dataType, ok := ci.DataTypeForValue(dst); ok {
			dt = dataType
		}
	} else {
		if dataType, ok := ci.DataTypeForOID(oid); ok {
			dt = dataType
		}
	}

	if dt != nil {
		if isScanner(dst) {
			return (*scanPlanDataTypeSQLScanner)(dt)
		}
		return (*scanPlanDataTypeAssignTo)(dt)
	}

	if isScanner(dst) {
		return scanPlanSQLScanner{}
	}

	return scanPlanReflection{}
}

func (ci *ConnInfo) Scan(oid uint32, formatCode int16, src []byte, dst interface{}) error {
	if dst == nil {
		return nil
	}

	plan := ci.PlanScan(oid, formatCode, dst)
	return plan.Scan(ci, oid, formatCode, src, dst)
}

func scanUnknownType(oid uint32, formatCode int16, buf []byte, dest interface{}) error {
	switch dest := dest.(type) {
	case *string:
		if formatCode == BinaryFormatCode {
			return fmt.Errorf("unknown oid %d in binary format cannot be scanned into %T", oid, dest)
		}
		*dest = string(buf)
		return nil
	case *[]byte:
		*dest = buf
		return nil
	default:
		if nextDst, retry := GetAssignToDstType(dest); retry {
			return scanUnknownType(oid, formatCode, buf, nextDst)
		}
		return fmt.Errorf("unknown oid %d cannot be scanned into %T", oid, dest)
	}
}

// NewValue returns a new instance of the same type as v.
func NewValue(v Value) Value {
	if tv, ok := v.(TypeValue); ok {
		return tv.NewTypeValue()
	} else {
		return reflect.New(reflect.ValueOf(v).Elem().Type()).Interface().(Value)
	}
}

var nameValues map[string]Value

func init() {
	nameValues = map[string]Value{
		"_aclitem":       &ACLItemArray{},
		"_bool":          &BoolArray{},
		"_bpchar":        &BPCharArray{},
		"_bytea":         &ByteaArray{},
		"_cidr":          &CIDRArray{},
		"_date":          &DateArray{},
		"_float4":        &Float4Array{},
		"_float8":        &Float8Array{},
		"_inet":          &InetArray{},
		"_int2":          &Int2Array{},
		"_int4":          &Int4Array{},
		"_int8":          &Int8Array{},
		"_numeric":       &NumericArray{},
		"_text":          &TextArray{},
		"_timestamp":     &TimestampArray{},
		"_timestamptz":   &TimestamptzArray{},
		"_uuid":          &UUIDArray{},
		"_varchar":       &VarcharArray{},
		"_json":          &JSONArray{},
		"_jsonb":         &JSONBArray{},
		"aclitem":        &ACLItem{},
		"bit":            &Bit{},
		"bool":           &Bool{},
		"box":            &Box{},
		"bpchar":         &BPChar{},
		"bytea":          &Bytea{},
		"char":           &QChar{},
		"cid":            &CID{},
		"cidr":           &CIDR{},
		"circle":         &Circle{},
		"date":           &Date{},
		"daterange":      &Daterange{},
		"float4":         &Float4{},
		"float8":         &Float8{},
		"hstore":         &Hstore{},
		"inet":           &Inet{},
		"int2":           &Int2{},
		"int4":           &Int4{},
		"int4range":      &Int4range{},
		"int4multirange": &Int4multirange{},
		"int8":           &Int8{},
		"int8range":      &Int8range{},
		"int8multirange": &Int8multirange{},
		"interval":       &Interval{},
		"json":           &JSON{},
		"jsonb":          &JSONB{},
		"line":           &Line{},
		"lseg":           &Lseg{},
		"ltree":          &Ltree{},
		"macaddr":        &Macaddr{},
		"name":           &Name{},
		"numeric":        &Numeric{},
		"numrange":       &Numrange{},
		"nummultirange":  &Nummultirange{},
		"oid":            &OIDValue{},
		"path":           &Path{},
		"point":          &Point{},
		"polygon":        &Polygon{},
		"record":         &Record{},
		"text":           &Text{},
		"tid":            &TID{},
		"timestamp":      &Timestamp{},
		"timestamptz":    &Timestamptz{},
		"tsrange":        &Tsrange{},
		"_tsrange":       &TsrangeArray{},
		"tstzrange":      &Tstzrange{},
		"_tstzrange":     &TstzrangeArray{},
		"unknown":        &Unknown{},
		"uuid":           &UUID{},
		"varbit":         &Varbit{},
		"varchar":        &Varchar{},
		"xid":            &XID{},
	}
}
