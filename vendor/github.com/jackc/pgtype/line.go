package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jackc/pgio"
)

type Line struct {
	A, B, C float64
	Status  Status
}

func (dst *Line) Set(src interface{}) error {
	return fmt.Errorf("cannot convert %v to Line", src)
}

func (dst Line) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Line) AssignTo(dst interface{}) error {
	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Line) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	if len(src) < 7 {
		return fmt.Errorf("invalid length for Line: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid format for line")
	}

	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}

	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	c, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return err
	}

	*dst = Line{A: a, B: b, C: c, Status: Present}
	return nil
}

func (dst *Line) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	if len(src) != 24 {
		return fmt.Errorf("invalid length for Line: %v", len(src))
	}

	a := binary.BigEndian.Uint64(src)
	b := binary.BigEndian.Uint64(src[8:])
	c := binary.BigEndian.Uint64(src[16:])

	*dst = Line{
		A:      math.Float64frombits(a),
		B:      math.Float64frombits(b),
		C:      math.Float64frombits(c),
		Status: Present,
	}
	return nil
}

func (src Line) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`{%s,%s,%s}`,
		strconv.FormatFloat(src.A, 'f', -1, 64),
		strconv.FormatFloat(src.B, 'f', -1, 64),
		strconv.FormatFloat(src.C, 'f', -1, 64),
	)...)

	return buf, nil
}

func (src Line) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.A))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.B))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.C))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Line) Scan(src interface{}) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Line) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
