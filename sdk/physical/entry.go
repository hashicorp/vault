package physical

import (
	"encoding/hex"
	"fmt"
)

// Entry is used to represent data stored by the physical backend
type Entry struct {
	Key      string
	Value    []byte
	SealWrap bool `json:"seal_wrap,omitempty"`

	// Only used in replication
	ValueHash []byte
}

func (e *Entry) String() string {
	return fmt.Sprintf("Key: %s. SealWrap: %t. Value: %s. ValueHash: %s", e.Key, e.SealWrap, hex.EncodeToString(e.Value), hex.EncodeToString(e.ValueHash))
}
