package pgtype

type Decimal Numeric

func (dst *Decimal) Set(src interface{}) error {
	return (*Numeric)(dst).Set(src)
}

func (dst *Decimal) Get() interface{} {
	return (*Numeric)(dst).Get()
}

func (src *Decimal) AssignTo(dst interface{}) error {
	return (*Numeric)(src).AssignTo(dst)
}

func (dst *Decimal) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeText(ci, src)
}

func (dst *Decimal) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeBinary(ci, src)
}

func (src *Decimal) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeText(ci, buf)
}

func (src *Decimal) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeBinary(ci, buf)
}
