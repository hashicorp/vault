package msgraphsdkgo

import (
	nethttp "net/http"

	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	core "github.com/microsoftgraph/msgraph-sdk-go-core"
)

var clientOptions = core.GraphClientOptions{
	GraphServiceVersion: "", //v1 doesn't include the service version in the telemetry header
	/** The SDK version */
	// x-release-please-start-version
	GraphServiceLibraryVersion: "1.51.0",
	// x-release-please-end
}

// GetDefaultClientOptions returns the default client options used by the GraphRequestAdapterBase and the middleware.
func GetDefaultClientOptions() core.GraphClientOptions {
	return clientOptions
}

// GraphRequestAdapter is the core service used by GraphBaseServiceClient to make requests to Microsoft Graph.
type GraphRequestAdapter struct {
	core.GraphRequestAdapterBase
}

// NewGraphRequestAdapter creates a new GraphRequestAdapter with the given parameters
// Parameters:
// authenticationProvider: the provider used to authenticate requests
// Returns:
// a new GraphRequestAdapter
func NewGraphRequestAdapter(authenticationProvider absauth.AuthenticationProvider) (*GraphRequestAdapter, error) {
	return NewGraphRequestAdapterWithParseNodeFactory(authenticationProvider, nil)
}

// NewGraphRequestAdapterWithParseNodeFactory creates a new GraphRequestAdapter with the given parameters
// Parameters:
// authenticationProvider: the provider used to authenticate requests
// parseNodeFactory: the factory used to create parse nodes
// Returns:
// a new GraphRequestAdapter
func NewGraphRequestAdapterWithParseNodeFactory(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory) (*GraphRequestAdapter, error) {
	return NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider, parseNodeFactory, nil)
}

// NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory creates a new GraphRequestAdapter with the given parameters
// Parameters:
// authenticationProvider: the provider used to authenticate requests
// parseNodeFactory: the factory used to create parse nodes
// serializationWriterFactory: the factory used to create serialization writers
// Returns:
// a new GraphRequestAdapter
func NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory) (*GraphRequestAdapter, error) {
	return NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider, parseNodeFactory, serializationWriterFactory, nil)
}

// NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient creates a new GraphRequestAdapter with the given parameters
// Parameters:
// authenticationProvider: the provider used to authenticate requests
// parseNodeFactory: the factory used to create parse nodes
// serializationWriterFactory: the factory used to create serialization writers
// httpClient: the client used to send requests
// Returns:
// a new GraphRequestAdapter
func NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory, httpClient *nethttp.Client) (*GraphRequestAdapter, error) {
	baseAdapter, err := core.NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider, clientOptions, parseNodeFactory, serializationWriterFactory, httpClient)
	if err != nil {
		return nil, err
	}
	result := &GraphRequestAdapter{
		GraphRequestAdapterBase: *baseAdapter,
	}

	return result, nil
}
