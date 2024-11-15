package nethttplibrary

import (
	"bytes"
	"context"
	"errors"
	"io"
	nethttp "net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	abs "github.com/microsoft/kiota-abstractions-go"
	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// nopCloser is an alternate io.nopCloser implementation which
// provides io.ReadSeekCloser instead of io.ReadCloser as we need
// Seek for retries
type nopCloser struct {
	io.ReadSeeker
}

func NopCloser(r io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{r}
}

func (nopCloser) Close() error { return nil }

// NetHttpRequestAdapter implements the RequestAdapter interface using net/http
type NetHttpRequestAdapter struct {
	// serializationWriterFactory is the factory used to create serialization writers
	serializationWriterFactory absser.SerializationWriterFactory
	// parseNodeFactory is the factory used to create parse nodes
	parseNodeFactory absser.ParseNodeFactory
	// httpClient is the client used to send requests
	httpClient *nethttp.Client
	// authenticationProvider is the provider used to authenticate requests
	authenticationProvider absauth.AuthenticationProvider
	// The base url for every request.
	baseUrl string
	// The observation options for the request adapter.
	observabilityOptions ObservabilityOptions
}

// NewNetHttpRequestAdapter creates a new NetHttpRequestAdapter with the given parameters
func NewNetHttpRequestAdapter(authenticationProvider absauth.AuthenticationProvider) (*NetHttpRequestAdapter, error) {
	return NewNetHttpRequestAdapterWithParseNodeFactory(authenticationProvider, nil)
}

// NewNetHttpRequestAdapterWithParseNodeFactory creates a new NetHttpRequestAdapter with the given parameters
func NewNetHttpRequestAdapterWithParseNodeFactory(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory) (*NetHttpRequestAdapter, error) {
	return NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider, parseNodeFactory, nil)
}

// NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory creates a new NetHttpRequestAdapter with the given parameters
func NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactory(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory) (*NetHttpRequestAdapter, error) {
	return NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider, parseNodeFactory, serializationWriterFactory, nil)
}

// NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient creates a new NetHttpRequestAdapter with the given parameters
func NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory, httpClient *nethttp.Client) (*NetHttpRequestAdapter, error) {
	return NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClientAndObservabilityOptions(authenticationProvider, parseNodeFactory, serializationWriterFactory, httpClient, ObservabilityOptions{})
}

// NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClientAndObservabilityOptions creates a new NetHttpRequestAdapter with the given parameters
func NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClientAndObservabilityOptions(authenticationProvider absauth.AuthenticationProvider, parseNodeFactory absser.ParseNodeFactory, serializationWriterFactory absser.SerializationWriterFactory, httpClient *nethttp.Client, observabilityOptions ObservabilityOptions) (*NetHttpRequestAdapter, error) {
	if authenticationProvider == nil {
		return nil, errors.New("authenticationProvider cannot be nil")
	}
	result := &NetHttpRequestAdapter{
		serializationWriterFactory: serializationWriterFactory,
		parseNodeFactory:           parseNodeFactory,
		httpClient:                 httpClient,
		authenticationProvider:     authenticationProvider,
		baseUrl:                    "",
		observabilityOptions:       observabilityOptions,
	}
	if result.httpClient == nil {
		defaultClient := GetDefaultClient()
		result.httpClient = defaultClient
	}
	if result.serializationWriterFactory == nil {
		result.serializationWriterFactory = absser.DefaultSerializationWriterFactoryInstance
	}
	if result.parseNodeFactory == nil {
		result.parseNodeFactory = absser.DefaultParseNodeFactoryInstance
	}
	return result, nil
}

// GetSerializationWriterFactory returns the serialization writer factory currently in use for the request adapter service.
func (a *NetHttpRequestAdapter) GetSerializationWriterFactory() absser.SerializationWriterFactory {
	return a.serializationWriterFactory
}

