package protocol

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

type columnOptions int8

const (
	coMandatory columnOptions = 0x01
	coOptional  columnOptions = 0x02
)

const (
	coMandatoryText = "mandatory"
	coOptionalText  = "optional"
)

func (k columnOptions) String() string {
	var s []string
	if k&coMandatory != 0 {
		s = append(s, coMandatoryText)
	}
	if k&coOptional != 0 {
		s = append(s, coOptionalText)
	}
	return fmt.Sprintf("%v", s)
}

// ResultsetID represents a resultset id.
type ResultsetID uint64

func (id ResultsetID) String() string { return fmt.Sprintf("%d", id) }
func (id *ResultsetID) decode(dec *encoding.Decoder) error {
	*id = ResultsetID(dec.Uint64())
	return dec.Error()
}
func (id ResultsetID) encode(enc *encoding.Encoder) error { enc.Uint64(uint64(id)); return nil }

func newResultFields(size int) []*ResultField {
	return make([]*ResultField, size)
}

// ResultField represents a database result field.
type ResultField struct {
	names                *fieldNames
	tableNameOfs         uint32
	schemaNameOfs        uint32
	columnNameOfs        uint32
	columnDisplayNameOfs uint32
	prec                 int // length
	scale                int // fraction
	columnOptions        columnOptions
	tc                   typeCode
}

// String implements the Stringer interface.
func (f *ResultField) String() string {
	return fmt.Sprintf("columnsOptions %s typeCode %s precision %d scale %d tablename %s schemaname %s columnname %s columnDisplayname %s",
		f.columnOptions,
		f.tc,
		f.prec,
		f.scale,
		f.names.name(f.tableNameOfs),
		f.names.name(f.schemaNameOfs),
		f.names.name(f.columnNameOfs),
		f.names.name(f.columnDisplayNameOfs),
	)
}

func (f *ResultField) isNullable() bool { return f.columnOptions == coOptional }

// DatabaseTypeName returns the type name of the field.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) DatabaseTypeName() string { return f.tc.typeName() }

// DecimalSize returns the type precision and scale of the field.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) DecimalSize() (int64, int64, bool) {
	if f.tc.isDecimalType() {
		return int64(f.prec), int64(f.scale), true
	}
	return 0, 0, false
}

// Length returns the type length of the field.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) Length() (int64, bool) {
	if f.tc.isVariableLength() {
		return int64(f.prec), true
	}
	return 0, false
}

// Name returns the result field name.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) Name() string { return f.names.name(f.columnDisplayNameOfs) }

// Nullable returns true if the field may be null, false otherwise.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) Nullable() (bool, bool) { return f.isNullable(), true }

// ScanType returns the scan type of the field.
// It implements the go-hdb driver ColumnType interface.
func (f *ResultField) ScanType() reflect.Type { return f.tc.dataType().ScanType(f.isNullable()) }

func (f *ResultField) decode(dec *encoding.Decoder) {
	f.columnOptions = columnOptions(dec.Int8())
	f.tc = typeCode(dec.Int8())
	f.scale = int(dec.Int16())
	f.prec = int(dec.Int16())
	dec.Skip(2) // filler
	f.tableNameOfs = dec.Uint32()
	f.schemaNameOfs = dec.Uint32()
	f.columnNameOfs = dec.Uint32()
	f.columnDisplayNameOfs = dec.Uint32()

	f.names.insert(f.tableNameOfs)
	f.names.insert(f.schemaNameOfs)
	f.names.insert(f.columnNameOfs)
	f.names.insert(f.columnDisplayNameOfs)
}

func (f *ResultField) decodeResult(dec *encoding.Decoder, lobReader LobReader, lobChunkSize int) (any, error) {
	return decodeResult(f.tc, dec, lobReader, lobChunkSize, f.scale)
}

// ResultMetadata represents the metadata of a set of database result fields.
type ResultMetadata struct {
	ResultFields []*ResultField
}

func (r *ResultMetadata) String() string {
	return fmt.Sprintf("result fields %v", r.ResultFields)
}

func (r *ResultMetadata) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	r.ResultFields = newResultFields(numArg)
	names := &fieldNames{}
	for i := range len(r.ResultFields) {
		f := &ResultField{names: names}
		f.decode(dec)
		r.ResultFields[i] = f
	}
	if err := names.decode(dec); err != nil {
		return err
	}
	return dec.Error()
}

// Resultset represents a database result set.
type Resultset struct {
	ResultFields []*ResultField
	FieldValues  []driver.Value
	DecodeErrors DecodeErrors
}

func (r *Resultset) String() string {
	return fmt.Sprintf("result fields %v field values %v", r.ResultFields, r.FieldValues)
}

func (r *Resultset) decodeResult(dec *encoding.Decoder, numArg int, lobReader LobReader, lobChunkSize int) error {
	cols := len(r.ResultFields)
	r.FieldValues = resizeSlice(r.FieldValues, numArg*cols)

	for i := range numArg {
		for j, f := range r.ResultFields {
			var err error
			if r.FieldValues[i*cols+j], err = f.decodeResult(dec, lobReader, lobChunkSize); err != nil {
				r.DecodeErrors = append(r.DecodeErrors, &DecodeError{row: i, fieldName: f.Name(), err: err}) // collect decode / conversion errors
			}
		}
	}
	return dec.Error()
}
