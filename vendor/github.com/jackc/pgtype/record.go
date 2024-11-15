package pgtype

import (
	"fmt"
	"reflect"
)

// Record is the generic PostgreSQL record type such as is created with the
// "row" function. Record only implements BinaryDecoder and Value. The text
// format output format from PostgreSQL does not include type information and is
// therefore impossible to decode. No encoders are implemented because
// PostgreSQL does not support input of generic records.
type Record struct {
	Fields []Value
	Status Status
}

func (dst *Record) Set(src interface{}) error {
	if src == nil {
		*dst = Record{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case []Value:
		*dst = Record{Fields: value, Status: Present}
	default:
		return fmt.Errorf("cannot convert %v to Record", src)
	}

	return nil
}

func (dst Record) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Fields
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Record) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *[]Value:
			*v = make([]Value, len(src.Fields))
			copy(*v, src.Fields)
			return nil
		case *[]interface{}:
			*v = make([]interface{}, len(src.Fields))
			for i := range *v {
				(*v)[i] = src.Fields[i].Get()
			}
			return nil
		default:
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

func prepareNewBinaryDecoder(ci *ConnInfo, fieldOID uint32, v *Value) (BinaryDecoder, error) {
	var binaryDecoder BinaryDecoder

	if dt, ok := ci.DataTypeForOID(fieldOID); ok {
		binaryDecoder, _ = dt.Value.(BinaryDecoder)
	} else {
		return nil, fmt.Errorf("unknown oid while decoding record: %v", fieldOID)
	}

	if binaryDecoder == nil {
		return nil, fmt.Errorf("no binary decoder registered for: %v", fieldOID)
	}

	// Duplicate struct to scan into
	binaryDecoder = reflect.New(reflect.ValueOf(binaryDecoder).Elem().Type()).Interface().(BinaryDecoder)
	*v = binaryDecoder.(Value)
	return binaryDecoder, nil
}

func (dst *Record) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Record{Status: Null}
		return nil
	}

	scanner := NewCompositeBinaryScanner(ci, src)

	fields := make([]Value, scanner.FieldCount())

	for i := 0; scanner.Next(); i++ {
		binaryDecoder, err := prepareNewBinaryDecoder(ci, scanner.OID(), &fields[i])
		if err != nil {
			return err
		}

		if err = binaryDecoder.DecodeBinary(ci, scanner.Bytes()); err != nil {
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	*dst = Record{Fields: fields, Status: Present}

	return nil
}
