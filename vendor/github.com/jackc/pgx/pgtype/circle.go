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

type Circle struct {
	P      Vec2
	R      float64
	Status Status
}

func (dst *Circle) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Circle", src)
}

func (dst *Circle) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Circle) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Circle) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Circle{Status: Null}
		return nil
	}

	if len(src) < 9 {
		return errors.Errorf("invalid length for Circle: %v", len(src))
	}

	str := string(src[2:])
	end := strings.IndexByte(str, ',')
	x, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+1:]
	end = strings.IndexByte(str, ')')

	y, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+2 : len(str)-1]

	r, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*dst = Circle{P: Vec2{x, y}, R: r, Status: Present}
	return nil
}

func (dst *Circle) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Circle{Status: Null}
		return nil
	}

	if len(src) != 24 {
		return errors.Errorf("invalid length for Circle: %v", len(src))
	}

	x := binary.BigEndian.Uint64(src)
	y := binary.BigEndian.Uint64(src[8:])
	r := binary.BigEndian.Uint64(src[16:])

	*dst = Circle{
		P:      Vec2{math.Float64frombits(x), math.Float64frombits(y)},
		R:      math.Float64frombits(r),
		Status: Present,
	}
	return nil
}

func (src *Circle) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`<(%s,%s),%s>`,
		strconv.FormatFloat(src.P.X, 'f', -1, 64),
		strconv.FormatFloat(src.P.Y, 'f', -1, 64),
		strconv.FormatFloat(src.R, 'f', -1, 64),
	)...)

	return buf, nil
}

func (src *Circle) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.P.X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P.Y))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.R))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Circle) Scan(src interface{}) error {
	if src == nil {
		*dst = Circle{Status: Null}
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
func (src *Circle) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
