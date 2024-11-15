package pgtype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgio"
)

type CompositeTypeField struct {
	Name string
	OID  uint32
}

type CompositeType struct {
	status Status

	typeName string

	fields           []CompositeTypeField
	valueTranscoders []ValueTranscoder
}

// NewCompositeType creates a CompositeType from fields and ci. ci is used to find the ValueTranscoders used
// for fields. All field OIDs must be previously registered in ci.
func NewCompositeType(typeName string, fields []CompositeTypeField, ci *ConnInfo) (*CompositeType, error) {
	valueTranscoders := make([]ValueTranscoder, len(fields))

	for i := range fields {
		dt, ok := ci.DataTypeForOID(fields[i].OID)
		if !ok {
			return nil, fmt.Errorf("no data type registered for oid: %d", fields[i].OID)
		}

		value := NewValue(dt.Value)
		valueTranscoder, ok := value.(ValueTranscoder)
		if !ok {
			return nil, fmt.Errorf("data type for oid does not implement ValueTranscoder: %d", fields[i].OID)
		}

		valueTranscoders[i] = valueTranscoder
	}

	return &CompositeType{typeName: typeName, fields: fields, valueTranscoders: valueTranscoders}, nil
}

// NewCompositeTypeValues creates a CompositeType from fields and values. fields and values must have the same length.
// Prefer NewCompositeType unless overriding the transcoding of fields is required.
func NewCompositeTypeValues(typeName string, fields []CompositeTypeField, values []ValueTranscoder) (*CompositeType, error) {
	if len(fields) != len(values) {
		return nil, errors.New("fields and valueTranscoders must have same length")
	}

	return &CompositeType{typeName: typeName, fields: fields, valueTranscoders: values}, nil
}

func (src CompositeType) Get() interface{} {
	switch src.status {
	case Present:
		results := make(map[string]interface{}, len(src.valueTranscoders))
		for i := range src.valueTranscoders {
			results[src.fields[i].Name] = src.valueTranscoders[i].Get()
		}
		return results
	case Null:
		return nil
	default:
		return src.status
	}
}

func (ct *CompositeType) NewTypeValue() Value {
	a := &CompositeType{
		typeName:         ct.typeName,
		fields:           ct.fields,
		valueTranscoders: make([]ValueTranscoder, len(ct.valueTranscoders)),
	}

	for i := range ct.valueTranscoders {
		a.valueTranscoders[i] = NewValue(ct.valueTranscoders[i]).(ValueTranscoder)
	}

	return a
}

func (ct *CompositeType) TypeName() string {
	return ct.typeName
}

func (ct *CompositeType) Fields() []CompositeTypeField {
	return ct.fields
}

func (dst *CompositeType) Set(src interface{}) error {
	if src == nil {
		dst.status = Null
		return nil
	}

	switch value := src.(type) {
	case []interface{}:
		if len(value) != len(dst.valueTranscoders) {
			return fmt.Errorf("Number of fields don't match. CompositeType has %d fields", len(dst.valueTranscoders))
		}
		for i, v := range value {
			if err := dst.valueTranscoders[i].Set(v); err != nil {
				return err
			}
		}
		dst.status = Present
	case *[]interface{}:
		if value == nil {
			dst.status = Null
			return nil
		}
		return dst.Set(*value)
	default:
		return fmt.Errorf("Can not convert %v to Composite", src)
	}

	return nil
}

// AssignTo should never be called on composite value directly
func (src CompositeType) AssignTo(dst interface{}) error {
	switch src.status {
	case Present:
		switch v := dst.(type) {
		case []interface{}:
			if len(v) != len(src.valueTranscoders) {
				return fmt.Errorf("Number of fields don't match. CompositeType has %d fields", len(src.valueTranscoders))
			}
			for i := range src.valueTranscoders {
				if v[i] == nil {
					continue
				}

				err := assignToOrSet(src.valueTranscoders[i], v[i])
				if err != nil {
					return fmt.Errorf("unable to assign to dst[%d]: %v", i, err)
				}
			}
			return nil
		case *[]interface{}:
			return src.AssignTo(*v)
		default:
			if isPtrStruct, err := src.assignToPtrStruct(dst); isPtrStruct {
				return err
			}

			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return fmt.Errorf("unable to assign to %T", dst)
		}
	case Null:
		return NullAssignTo(dst)
	}
	return fmt.Errorf("cannot decode %#v into %T", src, dst)
}

