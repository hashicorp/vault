package gosnowflake

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ObjectType Empty marker of an object used in column type ScanType function
type ObjectType struct {
}

var structuredObjectWriterType = reflect.TypeOf((*StructuredObjectWriter)(nil)).Elem()

// StructuredObject is a representation of structured object for reading.
type StructuredObject interface {
	GetString(fieldName string) (string, error)
	GetNullString(fieldName string) (sql.NullString, error)
	GetByte(fieldName string) (byte, error)
	GetNullByte(fieldName string) (sql.NullByte, error)
	GetInt16(fieldName string) (int16, error)
	GetNullInt16(fieldName string) (sql.NullInt16, error)
	GetInt32(fieldName string) (int32, error)
	GetNullInt32(fieldName string) (sql.NullInt32, error)
	GetInt64(fieldName string) (int64, error)
	GetNullInt64(fieldName string) (sql.NullInt64, error)
	GetBigInt(fieldName string) (*big.Int, error)
	GetFloat32(fieldName string) (float32, error)
	GetFloat64(fieldName string) (float64, error)
	GetNullFloat64(fieldName string) (sql.NullFloat64, error)
	GetBigFloat(fieldName string) (*big.Float, error)
	GetBool(fieldName string) (bool, error)
	GetNullBool(fieldName string) (sql.NullBool, error)
	GetBytes(fieldName string) ([]byte, error)
	GetTime(fieldName string) (time.Time, error)
	GetNullTime(fieldName string) (sql.NullTime, error)
	GetStruct(fieldName string, scanner sql.Scanner) (sql.Scanner, error)
	GetRaw(fieldName string) (any, error)
	ScanTo(sc sql.Scanner) error
}

// StructuredObjectWriter is an interface to implement, when binding structured objects.
type StructuredObjectWriter interface {
	Write(sowc StructuredObjectWriterContext) error
}

// StructuredObjectWriterContext is a helper interface to write particular fields of structured object.
type StructuredObjectWriterContext interface {
	WriteString(fieldName string, value string) error
	WriteNullString(fieldName string, value sql.NullString) error
	WriteByt(fieldName string, value byte) error // WriteByte name is prohibited by go vet
	WriteNullByte(fieldName string, value sql.NullByte) error
	WriteInt16(fieldName string, value int16) error
	WriteNullInt16(fieldName string, value sql.NullInt16) error
	WriteInt32(fieldName string, value int32) error
	WriteNullInt32(fieldName string, value sql.NullInt32) error
	WriteInt64(fieldName string, value int64) error
	WriteNullInt64(fieldName string, value sql.NullInt64) error
	WriteFloat32(fieldName string, value float32) error
	WriteFloat64(fieldName string, value float64) error
	WriteNullFloat64(fieldName string, value sql.NullFloat64) error
	WriteBytes(fieldName string, value []byte) error
	WriteBool(fieldName string, value bool) error
	WriteNullBool(fieldName string, value sql.NullBool) error
	WriteTime(fieldName string, value time.Time, tsmode []byte) error
	WriteNullTime(fieldName string, value sql.NullTime, tsmode []byte) error
	WriteStruct(fieldName string, value StructuredObjectWriter) error
	WriteNullableStruct(fieldName string, value StructuredObjectWriter, typ reflect.Type) error
	// WriteRaw is used for inserting slices and maps only.
	WriteRaw(fieldName string, value any, tsmode ...[]byte) error
	// WriteNullRaw is used for inserting nil slices and maps only.
	WriteNullRaw(fieldName string, typ reflect.Type, tsmode ...[]byte) error
	WriteAll(sow StructuredObjectWriter) error
}

// NilMapTypes is used to define types when binding nil maps.
type NilMapTypes struct {
	Key   reflect.Type
	Value reflect.Type
}

type structuredObjectWriterEntry struct {
	name      string
	typ       string
	nullable  bool
	length    int
	scale     int
	precision int
	fields    []fieldMetadata
}

func (e *structuredObjectWriterEntry) toFieldMetadata() fieldMetadata {
	return fieldMetadata{
		Name:      e.name,
		Type:      e.typ,
		Nullable:  e.nullable,
		Length:    e.length,
		Scale:     e.scale,
		Precision: e.precision,
		Fields:    e.fields,
	}
}

type structuredObjectWriterContext struct {
	values  map[string]any
	entries []structuredObjectWriterEntry
	params  map[string]*string
}

func (sowc *structuredObjectWriterContext) init(params map[string]*string) {
	sowc.values = make(map[string]any)
	sowc.params = params
}

