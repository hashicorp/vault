package hcp_link

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	link "github.com/hashicorp/hcp-link"
	linkConfig "github.com/hashicorp/hcp-link/pkg/config"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
)

// SessionStatus is used to express the current status of the SCADA session.
type SessionStatus = string

const (
	// Connected HCP link connection status when it is connected
	Connected = SessionStatus("connected")
	// Disconnected HCP link connection status when it is disconnected
	Disconnected = SessionStatus("disconnected")
)

type HCPLinkVault struct {
	l            sync.Mutex
	LinkStatus   internal.WrappedCoreHCPLinkStatus
	scadaConfig  *scada.Config
	linkConfig   *linkConfig.Config
	link         link.Link
	logger       hclog.Logger
	capabilities map[string]capabilities.Capability
	stopCh       chan struct{}
	running      bool

	// TODO: remove after testing
	URL string
}
