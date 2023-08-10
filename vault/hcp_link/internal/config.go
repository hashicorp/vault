// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	linkConfig "github.com/hashicorp/hcp-link/pkg/config"
	"github.com/hashicorp/hcp-link/pkg/nodestatus"
	scada "github.com/hashicorp/hcp-scada-provider"
	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
	sdkConfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/vault/internalshared/configutil"
)

const ServiceName = "vault-link"

func NewScadaConfig(linkConf *configutil.HCPLinkConfig, logger hclog.Logger) (*scada.Config, error) {
	// getting models.HashicorpCloudLocationLink to be passed in the
	// scada.config
	res := linkConf.Resource.Link()

	// creating a base from the env allows for overriding the following for dev purposes:
	// - auth URL:      HCP_AUTH_URL
	// - SCADA address: HCP_SCADA_ADDRESS
	// - API address:   HCP_API_ADDRESS
	opts := []sdkConfig.HCPConfigOption{sdkConfig.FromEnv()}

	// client ID and client secret from config takes precedence despite
	// sdkConfig.FromEnv allowing to set from env
	opts = append(opts, sdkConfig.WithClientCredentials(linkConf.ClientID, linkConf.ClientSecret))

	hcpConfig, err := sdkConfig.NewHCPConfig(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HCP config: %w", err)
	}

	// Compile SCADA config
	scadaConfig := &scada.Config{
		Service:   ServiceName,
		HCPConfig: hcpConfig,
		Resource:  *res,
		Logger:    logger,
	}
	return scadaConfig, nil
}

// NewLinkConfig validates the provided values and constructs an instance of a Config.
func NewLinkConfig(nodeID string, nodeVersion string, resource cloud.HashicorpCloudLocationLink, scadaProvider scada.SCADAProvider, hcpConfig sdkConfig.HCPConfig, nodeStatusReporter nodestatus.Reporter, logger hclog.Logger) (*linkConfig.Config, error) {
	config := &linkConfig.Config{
		NodeID:             nodeID,
		NodeVersion:        nodeVersion,
		HCPConfig:          hcpConfig,
		Resource:           resource,
		NodeStatusReporter: nodeStatusReporter,
		SCADAProvider:      scadaProvider,
		Logger:             logger,
	}

	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to create the link config: %w", err)
	}

	return config, nil
}
