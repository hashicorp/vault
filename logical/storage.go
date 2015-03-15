package logical

// Storage is the way that logical backends are able read/write data.
type Storage interface {
	List(prefix string) ([]string, error)
	Get(string) (*StorageEntry, error)
	Put(*StorageEntry) error
	Delete(string) error
}

// StorageEntry is the entry for an item in a Storage implementation.
type StorageEntry struct {
	Key   string
	Value []byte
}
