package sarama

type alterPartitionReassignmentsBlock struct {
	replicas []int32
}

func (b *alterPartitionReassignmentsBlock) encode(pe packetEncoder) error {
	if err := pe.putNullableCompactInt32Array(b.replicas); err != nil {
		return err
	}

	pe.putEmptyTaggedFieldArray()
	return nil
}

func (b *alterPartitionReassignmentsBlock) decode(pd packetDecoder) (err error) {
	if b.replicas, err = pd.getCompactInt32Array(); err != nil {
		return err
	}
	return nil
}

type AlterPartitionReassignmentsRequest struct {
	TimeoutMs int32
	blocks    map[string]map[int32]*alterPartitionReassignmentsBlock
	Version   int16
}

func (r *AlterPartitionReassignmentsRequest) encode(pe packetEncoder) error {
	pe.putInt32(r.TimeoutMs)

	pe.putCompactArrayLength(len(r.blocks))

	for topic, partitions := range r.blocks {
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

func (r *AlterPartitionReassignmentsRequest) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if r.TimeoutMs, err = pd.getInt32(); err != nil {
		return err
	}

	topicCount, err := pd.getCompactArrayLength()
	if err != nil {
		return err
	}
	if topicCount > 0 {
		r.blocks = make(map[string]map[int32]*alterPartitionReassignmentsBlock)
		for i := 0; i < topicCount; i++ {
			topic, err := pd.getCompactString()
			if err != nil {
				return err
			}
			partitionCount, err := pd.getCompactArrayLength()
			if err != nil {
				return err
			}
			r.blocks[topic] = make(map[int32]*alterPartitionReassignmentsBlock)
			for j := 0; j < partitionCount; j++ {
				partition, err := pd.getInt32()
				if err != nil {
					return err
				}
				block := &alterPartitionReassignmentsBlock{}
				if err := block.decode(pd); err != nil {
					return err
				}
				r.blocks[topic][partition] = block

				if _, err := pd.getEmptyTaggedFieldArray(); err != nil {
					return err
				}
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

func (r *AlterPartitionReassignmentsRequest) key() int16 {
	return 45
}

func (r *AlterPartitionReassignmentsRequest) version() int16 {
	return r.Version
}

func (r *AlterPartitionReassignmentsRequest) headerVersion() int16 {
	return 2
}

func (r *AlterPartitionReassignmentsRequest) requiredVersion() KafkaVersion {
	return V2_4_0_0
}

func (r *AlterPartitionReassignmentsRequest) AddBlock(topic string, partitionID int32, replicas []int32) {
	if r.blocks == nil {
		r.blocks = make(map[string]map[int32]*alterPartitionReassignmentsBlock)
	}

	if r.blocks[topic] == nil {
		r.blocks[topic] = make(map[int32]*alterPartitionReassignmentsBlock)
	}

	r.blocks[topic][partitionID] = &alterPartitionReassignmentsBlock{replicas}
}