func (sowc *structuredObjectWriterContext) WriteString(fieldName string, value string) error {
	return sowc.writeString(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullString(fieldName string, value sql.NullString) error {
	if value.Valid {
		return sowc.WriteString(fieldName, value.String)
	}
	return sowc.writeString(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) writeString(fieldName string, value any) error {
	return sowc.write(value, structuredObjectWriterEntry{
		name:     fieldName,
		typ:      "text",
		nullable: true,
		length:   134217728,
	})
}

func (sowc *structuredObjectWriterContext) WriteByt(fieldName string, value byte) error {
	return sowc.writeFixed(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullByte(fieldName string, value sql.NullByte) error {
	if value.Valid {
		return sowc.writeFixed(fieldName, value.Byte)
	}
	return sowc.writeFixed(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) WriteInt16(fieldName string, value int16) error {
	return sowc.writeFixed(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullInt16(fieldName string, value sql.NullInt16) error {
	if value.Valid {
		return sowc.writeFixed(fieldName, value.Int16)
	}
	return sowc.writeFixed(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) WriteInt32(fieldName string, value int32) error {
	return sowc.writeFixed(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullInt32(fieldName string, value sql.NullInt32) error {
	if value.Valid {
		return sowc.writeFixed(fieldName, value.Int32)
	}
	return sowc.writeFixed(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) WriteInt64(fieldName string, value int64) error {
	return sowc.writeFixed(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullInt64(fieldName string, value sql.NullInt64) error {
	if value.Valid {
		return sowc.writeFixed(fieldName, value.Int64)
	}
	return sowc.writeFixed(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) WriteFloat32(fieldName string, value float32) error {
	return sowc.writeFloat(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteFloat64(fieldName string, value float64) error {
	return sowc.writeFloat(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullFloat64(fieldName string, value sql.NullFloat64) error {
	if value.Valid {
		return sowc.writeFloat(fieldName, value.Float64)
	}
	return sowc.writeFloat(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) WriteBool(fieldName string, value bool) error {
	return sowc.writeBool(fieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteNullBool(fieldName string, value sql.NullBool) error {
	if value.Valid {
		return sowc.writeBool(fieldName, value.Bool)
	}
	return sowc.writeBool(fieldName, nil)
}

func (sowc *structuredObjectWriterContext) writeBool(fieldName string, value any) error {
	return sowc.write(value, structuredObjectWriterEntry{
		name:     fieldName,
		typ:      "boolean",
		nullable: true,
	})
}

func (sowc *structuredObjectWriterContext) WriteBytes(fieldName string, value []byte) error {
	var res *string
	if value != nil {
		r := hex.EncodeToString(value)
		res = &r
	}
	return sowc.write(res, structuredObjectWriterEntry{
		name:     fieldName,
		typ:      "binary",
		nullable: true,
	})
}

func (sowc *structuredObjectWriterContext) WriteTime(fieldName string, value time.Time, tsmode []byte) error {
	snowflakeType, err := dataTypeMode(tsmode)
	if err != nil {
		return err
	}
	typ := driverTypeToSnowflake[snowflakeType]
	sfFormat, err := dateTimeInputFormatByType(typ, sowc.params)
	if err != nil {
		return err
	}
	goFormat, err := snowflakeFormatToGoFormat(sfFormat)
	if err != nil {
		return err
	}
	return sowc.writeTime(fieldName, value.Format(goFormat), typ)
}

func (sowc *structuredObjectWriterContext) WriteNullTime(fieldName string, value sql.NullTime, tsmode []byte) error {
	if value.Valid {
		return sowc.WriteTime(fieldName, value.Time, tsmode)
	}
	snowflakeType, err := dataTypeMode(tsmode)
	if err != nil {
		return err
	}
	typ := driverTypeToSnowflake[snowflakeType]
	return sowc.writeTime(fieldName, nil, typ)
}

func (sowc *structuredObjectWriterContext) writeTime(fieldName string, value any, typ string) error {
	return sowc.write(value, structuredObjectWriterEntry{
		name:     fieldName,
		typ:      strings.ToLower(typ),
		nullable: true,
		scale:    9,
	})
}

func (sowc *structuredObjectWriterContext) WriteStruct(fieldName string, value StructuredObjectWriter) error {
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("%s is nil, use WriteNullableStruct instead", fieldName)
	}
	childSowc := structuredObjectWriterContext{}
	childSowc.init(sowc.params)
	err := value.Write(&childSowc)
	if err != nil {
		return err
	}
	return sowc.write(childSowc.values, structuredObjectWriterEntry{
		name:     fieldName,
		typ:      "object",
		nullable: true,
		fields:   childSowc.toFields(),
	})
}

func (sowc *structuredObjectWriterContext) WriteNullableStruct(structFieldName string, value StructuredObjectWriter, typ reflect.Type) error {
	if value == nil || reflect.ValueOf(value).IsNil() {
		childSowc, err := buildSowcFromType(sowc.params, typ)
		if err != nil {
			return err
		}
		return sowc.write(nil, structuredObjectWriterEntry{
			name:     structFieldName,
			typ:      "OBJECT",
			nullable: true,
			fields:   childSowc.toFields(),
		})
	}
	return sowc.WriteStruct(structFieldName, value)
}

func (sowc *structuredObjectWriterContext) WriteRaw(fieldName string, value any, dataTypeModes ...[]byte) error {
	dataTypeModeSingle := DataTypeArray
	if len(dataTypeModes) == 1 && dataTypeModes[0] != nil {
		dataTypeModeSingle = dataTypeModes[0]
	}
	tsmode, err := dataTypeMode(dataTypeModeSingle)
	if err != nil {
		return err
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Slice:
		metadata, err := goTypeToFieldMetadata(reflect.TypeOf(value).Elem(), tsmode, sowc.params)
		if err != nil {
			return err
		}
		return sowc.write(value, structuredObjectWriterEntry{
			name:     fieldName,
			typ:      "ARRAY",
			nullable: true,
			fields:   []fieldMetadata{metadata},
		})
	case reflect.Map:
		keyMetadata, err := goTypeToFieldMetadata(reflect.TypeOf(value).Key(), tsmode, sowc.params)
		if err != nil {
			return err
		}
		valueMetadata, err := goTypeToFieldMetadata(reflect.TypeOf(value).Elem(), tsmode, sowc.params)
		if err != nil {
			return err
		}
		return sowc.write(value, structuredObjectWriterEntry{
			name:     fieldName,
			typ:      "MAP",
			nullable: true,
			fields:   []fieldMetadata{keyMetadata, valueMetadata},
		})
	}
	return fmt.Errorf("unsupported raw type: %T", value)
}

func (sowc *structuredObjectWriterContext) WriteNullRaw(fieldName string, typ reflect.Type, dataTypeModes ...[]byte) error {
	dataTypeModeSingle := DataTypeArray
	if len(dataTypeModes) == 1 && dataTypeModes[0] != nil {
		dataTypeModeSingle = dataTypeModes[0]
	}
	tsmode, err := dataTypeMode(dataTypeModeSingle)
	if err != nil {
		return err
	}

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
		metadata, err := goTypeToFieldMetadata(typ, tsmode, sowc.params)
		if err != nil {
			return err
		}
		if err := sowc.write(nil, structuredObjectWriterEntry{
			name:     fieldName,
			typ:      metadata.Type,
			nullable: true,
			fields:   metadata.Fields,
		}); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cannot use %v as nillable field", typ.Kind().String())
}

func buildSowcFromType(params map[string]*string, typ reflect.Type) (*structuredObjectWriterContext, error) {
	childSowc := &structuredObjectWriterContext{}
	childSowc.init(params)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := getSfFieldName(field)
		if field.Type.Kind() == reflect.String {
			if err := childSowc.writeString(fieldName, nil); err != nil {
				return nil, err
			}
		} else if field.Type.Kind() == reflect.Uint8 || field.Type.Kind() == reflect.Int16 || field.Type.Kind() == reflect.Int32 || field.Type.Kind() == reflect.Int64 {
			if err := childSowc.writeFixed(fieldName, nil); err != nil {
				return nil, err
			}
		} else if field.Type.Kind() == reflect.Float32 || field.Type.Kind() == reflect.Float64 {
			if err := childSowc.writeFloat(fieldName, nil); err != nil {
				return nil, err
			}
		} else if field.Type.Kind() == reflect.Bool {
			if err := childSowc.writeBool(fieldName, nil); err != nil {
				return nil, err
			}
		} else if (field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && field.Type.Elem().Kind() == reflect.Uint8 {
			if err := childSowc.WriteBytes(fieldName, nil); err != nil {
				return nil, err
			}
		} else if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer {
			t := field.Type
			if field.Type.Kind() == reflect.Pointer {
				t = field.Type.Elem()
			}
			if t.AssignableTo(reflect.TypeOf(sql.NullString{})) {
				if err := childSowc.WriteNullString(fieldName, sql.NullString{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullByte{})) {
				if err := childSowc.WriteNullByte(fieldName, sql.NullByte{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullInt16{})) {
				if err := childSowc.WriteNullInt16(fieldName, sql.NullInt16{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullInt32{})) {
				if err := childSowc.WriteNullInt32(fieldName, sql.NullInt32{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullInt64{})) {
				if err := childSowc.WriteNullInt64(fieldName, sql.NullInt64{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullFloat64{})) {
				if err := childSowc.WriteNullFloat64(fieldName, sql.NullFloat64{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullBool{})) {
				if err := childSowc.WriteNullBool(fieldName, sql.NullBool{}); err != nil {
					return nil, err
				}
			} else if t.AssignableTo(reflect.TypeOf(sql.NullTime{})) || t.AssignableTo(reflect.TypeOf(time.Time{})) {
				timeSnowflakeType, err := getTimeSnowflakeType(field)
				if err != nil {
					return nil, err
				}
				if timeSnowflakeType == nil {
					return nil, fmt.Errorf("field %v does not have proper sf tag", fieldName)
				}
				if err := childSowc.WriteNullTime(fieldName, sql.NullTime{}, timeSnowflakeType); err != nil {
					return nil, err
				}
			} else if field.Type.AssignableTo(structuredObjectWriterType) {
				if err := childSowc.WriteNullableStruct(fieldName, nil, field.Type); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("field %s has unsupported type", field.Name)
			}
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Map {
			timeSnowflakeType, err := getTimeSnowflakeType(field)
			if err != nil {
				return nil, err
			}
			if err := childSowc.WriteNullRaw(fieldName, field.Type, timeSnowflakeType); err != nil {
				return nil, err
			}
		}
	}
	return childSowc, nil
}

func (sowc *structuredObjectWriterContext) writeFixed(fieldName string, value any) error {
	return sowc.write(value, structuredObjectWriterEntry{
		name:      fieldName,
		typ:       "fixed",
		nullable:  true,
		precision: 38,
		scale:     0,
	})
}

func (sowc *structuredObjectWriterContext) writeFloat(fieldName string, value any) error {
	return sowc.write(value, structuredObjectWriterEntry{
		name:      fieldName,
		typ:       "real",
		nullable:  true,
		precision: 38,
		scale:     0,
	})
}

func (sowc *structuredObjectWriterContext) write(value any, entry structuredObjectWriterEntry) error {
	sowc.values[entry.name] = value
	sowc.entries = append(sowc.entries, entry)
	return nil
}

func (sowc *structuredObjectWriterContext) WriteAll(sow StructuredObjectWriter) error {
	typ := reflect.TypeOf(sow)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(sow))
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if shouldIgnoreField(field) {
			continue
		}
		fieldName := getSfFieldName(field)
		if field.Type.Kind() == reflect.String {
			if err := sowc.WriteString(fieldName, val.Field(i).String()); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Uint8 {
			if err := sowc.WriteByt(fieldName, byte(val.Field(i).Uint())); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Int16 {
			if err := sowc.WriteInt16(fieldName, int16(val.Field(i).Int())); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Int32 {
			if err := sowc.WriteInt32(fieldName, int32(val.Field(i).Int())); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Int64 {
			if err := sowc.WriteInt64(fieldName, val.Field(i).Int()); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Float32 {
			if err := sowc.WriteFloat32(fieldName, float32(val.Field(i).Float())); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Float64 {
			if err := sowc.WriteFloat64(fieldName, val.Field(i).Float()); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Bool {
			if err := sowc.WriteBool(fieldName, val.Field(i).Bool()); err != nil {
				return err
			}
		} else if (field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice) && field.Type.Elem().Kind() == reflect.Uint8 {
			if err := sowc.WriteBytes(fieldName, val.Field(i).Bytes()); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer {
			if v, ok := val.Field(i).Interface().(time.Time); ok {
				timeSnowflakeType, err := getTimeSnowflakeType(typ.Field(i))
				if err != nil {
					return err
				}
				if timeSnowflakeType == nil {
					return fmt.Errorf("field %v does not have a proper sf tag", fieldName)
				}
				if err := sowc.WriteTime(fieldName, v, timeSnowflakeType); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullString); ok {
				if err := sowc.WriteNullString(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullByte); ok {
				if err := sowc.WriteNullByte(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullInt16); ok {
				if err := sowc.WriteNullInt16(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullInt32); ok {
				if err := sowc.WriteNullInt32(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullInt64); ok {
				if err := sowc.WriteNullInt64(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullFloat64); ok {
				if err := sowc.WriteNullFloat64(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullBool); ok {
				if err := sowc.WriteNullBool(fieldName, v); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(sql.NullTime); ok {
				timeSnowflakeType, err := getTimeSnowflakeType(typ.Field(i))
				if err != nil {
					return err
				}
				if timeSnowflakeType == nil {
					return fmt.Errorf("field %v does not have a proper sf tag", fieldName)
				}
				if err := sowc.WriteNullTime(fieldName, v, timeSnowflakeType); err != nil {
					return err
				}
			} else if v, ok := val.Field(i).Interface().(StructuredObjectWriter); ok {
				if reflect.ValueOf(v).IsNil() {
					if err := sowc.WriteNullableStruct(fieldName, nil, reflect.TypeOf(v)); err != nil {
						return err
					}
				} else {
					childSowc := &structuredObjectWriterContext{}
					childSowc.init(sowc.params)
					if err := v.Write(childSowc); err != nil {
						return err
					}
					if err := sowc.write(childSowc.values, structuredObjectWriterEntry{
						name:     fieldName,
						typ:      "OBJECT",
						nullable: true,
						fields:   childSowc.toFields(),
					}); err != nil {
						return err
					}
				}
			}
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Map {
			var timeSfType []byte
			var err error
			if field.Type.Elem().AssignableTo(reflect.TypeOf(time.Time{})) || field.Type.Elem().AssignableTo(reflect.TypeOf(sql.NullTime{})) {
				timeSfType, err = getTimeSnowflakeType(typ.Field(i))
				if err != nil {
					return err
				}
			}
			if err := sowc.WriteRaw(fieldName, val.Field(i).Interface(), timeSfType); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("field %s has unsupported type", field.Name)
		}
	}
	return nil
}

func (sowc *structuredObjectWriterContext) toFields() []fieldMetadata {
	fieldMetadatas := make([]fieldMetadata, len(sowc.entries))
	for i, entry := range sowc.entries {
		fieldMetadatas[i] = entry.toFieldMetadata()
	}
	return fieldMetadatas
}

// ArrayOfScanners Helper type for scanning array of sql.Scanner values.
type ArrayOfScanners[T sql.Scanner] []T

func (st *ArrayOfScanners[T]) Scan(val any) error {
	if val == nil {
		return nil
	}
	sts := val.([]*structuredType)
	*st = make([]T, len(sts))
	var t T
	for i, s := range sts {
		(*st)[i] = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		if err := (*st)[i].Scan(s); err != nil {
			return err
		}
	}
	return nil
}

// ScanArrayOfScanners is a helper function for scanning arrays of sql.Scanner values.
// Example:
//
//	var res []*simpleObject
//	err := rows.Scan(ScanArrayOfScanners(&res))
func ScanArrayOfScanners[T sql.Scanner](value *[]T) *ArrayOfScanners[T] {
	return (*ArrayOfScanners[T])(value)
}

// MapOfScanners Helper type for scanning map of sql.Scanner values.
type MapOfScanners[K comparable, V sql.Scanner] map[K]V

func (st *MapOfScanners[K, V]) Scan(val any) error {
	if val == nil {
		return nil
	}
	sts := val.(map[K]*structuredType)
	*st = make(map[K]V)
	var someV V
	for k, v := range sts {
		if v != nil && !reflect.ValueOf(v).IsNil() {
			(*st)[k] = reflect.New(reflect.TypeOf(someV).Elem()).Interface().(V)
			if err := (*st)[k].Scan(sts[k]); err != nil {
				return err
			}
		} else {
			(*st)[k] = reflect.Zero(reflect.TypeOf(someV)).Interface().(V)
		}
	}
	return nil
}

// ScanMapOfScanners is a helper function for scanning maps of sql.Scanner values.
// Example:
//
//	var res map[string]*simpleObject
//	err := rows.Scan(ScanMapOfScanners(&res))
func ScanMapOfScanners[K comparable, V sql.Scanner](m *map[K]V) *MapOfScanners[K, V] {
	return (*MapOfScanners[K, V])(m)
}

type structuredType struct {
	values        map[string]any
	fieldMetadata []fieldMetadata
	params        map[string]*string
}

func getType[T any](st *structuredType, fieldName string, emptyValue T) (T, bool, error) {
	v, ok := st.values[fieldName]
	if !ok {
		return emptyValue, false, errors.New("field " + fieldName + " does not exist")
	}
	if v == nil {
		return emptyValue, true, nil
	}
	v, ok = v.(T)
	if !ok {
		return emptyValue, false, fmt.Errorf("cannot convert field %v to %T", fieldName, emptyValue)
	}
	return v.(T), false, nil
}

func (st *structuredType) GetString(fieldName string) (string, error) {
	nullString, err := st.GetNullString(fieldName)
	if err != nil {
		return "", err
	}
	if !nullString.Valid {
		return "", fmt.Errorf("nil value for %v, use GetNullString instead", fieldName)
	}
	return nullString.String, nil
}

func (st *structuredType) GetNullString(fieldName string) (sql.NullString, error) {
	s, wasNil, err := getType[string](st, fieldName, "")
	if err != nil {
		return sql.NullString{}, err
	}
	if wasNil {
		return sql.NullString{Valid: false}, err
	}
	return sql.NullString{Valid: true, String: s}, nil
}

func (st *structuredType) GetByte(fieldName string) (byte, error) {
	nullByte, err := st.GetNullByte(fieldName)
	if err != nil {
		return 0, err
	}
	if !nullByte.Valid {
		return 0, fmt.Errorf("nil value for %v, use GetNullByte instead", fieldName)
	}
	return nullByte.Byte, nil
}

func (st *structuredType) GetNullByte(fieldName string) (sql.NullByte, error) {
	b, err := st.GetNullInt64(fieldName)
	if err != nil {
		return sql.NullByte{}, err
	}
	if !b.Valid {
		return sql.NullByte{Valid: false}, nil
	}
	return sql.NullByte{Valid: true, Byte: byte(b.Int64)}, nil
}

func (st *structuredType) GetInt16(fieldName string) (int16, error) {
	nullInt16, err := st.GetNullInt16(fieldName)
	if err != nil {
		return 0, err
	}
	if !nullInt16.Valid {
		return 0, fmt.Errorf("nil value for %v, use GetNullInt16 instead", fieldName)
	}
	return nullInt16.Int16, nil
}

func (st *structuredType) GetNullInt16(fieldName string) (sql.NullInt16, error) {
	b, err := st.GetNullInt64(fieldName)
	if err != nil {
		return sql.NullInt16{}, err
	}
	if !b.Valid {
		return sql.NullInt16{Valid: false}, nil
	}
	return sql.NullInt16{Valid: true, Int16: int16(b.Int64)}, nil
}

func (st *structuredType) GetInt32(fieldName string) (int32, error) {
	nullInt32, err := st.GetNullInt32(fieldName)
	if err != nil {
		return 0, err
	}
	if !nullInt32.Valid {
		return 0, fmt.Errorf("nil value for %v, use GetNullInt32 instead", fieldName)
	}
	return nullInt32.Int32, nil
}

func (st *structuredType) GetNullInt32(fieldName string) (sql.NullInt32, error) {
	b, err := st.GetNullInt64(fieldName)
	if err != nil {
		return sql.NullInt32{}, err
	}
	if !b.Valid {
		return sql.NullInt32{Valid: false}, nil
	}
	return sql.NullInt32{Valid: true, Int32: int32(b.Int64)}, nil
}

func (st *structuredType) GetInt64(fieldName string) (int64, error) {
	nullInt64, err := st.GetNullInt64(fieldName)
	if err != nil {
		return 0, err
	}
	if !nullInt64.Valid {
		return 0, fmt.Errorf("nil value for %v, use GetNullInt64 instead", fieldName)
	}
	return nullInt64.Int64, nil
}

func (st *structuredType) GetNullInt64(fieldName string) (sql.NullInt64, error) {
	i64, wasNil, err := getType[int64](st, fieldName, 0)
	if wasNil {
		return sql.NullInt64{Valid: false}, err
	}
	if err == nil {
		return sql.NullInt64{Valid: true, Int64: i64}, nil
	}
	if s, _, err := getType[string](st, fieldName, ""); err == nil {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return sql.NullInt64{Valid: false}, err
		}
		return sql.NullInt64{Valid: true, Int64: i}, nil
	} else if b, _, err := getType[float64](st, fieldName, 0); err == nil {
		return sql.NullInt64{Valid: true, Int64: int64(b)}, nil
	} else if b, _, err := getType[json.Number](st, fieldName, ""); err == nil {
		i, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return sql.NullInt64{Valid: false}, err
		}
		return sql.NullInt64{Valid: true, Int64: i}, err
	} else {
		return sql.NullInt64{Valid: false}, fmt.Errorf("cannot cast column %v to byte", fieldName)
	}
}

func (st *structuredType) GetBigInt(fieldName string) (*big.Int, error) {
	b, wasNull, err := getType[*big.Int](st, fieldName, new(big.Int))
	if wasNull {
		return nil, nil
	}
	return b, err
}

func (st *structuredType) GetFloat32(fieldName string) (float32, error) {
	f32, err := st.GetFloat64(fieldName)
	if err != nil {
		return 0, err
	}
	return float32(f32), err
}

func (st *structuredType) GetFloat64(fieldName string) (float64, error) {
	nullFloat64, err := st.GetNullFloat64(fieldName)
	if err != nil {
		return 0, nil
	}
	if !nullFloat64.Valid {
		return 0, fmt.Errorf("nil value for %v, use GetNullFloat64 instead", fieldName)
	}
	return nullFloat64.Float64, nil
}

func (st *structuredType) GetNullFloat64(fieldName string) (sql.NullFloat64, error) {
	float64, wasNull, err := getType[float64](st, fieldName, 0)
	if wasNull {
		return sql.NullFloat64{Valid: false}, nil
	}
	if err == nil {
		return sql.NullFloat64{Valid: true, Float64: float64}, nil
	}
	s, _, err := getType[string](st, fieldName, "")
	if err == nil {
		f64, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return sql.NullFloat64{}, err
		}
		return sql.NullFloat64{Valid: true, Float64: f64}, err
	}
	jsonNumber, _, err := getType[json.Number](st, fieldName, "")
	if err != nil {
		return sql.NullFloat64{}, err
	}
	f64, err := strconv.ParseFloat(string(jsonNumber), 64)
	if err != nil {
		return sql.NullFloat64{}, err
	}
	return sql.NullFloat64{Valid: true, Float64: f64}, nil
}

func (st *structuredType) GetBigFloat(fieldName string) (*big.Float, error) {
	float, wasNull, err := getType[*big.Float](st, fieldName, new(big.Float))
	if wasNull {
		return nil, nil
	}
	return float, err
}

func (st *structuredType) GetBool(fieldName string) (bool, error) {
	nullBool, err := st.GetNullBool(fieldName)
	if err != nil {
		return false, err
	}
	if !nullBool.Valid {
		return false, fmt.Errorf("nil value for %v, use GetNullBool instead", fieldName)
	}
	return nullBool.Bool, err
}

func (st *structuredType) GetNullBool(fieldName string) (sql.NullBool, error) {
	b, wasNull, err := getType[bool](st, fieldName, false)
	if wasNull {
		return sql.NullBool{Valid: false}, nil
	}
	if err != nil {
		return sql.NullBool{}, err
	}
	return sql.NullBool{Valid: true, Bool: b}, nil
}

func (st *structuredType) GetBytes(fieldName string) ([]byte, error) {
	if bi, _, err := getType[[]byte](st, fieldName, nil); err == nil {
		return bi, nil
	} else if bi, _, err := getType[string](st, fieldName, ""); err == nil {
		return hex.DecodeString(bi)
	}
	bytes, _, err := getType[[]byte](st, fieldName, []byte{})
	return bytes, err
}

func (st *structuredType) GetTime(fieldName string) (time.Time, error) {
	nullTime, err := st.GetNullTime(fieldName)
	if err != nil {
		return time.Time{}, err
	}
	if !nullTime.Valid {
		return time.Time{}, fmt.Errorf("nil value for %v, use GetNullBool instead", fieldName)
	}
	return nullTime.Time, nil
}

func (st *structuredType) GetNullTime(fieldName string) (sql.NullTime, error) {
	s, wasNull, err := getType[string](st, fieldName, "")
	if wasNull {
		return sql.NullTime{Valid: false}, nil
	}
	if err == nil {
		fieldMetadata, err := st.fieldMetadataByFieldName(fieldName)
		if err != nil {
			return sql.NullTime{}, err
		}
		format, err := dateTimeOutputFormatByType(fieldMetadata.Type, st.params)
		if err != nil {
			return sql.NullTime{}, err
		}
		goFormat, err := snowflakeFormatToGoFormat(format)
		if err != nil {
			return sql.NullTime{}, err
		}
		time, err := time.Parse(goFormat, s)
		return sql.NullTime{Valid: true, Time: time}, err
	}
	time, _, err := getType[time.Time](st, fieldName, time.Time{})
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Valid: true, Time: time}, nil
}

func (st *structuredType) GetStruct(fieldName string, scanner sql.Scanner) (sql.Scanner, error) {
	childSt, wasNull, err := getType[*structuredType](st, fieldName, &structuredType{})
	if wasNull {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = scanner.Scan(childSt)
	return scanner, err
}
func (st *structuredType) GetRaw(fieldName string) (any, error) {
	return st.values[fieldName], nil
}

func (st *structuredType) ScanTo(sc sql.Scanner) error {
	v := reflect.Indirect(reflect.ValueOf(sc))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if shouldIgnoreField(field) {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			s, err := st.GetString(getSfFieldName(field))
			if err != nil {
				return err
			}
			v.FieldByName(field.Name).SetString(s)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := st.GetInt64(getSfFieldName(field))
			if err != nil {
				return err
			}
			v.FieldByName(field.Name).SetInt(i)
		case reflect.Uint8:
			b, err := st.GetByte(getSfFieldName(field))
			if err != nil {
				return err
			}
			v.FieldByName(field.Name).SetUint(uint64(int64(b)))
		case reflect.Float32, reflect.Float64:
			f, err := st.GetFloat64(getSfFieldName(field))
			if err != nil {
				return err
			}
			v.FieldByName(field.Name).SetFloat(f)
		case reflect.Bool:
			b, err := st.GetBool(getSfFieldName(field))
			if err != nil {
				return err
			}
			v.FieldByName(field.Name).SetBool(b)
		case reflect.Slice, reflect.Array:
			switch field.Type.Elem().Kind() {
			case reflect.Uint8:
				b, err := st.GetBytes(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).SetBytes(b)
			default:
				raw, err := st.GetRaw(getSfFieldName(field))
				if err != nil {
					return err
				}
				if raw != nil {
					v.FieldByName(field.Name).Set(reflect.ValueOf(raw))
				}
			}
		case reflect.Map:
			raw, err := st.GetRaw(getSfFieldName(field))
			if err != nil {
				return err
			}
			if raw != nil {
				v.FieldByName(field.Name).Set(reflect.ValueOf(raw))
			}
		case reflect.Struct:
			a := v.FieldByName(field.Name).Interface()
			if _, ok := a.(time.Time); ok {
				time, err := st.GetTime(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(time))
			} else if _, ok := a.(sql.Scanner); ok {
				scanner := reflect.New(reflect.TypeOf(a)).Interface().(sql.Scanner)
				s, err := st.GetStruct(getSfFieldName(field), scanner)
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.Indirect(reflect.ValueOf(s)))
			} else if _, ok := a.(sql.NullString); ok {
				ns, err := st.GetNullString(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(ns))
			} else if _, ok := a.(sql.NullByte); ok {
				nb, err := st.GetNullByte(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(nb))
			} else if _, ok := a.(sql.NullBool); ok {
				nb, err := st.GetNullBool(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(nb))
			} else if _, ok := a.(sql.NullInt16); ok {
				ni, err := st.GetNullInt16(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(ni))
			} else if _, ok := a.(sql.NullInt32); ok {
				ni, err := st.GetNullInt32(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(ni))
			} else if _, ok := a.(sql.NullInt64); ok {
				ni, err := st.GetNullInt64(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(ni))
			} else if _, ok := a.(sql.NullFloat64); ok {
				nf, err := st.GetNullFloat64(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(nf))
			} else if _, ok := a.(sql.NullTime); ok {
				nt, err := st.GetNullTime(getSfFieldName(field))
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(nt))
			}
		case reflect.Pointer:
			switch field.Type.Elem().Kind() {
			case reflect.Struct:
				a := reflect.New(field.Type.Elem()).Interface()
				s, err := st.GetStruct(getSfFieldName(field), a.(sql.Scanner))
				if err != nil {
					return err
				}
				if s != nil {
					v.FieldByName(field.Name).Set(reflect.ValueOf(s))
				}
			default:
				return errors.New("only struct pointers are supported")
			}
		}
	}
	return nil
}

func (st *structuredType) fieldMetadataByFieldName(fieldName string) (fieldMetadata, error) {
	for _, fm := range st.fieldMetadata {
		if fm.Name == fieldName {
			return fm, nil
		}
	}
	return fieldMetadata{}, errors.New("no metadata for field " + fieldName)
}

func structuredTypesEnabled(ctx context.Context) bool {
	v := ctx.Value(enableStructuredTypes)
	if v == nil {
		return false
	}
	d, ok := v.(bool)
	return ok && d
}

func mapValuesNullableEnabled(ctx context.Context) bool {
	v := ctx.Value(mapValuesNullable)
	if v == nil {
		return false
	}
	d, ok := v.(bool)
	return ok && d
}

func arrayValuesNullableEnabled(ctx context.Context) bool {
	v := ctx.Value(arrayValuesNullable)
	if v == nil {
		return false
	}
	d, ok := v.(bool)
	return ok && d
}

func getSfFieldName(field reflect.StructField) string {
	sfTag := field.Tag.Get("sf")
	if sfTag != "" {
		return strings.Split(sfTag, ",")[0]
	}
	r := []rune(field.Name)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func shouldIgnoreField(field reflect.StructField) bool {
	sfTag := strings.ToLower(field.Tag.Get("sf"))
	if sfTag == "" {
		return false
	}
	return contains(strings.Split(sfTag, ",")[1:], "ignore")
}

func getTimeSnowflakeType(field reflect.StructField) ([]byte, error) {
	sfTag := strings.ToLower(field.Tag.Get("sf"))
	if sfTag == "" {
		return nil, nil
	}
	values := strings.Split(sfTag, ",")[1:]
	if contains(values, "time") {
		return DataTypeTime, nil
	} else if contains(values, "date") {
		return DataTypeDate, nil
	} else if contains(values, "ltz") {
		return DataTypeTimestampLtz, nil
	} else if contains(values, "ntz") {
		return DataTypeTimestampNtz, nil
	} else if contains(values, "tz") {
		return DataTypeTimestampTz, nil
	}
	return nil, nil
}
