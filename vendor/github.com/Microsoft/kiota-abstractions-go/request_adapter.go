// Package abstractions provides the base infrastructure for the Kiota-generated SDKs to function.
// It defines multiple concepts related to abstract HTTP requests, serialization, and authentication.
// These concepts can then be implemented independently without tying the SDKs to any specific implementation.
// Kiota also provides default implementations for these concepts.
// Checkout:
// - github.com/microsoft/kiota/authentication/go/azure
// - github.com/microsoft/kiota/http/go/nethttp
// - github.com/microsoft/kiota/serialization/go/json
package abstractions

import (
	"context"
	"github.com/microsoft/kiota-abstractions-go/store"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ErrorMappings is a mapping of status codes to error types factories.
type ErrorMappings map[string]s.ParsableFactory

// RequestAdapter is the service responsible for translating abstract RequestInformation into native HTTP requests.
type RequestAdapter interface {
	// Send executes the HTTP request specified by the given RequestInformation and returns the deserialized response model.
	Send(context context.Context, requestInfo *RequestInformation, constructor s.ParsableFactory, errorMappings ErrorMappings) (s.Parsable, error)
	// SendEnum executes the HTTP request specified by the given RequestInformation and returns the deserialized response model.
	SendEnum(context context.Context, requestInfo *RequestInformation, parser s.EnumFactory, errorMappings ErrorMappings) (any, error)
	// SendCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized response model collection.
	SendCollection(context context.Context, requestInfo *RequestInformation, constructor s.ParsableFactory, errorMappings ErrorMappings) ([]s.Parsable, error)
	// SendEnumCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized response model collection.
	SendEnumCollection(context context.Context, requestInfo *RequestInformation, parser s.EnumFactory, errorMappings ErrorMappings) ([]any, error)
	// SendPrimitive executes the HTTP request specified by the given RequestInformation and returns the deserialized primitive response model.
	SendPrimitive(context context.Context, requestInfo *RequestInformation, typeName string, errorMappings ErrorMappings) (any, error)
	// SendPrimitiveCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized primitive response model collection.
	SendPrimitiveCollection(context context.Context, requestInfo *RequestInformation, typeName string, errorMappings ErrorMappings) ([]any, error)
	// SendNoContent executes the HTTP request specified by the given RequestInformation with no return content.
	SendNoContent(context context.Context, requestInfo *RequestInformation, errorMappings ErrorMappings) error
	// GetSerializationWriterFactory returns the serialization writer factory currently in use for the request adapter service.
	GetSerializationWriterFactory() s.SerializationWriterFactory
	// EnableBackingStore enables the backing store proxies for the SerializationWriters and ParseNodes in use.
	EnableBackingStore(factory store.BackingStoreFactory)
	// SetBaseUrl sets the base url for every request.
	SetBaseUrl(baseUrl string)
	// GetBaseUrl gets the base url for every request.
	GetBaseUrl() string
	// ConvertToNativeRequest converts the given RequestInformation into a native HTTP request.
	ConvertToNativeRequest(context context.Context, requestInfo *RequestInformation) (any, error)
}
