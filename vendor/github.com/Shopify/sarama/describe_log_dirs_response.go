package sarama

import "time"

type DescribeLogDirsResponse struct {
	ThrottleTime time.Duration

	// Version 0 and 1 are equal
	// The version number is bumped to indicate that on quota violation brokers send out responses before throttling.
	Version int16

	LogDirs []DescribeLogDirsResponseDirMetadata
}

func (r *DescribeLogDirsResponse) encode(pe packetEncoder) error {
	pe.putInt32(int32(r.ThrottleTime / time.Millisecond))

	if err := pe.putArrayLength(len(r.LogDirs)); err != nil {
		return err
	}

	for _, dir := range r.LogDirs {
		if err := dir.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (r *DescribeLogDirsResponse) decode(pd packetDecoder, version int16) error {
	throttleTime, err := pd.getInt32()
	if err != nil {
		return err
	}
	r.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	// Decode array of DescribeLogDirsResponseDirMetadata
	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.LogDirs = make([]DescribeLogDirsResponseDirMetadata, n)
	for i := 0; i < n; i++ {
		dir := DescribeLogDirsResponseDirMetadata{}
		if err := dir.decode(pd, version); err != nil {
			return err
		}
		r.LogDirs[i] = dir
	}

	return nil
}

func (r *DescribeLogDirsResponse) key() int16 {
	return 35
}

func (r *DescribeLogDirsResponse) version() int16 {
	return r.Version
}

func (r *DescribeLogDirsResponse) headerVersion() int16 {
	return 0
}

func (r *DescribeLogDirsResponse) requiredVersion() KafkaVersion {
	return V1_0_0_0
}

type DescribeLogDirsResponseDirMetadata struct {
	ErrorCode KError

	// The absolute log directory path
	Path   string
	Topics []DescribeLogDirsResponseTopic
}

func (r *DescribeLogDirsResponseDirMetadata) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.ErrorCode))

	if err := pe.putString(r.Path); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.Topics)); err != nil {
		return err
	}
	for _, topic := range r.Topics {
		if err := topic.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (r *DescribeLogDirsResponseDirMetadata) decode(pd packetDecoder, version int16) error {
	errCode, err := pd.getInt16()
	if err != nil {
		return err
	}
	r.ErrorCode = KError(errCode)

	path, err := pd.getString()
	if err != nil {
		return err
	}
	r.Path = path

	// Decode array of DescribeLogDirsResponseTopic
	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.Topics = make([]DescribeLogDirsResponseTopic, n)
	for i := 0; i < n; i++ {
		t := DescribeLogDirsResponseTopic{}

		if err := t.decode(pd, version); err != nil {
			return err
		}

		r.Topics[i] = t
	}

	return nil
}

// DescribeLogDirsResponseTopic contains a topic's partitions descriptions
type DescribeLogDirsResponseTopic struct {
	Topic      string
	Partitions []DescribeLogDirsResponsePartition
}

func (r *DescribeLogDirsResponseTopic) encode(pe packetEncoder) error {
	if err := pe.putString(r.Topic); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.Partitions)); err != nil {
		return err
	}
	for _, partition := range r.Partitions {
		if err := partition.encode(pe); err != nil {
			return err
		}
	}

	return nil
}

func (r *DescribeLogDirsResponseTopic) decode(pd packetDecoder, version int16) error {
	t, err := pd.getString()
	if err != nil {
		return err
	}
	r.Topic = t

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	r.Partitions = make([]DescribeLogDirsResponsePartition, n)
	for i := 0; i < n; i++ {
		p := DescribeLogDirsResponsePartition{}
		if err := p.decode(pd, version); err != nil {
			return err
		}
		r.Partitions[i] = p
	}

	return nil
}

// DescribeLogDirsResponsePartition describes a partition's log directory
type DescribeLogDirsResponsePartition struct {
	PartitionID int32

	// The size of the log segments of the partition in bytes.
	Size int64

	// The lag of the log's LEO w.r.t. partition's HW (if it is the current log for the partition) or
	// current replica's LEO (if it is the future log for the partition)
	OffsetLag int64

	// True if this log is created by AlterReplicaLogDirsRequest and will replace the current log of
	// the replica in the future.
	IsTemporary bool
}

func (r *DescribeLogDirsResponsePartition) encode(pe packetEncoder) error {
	pe.putInt32(r.PartitionID)
	pe.putInt64(r.Size)
	pe.putInt64(r.OffsetLag)
	pe.putBool(r.IsTemporary)

	return nil
}

func (r *DescribeLogDirsResponsePartition) decode(pd packetDecoder, version int16) error {
	pID, err := pd.getInt32()
	if err != nil {
		return err
	}
	r.PartitionID = pID

	size, err := pd.getInt64()
	if err != nil {
		return err
	}
	r.Size = size

	lag, err := pd.getInt64()
	if err != nil {
		return err
	}
	r.OffsetLag = lag

	isTemp, err := pd.getBool()
	if err != nil {
		return err
	}
	r.IsTemporary = isTemp

	return nil
}
