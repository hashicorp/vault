package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"strconv"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

// OID (Object Identifier Type) is, according to
// https://www.postgresql.org/docs/current/static/datatype-oid.html, used
// internally by PostgreSQL as a primary key for various system tables. It is
// currently implemented as an unsigned four-byte integer. Its definition can be
// found in src/include/postgres_ext.h in the PostgreSQL sources. Because it is
// so frequently required to be in a NOT NULL condition OID cannot be NULL. To
// allow for NULL OIDs use OIDValue.
type OID uint32

func (dst *OID) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		return errors.Errorf("cannot decode nil into OID")
	}

	n, err := strconv.ParseUint(string(src), 10, 32)
	if err != nil {
		return err
	}

	*dst = OID(n)
	return nil
}

func (dst *OID) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		return errors.Errorf("cannot decode nil into OID")
	}

	if len(src) != 4 {
		return errors.Errorf("invalid length: %v", len(src))
	}

	n := binary.BigEndian.Uint32(src)
	*dst = OID(n)
	return nil
}

func (src OID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return append(buf, strconv.FormatUint(uint64(src), 10)...), nil
}

func (src OID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return pgio.AppendUint32(buf, uint32(src)), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *OID) Scan(src interface{}) error {
	if src == nil {
		return errors.Errorf("cannot scan NULL into %T", src)
	}

	switch src := src.(type) {
	case int64:
		*dst = OID(src)
		return nil
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
func (src OID) Value() (driver.Value, error) {
	return int64(src), nil
}
