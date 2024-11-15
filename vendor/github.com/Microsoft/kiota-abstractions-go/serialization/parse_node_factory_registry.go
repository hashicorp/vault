package serialization

import (
	"errors"
	re "regexp"
	"strings"
	"sync"
)

// ParseNodeFactoryRegistry holds a list of all the registered factories for the various types of nodes.
type ParseNodeFactoryRegistry struct {
	lock *sync.Mutex

	// ContentTypeAssociatedFactories maps content types onto the relevant factory.
	//
	// When interacting with this field, please make use of Lock and Unlock methods to ensure thread safety.
	ContentTypeAssociatedFactories map[string]ParseNodeFactory
}

func NewParseNodeFactoryRegistry() *ParseNodeFactoryRegistry {
	return &ParseNodeFactoryRegistry{
		lock:                           &sync.Mutex{},
		ContentTypeAssociatedFactories: make(map[string]ParseNodeFactory),
	}
}

// DefaultParseNodeFactoryInstance is the default singleton instance of the registry to be used when registering new factories that should be available by default.
var DefaultParseNodeFactoryInstance = NewParseNodeFactoryRegistry()

// GetValidContentType returns the valid content type for the ParseNodeFactoryRegistry
func (m *ParseNodeFactoryRegistry) GetValidContentType() (string, error) {
	return "", errors.New("the registry supports multiple content types. Get the registered factory instead")
}

var contentTypeVendorCleanupPattern = re.MustCompile("[^/]+\\+")

// GetRootParseNode returns a new ParseNode instance that is the root of the content
func (m *ParseNodeFactoryRegistry) GetRootParseNode(contentType string, content []byte) (ParseNode, error) {
	if contentType == "" {
		return nil, errors.New("contentType is required")
	}
	if content == nil {
		return nil, errors.New("content is required")
	}
	vendorSpecificContentType := strings.Split(contentType, ";")[0]
	factory, ok := m.ContentTypeAssociatedFactories[vendorSpecificContentType]
	if ok {
		return factory.GetRootParseNode(vendorSpecificContentType, content)
	}
	cleanedContentType := contentTypeVendorCleanupPattern.ReplaceAllString(vendorSpecificContentType, "")
	factory, ok = m.ContentTypeAssociatedFactories[cleanedContentType]
	if ok {
		return factory.GetRootParseNode(cleanedContentType, content)
	}
	return nil, errors.New("content type " + cleanedContentType + " does not have a factory registered to be parsed")
}

func (m *ParseNodeFactoryRegistry) Lock() {
	m.lock.Lock()
}

func (m *ParseNodeFactoryRegistry) Unlock() {
	m.lock.Unlock()
}