func assignToOrSet(src Value, dst interface{}) error {
	assignToErr := src.AssignTo(dst)
	if assignToErr != nil {
		// Try to use get / set instead -- this avoids every type having to be able to AssignTo type of self.
		setSucceeded := false
		if setter, ok := dst.(Value); ok {
			err := setter.Set(src.Get())
			setSucceeded = err == nil
		}
		if !setSucceeded {
			return assignToErr
		}
	}

	return nil
}

func (src CompositeType) assignToPtrStruct(dst interface{}) (bool, error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return false, nil
	}

	if dstValue.IsNil() {
		return false, nil
	}

	dstElemValue := dstValue.Elem()
	dstElemType := dstElemValue.Type()

	if dstElemType.Kind() != reflect.Struct {
		return false, nil
	}

	exportedFields := make([]int, 0, dstElemType.NumField())
	for i := 0; i < dstElemType.NumField(); i++ {
		sf := dstElemType.Field(i)
		if sf.PkgPath == "" {
			exportedFields = append(exportedFields, i)
		}
	}

	if len(exportedFields) != len(src.valueTranscoders) {
		return false, nil
	}

	for i := range exportedFields {
		err := assignToOrSet(src.valueTranscoders[i], dstElemValue.Field(exportedFields[i]).Addr().Interface())
		if err != nil {
			return true, fmt.Errorf("unable to assign to field %s: %v", dstElemType.Field(exportedFields[i]).Name, err)
		}
	}

	return true, nil
}

func (src CompositeType) EncodeBinary(ci *ConnInfo, buf []byte) (newBuf []byte, err error) {
	switch src.status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	b := NewCompositeBinaryBuilder(ci, buf)
	for i := range src.valueTranscoders {
		b.AppendEncoder(src.fields[i].OID, src.valueTranscoders[i])
	}

	return b.Finish()
}

