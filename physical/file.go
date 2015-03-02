package physical

// FileBackend is a physical backend that stores data on disk
// at a given file path. It can be used for durable single server
// situations, or to develop locally where durability is not critical.
type FileBackend struct {
}

// NewFileBackend constructs a Filebackend using the given directory
func NewFileBackend(dir string) (*FileBackend, error) {
	// TODO:
	f := &FileBackend{}
	return f, nil
}
