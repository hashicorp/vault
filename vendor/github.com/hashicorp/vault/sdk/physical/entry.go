package physical

// Entry is used to represent data stored by the physical backend
type Entry struct {
	Key      string
	Value    []byte
	SealWrap bool `json:"seal_wrap,omitempty"`

	// Only used in replication
	ValueHash []byte

	// The bool above is an easy control for whether it should be enabled; it
	// is used to carry information about whether seal wrapping is *desired*
	// regardless of whether it's currently available. The struct below stores
	// needed information when it's actually performed.
	SealWrapInfo *EncryptedBlobInfo `json:"seal_wrap_info,omitempty"`
}
