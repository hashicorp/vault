package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type CopyBothResponse struct {
	OverallFormat     byte
	ColumnFormatCodes []uint16
}

func (*CopyBothResponse) Backend() {}

func (dst *CopyBothResponse) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 3 {
		return &invalidMessageFormatErr{messageType: "CopyBothResponse"}
	}

	overallFormat := buf.Next(1)[0]

	columnCount := int(binary.BigEndian.Uint16(buf.Next(2)))
	if buf.Len() != columnCount*2 {
		return &invalidMessageFormatErr{messageType: "CopyBothResponse"}
	}

	columnFormatCodes := make([]uint16, columnCount)
	for i := 0; i < columnCount; i++ {
		columnFormatCodes[i] = binary.BigEndian.Uint16(buf.Next(2))
	}

	*dst = CopyBothResponse{OverallFormat: overallFormat, ColumnFormatCodes: columnFormatCodes}

	return nil
}

func (src *CopyBothResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'W')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendUint16(dst, uint16(len(src.ColumnFormatCodes)))
	for _, fc := range src.ColumnFormatCodes {
		dst = pgio.AppendUint16(dst, fc)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *CopyBothResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type              string
		ColumnFormatCodes []uint16
	}{
		Type:              "CopyBothResponse",
		ColumnFormatCodes: src.ColumnFormatCodes,
	})
}