// EnableBackingStore enables the backing store proxies for the SerializationWriters and ParseNodes in use.
func (a *NetHttpRequestAdapter) EnableBackingStore(factory store.BackingStoreFactory) {
	a.parseNodeFactory = abs.EnableBackingStoreForParseNodeFactory(a.parseNodeFactory)
	a.serializationWriterFactory = abs.EnableBackingStoreForSerializationWriterFactory(a.serializationWriterFactory)
	if factory != nil {
		store.BackingStoreFactoryInstance = factory
	}
}

// SetBaseUrl sets the base url for every request.
func (a *NetHttpRequestAdapter) SetBaseUrl(baseUrl string) {
	a.baseUrl = baseUrl
}

// GetBaseUrl gets the base url for every request.
func (a *NetHttpRequestAdapter) GetBaseUrl() string {
	return a.baseUrl
}

func (a *NetHttpRequestAdapter) getHttpResponseMessage(ctx context.Context, requestInfo *abs.RequestInformation, claims string, spanForAttributes trace.Span) (*nethttp.Response, error) {
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "getHttpResponseMessage")
	defer span.End()
	if ctx == nil {
		ctx = context.Background()
	}
	a.setBaseUrlForRequestInformation(requestInfo)
	additionalContext := make(map[string]any)
	if claims != "" {
		additionalContext[claimsKey] = claims
	}
	err := a.authenticationProvider.AuthenticateRequest(ctx, requestInfo, additionalContext)
	if err != nil {
		return nil, err
	}
	request, err := a.getRequestFromRequestInformation(ctx, requestInfo, spanForAttributes)
	if err != nil {
		return nil, err
	}
	response, err := (*a.httpClient).Do(request)
	if err != nil {
		spanForAttributes.RecordError(err)
		return nil, err
	}
	if response != nil {
		contentLenHeader := response.Header.Get("Content-Length")
		if contentLenHeader != "" {
			contentLen, _ := strconv.Atoi(contentLenHeader)
			spanForAttributes.SetAttributes(attribute.Int("http.response_content_length", contentLen))
		}
		contentTypeHeader := response.Header.Get("Content-Type")
		if contentTypeHeader != "" {
			spanForAttributes.SetAttributes(attribute.String("http.response_content_type", contentTypeHeader))
		}
		spanForAttributes.SetAttributes(
			attribute.Int("http.status_code", response.StatusCode),
			attribute.String("http.flavor", response.Proto),
		)
	}
	return a.retryCAEResponseIfRequired(ctx, response, requestInfo, claims, spanForAttributes)
}

const claimsKey = "claims"

var reBearer = regexp.MustCompile(`(?i)^Bearer\s`)
var reClaims = regexp.MustCompile(`\"([^\"]*)\"`)

const AuthenticateChallengedEventKey = "com.microsoft.kiota.authenticate_challenge_received"

func (a *NetHttpRequestAdapter) retryCAEResponseIfRequired(ctx context.Context, response *nethttp.Response, requestInfo *abs.RequestInformation, claims string, spanForAttributes trace.Span) (*nethttp.Response, error) {
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "retryCAEResponseIfRequired")
	defer span.End()
	if response.StatusCode == 401 &&
		claims == "" { //avoid infinite loop, we only retry once
		authenticateHeaderVal := response.Header.Get("WWW-Authenticate")
		if authenticateHeaderVal != "" && reBearer.Match([]byte(authenticateHeaderVal)) {
			span.AddEvent(AuthenticateChallengedEventKey)
			spanForAttributes.SetAttributes(attribute.Int("http.retry_count", 1))
			responseClaims := ""
			parametersRaw := string(reBearer.ReplaceAll([]byte(authenticateHeaderVal), []byte("")))
			parameters := strings.Split(parametersRaw, ",")
			for _, parameter := range parameters {
				if strings.HasPrefix(strings.Trim(parameter, " "), claimsKey) {
					responseClaims = reClaims.FindStringSubmatch(parameter)[1]
					break
				}
			}
			if responseClaims != "" {
				defer a.purge(response)
				return a.getHttpResponseMessage(ctx, requestInfo, responseClaims, spanForAttributes)
			}
		}
	}
	return response, nil
}

