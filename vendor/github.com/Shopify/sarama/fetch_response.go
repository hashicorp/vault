package sarama

import (
	"sort"
	"time"
)

type AbortedTransaction struct {
	ProducerID  int64
	FirstOffset int64
}

func (t *AbortedTransaction) decode(pd packetDecoder) (err error) {
	if t.ProducerID, err = pd.getInt64(); err != nil {
		return err
	}

	if t.FirstOffset, err = pd.getInt64(); err != nil {
		return err
	}

	return nil
}

func (t *AbortedTransaction) encode(pe packetEncoder) (err error) {
	pe.putInt64(t.ProducerID)
	pe.putInt64(t.FirstOffset)

	return nil
}

type FetchResponseBlock struct {
	Err                  KError
	HighWaterMarkOffset  int64
	LastStableOffset     int64
	LogStartOffset       int64
	AbortedTransactions  []*AbortedTransaction
	PreferredReadReplica int32
	Records              *Records // deprecated: use FetchResponseBlock.RecordsSet
	RecordsSet           []*Records
	Partial              bool
}

func (b *FetchResponseBlock) decode(pd packetDecoder, version int16) (err error) {
	tmp, err := pd.getInt16()
	if err != nil {
		return err
	}
	b.Err = KError(tmp)

	b.HighWaterMarkOffset, err = pd.getInt64()
	if err != nil {
		return err
	}

	if version >= 4 {
		b.LastStableOffset, err = pd.getInt64()
		if err != nil {
			return err
		}

		if version >= 5 {
			b.LogStartOffset, err = pd.getInt64()
			if err != nil {
				return err
			}
		}

		numTransact, err := pd.getArrayLength()
		if err != nil {
			return err
		}

		if numTransact >= 0 {
			b.AbortedTransactions = make([]*AbortedTransaction, numTransact)
		}

		for i := 0; i < numTransact; i++ {
			transact := new(AbortedTransaction)
			if err = transact.decode(pd); err != nil {
				return err
			}
			b.AbortedTransactions[i] = transact
		}
	}

	if version >= 11 {
		b.PreferredReadReplica, err = pd.getInt32()
		if err != nil {
			return err
		}
	}

	recordsSize, err := pd.getInt32()
	if err != nil {
		return err
	}

	recordsDecoder, err := pd.getSubset(int(recordsSize))
	if err != nil {
		return err
	}

	b.RecordsSet = []*Records{}

	for recordsDecoder.remaining() > 0 {
		records := &Records{}
		if err := records.decode(recordsDecoder); err != nil {
			// If we have at least one decoded records, this is not an error
			if err == ErrInsufficientData {
				if len(b.RecordsSet) == 0 {
					b.Partial = true
				}
				break
			}
			return err
		}

		partial, err := records.isPartial()
		if err != nil {
			return err
		}

		n, err := records.numRecords()
		if err != nil {
			return err
		}

		if n > 0 || (partial && len(b.RecordsSet) == 0) {
			b.RecordsSet = append(b.RecordsSet, records)

			if b.Records == nil {
				b.Records = records
			}
		}

		overflow, err := records.isOverflow()
		if err != nil {
			return err
		}

		if partial || overflow {
			break
		}
	}

	return nil
}

func (b *FetchResponseBlock) numRecords() (int, error) {
	sum := 0

	for _, records := range b.RecordsSet {
		count, err := records.numRecords()
		if err != nil {
			return 0, err
		}

		sum += count
	}

	return sum, nil
}

func (b *FetchResponseBlock) isPartial() (bool, error) {
	if b.Partial {
		return true, nil
	}

	if len(b.RecordsSet) == 1 {
		return b.RecordsSet[0].isPartial()
	}

	return false, nil
}

func (b *FetchResponseBlock) encode(pe packetEncoder, version int16) (err error) {
	pe.putInt16(int16(b.Err))

	pe.putInt64(b.HighWaterMarkOffset)

	if version >= 4 {
		pe.putInt64(b.LastStableOffset)

		if version >= 5 {
			pe.putInt64(b.LogStartOffset)
		}

		if err = pe.putArrayLength(len(b.AbortedTransactions)); err != nil {
			return err
		}
		for _, transact := range b.AbortedTransactions {
			if err = transact.encode(pe); err != nil {
				return err
			}
		}
	}

	if version >= 11 {
		pe.putInt32(b.PreferredReadReplica)
	}

	pe.push(&lengthField{})
	for _, records := range b.RecordsSet {
		err = records.encode(pe)
		if err != nil {
			return err
		}
	}
	return pe.pop()
}

func (b *FetchResponseBlock) getAbortedTransactions() []*AbortedTransaction {
	// I can't find any doc that guarantee the field `fetchResponse.AbortedTransactions` is ordered
	// plus Java implementation use a PriorityQueue based on `FirstOffset`. I guess we have to order it ourself
	at := b.AbortedTransactions
	sort.Slice(
		at,
		func(i, j int) bool { return at[i].FirstOffset < at[j].FirstOffset },
	)
	return at
}

