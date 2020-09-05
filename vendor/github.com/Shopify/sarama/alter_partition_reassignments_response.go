package sarama

type alterPartitionReassignmentsErrorBlock struct {
	errorCode    KError
	errorMessage *string
}

func (b *alterPartitionReassignmentsErrorBlock) encode(pe packetEncoder) error {
	pe.putInt16(int16(b.errorCode))
	if err := pe.putNullableCompactString(b.errorMessage); err != nil {
		return err
	}
	pe.putEmptyTaggedFieldArray()

	return nil
}

func (b *alterPartitionReassignmentsErrorBlock) decode(pd packetDecoder) (err error) {
	errorCode, err := pd.getInt16()
	if err != nil {
		return err
	}
	b.errorCode = KError(errorCode)
	b.errorMessage, err = pd.getCompactNullableString()

	if _, err := pd.getEmptyTaggedFieldArray(); err != nil {
		return err
	}
	return err
}

type AlterPartitionReassignmentsResponse struct {
	Version        int16
	ThrottleTimeMs int32
	ErrorCode      KError
	ErrorMessage   *string
	Errors         map[string]map[int32]*alterPartitionReassignmentsErrorBlock
}

func (r *AlterPartitionReassignmentsResponse) AddError(topic string, partition int32, kerror KError, message *string) {
	if r.Errors == nil {
		r.Errors = make(map[string]map[int32]*alterPartitionReassignmentsErrorBlock)
	}
	partitions := r.Errors[topic]
	if partitions == nil {
		partitions = make(map[int32]*alterPartitionReassignmentsErrorBlock)
		r.Errors[topic] = partitions
	}

	partitions[partition] = &alterPartitionReassignmentsErrorBlock{errorCode: kerror, errorMessage: message}
}

func (r *AlterPartitionReassignmentsResponse) encode(pe packetEncoder) error {
	pe.putInt32(r.ThrottleTimeMs)
	pe.putInt16(int16(r.ErrorCode))
	if err := pe.putNullableCompactString(r.ErrorMessage); err != nil {
		return err
	}

	pe.putCompactArrayLength(len(r.Errors))
	for topic, partitions := range r.Errors {
		if err := pe.putCompactString(topic); err != nil {
			return err
		}
		pe.putCompactArrayLength(len(partitions))
		for partition, block := range partitions {
			pe.putInt32(partition)

			if err := block.encode(pe); err != nil {
				return err
			}
		}
		pe.putEmptyTaggedFieldArray()
	}

	pe.putEmptyTaggedFieldArray()
	return nil
}

func (r *AlterPartitionReassignmentsResponse) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if r.ThrottleTimeMs, err = pd.getInt32(); err != nil {
		return err
	}

	kerr, err := pd.getInt16()
	if err != nil {
		return err
	}

	r.ErrorCode = KError(kerr)

	if r.ErrorMessage, err = pd.getCompactNullableString(); err != nil {
		return err
	}

	numTopics, err := pd.getCompactArrayLength()
	if err != nil {
		return err
	}

	if numTopics > 0 {
		r.Errors = make(map[string]map[int32]*alterPartitionReassignmentsErrorBlock, numTopics)
		for i := 0; i < numTopics; i++ {
			topic, err := pd.getCompactString()
			if err != nil {
				return err
			}

			ongoingPartitionReassignments, err := pd.getCompactArrayLength()
			if err != nil {
				return err
			}

			r.Errors[topic] = make(map[int32]*alterPartitionReassignmentsErrorBlock, ongoingPartitionReassignments)

			for j := 0; j < ongoingPartitionReassignments; j++ {
				partition, err := pd.getInt32()
				if err != nil {
					return err
				}
				block := &alterPartitionReassignmentsErrorBlock{}
				if err := block.decode(pd); err != nil {
					return err
				}

				r.Errors[topic][partition] = block
			}
			if _, err = pd.getEmptyTaggedFieldArray(); err != nil {
				return err
			}
		}
	}

	if _, err = pd.getEmptyTaggedFieldArray(); err != nil {
		return err
	}

	return nil
}

func (r *AlterPartitionReassignmentsResponse) key() int16 {
	return 45
}

func (r *AlterPartitionReassignmentsResponse) version() int16 {
	return r.Version
}

func (r *AlterPartitionReassignmentsResponse) headerVersion() int16 {
	return 1
}

func (r *AlterPartitionReassignmentsResponse) requiredVersion() KafkaVersion {
	return V2_4_0_0
}