func (a *NetHttpRequestAdapter) getResponsePrimaryContentType(response *nethttp.Response) string {
	if response.Header == nil {
		return ""
	}
	rawType := response.Header.Get("Content-Type")
	splat := strings.Split(rawType, ";")
	return strings.ToLower(splat[0])
}

func (a *NetHttpRequestAdapter) setBaseUrlForRequestInformation(requestInfo *abs.RequestInformation) {
	requestInfo.PathParameters["baseurl"] = a.GetBaseUrl()
}

func (a *NetHttpRequestAdapter) prepareContext(ctx context.Context, requestInfo *abs.RequestInformation) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	// set deadline if not set in receiving context
	// ignore if timeout is 0 as it means no timeout
	if _, deadlineSet := ctx.Deadline(); !deadlineSet && a.httpClient.Timeout != 0 {
		ctx, _ = context.WithTimeout(ctx, a.httpClient.Timeout)
	}

	for _, value := range requestInfo.GetRequestOptions() {
		ctx = context.WithValue(ctx, value.GetKey(), value)
	}
	obsOptionsSet := false
	if reqObsOpt := ctx.Value(observabilityOptionsKeyValue); reqObsOpt != nil {
		if _, ok := reqObsOpt.(ObservabilityOptionsInt); ok {
			obsOptionsSet = true
		}
	}
	if !obsOptionsSet {
		ctx = context.WithValue(ctx, observabilityOptionsKeyValue, &a.observabilityOptions)
	}
	return ctx
}

// ConvertToNativeRequest converts the given RequestInformation into a native HTTP request.
func (a *NetHttpRequestAdapter) ConvertToNativeRequest(context context.Context, requestInfo *abs.RequestInformation) (any, error) {
	err := a.authenticationProvider.AuthenticateRequest(context, requestInfo, nil)
	if err != nil {
		return nil, err
	}
	request, err := a.getRequestFromRequestInformation(context, requestInfo, nil)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (a *NetHttpRequestAdapter) getRequestFromRequestInformation(ctx context.Context, requestInfo *abs.RequestInformation, spanForAttributes trace.Span) (*nethttp.Request, error) {
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "getRequestFromRequestInformation")
	defer span.End()
	if spanForAttributes == nil {
		spanForAttributes = span
	}
	spanForAttributes.SetAttributes(attribute.String("http.method", requestInfo.Method.String()))
	uri, err := requestInfo.GetUri()
	if err != nil {
		spanForAttributes.RecordError(err)
		return nil, err
	}
	spanForAttributes.SetAttributes(
		attribute.String("http.scheme", uri.Scheme),
		attribute.String("http.host", uri.Host),
	)

	if a.observabilityOptions.IncludeEUIIAttributes {
		spanForAttributes.SetAttributes(attribute.String("http.uri", uri.String()))
	}

	request, err := nethttp.NewRequestWithContext(ctx, requestInfo.Method.String(), uri.String(), nil)

	if err != nil {
		spanForAttributes.RecordError(err)
		return nil, err
	}
	if len(requestInfo.Content) > 0 {
		reader := bytes.NewReader(requestInfo.Content)
		request.Body = NopCloser(reader)
	}
	if request.Header == nil {
		request.Header = make(nethttp.Header)
	}
	if requestInfo.Headers != nil {
		for _, key := range requestInfo.Headers.ListKeys() {
			values := requestInfo.Headers.Get(key)
			for _, v := range values {
				request.Header.Add(key, v)
			}
		}
		if request.Header.Get("Content-Type") != "" {
			spanForAttributes.SetAttributes(
				attribute.String("http.request_content_type", request.Header.Get("Content-Type")),
			)
		}
		if request.Header.Get("Content-Length") != "" {
			contentLenVal, _ := strconv.Atoi(request.Header.Get("Content-Length"))
			request.ContentLength = int64(contentLenVal)
			spanForAttributes.SetAttributes(
				attribute.Int("http.request_content_length", contentLenVal),
			)
		}
	}

	return request, nil
}

