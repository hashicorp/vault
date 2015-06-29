package audit

import (
	"io"

	"github.com/hashicorp/vault/logical"
)

// Formatter is an interface that is responsible for formating a
// request/response into some format. Formatters write their output
// to an io.Writer.
//
// It is recommended that you pass data through Hash prior to formatting it.
type Formatter interface {
	FormatRequest(io.Writer, *logical.Auth, *logical.Request, error) error
	FormatResponse(io.Writer, *logical.Auth, *logical.Request, *logical.Response, error) error
}
