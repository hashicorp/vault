package pgx

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/jackc/pgtype"
)

type extendedQueryBuilder struct {
	paramValues     [][]byte
	paramValueBytes []byte
	paramFormats    []int16
	resultFormats   []int16
}

func (eqb *extendedQueryBuilder) AppendParam(ci *pgtype.ConnInfo, oid uint32, arg interface{}) error {
	f := chooseParameterFormatCode(ci, oid, arg)
	eqb.paramFormats = append(eqb.paramFormats, f)

	v, err := eqb.encodeExtendedParamValue(ci, oid, f, arg)
	if err != nil {
		return err
	}
	eqb.paramValues = append(eqb.paramValues, v)

	return nil
}

func (eqb *extendedQueryBuilder) AppendResultFormat(f int16) {
	eqb.resultFormats = append(eqb.resultFormats, f)
}

// Reset readies eqb to build another query.
func (eqb *extendedQueryBuilder) Reset() {
	eqb.paramValues = eqb.paramValues[0:0]
	eqb.paramValueBytes = eqb.paramValueBytes[0:0]
	eqb.paramFormats = eqb.paramFormats[0:0]
	eqb.resultFormats = eqb.resultFormats[0:0]

	if cap(eqb.paramValues) > 64 {
		eqb.paramValues = make([][]byte, 0, 64)
	}

	if cap(eqb.paramValueBytes) > 256 {
		eqb.paramValueBytes = make([]byte, 0, 256)
	}

	if cap(eqb.paramFormats) > 64 {
		eqb.paramFormats = make([]int16, 0, 64)
	}
	if cap(eqb.resultFormats) > 64 {
		eqb.resultFormats = make([]int16, 0, 64)
	}
}

func (eqb *extendedQueryBuilder) encodeExtendedParamValue(ci *pgtype.ConnInfo, oid uint32, formatCode int16, arg interface{}) ([]byte, error) {
	if arg == nil {
		return nil, nil
	}

	refVal := reflect.ValueOf(arg)
	argIsPtr := refVal.Kind() == reflect.Ptr

	if argIsPtr && refVal.IsNil() {
		return nil, nil
	}

	if eqb.paramValueBytes == nil {
		eqb.paramValueBytes = make([]byte, 0, 128)
	}

	var err error
	var buf []byte
	pos := len(eqb.paramValueBytes)

	if arg, ok := arg.(string); ok {
		return []byte(arg), nil
	}

	if formatCode == TextFormatCode {
		if arg, ok := arg.(pgtype.TextEncoder); ok {
			buf, err = arg.EncodeText(ci, eqb.paramValueBytes)
			if err != nil {
				return nil, err
			}
			if buf == nil {
				return nil, nil
			}
			eqb.paramValueBytes = buf
			return eqb.paramValueBytes[pos:], nil
		}
	} else if formatCode == BinaryFormatCode {
		if arg, ok := arg.(pgtype.BinaryEncoder); ok {
			buf, err = arg.EncodeBinary(ci, eqb.paramValueBytes)
			if err != nil {
				return nil, err
			}
			if buf == nil {
				return nil, nil
			}
			eqb.paramValueBytes = buf
			return eqb.paramValueBytes[pos:], nil
		}
	}

	if argIsPtr {
		// We have already checked that arg is not pointing to nil,
		// so it is safe to dereference here.
		arg = refVal.Elem().Interface()
		return eqb.encodeExtendedParamValue(ci, oid, formatCode, arg)
	}

	if dt, ok := ci.DataTypeForOID(oid); ok {
		value := dt.Value
		err := value.Set(arg)
		if err != nil {
			{
				if arg, ok := arg.(driver.Valuer); ok {
					v, err := callValuerValue(arg)
					if err != nil {
						return nil, err
					}
					return eqb.encodeExtendedParamValue(ci, oid, formatCode, v)
				}
			}

			return nil, err
		}

		return eqb.encodeExtendedParamValue(ci, oid, formatCode, value)
	}

	// There is no data type registered for the destination OID, but maybe there is data type registered for the arg
	// type. If so use it's text encoder (if available).
	if dt, ok := ci.DataTypeForValue(arg); ok {
		value := dt.Value
		if textEncoder, ok := value.(pgtype.TextEncoder); ok {
			err := value.Set(arg)
			if err != nil {
				return nil, err
			}

			buf, err = textEncoder.EncodeText(ci, eqb.paramValueBytes)
			if err != nil {
				return nil, err
			}
			if buf == nil {
				return nil, nil
			}
			eqb.paramValueBytes = buf
			return eqb.paramValueBytes[pos:], nil
		}
	}

	if strippedArg, ok := stripNamedType(&refVal); ok {
		return eqb.encodeExtendedParamValue(ci, oid, formatCode, strippedArg)
	}
	return nil, SerializationError(fmt.Sprintf("Cannot encode %T into oid %v - %T must implement Encoder or be converted to a string", arg, oid, arg))
}