const EventResponseHandlerInvokedKey = "com.microsoft.kiota.response_handler_invoked"

var queryParametersCleanupRegex = regexp.MustCompile(`\{\?[^\}]+}`)

func (a *NetHttpRequestAdapter) startTracingSpan(ctx context.Context, requestInfo *abs.RequestInformation, methodName string) (context.Context, trace.Span) {
	decodedUriTemplate := decodeUriEncodedString(requestInfo.UrlTemplate, []byte{'-', '.', '~', '$'})
	telemetryPathValue := queryParametersCleanupRegex.ReplaceAll([]byte(decodedUriTemplate), []byte(""))
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, methodName+" - "+string(telemetryPathValue))
	span.SetAttributes(attribute.String("http.uri_template", decodedUriTemplate))
	return ctx, span
}

// Send executes the HTTP request specified by the given RequestInformation and returns the deserialized response model.
func (a *NetHttpRequestAdapter) Send(ctx context.Context, requestInfo *abs.RequestInformation, constructor absser.ParsableFactory, errorMappings abs.ErrorMappings) (absser.Parsable, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "Send")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.(absser.Parsable), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetObjectValue")
		defer deserializeSpan.End()
		result, err := parseNode.GetObjectValue(constructor)
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

func (a *NetHttpRequestAdapter) setResponseType(result any, span trace.Span) {
	if result != nil {
		span.SetAttributes(attribute.String("com.microsoft.kiota.response.type", reflect.TypeOf(result).String()))
	}
}

// SendEnum executes the HTTP request specified by the given RequestInformation and returns the deserialized response model.
func (a *NetHttpRequestAdapter) SendEnum(ctx context.Context, requestInfo *abs.RequestInformation, parser absser.EnumFactory, errorMappings abs.ErrorMappings) (any, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendEnum")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.(absser.Parsable), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetEnumValue")
		defer deserializeSpan.End()
		result, err := parseNode.GetEnumValue(parser)
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

// SendCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized response model collection.
func (a *NetHttpRequestAdapter) SendCollection(ctx context.Context, requestInfo *abs.RequestInformation, constructor absser.ParsableFactory, errorMappings abs.ErrorMappings) ([]absser.Parsable, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendCollection")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.([]absser.Parsable), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetCollectionOfObjectValues")
		defer deserializeSpan.End()
		result, err := parseNode.GetCollectionOfObjectValues(constructor)
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

// SendEnumCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized response model collection.
func (a *NetHttpRequestAdapter) SendEnumCollection(ctx context.Context, requestInfo *abs.RequestInformation, parser absser.EnumFactory, errorMappings abs.ErrorMappings) ([]any, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendEnumCollection")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.([]any), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetCollectionOfEnumValues")
		defer deserializeSpan.End()
		result, err := parseNode.GetCollectionOfEnumValues(parser)
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

func getResponseHandler(ctx context.Context) abs.ResponseHandler {
	var handlerOption = ctx.Value(abs.ResponseHandlerOptionKey)
	if handlerOption != nil {
		return handlerOption.(abs.RequestHandlerOption).GetResponseHandler()
	}
	return nil
}

