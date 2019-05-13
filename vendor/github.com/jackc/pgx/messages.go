package pgx

import (
	"database/sql/driver"
	"math"
	"reflect"
	"time"

	"github.com/jackc/pgx/pgio"
	"github.com/jackc/pgx/pgtype"
)

const (
	copyData      = 'd'
	copyFail      = 'f'
	copyDone      = 'c'
	varHeaderSize = 4
)

type FieldDescription struct {
	Name            string
	Table           pgtype.OID
	AttributeNumber uint16
	DataType        pgtype.OID
	DataTypeSize    int16
	DataTypeName    string
	Modifier        int32
	FormatCode      int16
}

func (fd FieldDescription) Length() (int64, bool) {
	switch fd.DataType {
	case pgtype.TextOID, pgtype.ByteaOID:
		return math.MaxInt64, true
	case pgtype.VarcharOID, pgtype.BPCharArrayOID:
		return int64(fd.Modifier - varHeaderSize), true
	default:
		return 0, false
	}
}

func (fd FieldDescription) PrecisionScale() (precision, scale int64, ok bool) {
	switch fd.DataType {
	case pgtype.NumericOID:
		mod := fd.Modifier - varHeaderSize
		precision = int64((mod >> 16) & 0xffff)
		scale = int64(mod & 0xffff)
		return precision, scale, true
	default:
		return 0, 0, false
	}
}

func (fd FieldDescription) Type() reflect.Type {
	switch fd.DataType {
	case pgtype.Float8OID:
		return reflect.TypeOf(float64(0))
	case pgtype.Float4OID:
		return reflect.TypeOf(float32(0))
	case pgtype.Int8OID:
		return reflect.TypeOf(int64(0))
	case pgtype.Int4OID:
		return reflect.TypeOf(int32(0))
	case pgtype.Int2OID:
		return reflect.TypeOf(int16(0))
	case pgtype.VarcharOID, pgtype.BPCharArrayOID, pgtype.TextOID:
		return reflect.TypeOf("")
	case pgtype.BoolOID:
		return reflect.TypeOf(false)
	case pgtype.NumericOID:
		return reflect.TypeOf(float64(0))
	case pgtype.DateOID, pgtype.TimestampOID, pgtype.TimestamptzOID:
		return reflect.TypeOf(time.Time{})
	case pgtype.ByteaOID:
		return reflect.TypeOf([]byte(nil))
	default:
		return reflect.TypeOf(new(interface{})).Elem()
	}
}

// PgError represents an error reported by the PostgreSQL server. See
// http://www.postgresql.org/docs/9.3/static/protocol-error-fields.html for
// detailed field description.
type PgError struct {
	Severity         string
	Code             string
	Message          string
	Detail           string
	Hint             string
	Position         int32
	InternalPosition int32
	InternalQuery    string
	Where            string
	SchemaName       string
	TableName        string
	ColumnName       string
	DataTypeName     string
	ConstraintName   string
	File             string
	Line             int32
	Routine          string
}

func (pe PgError) Error() string {
	return pe.Severity + ": " + pe.Message + " (SQLSTATE " + pe.Code + ")"
}

// Notice represents a notice response message reported by the PostgreSQL
// server. Be aware that this is distinct from LISTEN/NOTIFY notification.
type Notice PgError

// appendParse appends a PostgreSQL wire protocol parse message to buf and returns it.
func appendParse(buf []byte, name string, query string, parameterOIDs []pgtype.OID) []byte {
	buf = append(buf, 'P')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, name...)
	buf = append(buf, 0)
	buf = append(buf, query...)
	buf = append(buf, 0)

	buf = pgio.AppendInt16(buf, int16(len(parameterOIDs)))
	for _, oid := range parameterOIDs {
		buf = pgio.AppendUint32(buf, uint32(oid))
	}
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendDescribe appends a PostgreSQL wire protocol describe message to buf and returns it.
func appendDescribe(buf []byte, objectType byte, name string) []byte {
	buf = append(buf, 'D')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, objectType)
	buf = append(buf, name...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendSync appends a PostgreSQL wire protocol sync message to buf and returns it.
func appendSync(buf []byte) []byte {
	buf = append(buf, 'S')
	buf = pgio.AppendInt32(buf, 4)

	return buf
}

// appendBind appends a PostgreSQL wire protocol bind message to buf and returns it.
func appendBind(
	buf []byte,
	destinationPortal,
	preparedStatement string,
	connInfo *pgtype.ConnInfo,
	parameterOIDs []pgtype.OID,
	arguments []interface{},
	resultFormatCodes []int16,
) ([]byte, error) {
	buf = append(buf, 'B')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, destinationPortal...)
	buf = append(buf, 0)
	buf = append(buf, preparedStatement...)
	buf = append(buf, 0)

	var err error
	arguments, err = convertDriverValuers(arguments)
	if err != nil {
		return nil, err
	}

	buf = pgio.AppendInt16(buf, int16(len(parameterOIDs)))
	for i, oid := range parameterOIDs {
		buf = pgio.AppendInt16(buf, chooseParameterFormatCode(connInfo, oid, arguments[i]))
	}

	buf = pgio.AppendInt16(buf, int16(len(arguments)))
	for i, oid := range parameterOIDs {
		var err error
		buf, err = encodePreparedStatementArgument(connInfo, buf, oid, arguments[i])
		if err != nil {
			return nil, err
		}
	}

	buf = pgio.AppendInt16(buf, int16(len(resultFormatCodes)))
	for _, fc := range resultFormatCodes {
		buf = pgio.AppendInt16(buf, fc)
	}
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf, nil
}

func convertDriverValuers(args []interface{}) ([]interface{}, error) {
	for i, arg := range args {
		switch arg := arg.(type) {
		case pgtype.BinaryEncoder:
		case pgtype.TextEncoder:
		case driver.Valuer:
			v, err := callValuerValue(arg)
			if err != nil {
				return nil, err
			}
			args[i] = v
		}
	}
	return args, nil
}

// appendExecute appends a PostgreSQL wire protocol execute message to buf and returns it.
func appendExecute(buf []byte, portal string, maxRows uint32) []byte {
	buf = append(buf, 'E')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)

	buf = append(buf, portal...)
	buf = append(buf, 0)
	buf = pgio.AppendUint32(buf, maxRows)

	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendQuery appends a PostgreSQL wire protocol query message to buf and returns it.
func appendQuery(buf []byte, query string) []byte {
	buf = append(buf, 'Q')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, query...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}
