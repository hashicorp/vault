package pgproto3

type NoticeResponse ErrorResponse

func (*NoticeResponse) Backend() {}

func (dst *NoticeResponse) Decode(src []byte) error {
	return (*ErrorResponse)(dst).Decode(src)
}

func (src *NoticeResponse) Encode(dst []byte) []byte {
	return append(dst, (*ErrorResponse)(src).marshalBinary('N')...)
}
