package abstractions

// BaseRequestBuilder is the base class for all request builders.
type BaseRequestBuilder struct {
	// Path parameters for the request
	PathParameters map[string]string
	// The request adapter to use to execute the requests.
	RequestAdapter RequestAdapter
	// Url template to use to build the URL for the current request builder
	UrlTemplate string
}

// NewBaseRequestBuilder creates a new BaseRequestBuilder instance.
func NewBaseRequestBuilder(requestAdapter RequestAdapter, urlTemplate string, pathParameters map[string]string) *BaseRequestBuilder {
	if requestAdapter == nil {
		panic("requestAdapter cannot be nil")
	}
	pathParametersCopy := make(map[string]string)
	if pathParameters != nil {
		for idx, item := range pathParameters {
			pathParametersCopy[idx] = item
		}
	}
	return &BaseRequestBuilder{
		RequestAdapter: requestAdapter,
		UrlTemplate:    urlTemplate,
		PathParameters: pathParametersCopy,
	}
}
