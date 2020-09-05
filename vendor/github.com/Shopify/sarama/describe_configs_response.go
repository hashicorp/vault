package sarama

import (
	"fmt"
	"time"
)

type ConfigSource int8

func (s ConfigSource) String() string {
	switch s {
	case SourceUnknown:
		return "Unknown"
	case SourceTopic:
		return "Topic"
	case SourceDynamicBroker:
		return "DynamicBroker"
	case SourceDynamicDefaultBroker:
		return "DynamicDefaultBroker"
	case SourceStaticBroker:
		return "StaticBroker"
	case SourceDefault:
		return "Default"
	}
	return fmt.Sprintf("Source Invalid: %d", int(s))
}

const (
	SourceUnknown ConfigSource = iota
	SourceTopic
	SourceDynamicBroker
	SourceDynamicDefaultBroker
	SourceStaticBroker
	SourceDefault
)

type DescribeConfigsResponse struct {
	Version      int16
	ThrottleTime time.Duration
	Resources    []*ResourceResponse
}

type ResourceResponse struct {
	ErrorCode int16
	ErrorMsg  string
	Type      ConfigResourceType
	Name      string
	Configs   []*ConfigEntry
}

type ConfigEntry struct {
	Name      string
	Value     string
	ReadOnly  bool
	Default   bool
	Source    ConfigSource
	Sensitive bool
	Synonyms  []*ConfigSynonym
}

type ConfigSynonym struct {
	ConfigName  string
	ConfigValue string
	Source      ConfigSource
}

func (r *DescribeConfigsResponse) encode(pe packetEncoder) (err error) {
	pe.putInt32(int32(r.ThrottleTime / time.Millisecond))
	if err = pe.putArrayLength(len(r.Resources)); err != nil {
		return err
	}

	for _, c := range r.Resources {
		if err = c.encode(pe, r.Version); err != nil {
			return err
		}
	}

	return nil
}

func (r *DescribeConfigsResponse) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version
	throttleTime, err := pd.getInt32()
	if err != nil {
		return err
	}
	r.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.Resources = make([]*ResourceResponse, n)
	for i := 0; i < n; i++ {
		rr := &ResourceResponse{}
		if err := rr.decode(pd, version); err != nil {
			return err
		}
		r.Resources[i] = rr
	}

	return nil
}

func (r *DescribeConfigsResponse) key() int16 {
	return 32
}

func (r *DescribeConfigsResponse) version() int16 {
	return r.Version
}

func (r *DescribeConfigsResponse) headerVersion() int16 {
	return 0
}

func (r *DescribeConfigsResponse) requiredVersion() KafkaVersion {
	switch r.Version {
	case 1:
		return V1_0_0_0
	case 2:
		return V2_0_0_0
	default:
		return V0_11_0_0
	}
}

func (r *ResourceResponse) encode(pe packetEncoder, version int16) (err error) {
	pe.putInt16(r.ErrorCode)

	if err = pe.putString(r.ErrorMsg); err != nil {
		return err
	}

	pe.putInt8(int8(r.Type))

	if err = pe.putString(r.Name); err != nil {
		return err
	}

	if err = pe.putArrayLength(len(r.Configs)); err != nil {
		return err
	}

	for _, c := range r.Configs {
		if err = c.encode(pe, version); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceResponse) decode(pd packetDecoder, version int16) (err error) {
	ec, err := pd.getInt16()
	if err != nil {
		return err
	}
	r.ErrorCode = ec

	em, err := pd.getString()
	if err != nil {
		return err
	}
	r.ErrorMsg = em

	t, err := pd.getInt8()
	if err != nil {
		return err
	}
	r.Type = ConfigResourceType(t)

	name, err := pd.getString()
	if err != nil {
		return err
	}
	r.Name = name

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.Configs = make([]*ConfigEntry, n)
	for i := 0; i < n; i++ {
		c := &ConfigEntry{}
		if err := c.decode(pd, version); err != nil {
			return err
		}
		r.Configs[i] = c
	}
	return nil
}

func (r *ConfigEntry) encode(pe packetEncoder, version int16) (err error) {
	if err = pe.putString(r.Name); err != nil {
		return err
	}

	if err = pe.putString(r.Value); err != nil {
		return err
	}

	pe.putBool(r.ReadOnly)

	if version <= 0 {
		pe.putBool(r.Default)
		pe.putBool(r.Sensitive)
	} else {
		pe.putInt8(int8(r.Source))
		pe.putBool(r.Sensitive)

		if err := pe.putArrayLength(len(r.Synonyms)); err != nil {
			return err
		}
		for _, c := range r.Synonyms {
			if err = c.encode(pe, version); err != nil {
				return err
			}
		}
	}

	return nil
}

//https://cwiki.apache.org/confluence/display/KAFKA/KIP-226+-+Dynamic+Broker+Configuration
func (r *ConfigEntry) decode(pd packetDecoder, version int16) (err error) {
	if version == 0 {
		r.Source = SourceUnknown
	}
	name, err := pd.getString()
	if err != nil {
		return err
	}
	r.Name = name

	value, err := pd.getString()
	if err != nil {
		return err
	}
	r.Value = value

	read, err := pd.getBool()
	if err != nil {
		return err
	}
	r.ReadOnly = read

	if version == 0 {
		defaultB, err := pd.getBool()
		if err != nil {
			return err
		}
		r.Default = defaultB
		if defaultB {
			r.Source = SourceDefault
		}
	} else {
		source, err := pd.getInt8()
		if err != nil {
			return err
		}
		r.Source = ConfigSource(source)
		r.Default = r.Source == SourceDefault
	}

	sensitive, err := pd.getBool()
	if err != nil {
		return err
	}
	r.Sensitive = sensitive

	if version > 0 {
		n, err := pd.getArrayLength()
		if err != nil {
			return err
		}
		r.Synonyms = make([]*ConfigSynonym, n)

		for i := 0; i < n; i++ {
			s := &ConfigSynonym{}
			if err := s.decode(pd, version); err != nil {
				return err
			}
			r.Synonyms[i] = s
		}
	}
	return nil
}

func (c *ConfigSynonym) encode(pe packetEncoder, version int16) (err error) {
	err = pe.putString(c.ConfigName)
	if err != nil {
		return err
	}

	err = pe.putString(c.ConfigValue)
	if err != nil {
		return err
	}

	pe.putInt8(int8(c.Source))

	return nil
}

func (c *ConfigSynonym) decode(pd packetDecoder, version int16) error {
	name, err := pd.getString()
	if err != nil {
		return nil
	}
	c.ConfigName = name

	value, err := pd.getString()
	if err != nil {
		return nil
	}
	c.ConfigValue = value

	source, err := pd.getInt8()
	if err != nil {
		return nil
	}
	c.Source = ConfigSource(source)
	return nil
}
