// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package link implements routines for self-managed resources to connect to HashiCorp Cloud Platform (HCP).
package link

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"google.golang.org/grpc"

	linkstatuspb "github.com/hashicorp/hcp-link/gen/proto/go/hashicorp/cloud/hcp_link/link_status/v1"
	nodestatuspb "github.com/hashicorp/hcp-link/gen/proto/go/hashicorp/cloud/hcp_link/node_status/v1"
	linkstatusinternal "github.com/hashicorp/hcp-link/internal/linkstatus"
	nodestatusinternal "github.com/hashicorp/hcp-link/internal/nodestatus"
	"github.com/hashicorp/hcp-link/pkg/config"
)

const (
	metaDataNodeId      = "link.node_id"
	metaDataNodeVersion = "link.node_version"

	// Capability defines the name of the SCADA capability that is used to expose
	// the Link status and node status gRPC services.
	Capability = "link"
)

type link struct {
	// Config contains all dependencies as well as information about the node Link is
	// running on.
	*config.Config

	// apiClient is a http.Client that can be used to call the HCP API.
	apiClient *http.Client

	// collector is used to collect node status information.
	collector *nodestatusinternal.Collector

	// listener is the listener of the Link SCADA capability.
	listener net.Listener

	// grpcServer is the gRPC server of the Link SCADA capability.
	grpcServer *grpc.Server

	// running is set true if Link is running
	running     bool
	runningLock sync.Mutex
}

// New creates a new instance of a Link interface, that allows access to
// functionality of linked HCP services.
func New(config *config.Config) (Link, error) {
	if config == nil {
		return nil, fmt.Errorf("failed to initialize link library: config must be provided")
	}

	// Use the HCP SDK to prepare a transport
	runtime, err := httpclient.New(httpclient.Config{HCPConfig: config.HCPConfig, SourceChannel: "link"})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare API client transport: %w", err)
	}

	// Configure a http client with the transport
	apiClient := cleanhttp.DefaultClient()
	apiClient.Transport = runtime.Transport

	return &link{
		Config:    config,
		apiClient: apiClient,
		collector: &nodestatusinternal.Collector{

			Config: config,
		},
	}, nil
}

// Start implements Link interface.
//
// It will set the Link specific meta-data values and expose the Link specific capability.
func (l *link) Start() error {
	l.runningLock.Lock()
	defer l.runningLock.Unlock()

	// Check if Link is already running
	if l.running {
		return nil
	}

	// Configure Link specific meta-data
	l.SCADAProvider.UpdateMeta(map[string]string{
		metaDataNodeId:      l.NodeID,
		metaDataNodeVersion: l.NodeVersion,
	})

	// Start listening on Link capability
	listener, err := l.SCADAProvider.Listen(Capability)
	if err != nil {
		return fmt.Errorf("failed to start listening on the %q capability: %w", Capability, err)
	}
	l.listener = listener

	// Setup gRPC server
	l.grpcServer = grpc.NewServer()

	// Handle LinkStatus requests
	linkstatuspb.RegisterLinkStatusServiceServer(l.grpcServer, &linkstatusinternal.Service{
		Config: l.Config,
	})

	// Handle NodeStatus requests, if a node status reporter has been registered
	if l.NodeStatusReporter != nil {
		nodestatuspb.RegisterNodeStatusServiceServer(l.grpcServer, &nodestatusinternal.Service{
			Collector: l.collector,
		})
	}

	// Start the gRPC server
	go func() {
		_ = l.grpcServer.Serve(listener)
	}()

	// Mark Link as running
	l.running = true

	return nil
}

// Stop implements Link interface.
//
// It will unset the Link specific meta-data value and stop to expose the Link capability.
func (l *link) Stop() error {
	l.runningLock.Lock()
	defer l.runningLock.Unlock()

	// Check if Link is already stopped
	if !l.running {
		return nil
	}

	// Stop the gRPC server
	l.grpcServer.Stop()

	// Stop listening on the Link capability
	err := l.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener for %q capability: %w", Capability, err)
	}

	// Reset listener
	l.listener = nil

	// Mark Link as stopped
	l.running = false

	return nil
}

// ReportNodeStatus will get the most recent node status information from the
// configured node status reporter and push it to HCP.
//
// This function only needs to be invoked in situations where it is important
// that the node status is reported right away. HCP will regularly poll for node
// status information.
func (l *link) ReportNodeStatus(ctx context.Context) error {
	// Get the node status
	status, err := l.collector.CollectPb(ctx)
	if err != nil {
		return fmt.Errorf("failed to collect node status: %w", err)
	}

	// Marshal the node status into a binary proto message
	requestMessage, err := proto.Marshal(&nodestatuspb.SetNodeStatusRequest{NodeStatus: status})
	if err != nil {
		return fmt.Errorf("failed to marshal node status: %w", err)
	}

	// Determine the scheme
	scheme := "https"
	if l.HCPConfig.APITLSConfig() == nil {
		scheme = "http"
	}

	// Build the URL
	requestURL := url.URL{
		Scheme: scheme,
		Host:   l.HCPConfig.APIAddress(),
		Path: path.Join(
			"link/2022-06-04/status/organization",
			l.Resource.Location.OrganizationID,
			"project",
			l.Resource.Location.ProjectID,
			l.Resource.Type,
			l.Resource.ID,
			"node",
			status.NodeId,
		),
	}

	// Prepare the request
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL.String(), bytes.NewReader(requestMessage))
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}
	request.Header.Set("Content-Type", "application/octet-stream")

	// Call the API
	response, err := l.apiClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to set node status: %w", err)
	}
	defer response.Body.Close()

	// Check if the call was successful
	if response.StatusCode != http.StatusOK {
		body := make([]byte, 0, 256)
		_, _ = response.Body.Read(body)

		return fmt.Errorf("received non 200 response: %d, %v", response.StatusCode, string(body))
	}

	return nil
}
