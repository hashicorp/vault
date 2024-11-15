package pgx

import (
	"database/sql/driver"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
)

// PostgreSQL format codes
const (
	TextFormatCode   = 0
	BinaryFormatCode = 1
)

// SerializationError occurs on failure to encode or decode a value
type SerializationError string

func (e SerializationError) Error() string {
	return string(e)
}

func convertSimpleArgument(ci *pgtype.ConnInfo, arg interface{}) (interface{}, error) {
	if arg == nil {
		return nil, nil
	}

	refVal := reflect.ValueOf(arg)
	if refVal.Kind() == reflect.Ptr && refVal.IsNil() {
		return nil, nil
	}

	switch arg := arg.(type) {

	// https://github.com/jackc/pgx/issues/409 Changed JSON and JSONB to surface
	// []byte to database/sql instead of string. But that caused problems with the
	// simple protocol because the driver.Valuer case got taken before the
	// pgtype.TextEncoder case. And driver.Valuer needed to be first in the usual
	// case because of https://github.com/jackc/pgx/issues/339. So instead we
	// special case JSON and JSONB.
	case *pgtype.JSON:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil
	case *pgtype.JSONB:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil

	case driver.Valuer:
		return callValuerValue(arg)
	case pgtype.TextEncoder:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil
	case float32:
		return float64(arg), nil
	case float64:
		return arg, nil
	case bool:
		return arg, nil
	case time.Duration:
		return fmt.Sprintf("%d microsecond", int64(arg)/1000), nil
	case time.Time:
		return arg, nil
	case string:
		return arg, nil
	case []byte:
		return arg, nil
	case int8:
		return int64(arg), nil
	case int16:
		return int64(arg), nil
	case int32:
		return int64(arg), nil
	case int64:
		return arg, nil
	case int:
		return int64(arg), nil
	case uint8:
		return int64(arg), nil
	case uint16:
		return int64(arg), nil
	case uint32:
		return int64(arg), nil
	case uint64:
		if arg > math.MaxInt64 {
			return nil, fmt.Errorf("arg too big for int64: %v", arg)
		}
		return int64(arg), nil
	case uint:
		if uint64(arg) > math.MaxInt64 {
			return nil, fmt.Errorf("arg too big for int64: %v", arg)
		}
		return int64(arg), nil
	}

	if dt, found := ci.DataTypeForValue(arg); found {
		v := dt.Value
		err := v.Set(arg)
		if err != nil {
			return nil, err
		}
		buf, err := v.(pgtype.TextEncoder).EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil
	}

	if refVal.Kind() == reflect.Ptr {
		arg = refVal.Elem().Interface()
		return convertSimpleArgument(ci, arg)
	}

	if strippedArg, ok := stripNamedType(&refVal); ok {
		return convertSimpleArgument(ci, strippedArg)
	}
	return nil, SerializationError(fmt.Sprintf("Cannot encode %T in simple protocol - %T must implement driver.Valuer, pgtype.TextEncoder, or be a native type", arg, arg))
}

func encodePreparedStatementArgument(ci *pgtype.ConnInfo, buf []byte, oid uint32, arg interface{}) ([]byte, error) {
	if arg == nil {
		return pgio.AppendInt32(buf, -1), nil
	}

	switch arg := arg.(type) {
	case pgtype.BinaryEncoder:
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := arg.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	case pgtype.TextEncoder:
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := arg.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	case string:
		buf = pgio.AppendInt32(buf, int32(len(arg)))
		buf = append(buf, arg...)
		return buf, nil
	}

	refVal := reflect.ValueOf(arg)

	if refVal.Kind() == reflect.Ptr {
		if refVal.IsNil() {
			return pgio.AppendInt32(buf, -1), nil
		}
		arg = refVal.Elem().Interface()
		return encodePreparedStatementArgument(ci, buf, oid, arg)
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
					return encodePreparedStatementArgument(ci, buf, oid, v)
				}
			}

			return nil, err
		}

		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := value.(pgtype.BinaryEncoder).EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	}

	if strippedArg, ok := stripNamedType(&refVal); ok {
		return encodePreparedStatementArgument(ci, buf, oid, strippedArg)
	}
	return nil, SerializationError(fmt.Sprintf("Cannot encode %T into oid %v - %T must implement Encoder or be converted to a string", arg, oid, arg))
}

// chooseParameterFormatCode determines the correct format code for an
// argument to a prepared statement. It defaults to TextFormatCode if no
// determination can be made.
func chooseParameterFormatCode(ci *pgtype.ConnInfo, oid uint32, arg interface{}) int16 {
	switch arg := arg.(type) {
	case pgtype.ParamFormatPreferrer:
		return arg.PreferredParamFormat()
	case pgtype.BinaryEncoder:
		return BinaryFormatCode
	case string, *string, pgtype.TextEncoder:
		return TextFormatCode
	}

	return ci.ParamFormatCodeForOID(oid)
}

func stripNamedType(val *reflect.Value) (interface{}, bool) {
	switch val.Kind() {
	case reflect.Int:
		convVal := int(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int8:
		convVal := int8(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int16:
		convVal := int16(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int32:
		convVal := int32(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int64:
		convVal := int64(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint:
		convVal := uint(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint8:
		convVal := uint8(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint16:
		convVal := uint16(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint32:
		convVal := uint32(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint64:
		convVal := uint64(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.String:
		convVal := val.String()
		return convVal, reflect.TypeOf(convVal) != val.Type()
	}

	return nil, false
}
