package xmlrpc

import (
	"fmt"
)

// xmlrpcError represents errors returned on xmlrpc request.
type XmlRpcError struct {
	Code           interface{}
	Err            string
	HttpStatusCode int
}

// Error() method implements Error interface
func (e *XmlRpcError) Error() string {
	return fmt.Sprintf(
		"error: %s, code: %v, http status code: %d",
		e.Err, e.Code, e.HttpStatusCode)
}

// Base64 represents value in base64 encoding
type Base64 string

type Params struct {
	Params []interface{}
}

type Struct map[string]interface{}
