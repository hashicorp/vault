package sarama

//ControlRecordType ...
type ControlRecordType int

const (
	//ControlRecordAbort is a control record for abort
	ControlRecordAbort ControlRecordType = iota
	//ControlRecordCommit is a control record for commit
	ControlRecordCommit
	//ControlRecordUnknown is a control record of unknown type
	ControlRecordUnknown
)

// Control records are returned as a record by fetchRequest
// However unlike "normal" records, they mean nothing application wise.
// They only serve internal logic for supporting transactions.
type ControlRecord struct {
	Version          int16
	CoordinatorEpoch int32
	Type             ControlRecordType
}

func (cr *ControlRecord) decode(key, value packetDecoder) error {
	var err error
	cr.Version, err = value.getInt16()
	if err != nil {
		return err
	}

	cr.CoordinatorEpoch, err = value.getInt32()
	if err != nil {
		return err
	}

	// There a version for the value part AND the key part. And I have no idea if they are supposed to match or not
	// Either way, all these version can only be 0 for now
	cr.Version, err = key.getInt16()
	if err != nil {
		return err
	}

	recordType, err := key.getInt16()
	if err != nil {
		return err
	}

	switch recordType {
	case 0:
		cr.Type = ControlRecordAbort
	case 1:
		cr.Type = ControlRecordCommit
	default:
		// from JAVA implementation:
		// UNKNOWN is used to indicate a control type which the client is not aware of and should be ignored
		cr.Type = ControlRecordUnknown
	}
	return nil
}

func (cr *ControlRecord) encode(key, value packetEncoder) {
	value.putInt16(cr.Version)
	value.putInt32(cr.CoordinatorEpoch)
	key.putInt16(cr.Version)

	switch cr.Type {
	case ControlRecordAbort:
		key.putInt16(0)
	case ControlRecordCommit:
		key.putInt16(1)
	}
}
