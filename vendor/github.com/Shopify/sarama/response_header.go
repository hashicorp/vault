package sarama

import "fmt"

const responseLengthSize = 4
const correlationIDSize = 4

type responseHeader struct {
	length        int32
	correlationID int32
}

func (r *responseHeader) decode(pd packetDecoder, version int16) (err error) {
	r.length, err = pd.getInt32()
	if err != nil {
		return err
	}
	if r.length <= 4 || r.length > MaxResponseSize {
		return PacketDecodingError{fmt.Sprintf("message of length %d too large or too small", r.length)}
	}

	r.correlationID, err = pd.getInt32()

	if version >= 1 {
		if _, err := pd.getEmptyTaggedFieldArray(); err != nil {
			return err
		}
	}

	return err
}
