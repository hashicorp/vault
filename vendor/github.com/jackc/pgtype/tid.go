package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgio"
)

// TID is PostgreSQL's Tuple Identifier type.
//
// When one does
//
// 	select ctid, * from some_table;
//
// it is the data type of the ctid hidden system column.
//
// It is currently implemented as a pair unsigned two byte integers.
// Its conversion functions can be found in src/backend/utils/adt/tid.c
// in the PostgreSQL sources.
type TID struct {
	BlockNumber  uint32
	OffsetNumber uint16
	Status       Status
}

func (dst *TID) Set(src interface{}) error {
	return fmt.Errorf("cannot convert %v to TID", src)
}

func (dst TID) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *TID) AssignTo(dst interface{}) error {
	if src.Status == Present {
		switch v := dst.(type) {
		case *string:
			*v = fmt.Sprintf(`(%d,%d)`, src.BlockNumber, src.OffsetNumber)
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return fmt.Errorf("unable to assign to %T", dst)
		}
	}

	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *TID) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TID{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return fmt.Errorf("invalid length for tid: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid format for tid")
	}

	blockNumber, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return err
	}

	offsetNumber, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return err
	}

	*dst = TID{BlockNumber: uint32(blockNumber), OffsetNumber: uint16(offsetNumber), Status: Present}
	return nil
}

func (dst *TID) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TID{Status: Null}
		return nil
	}

	if len(src) != 6 {
		return fmt.Errorf("invalid length for tid: %v", len(src))
	}

	*dst = TID{
		BlockNumber:  binary.BigEndian.Uint32(src),
		OffsetNumber: binary.BigEndian.Uint16(src[4:]),
		Status:       Present,
	}
	return nil
}

func (src TID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`(%d,%d)`, src.BlockNumber, src.OffsetNumber)...)
	return buf, nil
}

func (src TID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint32(buf, src.BlockNumber)
	buf = pgio.AppendUint16(buf, src.OffsetNumber)
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *TID) Scan(src interface{}) error {
	if src == nil {
		*dst = TID{Status: Null}
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
func (src TID) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
