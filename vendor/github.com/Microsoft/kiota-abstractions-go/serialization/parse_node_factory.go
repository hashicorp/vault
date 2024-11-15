package serialization

// ParseNodeFactory defines the contract for a factory that creates new ParseNode.
type ParseNodeFactory interface {
	// GetValidContentType returns the content type this factory's parse nodes can deserialize.
	GetValidContentType() (string, error)
	// GetRootParseNode return a new ParseNode instance that is the root of the content
	GetRootParseNode(contentType string, content []byte) (ParseNode, error)
}
