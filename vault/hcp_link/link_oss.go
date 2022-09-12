//go:build !enterprise

package hcp_link

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func NewHCPLink(linkConf *configutil.HCPLinkConfig, core *vault.Core, logger hclog.Logger) (*WrappedHCPLinkVault, error) {
	return nil, nil
}

func (h *WrappedHCPLinkVault) Shutdown() error {
	return nil
}

func (h *WrappedHCPLinkVault) GetScadaSessionStatus() string { return Disconnected }

func (h *WrappedHCPLinkVault) GetConnectionStatusMessage(scadaStatus string) string {
	return scadaStatus
}
