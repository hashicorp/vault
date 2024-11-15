// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package config contains the configuration required for a self-managed resource to connect to HashiCorp Cloud
// Platform (HCP).
package config

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
	sdk "github.com/hashicorp/hcp-sdk-go/config"

	"github.com/hashicorp/hcp-link/internal/resource"
	"github.com/hashicorp/hcp-link/pkg/nodestatus"
)

// Config contains information about the node, the linked resource and the SCADA connection.
type Config struct {
	// NodeID is an ID that uniquely identifies the node within the Resource
	// (e.g. within the Vault cluster).
	//
	// This ID should ideally persist through restarts of the node.
	NodeID string

	// NodeVersion is the semantic versioning formatted version of the node.
	NodeVersion string

	// Resource is the Resource the Link client should identify as, the Resource
	// will have to be created on HCP before it can be provided to the library.
	Resource cloud.HashicorpCloudLocationLink

	// HCPConfig is the HCP specific configuration, it provides information
	// necessary to talk to HCP APIs.
	HCPConfig sdk.HCPConfig

	// NodeStatusReporter is used as a callback to retrieve the node's current
	// status information.
	NodeStatusReporter nodestatus.Reporter

	// SCADAProvider is a SCADA provider that is registered on HCP's SCADA broker.
	SCADAProvider scada.SCADAProvider

	// Logger is HCLog Logger instance that will be used to log debug information.
	Logger hclog.Logger
}

// Validate will validate the Link configuration. It ensures that the specified Resource is valid and that all
// fields are provided.
func (c *Config) Validate() error {
	err := resource.Validate(c.Resource)
	if err != nil {
		return fmt.Errorf("resource link is invalid: %w", err)
	}

	if c.NodeID == "" {
		return fmt.Errorf("node ID must be provided")
	}
	if c.NodeVersion == "" {
		return fmt.Errorf("node version must be provided")
	}
	if c.HCPConfig == nil {
		return fmt.Errorf("HCP config must be provided")
	}
	if c.SCADAProvider == nil {
		return fmt.Errorf("SCADA provider must be provided")
	}
	if c.Logger == nil {
		return fmt.Errorf("logger must be provided")
	}

	return nil
}
