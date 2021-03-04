package sarama

//DeleteAclsRequest is a delete acl request
type DeleteAclsRequest struct {
	Version int
	Filters []*AclFilter
}

func (d *DeleteAclsRequest) encode(pe packetEncoder) error {
	if err := pe.putArrayLength(len(d.Filters)); err != nil {
		return err
	}

	for _, filter := range d.Filters {
		filter.Version = d.Version
		if err := filter.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeleteAclsRequest) decode(pd packetDecoder, version int16) (err error) {
	d.Version = int(version)
	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	d.Filters = make([]*AclFilter, n)
	for i := 0; i < n; i++ {
		d.Filters[i] = new(AclFilter)
		d.Filters[i].Version = int(version)
		if err := d.Filters[i].decode(pd, version); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeleteAclsRequest) key() int16 {
	return 31
}

func (d *DeleteAclsRequest) version() int16 {
	return int16(d.Version)
}

func (c *DeleteAclsRequest) headerVersion() int16 {
	return 1
}

func (d *DeleteAclsRequest) requiredVersion() KafkaVersion {
	switch d.Version {
	case 1:
		return V2_0_0_0
	default:
		return V0_11_0_0
	}
}
