// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/hashicorp/go-plugin/internal/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func dialGRPCConn(tls *tls.Config, dialer func(string, time.Duration) (net.Conn, error), dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Build dialing options.
	opts := make([]grpc.DialOption, 0)

	// We use a custom dialer so that we can connect over unix domain sockets.
	opts = append(opts, grpc.WithDialer(dialer))

	// Fail right away
	opts = append(opts, grpc.FailOnNonTempDialError(true))

	// If we have no TLS configuration set, we need to explicitly tell grpc
	// that we're connecting with an insecure connection.
	if tls == nil {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(tls)))
	}

	opts = append(opts,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)))

	// Add our custom options if we have any
	opts = append(opts, dialOpts...)

	// Connect. Note the first parameter is unused because we use a custom
	// dialer that has the state to see the address.
	conn, err := grpc.Dial("unused", opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// newGRPCClient creates a new GRPCClient. The Client argument is expected
// to be successfully started already with a lock held.
func newGRPCClient(doneCtx context.Context, c *Client) (*GRPCClient, error) {
	conn, err := dialGRPCConn(c.config.TLSConfig, c.dialer, c.config.GRPCDialOptions...)
	if err != nil {
		return nil, err
	}

	muxer, err := c.getGRPCMuxer(c.address)
	if err != nil {
		return nil, err
	}

	// Start the broker.
	brokerGRPCClient := newGRPCBrokerClient(conn)
	broker := newGRPCBroker(brokerGRPCClient, c.config.TLSConfig, c.unixSocketCfg, c.runner, muxer)
	go broker.Run()
	go brokerGRPCClient.StartStream()

	// Start the stdio client
	stdioClient, err := newGRPCStdioClient(doneCtx, c.logger.Named("stdio"), conn)
	if err != nil {
		return nil, err
	}
	go stdioClient.Run(c.config.SyncStdout, c.config.SyncStderr)

	cl := &GRPCClient{
		Conn:       conn,
		Plugins:    c.config.Plugins,
		doneCtx:    doneCtx,
		broker:     broker,
		controller: plugin.NewGRPCControllerClient(conn),
	}

	return cl, nil
}

// GRPCClient connects to a GRPCServer over gRPC to dispense plugin types.
type GRPCClient struct {
	Conn    *grpc.ClientConn
	Plugins map[string]Plugin

	doneCtx context.Context
	broker  *GRPCBroker

	controller plugin.GRPCControllerClient
}

// ClientProtocol impl.
func (c *GRPCClient) Close() error {
	c.broker.Close()
	c.controller.Shutdown(c.doneCtx, &plugin.Empty{})
	return c.Conn.Close()
}

// ClientProtocol impl.
func (c *GRPCClient) Dispense(name string) (interface{}, error) {
	raw, ok := c.Plugins[name]
	if !ok {
		return nil, fmt.Errorf("unknown plugin type: %s", name)
	}

	p, ok := raw.(GRPCPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin %q doesn't support gRPC", name)
	}

	return p.GRPCClient(c.doneCtx, c.broker, c.Conn)
}

// ClientProtocol impl.
func (c *GRPCClient) Ping() error {
	client := grpc_health_v1.NewHealthClient(c.Conn)
	_, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
		Service: GRPCServiceName,
	})

	return err
}
