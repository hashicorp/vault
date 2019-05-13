package pgtype

type CIDR Inet

func (dst *CIDR) Set(src interface{}) error {
	return (*Inet)(dst).Set(src)
}

func (dst *CIDR) Get() interface{} {
	return (*Inet)(dst).Get()
}

func (src *CIDR) AssignTo(dst interface{}) error {
	return (*Inet)(src).AssignTo(dst)
}

func (dst *CIDR) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Inet)(dst).DecodeText(ci, src)
}

func (dst *CIDR) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Inet)(dst).DecodeBinary(ci, src)
}

func (src *CIDR) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Inet)(src).EncodeText(ci, buf)
}

func (src *CIDR) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Inet)(src).EncodeBinary(ci, buf)
}
