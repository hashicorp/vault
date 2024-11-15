package gocb

type internalProvider interface {
	GetNodesMetadata(opts *GetNodesMetadataOptions) ([]NodeMetadata, error)
}
