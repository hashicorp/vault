package sarama

type ListPartitionReassignmentsRequest struct {
	TimeoutMs int32
	blocks    map[string][]int32
	Version   int16
}

func (r *ListPartitionReassignmentsRequest) encode(pe packetEncoder) error {
	pe.putInt32(r.TimeoutMs)

	pe.putCompactArrayLength(len(r.blocks))

	for topic, partitions := range r.blocks {
		if err := pe.putCompactString(topic); err != nil {
			return err
		}

		if err := pe.putCompactInt32Array(partitions); err != nil {
			return err
		}

		pe.putEmptyTaggedFieldArray()
	}

	pe.putEmptyTaggedFieldArray()

	return nil
}

func (r *ListPartitionReassignmentsRequest) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if r.TimeoutMs, err = pd.getInt32(); err != nil {
		return err
	}

	topicCount, err := pd.getCompactArrayLength()
	if err != nil {
		return err
	}
	if topicCount > 0 {
		r.blocks = make(map[string][]int32)
		for i := 0; i < topicCount; i++ {
			topic, err := pd.getCompactString()
			if err != nil {
				return err
			}
			partitionCount, err := pd.getCompactArrayLength()
			if err != nil {
				return err
			}
			r.blocks[topic] = make([]int32, partitionCount)
			for j := 0; j < partitionCount; j++ {
				partition, err := pd.getInt32()
				if err != nil {
					return err
				}
				r.blocks[topic][j] = partition
			}
			if _, err := pd.getEmptyTaggedFieldArray(); err != nil {
				return err
			}
		}
	}

	if _, err := pd.getEmptyTaggedFieldArray(); err != nil {
		return err
	}

	return
}

func (r *ListPartitionReassignmentsRequest) key() int16 {
	return 46
}

func (r *ListPartitionReassignmentsRequest) version() int16 {
	return r.Version
}

func (r *ListPartitionReassignmentsRequest) headerVersion() int16 {
	return 2
}

func (r *ListPartitionReassignmentsRequest) requiredVersion() KafkaVersion {
	return V2_4_0_0
}

func (r *ListPartitionReassignmentsRequest) AddBlock(topic string, partitionIDs []int32) {
	if r.blocks == nil {
		r.blocks = make(map[string][]int32)
	}

	if r.blocks[topic] == nil {
		r.blocks[topic] = partitionIDs
	}
}