// SendPrimitive executes the HTTP request specified by the given RequestInformation and returns the deserialized primitive response model.
func (a *NetHttpRequestAdapter) SendPrimitive(ctx context.Context, requestInfo *abs.RequestInformation, typeName string, errorMappings abs.ErrorMappings) (any, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendPrimitive")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.(absser.Parsable), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		if typeName == "[]byte" {
			res, err := io.ReadAll(response.Body)
			if err != nil {
				span.RecordError(err)
				return nil, err
			} else if len(res) == 0 {
				return nil, nil
			}
			return res, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "Get"+typeName+"Value")
		defer deserializeSpan.End()
		var result any
		switch typeName {
		case "string":
			result, err = parseNode.GetStringValue()
		case "float32":
			result, err = parseNode.GetFloat32Value()
		case "float64":
			result, err = parseNode.GetFloat64Value()
		case "int32":
			result, err = parseNode.GetInt32Value()
		case "int64":
			result, err = parseNode.GetInt64Value()
		case "bool":
			result, err = parseNode.GetBoolValue()
		case "Time":
			result, err = parseNode.GetTimeValue()
		case "UUID":
			result, err = parseNode.GetUUIDValue()
		default:
			return nil, errors.New("unsupported type")
		}
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

// SendPrimitiveCollection executes the HTTP request specified by the given RequestInformation and returns the deserialized primitive response model collection.
func (a *NetHttpRequestAdapter) SendPrimitiveCollection(ctx context.Context, requestInfo *abs.RequestInformation, typeName string, errorMappings abs.ErrorMappings) ([]any, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendPrimitiveCollection")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return nil, err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		result, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return result.([]any), nil
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return nil, err
		}
		if a.shouldReturnNil(response) {
			return nil, nil
		}
		parseNode, _, err := a.getRootParseNode(ctx, response, span)
		if err != nil {
			return nil, err
		}
		if parseNode == nil {
			return nil, nil
		}
		_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetCollectionOfPrimitiveValues")
		defer deserializeSpan.End()
		result, err := parseNode.GetCollectionOfPrimitiveValues(typeName)
		a.setResponseType(result, span)
		if err != nil {
			span.RecordError(err)
		}
		return result, err
	} else {
		return nil, errors.New("response is nil")
	}
}

// SendNoContent executes the HTTP request specified by the given RequestInformation with no return content.
func (a *NetHttpRequestAdapter) SendNoContent(ctx context.Context, requestInfo *abs.RequestInformation, errorMappings abs.ErrorMappings) error {
	if requestInfo == nil {
		return errors.New("requestInfo cannot be nil")
	}
	ctx = a.prepareContext(ctx, requestInfo)
	ctx, span := a.startTracingSpan(ctx, requestInfo, "SendNoContent")
	defer span.End()
	response, err := a.getHttpResponseMessage(ctx, requestInfo, "", span)
	if err != nil {
		return err
	}

	responseHandler := getResponseHandler(ctx)
	if responseHandler != nil {
		span.AddEvent(EventResponseHandlerInvokedKey)
		_, err := responseHandler(response, errorMappings)
		if err != nil {
			span.RecordError(err)
		}
		return err
	} else if response != nil {
		defer a.purge(response)
		err = a.throwIfFailedResponse(ctx, response, errorMappings, span)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("response is nil")
	}
}

