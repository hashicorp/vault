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

type Box struct {
	P      [2]Vec2
	Status Status
}

func (dst *Box) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Box", src)
}

func (dst *Box) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Box) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Box) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Box{Status: Null}
		return nil
	}

	if len(src) < 11 {
		return errors.Errorf("invalid length for Box: %v", len(src))
	}

	str := string(src[1:])

	var end int
	end = strings.IndexByte(str, ',')

	x1, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+1:]
	end = strings.IndexByte(str, ')')

	y1, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+3:]
	end = strings.IndexByte(str, ',')

	x2, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+1 : len(str)-1]

	y2, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*dst = Box{P: [2]Vec2{{x1, y1}, {x2, y2}}, Status: Present}
	return nil
}

func (dst *Box) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Box{Status: Null}
		return nil
	}

	if len(src) != 32 {
		return errors.Errorf("invalid length for Box: %v", len(src))
	}

	x1 := binary.BigEndian.Uint64(src)
	y1 := binary.BigEndian.Uint64(src[8:])
	x2 := binary.BigEndian.Uint64(src[16:])
	y2 := binary.BigEndian.Uint64(src[24:])

	*dst = Box{
		P: [2]Vec2{
			{math.Float64frombits(x1), math.Float64frombits(y1)},
			{math.Float64frombits(x2), math.Float64frombits(y2)},
		},
		Status: Present,
	}
	return nil
}

func (src *Box) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`(%s,%s),(%s,%s)`,
		strconv.FormatFloat(src.P[0].X, 'f', -1, 64),
		strconv.FormatFloat(src.P[0].Y, 'f', -1, 64),
		strconv.FormatFloat(src.P[1].X, 'f', -1, 64),
		strconv.FormatFloat(src.P[1].Y, 'f', -1, 64),
	)...)
	return buf, nil
}

func (src *Box) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[0].X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[0].Y))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[1].X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[1].Y))

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Box) Scan(src interface{}) error {
	if src == nil {
		*dst = Box{Status: Null}
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
func (src *Box) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
