package sarama

type topicPartitionAssignment struct {
	Topic     string
	Partition int32
}

type StickyAssignorUserData interface {
	partitions() []topicPartitionAssignment
	hasGeneration() bool
	generation() int
}

//StickyAssignorUserDataV0 holds topic partition information for an assignment
type StickyAssignorUserDataV0 struct {
	Topics map[string][]int32

	topicPartitions []topicPartitionAssignment
}

func (m *StickyAssignorUserDataV0) encode(pe packetEncoder) error {
	if err := pe.putArrayLength(len(m.Topics)); err != nil {
		return err
	}

	for topic, partitions := range m.Topics {
		if err := pe.putString(topic); err != nil {
			return err
		}
		if err := pe.putInt32Array(partitions); err != nil {
			return err
		}
	}
	return nil
}

func (m *StickyAssignorUserDataV0) decode(pd packetDecoder) (err error) {
	var topicLen int
	if topicLen, err = pd.getArrayLength(); err != nil {
		return
	}

	m.Topics = make(map[string][]int32, topicLen)
	for i := 0; i < topicLen; i++ {
		var topic string
		if topic, err = pd.getString(); err != nil {
			return
		}
		if m.Topics[topic], err = pd.getInt32Array(); err != nil {
			return
		}
	}
	m.topicPartitions = populateTopicPartitions(m.Topics)
	return nil
}

func (m *StickyAssignorUserDataV0) partitions() []topicPartitionAssignment { return m.topicPartitions }
func (m *StickyAssignorUserDataV0) hasGeneration() bool                    { return false }
func (m *StickyAssignorUserDataV0) generation() int                        { return defaultGeneration }

//StickyAssignorUserDataV1 holds topic partition information for an assignment
type StickyAssignorUserDataV1 struct {
	Topics     map[string][]int32
	Generation int32

	topicPartitions []topicPartitionAssignment
}

func (m *StickyAssignorUserDataV1) encode(pe packetEncoder) error {
	if err := pe.putArrayLength(len(m.Topics)); err != nil {
		return err
	}

	for topic, partitions := range m.Topics {
		if err := pe.putString(topic); err != nil {
			return err
		}
		if err := pe.putInt32Array(partitions); err != nil {
			return err
		}
	}

	pe.putInt32(m.Generation)
	return nil
}

func (m *StickyAssignorUserDataV1) decode(pd packetDecoder) (err error) {
	var topicLen int
	if topicLen, err = pd.getArrayLength(); err != nil {
		return
	}

	m.Topics = make(map[string][]int32, topicLen)
	for i := 0; i < topicLen; i++ {
		var topic string
		if topic, err = pd.getString(); err != nil {
			return
		}
		if m.Topics[topic], err = pd.getInt32Array(); err != nil {
			return
		}
	}

	m.Generation, err = pd.getInt32()
	if err != nil {
		return err
	}
	m.topicPartitions = populateTopicPartitions(m.Topics)
	return nil
}

func (m *StickyAssignorUserDataV1) partitions() []topicPartitionAssignment { return m.topicPartitions }
func (m *StickyAssignorUserDataV1) hasGeneration() bool                    { return true }
func (m *StickyAssignorUserDataV1) generation() int                        { return int(m.Generation) }

func populateTopicPartitions(topics map[string][]int32) []topicPartitionAssignment {
	topicPartitions := make([]topicPartitionAssignment, 0)
	for topic, partitions := range topics {
		for _, partition := range partitions {
			topicPartitions = append(topicPartitions, topicPartitionAssignment{Topic: topic, Partition: partition})
		}
	}
	return topicPartitions
}
