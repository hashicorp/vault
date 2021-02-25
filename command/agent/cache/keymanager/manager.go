package keymanager

import wrapping "github.com/hashicorp/go-kms-wrapping"

const (
	KeyID = "root"
)

type KeyManager interface {
	// Returns a wrapping.Wrapper which can be used to perform key-related operations.
	Wrapper() wrapping.Wrapper
	// RetrievalToken is the material returned which can be used to source back the
	// encryption key. Depending on the implementation, the token can be the
	// encryption key itself or a token/identifier used to exchange the token.
	RetrievalToken() ([]byte, error)
}