// DecodeBinary implements BinaryDecoder interface.
// Opposite to Record, fields in a composite act as a "schema"
// and decoding fails if SQL value can't be assigned due to
// type mismatch
func (dst *CompositeType) DecodeBinary(ci *ConnInfo, buf []byte) error {
	if buf == nil {
		dst.status = Null
		return nil
	}

	scanner := NewCompositeBinaryScanner(ci, buf)

	for _, f := range dst.valueTranscoders {
		scanner.ScanDecoder(f)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	dst.status = Present

	return nil
}

func (dst *CompositeType) DecodeText(ci *ConnInfo, buf []byte) error {
	if buf == nil {
		dst.status = Null
		return nil
	}

	scanner := NewCompositeTextScanner(ci, buf)

	for _, f := range dst.valueTranscoders {
		scanner.ScanDecoder(f)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	dst.status = Present

	return nil
}

func (src CompositeType) EncodeText(ci *ConnInfo, buf []byte) (newBuf []byte, err error) {
	switch src.status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	b := NewCompositeTextBuilder(ci, buf)
	for _, f := range src.valueTranscoders {
		b.AppendEncoder(f)
	}

	return b.Finish()
}

type CompositeBinaryScanner struct {
	ci  *ConnInfo
	rp  int
	src []byte

	fieldCount int32
	fieldBytes []byte
	fieldOID   uint32
	err        error
}

// NewCompositeBinaryScanner a scanner over a binary encoded composite balue.
func NewCompositeBinaryScanner(ci *ConnInfo, src []byte) *CompositeBinaryScanner {
	rp := 0
	if len(src[rp:]) < 4 {
		return &CompositeBinaryScanner{err: fmt.Errorf("Record incomplete %v", src)}
	}

	fieldCount := int32(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	return &CompositeBinaryScanner{
		ci:         ci,
		rp:         rp,
		src:        src,
		fieldCount: fieldCount,
	}
}

// ScanDecoder calls Next and decodes the result with d.
func (cfs *CompositeBinaryScanner) ScanDecoder(d BinaryDecoder) {
	if cfs.err != nil {
		return
	}

	if cfs.Next() {
		cfs.err = d.DecodeBinary(cfs.ci, cfs.fieldBytes)
	} else {
		cfs.err = errors.New("read past end of composite")
	}
}

// ScanDecoder calls Next and scans the result into d.
func (cfs *CompositeBinaryScanner) ScanValue(d interface{}) {
	if cfs.err != nil {
		return
	}

	if cfs.Next() {
		cfs.err = cfs.ci.Scan(cfs.OID(), BinaryFormatCode, cfs.Bytes(), d)
	} else {
		cfs.err = errors.New("read past end of composite")
	}
}

// Next advances the scanner to the next field. It returns false after the last field is read or an error occurs. After
// Next returns false, the Err method can be called to check if any errors occurred.
func (cfs *CompositeBinaryScanner) Next() bool {
	if cfs.err != nil {
		return false
	}

	if cfs.rp == len(cfs.src) {
		return false
	}

	if len(cfs.src[cfs.rp:]) < 8 {
		cfs.err = fmt.Errorf("Record incomplete %v", cfs.src)
		return false
	}
	cfs.fieldOID = binary.BigEndian.Uint32(cfs.src[cfs.rp:])
	cfs.rp += 4

	fieldLen := int(int32(binary.BigEndian.Uint32(cfs.src[cfs.rp:])))
	cfs.rp += 4

	if fieldLen >= 0 {
		if len(cfs.src[cfs.rp:]) < fieldLen {
			cfs.err = fmt.Errorf("Record incomplete rp=%d src=%v", cfs.rp, cfs.src)
			return false
		}
		cfs.fieldBytes = cfs.src[cfs.rp : cfs.rp+fieldLen]
		cfs.rp += fieldLen
	} else {
		cfs.fieldBytes = nil
	}

	return true
}

func (cfs *CompositeBinaryScanner) FieldCount() int {
	return int(cfs.fieldCount)
}

// Bytes returns the bytes of the field most recently read by Scan().
func (cfs *CompositeBinaryScanner) Bytes() []byte {
	return cfs.fieldBytes
}

// OID returns the OID of the field most recently read by Scan().
func (cfs *CompositeBinaryScanner) OID() uint32 {
	return cfs.fieldOID
}

// Err returns any error encountered by the scanner.
func (cfs *CompositeBinaryScanner) Err() error {
	return cfs.err
}

type CompositeTextScanner struct {
	ci  *ConnInfo
	rp  int
	src []byte

	fieldBytes []byte
	err        error
}

// NewCompositeTextScanner a scanner over a text encoded composite value.
func NewCompositeTextScanner(ci *ConnInfo, src []byte) *CompositeTextScanner {
	if len(src) < 2 {
		return &CompositeTextScanner{err: fmt.Errorf("Record incomplete %v", src)}
	}

	if src[0] != '(' {
		return &CompositeTextScanner{err: fmt.Errorf("composite text format must start with '('")}
	}

	if src[len(src)-1] != ')' {
		return &CompositeTextScanner{err: fmt.Errorf("composite text format must end with ')'")}
	}

	return &CompositeTextScanner{
		ci:  ci,
		rp:  1,
		src: src,
	}
}

// ScanDecoder calls Next and decodes the result with d.
func (cfs *CompositeTextScanner) ScanDecoder(d TextDecoder) {
	if cfs.err != nil {
		return
	}

	if cfs.Next() {
		cfs.err = d.DecodeText(cfs.ci, cfs.fieldBytes)
	} else {
		cfs.err = errors.New("read past end of composite")
	}
}

// ScanDecoder calls Next and scans the result into d.
func (cfs *CompositeTextScanner) ScanValue(d interface{}) {
	if cfs.err != nil {
		return
	}

	if cfs.Next() {
		cfs.err = cfs.ci.Scan(0, TextFormatCode, cfs.Bytes(), d)
	} else {
		cfs.err = errors.New("read past end of composite")
	}
}

// Next advances the scanner to the next field. It returns false after the last field is read or an error occurs. After
// Next returns false, the Err method can be called to check if any errors occurred.
func (cfs *CompositeTextScanner) Next() bool {
	if cfs.err != nil {
		return false
	}

	if cfs.rp == len(cfs.src) {
		return false
	}

	switch cfs.src[cfs.rp] {
	case ',', ')': // null
		cfs.rp++
		cfs.fieldBytes = nil
		return true
	case '"': // quoted value
		cfs.rp++
		cfs.fieldBytes = make([]byte, 0, 16)
		for {
			ch := cfs.src[cfs.rp]

			if ch == '"' {
				cfs.rp++
				if cfs.src[cfs.rp] == '"' {
					cfs.fieldBytes = append(cfs.fieldBytes, '"')
					cfs.rp++
				} else {
					break
				}
			} else if ch == '\\' {
				cfs.rp++
				cfs.fieldBytes = append(cfs.fieldBytes, cfs.src[cfs.rp])
				cfs.rp++
			} else {
				cfs.fieldBytes = append(cfs.fieldBytes, ch)
				cfs.rp++
			}
		}
		cfs.rp++
		return true
	default: // unquoted value
		start := cfs.rp
		for {
			ch := cfs.src[cfs.rp]
			if ch == ',' || ch == ')' {
				break
			}
			cfs.rp++
		}
		cfs.fieldBytes = cfs.src[start:cfs.rp]
		cfs.rp++
		return true
	}
}

// Bytes returns the bytes of the field most recently read by Scan().
func (cfs *CompositeTextScanner) Bytes() []byte {
	return cfs.fieldBytes
}

// Err returns any error encountered by the scanner.
func (cfs *CompositeTextScanner) Err() error {
	return cfs.err
}

type CompositeBinaryBuilder struct {
	ci         *ConnInfo
	buf        []byte
	startIdx   int
	fieldCount uint32
	err        error
}

func NewCompositeBinaryBuilder(ci *ConnInfo, buf []byte) *CompositeBinaryBuilder {
	startIdx := len(buf)
	buf = append(buf, 0, 0, 0, 0) // allocate room for number of fields
	return &CompositeBinaryBuilder{ci: ci, buf: buf, startIdx: startIdx}
}

func (b *CompositeBinaryBuilder) AppendValue(oid uint32, field interface{}) {
	if b.err != nil {
		return
	}

	dt, ok := b.ci.DataTypeForOID(oid)
	if !ok {
		b.err = fmt.Errorf("unknown data type for OID: %d", oid)
		return
	}

	err := dt.Value.Set(field)
	if err != nil {
		b.err = err
		return
	}

	binaryEncoder, ok := dt.Value.(BinaryEncoder)
	if !ok {
		b.err = fmt.Errorf("unable to encode binary for OID: %d", oid)
		return
	}

	b.AppendEncoder(oid, binaryEncoder)
}

func (b *CompositeBinaryBuilder) AppendEncoder(oid uint32, field BinaryEncoder) {
	if b.err != nil {
		return
	}

	b.buf = pgio.AppendUint32(b.buf, oid)
	lengthPos := len(b.buf)
	b.buf = pgio.AppendInt32(b.buf, -1)
	fieldBuf, err := field.EncodeBinary(b.ci, b.buf)
	if err != nil {
		b.err = err
		return
	}
	if fieldBuf != nil {
		binary.BigEndian.PutUint32(fieldBuf[lengthPos:], uint32(len(fieldBuf)-len(b.buf)))
		b.buf = fieldBuf
	}

	b.fieldCount++
}

func (b *CompositeBinaryBuilder) Finish() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}

	binary.BigEndian.PutUint32(b.buf[b.startIdx:], b.fieldCount)
	return b.buf, nil
}

type CompositeTextBuilder struct {
	ci       *ConnInfo
	buf      []byte
	err      error
	fieldBuf [32]byte
}

func NewCompositeTextBuilder(ci *ConnInfo, buf []byte) *CompositeTextBuilder {
	buf = append(buf, '(') // allocate room for number of fields
	return &CompositeTextBuilder{ci: ci, buf: buf}
}

func (b *CompositeTextBuilder) AppendValue(field interface{}) {
	if b.err != nil {
		return
	}

	if field == nil {
		b.buf = append(b.buf, ',')
		return
	}

	dt, ok := b.ci.DataTypeForValue(field)
	if !ok {
		b.err = fmt.Errorf("unknown data type for field: %v", field)
		return
	}

	err := dt.Value.Set(field)
	if err != nil {
		b.err = err
		return
	}

	textEncoder, ok := dt.Value.(TextEncoder)
	if !ok {
		b.err = fmt.Errorf("unable to encode text for value: %v", field)
		return
	}

	b.AppendEncoder(textEncoder)
}

func (b *CompositeTextBuilder) AppendEncoder(field TextEncoder) {
	if b.err != nil {
		return
	}

	fieldBuf, err := field.EncodeText(b.ci, b.fieldBuf[0:0])
	if err != nil {
		b.err = err
		return
	}
	if fieldBuf != nil {
		b.buf = append(b.buf, quoteCompositeFieldIfNeeded(string(fieldBuf))...)
	}

	b.buf = append(b.buf, ',')
}

func (b *CompositeTextBuilder) Finish() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.buf[len(b.buf)-1] = ')'
	return b.buf, nil
}

var quoteCompositeReplacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func quoteCompositeField(src string) string {
	return `"` + quoteCompositeReplacer.Replace(src) + `"`
}

func quoteCompositeFieldIfNeeded(src string) string {
	if src == "" || src[0] == ' ' || src[len(src)-1] == ' ' || strings.ContainsAny(src, `(),"\`) {
		return quoteCompositeField(src)
	}
	return src
}
