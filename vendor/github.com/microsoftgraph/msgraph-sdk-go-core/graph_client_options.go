package msgraphgocore

// GraphClientOptions represents a combination of GraphServiceVersion and GraphServiceLibraryVersion
//
// GraphServiceVersion is version of the targeted service.
// GraphServiceLibraryVersion is the version of the service library
type GraphClientOptions struct {
	GraphServiceVersion        string
	GraphServiceLibraryVersion string
}
