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

type Path struct {
	P      []Vec2
	Closed bool
	Status Status
}

func (dst *Path) Set(src interface{}) error {
	return fmt.Errorf("cannot convert %v to Path", src)
}

func (dst Path) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Path) AssignTo(dst interface{}) error {
	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Path) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Path{Status: Null}
		return nil
	}

	if len(src) < 7 {
		return fmt.Errorf("invalid length for Path: %v", len(src))
	}

	closed := src[0] == '('
	points := make([]Vec2, 0)

	str := string(src[2:])

	for {
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

		points = append(points, Vec2{x, y})

		if end+3 < len(str) {
			str = str[end+3:]
		} else {
			break
		}
	}

	*dst = Path{P: points, Closed: closed, Status: Present}
	return nil
}

func (dst *Path) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Path{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return fmt.Errorf("invalid length for Path: %v", len(src))
	}

	closed := src[0] == 1
	pointCount := int(binary.BigEndian.Uint32(src[1:]))

	rp := 5

	if 5+pointCount*16 != len(src) {
		return fmt.Errorf("invalid length for Path with %d points: %v", pointCount, len(src))
	}

	points := make([]Vec2, pointCount)
	for i := 0; i < len(points); i++ {
		x := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		y := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		points[i] = Vec2{math.Float64frombits(x), math.Float64frombits(y)}
	}

	*dst = Path{
		P:      points,
		Closed: closed,
		Status: Present,
	}
	return nil
}

func (src Path) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var startByte, endByte byte
	if src.Closed {
		startByte = '('
		endByte = ')'
	} else {
		startByte = '['
		endByte = ']'
	}
	buf = append(buf, startByte)

	for i, p := range src.P {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, fmt.Sprintf(`(%s,%s)`,
			strconv.FormatFloat(p.X, 'f', -1, 64),
			strconv.FormatFloat(p.Y, 'f', -1, 64),
		)...)
	}

	return append(buf, endByte), nil
}

func (src Path) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var closeByte byte
	if src.Closed {
		closeByte = 1
	}
	buf = append(buf, closeByte)

	buf = pgio.AppendInt32(buf, int32(len(src.P)))

	for _, p := range src.P {
		buf = pgio.AppendUint64(buf, math.Float64bits(p.X))
		buf = pgio.AppendUint64(buf, math.Float64bits(p.Y))
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Path) Scan(src interface{}) error {
	if src == nil {
		*dst = Path{Status: Null}
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
func (src Path) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
