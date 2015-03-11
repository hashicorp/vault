package api

import (
	"encoding/json"
	"net/http"
)

// Response is a raw response that wraps an HTTP response.
type Response struct {
	*http.Response
}

// DecodeJSON will decode the response body to a JSON structure. This
// will consume the response body, but will not close it. Close must
// still be called.
func (r *Response) DecodeJSON(out interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(out)
}
