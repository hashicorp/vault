package abstractions

//ResponseHeaders represents a collection of response headers
type ResponseHeaders struct {
	header
}

//NewResponseHeaders creates a new ResponseHeaders
func NewResponseHeaders() *ResponseHeaders {
	return &ResponseHeaders{
		header{make(map[string]map[string]struct{})},
	}
}

//AddAll adds all headers from the other headers
func (r *ResponseHeaders) AddAll(other *ResponseHeaders) {
	if other == nil || other.headers == nil {
		return
	}
	for k, v := range other.headers {
		for k2 := range v {
			r.Add(k, k2)
		}
	}
}
