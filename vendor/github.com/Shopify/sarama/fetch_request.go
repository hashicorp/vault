package sarama

type fetchRequestBlock struct {
	Version            int16
	currentLeaderEpoch int32
	fetchOffset        int64
	logStartOffset     int64
	maxBytes           int32
}

func (b *fetchRequestBlock) encode(pe packetEncoder, version int16) error {
	b.Version = version
	if b.Version >= 9 {
		pe.putInt32(b.currentLeaderEpoch)
	}
	pe.putInt64(b.fetchOffset)
	if b.Version >= 5 {
		pe.putInt64(b.logStartOffset)
	}
	pe.putInt32(b.maxBytes)
	return nil
}

func (b *fetchRequestBlock) decode(pd packetDecoder, version int16) (err error) {
	b.Version = version
	if b.Version >= 9 {
		if b.currentLeaderEpoch, err = pd.getInt32(); err != nil {
			return err
		}
	}
	if b.fetchOffset, err = pd.getInt64(); err != nil {
		return err
	}
	if b.Version >= 5 {
		if b.logStartOffset, err = pd.getInt64(); err != nil {
			return err
		}
	}
	if b.maxBytes, err = pd.getInt32(); err != nil {
		return err
	}
	return nil
}

// FetchRequest (API key 1) will fetch Kafka messages. Version 3 introduced the MaxBytes field. See
// https://issues.apache.org/jira/browse/KAFKA-2063 for a discussion of the issues leading up to that.  The KIP is at
// https://cwiki.apache.org/confluence/display/KAFKA/KIP-74%3A+Add+Fetch+Response+Size+Limit+in+Bytes
type FetchRequest struct {
	MaxWaitTime  int32
	MinBytes     int32
	MaxBytes     int32
	Version      int16
	Isolation    IsolationLevel
	SessionID    int32
	SessionEpoch int32
	blocks       map[string]map[int32]*fetchRequestBlock
	forgotten    map[string][]int32
	RackID       string
}

type IsolationLevel int8

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
)

func (r *FetchRequest) encode(pe packetEncoder) (err error) {
	pe.putInt32(-1) // replica ID is always -1 for clients
	pe.putInt32(r.MaxWaitTime)
	pe.putInt32(r.MinBytes)
	if r.Version >= 3 {
		pe.putInt32(r.MaxBytes)
	}
	if r.Version >= 4 {
		pe.putInt8(int8(r.Isolation))
	}
	if r.Version >= 7 {
		pe.putInt32(r.SessionID)
		pe.putInt32(r.SessionEpoch)
	}
	err = pe.putArrayLength(len(r.blocks))
	if err != nil {
		return err
	}
	for topic, blocks := range r.blocks {
		err = pe.putString(topic)
		if err != nil {
			return err
		}
		err = pe.putArrayLength(len(blocks))
		if err != nil {
			return err
		}
		for partition, block := range blocks {
			pe.putInt32(partition)
			err = block.encode(pe, r.Version)
			if err != nil {
				return err
			}
		}
	}
	if r.Version >= 7 {
		err = pe.putArrayLength(len(r.forgotten))
		if err != nil {
			return err
		}
		for topic, partitions := range r.forgotten {
			err = pe.putString(topic)
			if err != nil {
				return err
			}
			err = pe.putArrayLength(len(partitions))
			if err != nil {
				return err
			}
			for _, partition := range partitions {
				pe.putInt32(partition)
			}
		}
	}
	if r.Version >= 11 {
		err = pe.putString(r.RackID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FetchRequest) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if _, err = pd.getInt32(); err != nil {
		return err
	}
	if r.MaxWaitTime, err = pd.getInt32(); err != nil {
		return err
	}
	if r.MinBytes, err = pd.getInt32(); err != nil {
		return err
	}
	if r.Version >= 3 {
		if r.MaxBytes, err = pd.getInt32(); err != nil {
			return err
		}
	}
	if r.Version >= 4 {
		isolation, err := pd.getInt8()
		if err != nil {
			return err
		}
		r.Isolation = IsolationLevel(isolation)
	}
	if r.Version >= 7 {
		r.SessionID, err = pd.getInt32()
		if err != nil {
			return err
		}
		r.SessionEpoch, err = pd.getInt32()
		if err != nil {
			return err
		}
	}
	topicCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if topicCount == 0 {
		return nil
	}
	r.blocks = make(map[string]map[int32]*fetchRequestBlock)
	for i := 0; i < topicCount; i++ {
		topic, err := pd.getString()
		if err != nil {
			return err
		}
		partitionCount, err := pd.getArrayLength()
		if err != nil {
			return err
		}
		r.blocks[topic] = make(map[int32]*fetchRequestBlock)
		for j := 0; j < partitionCount; j++ {
			partition, err := pd.getInt32()
			if err != nil {
				return err
			}
			fetchBlock := &fetchRequestBlock{}
			if err = fetchBlock.decode(pd, r.Version); err != nil {
				return err
			}
			r.blocks[topic][partition] = fetchBlock
		}
	}

	if r.Version >= 7 {
		forgottenCount, err := pd.getArrayLength()
		if err != nil {
			return err
		}
		r.forgotten = make(map[string][]int32)
		for i := 0; i < forgottenCount; i++ {
			topic, err := pd.getString()
			if err != nil {
				return err
			}
			partitionCount, err := pd.getArrayLength()
			if err != nil {
				return err
			}
			r.forgotten[topic] = make([]int32, partitionCount)

			for j := 0; j < partitionCount; j++ {
				partition, err := pd.getInt32()
				if err != nil {
					return err
				}
				r.forgotten[topic][j] = partition
			}
		}
	}

	if r.Version >= 11 {
		r.RackID, err = pd.getString()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FetchRequest) key() int16 {
	return 1
}

func (r *FetchRequest) version() int16 {
	return r.Version
}

func (r *FetchRequest) headerVersion() int16 {
	return 1
}

func (r *FetchRequest) requiredVersion() KafkaVersion {
	switch r.Version {
	case 0:
		return MinVersion
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_10_1_0
	case 4, 5:
		return V0_11_0_0
	case 6:
		return V1_0_0_0
	case 7:
		return V1_1_0_0
	case 8:
		return V2_0_0_0
	case 9, 10:
		return V2_1_0_0
	case 11:
		return V2_3_0_0
	default:
		return MaxVersion
	}
}

func (r *FetchRequest) AddBlock(topic string, partitionID int32, fetchOffset int64, maxBytes int32) {
	if r.blocks == nil {
		r.blocks = make(map[string]map[int32]*fetchRequestBlock)
	}

	if r.Version >= 7 && r.forgotten == nil {
		r.forgotten = make(map[string][]int32)
	}

	if r.blocks[topic] == nil {
		r.blocks[topic] = make(map[int32]*fetchRequestBlock)
	}

	tmp := new(fetchRequestBlock)
	tmp.Version = r.Version
	tmp.maxBytes = maxBytes
	tmp.fetchOffset = fetchOffset
	if r.Version >= 9 {
		tmp.currentLeaderEpoch = int32(-1)
	}

	r.blocks[topic][partitionID] = tmp
}