func (a *NetHttpRequestAdapter) getRootParseNode(ctx context.Context, response *nethttp.Response, spanForAttributes trace.Span) (absser.ParseNode, context.Context, error) {
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "getRootParseNode")
	defer span.End()

	if response.ContentLength == 0 {
		return nil, ctx, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		spanForAttributes.RecordError(err)
		return nil, ctx, err
	}
	contentType := a.getResponsePrimaryContentType(response)
	if contentType == "" {
		return nil, ctx, nil
	}
	rootNode, err := a.parseNodeFactory.GetRootParseNode(contentType, body)
	if err != nil {
		spanForAttributes.RecordError(err)
	}
	return rootNode, ctx, err
}
func (a *NetHttpRequestAdapter) purge(response *nethttp.Response) error {
	_, _ = io.ReadAll(response.Body) //we don't care about errors comming from reading the body, just trying to purge anything that maybe left
	err := response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
func (a *NetHttpRequestAdapter) shouldReturnNil(response *nethttp.Response) bool {
	return response.StatusCode == 204
}

// ErrorMappingFoundAttributeName is the attribute name used to indicate whether an error code mapping was found.
const ErrorMappingFoundAttributeName = "com.microsoft.kiota.error.mapping_found"

// ErrorBodyFoundAttributeName is the attribute name used to indicate whether the error response contained a body
const ErrorBodyFoundAttributeName = "com.microsoft.kiota.error.body_found"

func (a *NetHttpRequestAdapter) throwIfFailedResponse(ctx context.Context, response *nethttp.Response, errorMappings abs.ErrorMappings, spanForAttributes trace.Span) error {
	ctx, span := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "throwIfFailedResponse")
	defer span.End()
	if response.StatusCode < 400 {
		return nil
	}
	spanForAttributes.SetStatus(codes.Error, "received_error_response")

	statusAsString := strconv.Itoa(response.StatusCode)
	responseHeaders := abs.NewResponseHeaders()
	for key, values := range response.Header {
		for i := range values {
			responseHeaders.Add(key, values[i])
		}
	}
	var errorCtor absser.ParsableFactory = nil
	if len(errorMappings) != 0 {
		if errorMappings[statusAsString] != nil {
			errorCtor = errorMappings[statusAsString]
		} else if response.StatusCode >= 400 && response.StatusCode < 500 && errorMappings["4XX"] != nil {
			errorCtor = errorMappings["4XX"]
		} else if response.StatusCode >= 500 && response.StatusCode < 600 && errorMappings["5XX"] != nil {
			errorCtor = errorMappings["5XX"]
		} else if errorMappings["XXX"] != nil && response.StatusCode >= 400 && response.StatusCode < 600 {
			errorCtor = errorMappings["XXX"]
		}
	}

	if errorCtor == nil {
		spanForAttributes.SetAttributes(attribute.Bool(ErrorMappingFoundAttributeName, false))
		err := &abs.ApiError{
			Message:            "The server returned an unexpected status code and no error factory is registered for this code: " + statusAsString,
			ResponseStatusCode: response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
		spanForAttributes.RecordError(err)
		return err
	}
	spanForAttributes.SetAttributes(attribute.Bool(ErrorMappingFoundAttributeName, true))

	rootNode, _, err := a.getRootParseNode(ctx, response, spanForAttributes)
	if err != nil {
		spanForAttributes.RecordError(err)
		return err
	}
	if rootNode == nil {
		spanForAttributes.SetAttributes(attribute.Bool(ErrorBodyFoundAttributeName, false))
		err := &abs.ApiError{
			Message:            "The server returned an unexpected status code with no response body: " + statusAsString,
			ResponseStatusCode: response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
		spanForAttributes.RecordError(err)
		return err
	}
	spanForAttributes.SetAttributes(attribute.Bool(ErrorBodyFoundAttributeName, true))

	_, deserializeSpan := otel.GetTracerProvider().Tracer(a.observabilityOptions.GetTracerInstrumentationName()).Start(ctx, "GetObjectValue")
	defer deserializeSpan.End()
	errValue, err := rootNode.GetObjectValue(errorCtor)
	if err != nil {
		spanForAttributes.RecordError(err)
		if apiErrorable, ok := err.(abs.ApiErrorable); ok {
			apiErrorable.SetResponseHeaders(responseHeaders)
			apiErrorable.SetStatusCode(response.StatusCode)
		}
		return err
	} else if errValue == nil {
		return &abs.ApiError{
			Message:            "The server returned an unexpected status code but the error could not be deserialized: " + statusAsString,
			ResponseStatusCode: response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
	}

	if apiErrorable, ok := errValue.(abs.ApiErrorable); ok {
		apiErrorable.SetResponseHeaders(responseHeaders)
		apiErrorable.SetStatusCode(response.StatusCode)
	}

	err = errValue.(error)

	spanForAttributes.RecordError(err)
	return err
}
