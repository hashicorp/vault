package sarama

// DescribeLogDirsRequest is a describe request to get partitions' log size
type DescribeLogDirsRequest struct {
	// Version 0 and 1 are equal
	// The version number is bumped to indicate that on quota violation brokers send out responses before throttling.
	Version int16

	// If this is an empty array, all topics will be queried
	DescribeTopics []DescribeLogDirsRequestTopic
}

// DescribeLogDirsRequestTopic is a describe request about the log dir of one or more partitions within a Topic
type DescribeLogDirsRequestTopic struct {
	Topic        string
	PartitionIDs []int32
}

func (r *DescribeLogDirsRequest) encode(pe packetEncoder) error {
	length := len(r.DescribeTopics)
	if length == 0 {
		// In order to query all topics we must send null
		length = -1
	}

	if err := pe.putArrayLength(length); err != nil {
		return err
	}

	for _, d := range r.DescribeTopics {
		if err := pe.putString(d.Topic); err != nil {
			return err
		}

		if err := pe.putInt32Array(d.PartitionIDs); err != nil {
			return err
		}
	}

	return nil
}

func (r *DescribeLogDirsRequest) decode(pd packetDecoder, version int16) error {
	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if n == -1 {
		n = 0
	}

	topics := make([]DescribeLogDirsRequestTopic, n)
	for i := 0; i < n; i++ {
		topics[i] = DescribeLogDirsRequestTopic{}

		topic, err := pd.getString()
		if err != nil {
			return err
		}
		topics[i].Topic = topic

		pIDs, err := pd.getInt32Array()
		if err != nil {
			return err
		}
		topics[i].PartitionIDs = pIDs
	}
	r.DescribeTopics = topics

	return nil
}

func (r *DescribeLogDirsRequest) key() int16 {
	return 35
}

func (r *DescribeLogDirsRequest) version() int16 {
	return r.Version
}

func (r *DescribeLogDirsRequest) headerVersion() int16 {
	return 1
}

func (r *DescribeLogDirsRequest) requiredVersion() KafkaVersion {
	return V1_0_0_0
}
