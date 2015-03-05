package physical

// FileBackend is a physical backend that stores data on disk
// at a given file path. It can be used for durable single server
// situations, or to develop locally where durability is not critical.
type FileBackend struct {
}

// newFileBackend constructs a Filebackend using the given directory
func newFileBackend(conf map[string]string) (Backend, error) {
	// TODO:
	return nil, nil
}
