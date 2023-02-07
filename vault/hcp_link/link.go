package hcp_link

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	link "github.com/hashicorp/hcp-link"
	linkConfig "github.com/hashicorp/hcp-link/pkg/config"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/internalshared/configutil"
	vaultVersion "github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities/api_capability"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities/link_control"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities/meta"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities/node_status"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
)

const (
	SetLinkStatusCadence = 5 * time.Second

	// metaDataNodeStatus is used to set the Scada provider metadata status
	// to indicate if Vault is in active or standby status
	metaDataNodeStatus = "link.node_status"

	standbyStatus     = "STANDBY"
	activeStatus      = "ACTIVE"
	perfStandbyStatus = "PERF-STANDBY"
)

var (
	// genericScadaConnectionError is used when Vault fails to fetch
	// last connection error from Scada Provider
	genericScadaConnectionError = errors.New("unable to establish a connection with HCP")
	invalidClientCredentials    = errors.New("failed to get access token: oauth2: cannot fetch token: 401 Unauthorized")
)

type HCPLinkVault struct {
	l            sync.RWMutex
	LinkStatus   internal.WrappedCoreHCPLinkStatus
	scadaConfig  *scada.Config
	linkConfig   *linkConfig.Config
	link         link.Link
	logger       hclog.Logger
	capabilities map[string]capabilities.Capability
	stopCh       chan struct{}
	running      bool
}

func NewHCPLink(linkConf *configutil.HCPLinkConfig, core *vault.Core, logger hclog.Logger) (*HCPLinkVault, error) {
	if linkConf == nil {
		return nil, nil
	}

	scadaLogger := logger.Named("scada")
	scadaConfig, err := internal.NewScadaConfig(linkConf, scadaLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate SCADA config, %w", err)
	}

	// setting the link status in core, as link config is not nil
	// At this point scada provider has not been started yet
	// After starting scada provider, we need to use
	// scadaProvider.SessionStatus() to get the status of the connection
	core.SetHCPLinkStatus(
		buildConnectionErrorMessage(scada.SessionStatusDisconnected, scada.ErrProviderNotStarted.Error(), time.Now()),
		scadaConfig.Resource.Location.ProjectID,
	)

	// Creating SCADA provider
	scadaProvider, err := scada.New(scadaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate SCADA provider: %w", err)
	}

	resource := scadaConfig.Resource
	hcpConfig := scadaConfig.HCPConfig
	version := vaultVersion.Version

	// initializing node status reporter. This capability is configured by link lib.
	statusReporter := &node_status.NodeStatusReporter{
		NodeStatusGetter: core,
	}
	nodeID, err := core.LoadNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodeID, %w", err)
	}

	// Compile the Link config
	var conf *linkConfig.Config
	conf, err = internal.NewLinkConfig(
		nodeID,
		version,
		resource,
		scadaProvider,
		hcpConfig,
		statusReporter,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Link library config: %w", err)
	}

	// Create a Link library instance
	hcpLink, err := link.New(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Link library: %w", err)
	}

	hcpLinkCaps, err := initializeCapabilities(linkConf, scadaConfig, scadaProvider, core, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize capabilities: %w", err)
	}

	hcpLinkVault := &HCPLinkVault{
		LinkStatus:   core,
		scadaConfig:  scadaConfig,
		linkConfig:   conf,
		link:         hcpLink,
		capabilities: hcpLinkCaps,
		stopCh:       make(chan struct{}),
		logger:       logger,
	}

	// Start hcpLink and ScadaProvider
	err = hcpLinkVault.start()
	if err != nil {
		return nil, fmt.Errorf("failed to start hcp link, %w", err)
	}

	return hcpLinkVault, nil
}

