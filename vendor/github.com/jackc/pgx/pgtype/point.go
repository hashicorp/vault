package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

type Vec2 struct {
	X float64
	Y float64
}

type Point struct {
	P      Vec2
	Status Status
}

func (dst *Point) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Point", src)
}

func (dst *Point) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Point) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Point) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Point{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return errors.Errorf("invalid length for point: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 2)
	if len(parts) < 2 {
		return errors.Errorf("invalid format for point")
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	*dst = Point{P: Vec2{x, y}, Status: Present}
	return nil
}

func (dst *Point) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Point{Status: Null}
		return nil
	}

	if len(src) != 16 {
		return errors.Errorf("invalid length for point: %v", len(src))
	}

	x := binary.BigEndian.Uint64(src)
	y := binary.BigEndian.Uint64(src[8:])

	*dst = Point{
		P:      Vec2{math.Float64frombits(x), math.Float64frombits(y)},
		Status: Present,
	}
	return nil
}

func (src *Point) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, fmt.Sprintf(`(%s,%s)`,
		strconv.FormatFloat(src.P.X, 'f', -1, 64),
		strconv.FormatFloat(src.P.Y, 'f', -1, 64),
	)...), nil
}

func (src *Point) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.P.X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P.Y))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Point) Scan(src interface{}) error {
	if src == nil {
		*dst = Point{Status: Null}
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

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Point) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
