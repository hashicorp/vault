package command

import (
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/meta"
)

// DefaultTokenHelper returns the token helper that is configured for Vault.
func DefaultTokenHelper(m *meta.Meta) (token.TokenHelper, error) {
	config, err := m.Config()
	if err != nil {
		return nil, err
	}

	path := config.TokenHelper
	if path == "" {
		return &token.InternalTokenHelper{}, nil
	}

	path, err = token.ExternalTokenHelperPath(path)
	if err != nil {
		return nil, err
	}
	return &token.ExternalTokenHelper{BinaryPath: path}, nil
}
