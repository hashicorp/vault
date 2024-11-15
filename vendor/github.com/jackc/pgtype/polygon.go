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

type Polygon struct {
	P      []Vec2
	Status Status
}

// Set converts src to dest.
//
// src can be nil, string, []float64, and []pgtype.Vec2.
//
// If src is string the format must be ((x1,y1),(x2,y2),...,(xn,yn)).
// Important that there are no spaces in it.
func (dst *Polygon) Set(src interface{}) error {
	if src == nil {
		dst.Status = Null
		return nil
	}
	err := fmt.Errorf("cannot convert %v to Polygon", src)
	var p *Polygon
	switch value := src.(type) {
	case string:
		p, err = stringToPolygon(value)
	case []Vec2:
		p = &Polygon{Status: Present, P: value}
		err = nil
	case []float64:
		p, err = float64ToPolygon(value)
	default:
		return err
	}
	if err != nil {
		return err
	}
	*dst = *p
	return nil
}

func stringToPolygon(src string) (*Polygon, error) {
	p := &Polygon{}
	err := p.DecodeText(nil, []byte(src))
	return p, err
}

func float64ToPolygon(src []float64) (*Polygon, error) {
	p := &Polygon{Status: Null}
	if len(src) == 0 {
		return p, nil
	}
	if len(src)%2 != 0 {
		p.Status = Undefined
		return p, fmt.Errorf("invalid length for polygon: %v", len(src))
	}
	p.Status = Present
	p.P = make([]Vec2, 0)
	for i := 0; i < len(src); i += 2 {
		p.P = append(p.P, Vec2{X: src[i], Y: src[i+1]})
	}
	return p, nil
}

func (dst Polygon) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Polygon) AssignTo(dst interface{}) error {
	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Polygon) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Polygon{Status: Null}
		return nil
	}

	if len(src) < 7 {
		return fmt.Errorf("invalid length for Polygon: %v", len(src))
	}

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

	*dst = Polygon{P: points, Status: Present}
	return nil
}

func (dst *Polygon) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Polygon{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return fmt.Errorf("invalid length for Polygon: %v", len(src))
	}

	pointCount := int(binary.BigEndian.Uint32(src))
	rp := 4

	if 4+pointCount*16 != len(src) {
		return fmt.Errorf("invalid length for Polygon with %d points: %v", pointCount, len(src))
	}

	points := make([]Vec2, pointCount)
	for i := 0; i < len(points); i++ {
		x := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		y := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		points[i] = Vec2{math.Float64frombits(x), math.Float64frombits(y)}
	}

	*dst = Polygon{
		P:      points,
		Status: Present,
	}
	return nil
}

func (src Polygon) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, '(')

	for i, p := range src.P {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, fmt.Sprintf(`(%s,%s)`,
			strconv.FormatFloat(p.X, 'f', -1, 64),
			strconv.FormatFloat(p.Y, 'f', -1, 64),
		)...)
	}

	return append(buf, ')'), nil
}

func (src Polygon) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendInt32(buf, int32(len(src.P)))

	for _, p := range src.P {
		buf = pgio.AppendUint64(buf, math.Float64bits(p.X))
		buf = pgio.AppendUint64(buf, math.Float64bits(p.Y))
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Polygon) Scan(src interface{}) error {
	if src == nil {
		*dst = Polygon{Status: Null}
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
func (src Polygon) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
