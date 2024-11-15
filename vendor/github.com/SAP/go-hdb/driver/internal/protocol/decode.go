package protocol

import (
	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

func decodeResult(tc typeCode, d *encoding.Decoder, lobReader LobReader, lobChunkSize, scale int) (any, error) { //nolint: gocyclo
	switch tc {
	case tcBoolean:
		return d.BooleanField()
	case tcTinyint:
		if !d.Bool() { // null value
			return nil, nil
		}
		return int64(d.Byte()), nil
	case tcSmallint:
		if !d.Bool() { // null value
			return nil, nil
		}
		return int64(d.Int16()), nil
	case tcInteger:
		if !d.Bool() { // null value
			return nil, nil
		}
		return int64(d.Int32()), nil
	case tcBigint:
		if !d.Bool() { // null value
			return nil, nil
		}
		return d.Int64(), nil
	case tcReal:
		return d.RealField()
	case tcDouble:
		return d.DoubleField()
	case tcDate:
		return d.DateField()
	case tcTime:
		return d.TimeField()
	case tcTimestamp:
		return d.TimestampField()
	case tcLongdate:
		return d.LongdateField()
	case tcSeconddate:
		return d.SeconddateField()
	case tcDaydate:
		return d.DaydateField()
	case tcSecondtime:
		return d.SecondtimeField()
	case tcDecimal:
		return d.DecimalField()
	case tcFixed8:
		return d.Fixed8Field(scale)
	case tcFixed12:
		return d.Fixed12Field(scale)
	case tcFixed16:
		return d.Fixed16Field(scale)
	case tcChar, tcVarchar, tcString, tcBinary, tcVarbinary:
		return d.VarField()
	case tcAlphanum:
		return d.AlphanumField()
	case tcNchar, tcNvarchar, tcNstring, tcShorttext:
		return d.Cesu8Field()
	case tcStPoint, tcStGeometry:
		return d.HexField()
	case tcBlob, tcClob, tcLocator, tcBintext:
		descr := newLobOutDescr(nil, lobReader, lobChunkSize)
		if descr.decode(d) {
			return nil, nil
		}
		return descr, nil
	case tcText, tcNclob, tcNlocator:
		descr := newLobOutDescr(d.Transformer(), lobReader, lobChunkSize)
		if descr.decode(d) {
			return nil, nil
		}
		return descr, nil
	default:
		panic("invalid type code")
	}
}

func decodeLobParameter(d *encoding.Decoder) (any, error) {
	// real decoding (sniffer) not yet supported
	// descr := &LobInDescr{}
	// descr.Opt = LobOptions(d.Byte())
	// descr._size = int(d.Int32())
	// descr.pos = int(d.Int32())
	d.Byte()
	d.Int32()
	d.Int32()
	return nil, nil
}

func decodeParameter(tc typeCode, d *encoding.Decoder, scale int) (any, error) {
	switch tc {
	case tcBoolean:
		return d.BooleanField()
	case tcTinyint:
		return int64(d.Byte()), nil
	case tcSmallint:
		return int64(d.Int16()), nil
	case tcInteger:
		return int64(d.Int32()), nil
	case tcBigint:
		return d.Int64(), nil
	case tcReal:
		return d.RealField()
	case tcDouble:
		return d.DoubleField()
	case tcDate:
		return d.DateField()
	case tcTime:
		return d.TimeField()
	case tcTimestamp:
		return d.TimestampField()
	case tcLongdate:
		return d.LongdateField()
	case tcSeconddate:
		return d.SeconddateField()
	case tcDaydate:
		return d.DaydateField()
	case tcSecondtime:
		return d.SecondtimeField()
	case tcDecimal:
		return d.DecimalField()
	case tcFixed8:
		return d.Fixed8Field(scale)
	case tcFixed12:
		return d.Fixed12Field(scale)
	case tcFixed16:
		return d.Fixed16Field(scale)
	case tcChar, tcVarchar, tcString, tcBinary, tcVarbinary:
		return d.VarField()
	case tcAlphanum:
		return d.AlphanumField()
	case tcNchar, tcNvarchar, tcNstring, tcShorttext:
		return d.Cesu8Field()
	case tcStPoint, tcStGeometry:
		return d.HexField()
	case tcBlob, tcClob, tcLocator, tcBintext:
		return decodeLobParameter(d)
	case tcText, tcNclob, tcNlocator:
		return decodeLobParameter(d)
	default:
		panic("invalid type code")
	}
}
