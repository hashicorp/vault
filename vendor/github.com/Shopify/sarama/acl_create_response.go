package sarama

import "time"

//CreateAclsResponse is a an acl response creation type
type CreateAclsResponse struct {
	ThrottleTime         time.Duration
	AclCreationResponses []*AclCreationResponse
}

func (c *CreateAclsResponse) encode(pe packetEncoder) error {
	pe.putInt32(int32(c.ThrottleTime / time.Millisecond))

	if err := pe.putArrayLength(len(c.AclCreationResponses)); err != nil {
		return err
	}

	for _, aclCreationResponse := range c.AclCreationResponses {
		if err := aclCreationResponse.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (c *CreateAclsResponse) decode(pd packetDecoder, version int16) (err error) {
	throttleTime, err := pd.getInt32()
	if err != nil {
		return err
	}
	c.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	c.AclCreationResponses = make([]*AclCreationResponse, n)
	for i := 0; i < n; i++ {
		c.AclCreationResponses[i] = new(AclCreationResponse)
		if err := c.AclCreationResponses[i].decode(pd, version); err != nil {
			return err
		}
	}

	return nil
}

func (c *CreateAclsResponse) key() int16 {
	return 30
}

func (c *CreateAclsResponse) version() int16 {
	return 0
}

func (c *CreateAclsResponse) headerVersion() int16 {
	return 0
}

func (c *CreateAclsResponse) requiredVersion() KafkaVersion {
	return V0_11_0_0
}

//AclCreationResponse is an acl creation response type
type AclCreationResponse struct {
	Err    KError
	ErrMsg *string
}

func (a *AclCreationResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(a.Err))

	if err := pe.putNullableString(a.ErrMsg); err != nil {
		return err
	}

	return nil
}

func (a *AclCreationResponse) decode(pd packetDecoder, version int16) (err error) {
	kerr, err := pd.getInt16()
	if err != nil {
		return err
	}
	a.Err = KError(kerr)

	if a.ErrMsg, err = pd.getNullableString(); err != nil {
		return err
	}

	return nil
}
