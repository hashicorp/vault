//go:build !enterprise

package hcp_link

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func NewHCPLink(linkConf *configutil.HCPLinkConfig, core *vault.Core, logger hclog.Logger) (*HCPLinkVault, error) {
	return nil, nil
}

func (h *HCPLinkVault) Shutdown() error {
	return nil
}

func (h *HCPLinkVault) GetScadaSessionStatus() string { return Disconnected }