func initializeCapabilities(linkConf *configutil.HCPLinkConfig, scadaConfig *scada.Config, scadaProvider scada.SCADAProvider, core *vault.Core, logger hclog.Logger) (map[string]capabilities.Capability, error) {
	hcpLinkCaps := make(map[string]capabilities.Capability, 0)

	metaCap := meta.NewHCPLinkMetaService(scadaProvider, core, logger)
	hcpLinkCaps[capabilities.MetaCapability] = metaCap

	// Initializing API and link-control capabilities
	var retErr *multierror.Error
	if linkConf.EnableAPICapability {
		apiCap, err := api_capability.NewAPICapability(scadaConfig, scadaProvider, core, logger)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("failed to instantiate API capability, %w", err))
		}
		hcpLinkCaps[capabilities.APICapability] = apiCap

		// link control capability is tied to api capability
		linkControlCap := link_control.NewHCPLinkControlService(scadaProvider, core, apiCap.PurgePolicy, logger)
		hcpLinkCaps[capabilities.LinkControlCapability] = linkControlCap

	}

	// Initializing Passthrough capability
	if linkConf.EnablePassThroughCapability {
		apiPassCap, err := api_capability.NewAPIPassThroughCapability(scadaProvider, core, logger)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("failed to instantiate PassThrough capability, %w", err))
		}
		hcpLinkCaps[capabilities.APIPassThroughCapability] = apiPassCap
	}

	return hcpLinkCaps, retErr.ErrorOrNil()
}

// Start the connection regardless if the node is in seal mode or not
func (h *HCPLinkVault) start() error {
	h.l.Lock()
	defer h.l.Unlock()

	if h.running {
		return nil
	}

	if h.linkConfig == nil {
		return fmt.Errorf("hcpLink config has not been provided")
	}

	scadaProvider := h.linkConfig.SCADAProvider
	if scadaProvider == nil {
		return fmt.Errorf("the reference to Scada provider in hcp link config is nil")
	}

	// Start both the Link functionality and the provider
	if err := h.link.Start(); err != nil {
		return fmt.Errorf("failed to start Link functionality, %w", err)
	}

	if err := scadaProvider.Start(); err != nil {
		return fmt.Errorf("failed to start SCADA provider, %w", err)
	}

	// The connection should have been established between Vault and HCP
	// Update core with the status
	h.LinkStatus.SetHCPLinkStatus(h.GetConnectionStatusMessage(h.GetScadaSessionStatus()), h.getResourceID())

	// Running capabilities
	err := h.RunCapabilities()
	if err != nil {
		h.logger.Error("failed to start HCP link capabilities", "error", err.Error())
	}

	go h.reportStatus()

	h.running = true

	h.logger.Info("established connection to HCP")

	return nil
}

// runs in a goroutine and in every 5 seconds, it sets the link status in Core
// such that a user could query the health of the connection via seal-status
// API. In addition, it checks replication status of Vault and sets that in
// Scada provider metadata status
func (h *HCPLinkVault) reportStatus() {
	h.l.RLock()
	stopCh := h.stopCh
	h.l.RUnlock()

	var currentNodeStatus string

	ticker := time.NewTicker(SetLinkStatusCadence)
	defer ticker.Stop()
	for {
		// Check for a shutdown
		select {
		case <-stopCh:
			h.logger.Trace("returning from reporting link/node status")
			return
		case <-ticker.C:
			// setting the HCP link status in core in this cadence
			h.LinkStatus.SetHCPLinkStatus(
				h.GetConnectionStatusMessage(h.GetScadaSessionStatus()),
				h.getResourceID(),
			)

			// if node is in standby mode, set Scada metadata to indicate that
			var nodeStatus string
			standby, perfStandby := h.LinkStatus.StandbyStates()
			switch {
			case perfStandby:
				nodeStatus = perfStandbyStatus
			case standby:
				nodeStatus = standbyStatus
			default:
				nodeStatus = activeStatus
			}

			// Only update SCADA session metadata if status has changed
			if currentNodeStatus != nodeStatus {
				currentNodeStatus = nodeStatus
				h.linkConfig.SCADAProvider.UpdateMeta(map[string]string{metaDataNodeStatus: currentNodeStatus})
			}
		}
	}
}