type FetchResponse struct {
	Blocks        map[string]map[int32]*FetchResponseBlock
	ThrottleTime  time.Duration
	ErrorCode     int16
	SessionID     int32
	Version       int16
	LogAppendTime bool
	Timestamp     time.Time
}

func (r *FetchResponse) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if r.Version >= 1 {
		throttle, err := pd.getInt32()
		if err != nil {
			return err
		}
		r.ThrottleTime = time.Duration(throttle) * time.Millisecond
	}

	if r.Version >= 7 {
		r.ErrorCode, err = pd.getInt16()
		if err != nil {
			return err
		}
		r.SessionID, err = pd.getInt32()
		if err != nil {
			return err
		}
	}

	numTopics, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.Blocks = make(map[string]map[int32]*FetchResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ {
		name, err := pd.getString()
		if err != nil {
			return err
		}

		numBlocks, err := pd.getArrayLength()
		if err != nil {
			return err
		}

		r.Blocks[name] = make(map[int32]*FetchResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ {
			id, err := pd.getInt32()
			if err != nil {
				return err
			}

			block := new(FetchResponseBlock)
			err = block.decode(pd, version)
			if err != nil {
				return err
			}
			r.Blocks[name][id] = block
		}
	}

	return nil
}

func (r *FetchResponse) encode(pe packetEncoder) (err error) {
	if r.Version >= 1 {
		pe.putInt32(int32(r.ThrottleTime / time.Millisecond))
	}

	if r.Version >= 7 {
		pe.putInt16(r.ErrorCode)
		pe.putInt32(r.SessionID)
	}

	err = pe.putArrayLength(len(r.Blocks))
	if err != nil {
		return err
	}

	for topic, partitions := range r.Blocks {
		err = pe.putString(topic)
		if err != nil {
			return err
		}

		err = pe.putArrayLength(len(partitions))
		if err != nil {
			return err
		}

		for id, block := range partitions {
			pe.putInt32(id)
			err = block.encode(pe, r.Version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *FetchResponse) key() int16 {
	return 1
}

func (r *FetchResponse) version() int16 {
	return r.Version
}

func (r *FetchResponse) headerVersion() int16 {
	return 0
}

func (r *FetchResponse) requiredVersion() KafkaVersion {
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

func (r *FetchResponse) GetBlock(topic string, partition int32) *FetchResponseBlock {
	if r.Blocks == nil {
		return nil
	}

	if r.Blocks[topic] == nil {
		return nil
	}

	return r.Blocks[topic][partition]
}

func (r *FetchResponse) AddError(topic string, partition int32, err KError) {
	if r.Blocks == nil {
		r.Blocks = make(map[string]map[int32]*FetchResponseBlock)
	}
	partitions, ok := r.Blocks[topic]
	if !ok {
		partitions = make(map[int32]*FetchResponseBlock)
		r.Blocks[topic] = partitions
	}
	frb, ok := partitions[partition]
	if !ok {
		frb = new(FetchResponseBlock)
		partitions[partition] = frb
	}
	frb.Err = err
}

func (r *FetchResponse) getOrCreateBlock(topic string, partition int32) *FetchResponseBlock {
	if r.Blocks == nil {
		r.Blocks = make(map[string]map[int32]*FetchResponseBlock)
	}
	partitions, ok := r.Blocks[topic]
	if !ok {
		partitions = make(map[int32]*FetchResponseBlock)
		r.Blocks[topic] = partitions
	}
	frb, ok := partitions[partition]
	if !ok {
		frb = new(FetchResponseBlock)
		partitions[partition] = frb
	}

	return frb
}

func encodeKV(key, value Encoder) ([]byte, []byte) {
	var kb []byte
	var vb []byte
	if key != nil {
		kb, _ = key.Encode()
	}
	if value != nil {
		vb, _ = value.Encode()
	}

	return kb, vb
}

func (r *FetchResponse) AddMessageWithTimestamp(topic string, partition int32, key, value Encoder, offset int64, timestamp time.Time, version int8) {
	frb := r.getOrCreateBlock(topic, partition)
	kb, vb := encodeKV(key, value)
	if r.LogAppendTime {
		timestamp = r.Timestamp
	}
	msg := &Message{Key: kb, Value: vb, LogAppendTime: r.LogAppendTime, Timestamp: timestamp, Version: version}
	msgBlock := &MessageBlock{Msg: msg, Offset: offset}
	if len(frb.RecordsSet) == 0 {
		records := newLegacyRecords(&MessageSet{})
		frb.RecordsSet = []*Records{&records}
	}
	set := frb.RecordsSet[0].MsgSet
	set.Messages = append(set.Messages, msgBlock)
}

func (r *FetchResponse) AddRecordWithTimestamp(topic string, partition int32, key, value Encoder, offset int64, timestamp time.Time) {
	frb := r.getOrCreateBlock(topic, partition)
	kb, vb := encodeKV(key, value)
	if len(frb.RecordsSet) == 0 {
		records := newDefaultRecords(&RecordBatch{Version: 2, LogAppendTime: r.LogAppendTime, FirstTimestamp: timestamp, MaxTimestamp: r.Timestamp})
		frb.RecordsSet = []*Records{&records}
	}
	batch := frb.RecordsSet[0].RecordBatch
	rec := &Record{Key: kb, Value: vb, OffsetDelta: offset, TimestampDelta: timestamp.Sub(batch.FirstTimestamp)}
	batch.addRecord(rec)
}

// AddRecordBatchWithTimestamp is similar to AddRecordWithTimestamp
// But instead of appending 1 record to a batch, it append a new batch containing 1 record to the fetchResponse
// Since transaction are handled on batch level (the whole batch is either committed or aborted), use this to test transactions
func (r *FetchResponse) AddRecordBatchWithTimestamp(topic string, partition int32, key, value Encoder, offset int64, producerID int64, isTransactional bool, timestamp time.Time) {
	frb := r.getOrCreateBlock(topic, partition)
	kb, vb := encodeKV(key, value)

	records := newDefaultRecords(&RecordBatch{Version: 2, LogAppendTime: r.LogAppendTime, FirstTimestamp: timestamp, MaxTimestamp: r.Timestamp})
	batch := &RecordBatch{
		Version:         2,
		LogAppendTime:   r.LogAppendTime,
		FirstTimestamp:  timestamp,
		MaxTimestamp:    r.Timestamp,
		FirstOffset:     offset,
		LastOffsetDelta: 0,
		ProducerID:      producerID,
		IsTransactional: isTransactional,
	}
	rec := &Record{Key: kb, Value: vb, OffsetDelta: 0, TimestampDelta: timestamp.Sub(batch.FirstTimestamp)}
	batch.addRecord(rec)
	records.RecordBatch = batch

	frb.RecordsSet = append(frb.RecordsSet, &records)
}

func (r *FetchResponse) AddControlRecordWithTimestamp(topic string, partition int32, offset int64, producerID int64, recordType ControlRecordType, timestamp time.Time) {
	frb := r.getOrCreateBlock(topic, partition)

	// batch
	batch := &RecordBatch{
		Version:         2,
		LogAppendTime:   r.LogAppendTime,
		FirstTimestamp:  timestamp,
		MaxTimestamp:    r.Timestamp,
		FirstOffset:     offset,
		LastOffsetDelta: 0,
		ProducerID:      producerID,
		IsTransactional: true,
		Control:         true,
	}

	// records
	records := newDefaultRecords(nil)
	records.RecordBatch = batch

	// record
	crAbort := ControlRecord{
		Version: 0,
		Type:    recordType,
	}
	crKey := &realEncoder{raw: make([]byte, 4)}
	crValue := &realEncoder{raw: make([]byte, 6)}
	crAbort.encode(crKey, crValue)
	rec := &Record{Key: ByteEncoder(crKey.raw), Value: ByteEncoder(crValue.raw), OffsetDelta: 0, TimestampDelta: timestamp.Sub(batch.FirstTimestamp)}
	batch.addRecord(rec)

	frb.RecordsSet = append(frb.RecordsSet, &records)
}

func (r *FetchResponse) AddMessage(topic string, partition int32, key, value Encoder, offset int64) {
	r.AddMessageWithTimestamp(topic, partition, key, value, offset, time.Time{}, 0)
}

func (r *FetchResponse) AddRecord(topic string, partition int32, key, value Encoder, offset int64) {
	r.AddRecordWithTimestamp(topic, partition, key, value, offset, time.Time{})
}

func (r *FetchResponse) AddRecordBatch(topic string, partition int32, key, value Encoder, offset int64, producerID int64, isTransactional bool) {
	r.AddRecordBatchWithTimestamp(topic, partition, key, value, offset, producerID, isTransactional, time.Time{})
}

func (r *FetchResponse) AddControlRecord(topic string, partition int32, offset int64, producerID int64, recordType ControlRecordType) {
	// define controlRecord key and value
	r.AddControlRecordWithTimestamp(topic, partition, offset, producerID, recordType, time.Time{})
}

func (r *FetchResponse) SetLastOffsetDelta(topic string, partition int32, offset int32) {
	frb := r.getOrCreateBlock(topic, partition)
	if len(frb.RecordsSet) == 0 {
		records := newDefaultRecords(&RecordBatch{Version: 2})
		frb.RecordsSet = []*Records{&records}
	}
	batch := frb.RecordsSet[0].RecordBatch
	batch.LastOffsetDelta = offset
}

func (r *FetchResponse) SetLastStableOffset(topic string, partition int32, offset int64) {
	frb := r.getOrCreateBlock(topic, partition)
	frb.LastStableOffset = offset
}
