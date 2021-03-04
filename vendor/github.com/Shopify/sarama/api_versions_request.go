package sarama

//ApiVersionsRequest ...
type ApiVersionsRequest struct {
}

func (a *ApiVersionsRequest) encode(pe packetEncoder) error {
	return nil
}

func (a *ApiVersionsRequest) decode(pd packetDecoder, version int16) (err error) {
	return nil
}

func (a *ApiVersionsRequest) key() int16 {
	return 18
}

func (a *ApiVersionsRequest) version() int16 {
	return 0
}

func (a *ApiVersionsRequest) headerVersion() int16 {
	return 1
}

func (a *ApiVersionsRequest) requiredVersion() KafkaVersion {
	return V0_10_0_0
}
