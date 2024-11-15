package workload

import "fmt"

// CredentialTokenSource sources credentials via a directly supplied token.
type CredentialTokenSource struct {
	Token string `json:"token,omitempty"`
}

// Validate validates the config.
func (ct *CredentialTokenSource) Validate() error {
	if ct.Token == "" {
		return fmt.Errorf("token must be set")
	}

	return nil
}

// token returns the token.
func (ct *CredentialTokenSource) token() (string, error) {
	return ct.Token, nil
}