func buildConnectionErrorMessage(scadaStatus, errMsg string, errTime time.Time) string {
	return fmt.Sprintf("%s since %s; error: %v", scadaStatus, errTime.Format(time.RFC3339Nano), errMsg)
}

// GetConnectionStatusMessage returns a meaningful message about connection
// status. If Scada connection is anything other than "connected", it will
// get the LastError from ScadaProvider, and builds a message with the
// scada session status, error time and error message, and returns the message.
func (h *HCPLinkVault) GetConnectionStatusMessage(scadaStatus string) string {
	if scadaStatus == scada.SessionStatusConnected {
		return scadaStatus
	}

	// HCP connectivity team is going to unify "connecting" with "waiting"
	// statuses later. For simplicity, we unify the two until Scada
	// provider unifies them
	if scadaStatus == scada.SessionStatusWaiting {
		scadaStatus = scada.SessionStatusConnecting
	}

	// There are two other states "connecting" and "disconnected"
	// For those, there could have been an error with the connection
	var errToReturn string
	errTime, err := h.linkConfig.SCADAProvider.LastError()
	if err == nil {
		err = genericScadaConnectionError
		errTime = time.Now()
	}

	switch {
	case strings.Contains(err.Error(), scada.ErrPermissionDenied.Error()):
		errToReturn = scada.ErrPermissionDenied.Error()
	case strings.Contains(err.Error(), invalidClientCredentials.Error()), strings.Contains(err.Error(), scada.ErrInvalidCredentials.Error()):
		errToReturn = scada.ErrInvalidCredentials.Error()
	default:
		errToReturn = genericScadaConnectionError.Error()
	}

	return buildConnectionErrorMessage(scadaStatus, errToReturn, errTime)
}

func (h *HCPLinkVault) getResourceID() string {
	if h.scadaConfig != nil {
		return h.scadaConfig.Resource.ID
	}

	return ""
}

func (h *HCPLinkVault) GetScadaSessionStatus() string {
	if h.linkConfig != nil && h.linkConfig.SCADAProvider != nil {
		return h.linkConfig.SCADAProvider.SessionStatus()
	}
	return scada.SessionStatusDisconnected
}

func (h *HCPLinkVault) Shutdown() error {
	h.l.Lock()
	defer h.l.Unlock()

	if !h.running {
		return nil
	}

	if h.stopCh != nil {
		close(h.stopCh)
		h.stopCh = nil
	}

	h.logger.Info("tearing down connection to HCP")

	var retErr *multierror.Error

	// stopping capabilities
	for capName, capability := range h.capabilities {
		err := capability.Stop()
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("failed to close capability %s, %w", capName, err))
		}
	}

	// updating metaDataNodeStatus before stopping link
	h.linkConfig.SCADAProvider.UpdateMeta(map[string]string{metaDataNodeStatus: ""})

	// stopping hcp link
	err := h.link.Stop()
	if err != nil {
		retErr = multierror.Append(err, fmt.Errorf("failed to stop link %w", err))
	}
	h.link = nil

	// stopping scada provider
	err = h.linkConfig.SCADAProvider.Stop()
	if err != nil {
		retErr = multierror.Append(err, fmt.Errorf("failed to stop scada provider %w", err))
	}

	// setting the link status in Vault
	h.LinkStatus.SetHCPLinkStatus(h.GetConnectionStatusMessage(h.GetScadaSessionStatus()), h.getResourceID())

	h.running = false

	return retErr.ErrorOrNil()
}

func (h *HCPLinkVault) RunCapabilities() error {
	var retErr *multierror.Error
	for capName, capability := range h.capabilities {
		err := capability.Start()
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("failed to start capability %s, %v", capName, err))
		}
	}

	return retErr.ErrorOrNil()
}
