package msgraphgocore

import (
	"errors"
	nethttp "net/http"

	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	khttp "github.com/microsoft/kiota-http-go"
)

// GraphRequestAdapterBase is the core service used by GraphServiceClient to make requests to Microsoft Graph.
type GraphRequestAdapterBase struct {
	khttp.NetHttpRequestAdapter
}

// NewGraphRequestAdapterBase creates a new GraphRequestAdapterBase with the given parameters
func NewGraphRequestAdapterBase(authenticationProvider absauth.AuthenticationProvider, clientOptions GraphClientOptions) (*GraphRequestAdapterBase, error) {
	return NewGraphRequestAdapterBaseWithParseNodeFactory(authenticationProvider, clientOptions, nil)
}

// NewGraphRequestAdapterBaseWithParseNodeFactory creates a new GraphRequestAdapterBase with the given parameters
func NewGraphRequestAdapterBaseWithParseNodeFactory(authenticationProvider absauth.AuthenticationProvider, clientOptions GraphClientOptions, parseNodeFactory absser.ParseNodeFactory) (*GraphRequestAdapterBase, error) {
	return NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider, clientOptions, parseNodeFactory, nil)
}

// NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactory creates a new GraphRequestAdapterBase with the given parameters
func NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider absauth.AuthenticationProvider, clientOptions GraphClientOptions, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory) (*GraphRequestAdapterBase, error) {
	return NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider, clientOptions, parseNodeFactory, serializationWriterFactory, nil)
}

// NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient creates a new GraphRequestAdapterBase with the given parameters
func NewGraphRequestAdapterBaseWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider absauth.AuthenticationProvider, clientOptions GraphClientOptions, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory, httpClient *nethttp.Client) (*GraphRequestAdapterBase, error) {
	if authenticationProvider == nil {
		return nil, errors.New("authenticationProvider cannot be nil")
	}
	if httpClient == nil {
		httpClient = GetDefaultClient(&clientOptions)
	}
	if serializationWriterFactory == nil {
		serializationWriterFactory = absser.DefaultSerializationWriterFactoryInstance
	}
	if parseNodeFactory == nil {
		parseNodeFactory = absser.DefaultParseNodeFactoryInstance
	}
	baseAdapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider, parseNodeFactory, serializationWriterFactory, httpClient)
	if err != nil {
		return nil, err
	}
	result := &GraphRequestAdapterBase{
		NetHttpRequestAdapter: *baseAdapter,
	}

	return result, nil
}
