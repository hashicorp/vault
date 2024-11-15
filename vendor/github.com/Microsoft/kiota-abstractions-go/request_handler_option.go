package abstractions

// RequestHandlerOption represents an abstract provider for ResponseHandler
type RequestHandlerOption interface {
	GetResponseHandler() ResponseHandler
	SetResponseHandler(responseHandler ResponseHandler)
	GetKey() RequestOptionKey
}

var ResponseHandlerOptionKey = RequestOptionKey{
	Key: "ResponseHandlerOptionKey",
}

type requestHandlerOption struct {
	handler ResponseHandler
}

// NewRequestHandlerOption creates a new RequestInformation object with default values.
func NewRequestHandlerOption() RequestHandlerOption {
	return &requestHandlerOption{}
}

func (r *requestHandlerOption) GetResponseHandler() ResponseHandler {
	return r.handler
}

func (r *requestHandlerOption) SetResponseHandler(responseHandler ResponseHandler) {
	r.handler = responseHandler
}

func (r *requestHandlerOption) GetKey() RequestOptionKey {
	return ResponseHandlerOptionKey
}
