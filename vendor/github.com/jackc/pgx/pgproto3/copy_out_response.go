package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type CopyOutResponse struct {
	OverallFormat     byte
	ColumnFormatCodes []uint16
}

func (*CopyOutResponse) Backend() {}

func (dst *CopyOutResponse) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 3 {
		return &invalidMessageFormatErr{messageType: "CopyOutResponse"}
	}

	overallFormat := buf.Next(1)[0]

	columnCount := int(binary.BigEndian.Uint16(buf.Next(2)))
	if buf.Len() != columnCount*2 {
		return &invalidMessageFormatErr{messageType: "CopyOutResponse"}
	}

	columnFormatCodes := make([]uint16, columnCount)
	for i := 0; i < columnCount; i++ {
		columnFormatCodes[i] = binary.BigEndian.Uint16(buf.Next(2))
	}

	*dst = CopyOutResponse{OverallFormat: overallFormat, ColumnFormatCodes: columnFormatCodes}

	return nil
}

func (src *CopyOutResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'H')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendUint16(dst, uint16(len(src.ColumnFormatCodes)))
	for _, fc := range src.ColumnFormatCodes {
		dst = pgio.AppendUint16(dst, fc)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *CopyOutResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type              string
		ColumnFormatCodes []uint16
	}{
		Type:              "CopyOutResponse",
		ColumnFormatCodes: src.ColumnFormatCodes,
	})
}
