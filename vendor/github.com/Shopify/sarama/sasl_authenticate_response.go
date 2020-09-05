package sarama

type SaslAuthenticateResponse struct {
	Err           KError
	ErrorMessage  *string
	SaslAuthBytes []byte
}

func (r *SaslAuthenticateResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	if err := pe.putNullableString(r.ErrorMessage); err != nil {
		return err
	}
	return pe.putBytes(r.SaslAuthBytes)
}

func (r *SaslAuthenticateResponse) decode(pd packetDecoder, version int16) error {
	kerr, err := pd.getInt16()
	if err != nil {
		return err
	}

	r.Err = KError(kerr)

	if r.ErrorMessage, err = pd.getNullableString(); err != nil {
		return err
	}

	r.SaslAuthBytes, err = pd.getBytes()

	return err
}

func (r *SaslAuthenticateResponse) key() int16 {
	return APIKeySASLAuth
}

func (r *SaslAuthenticateResponse) version() int16 {
	return 0
}

func (r *SaslAuthenticateResponse) headerVersion() int16 {
	return 0
}

func (r *SaslAuthenticateResponse) requiredVersion() KafkaVersion {
	return V1_0_0_0
}
