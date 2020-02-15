package physical

// Entry is used to represent data stored by the physical backend
type Entry struct {
	Key      string
	Value    []byte
	SealWrap bool `json:"seal_wrap,omitempty"`

	// Only used in replication
	ValueHash []byte
}
