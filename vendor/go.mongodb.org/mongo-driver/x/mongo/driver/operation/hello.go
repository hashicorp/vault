// Copyright (C) MongoDB, Inc. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/bsonutil"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/handshake"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/version"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// maxClientMetadataSize is the maximum size of the client metadata document
// that can be sent to the server. Note that the maximum document size on
// standalone and replica servers is 1024, but the maximum document size on
// sharded clusters is 512.
const maxClientMetadataSize = 512

const driverName = "mongo-go-driver"

// Hello is used to run the handshake operation.
type Hello struct {
	authenticator      driver.Authenticator
	appname            string
	compressors        []string
	saslSupportedMechs string
	d                  driver.Deployment
	clock              *session.ClusterClock
	speculativeAuth    bsoncore.Document
	topologyVersion    *description.TopologyVersion
	maxAwaitTimeMS     *int64
	serverAPI          *driver.ServerAPIOptions
	loadBalanced       bool

	res bsoncore.Document
}

var _ driver.Handshaker = (*Hello)(nil)

// NewHello constructs a Hello.
func NewHello() *Hello { return &Hello{} }

// AppName sets the application name in the client metadata sent in this operation.
func (h *Hello) AppName(appname string) *Hello {
	h.appname = appname
	return h
}

// ClusterClock sets the cluster clock for this operation.
func (h *Hello) ClusterClock(clock *session.ClusterClock) *Hello {
	if h == nil {
		h = new(Hello)
	}

	h.clock = clock
	return h
}

// Compressors sets the compressors that can be used.
func (h *Hello) Compressors(compressors []string) *Hello {
	h.compressors = compressors
	return h
}

// SASLSupportedMechs retrieves the supported SASL mechanism for the given user when this operation
// is run.
func (h *Hello) SASLSupportedMechs(username string) *Hello {
	h.saslSupportedMechs = username
	return h
}

// Deployment sets the Deployment for this operation.
func (h *Hello) Deployment(d driver.Deployment) *Hello {
	h.d = d
	return h
}

// SpeculativeAuthenticate sets the document to be used for speculative authentication.
func (h *Hello) SpeculativeAuthenticate(doc bsoncore.Document) *Hello {
	h.speculativeAuth = doc
	return h
}

// TopologyVersion sets the TopologyVersion to be used for heartbeats.
func (h *Hello) TopologyVersion(tv *description.TopologyVersion) *Hello {
	h.topologyVersion = tv
	return h
}

// MaxAwaitTimeMS sets the maximum time for the server to wait for topology changes during a heartbeat.
func (h *Hello) MaxAwaitTimeMS(awaitTime int64) *Hello {
	h.maxAwaitTimeMS = &awaitTime
	return h
}

// ServerAPI sets the server API version for this operation.
func (h *Hello) ServerAPI(serverAPI *driver.ServerAPIOptions) *Hello {
	h.serverAPI = serverAPI
	return h
}

// LoadBalanced specifies whether or not this operation is being sent over a connection to a load balanced cluster.
func (h *Hello) LoadBalanced(lb bool) *Hello {
	h.loadBalanced = lb
	return h
}

// Result returns the result of executing this operation.
func (h *Hello) Result(addr address.Address) description.Server {
	return description.NewServer(addr, bson.Raw(h.res))
}

const dockerEnvPath = "/.dockerenv"

const (
	// Runtime names
	runtimeNameDocker = "docker"

	// Orchestrator names
	orchestratorNameK8s = "kubernetes"
)

// getFaasEnvName parses the FaaS environment variable name and returns the
// corresponding name used by the client. If none of the variables or variables
// for multiple names are populated the FaaS values MUST be entirely omitted.
// When variables for multiple "client.env.name" values are present, "vercel"
// takes precedence over "aws.lambda"; any other combination MUST cause FaaS
// values to be entirely omitted.
func getFaasEnvName() string {
	envVars := []string{
		driverutil.EnvVarAWSExecutionEnv,
		driverutil.EnvVarAWSLambdaRuntimeAPI,
		driverutil.EnvVarFunctionsWorkerRuntime,
		driverutil.EnvVarKService,
		driverutil.EnvVarFunctionName,
		driverutil.EnvVarVercel,
	}

	// If none of the variables are populated the client.env value MUST be
	// entirely omitted.
	names := make(map[string]struct{})

	for _, envVar := range envVars {
		val := os.Getenv(envVar)
		if val == "" {
			continue
		}

		var name string

		switch envVar {
		case driverutil.EnvVarAWSExecutionEnv:
			if !strings.HasPrefix(val, driverutil.AwsLambdaPrefix) {
				continue
			}

			name = driverutil.EnvNameAWSLambda
		case driverutil.EnvVarAWSLambdaRuntimeAPI:
			name = driverutil.EnvNameAWSLambda
		case driverutil.EnvVarFunctionsWorkerRuntime:
			name = driverutil.EnvNameAzureFunc
		case driverutil.EnvVarKService, driverutil.EnvVarFunctionName:
			name = driverutil.EnvNameGCPFunc
		case driverutil.EnvVarVercel:
			// "vercel" takes precedence over "aws.lambda".
			delete(names, driverutil.EnvNameAWSLambda)

			name = driverutil.EnvNameVercel
		}

		names[name] = struct{}{}
		if len(names) > 1 {
			// If multiple names are populated the client.env value
			// MUST be entirely omitted.
			names = nil

			break
		}
	}

	for name := range names {
		return name
	}

	return ""
}

type containerInfo struct {
	runtime      string
	orchestrator string
}

// getContainerEnvInfo returns runtime and orchestrator of a container.
// If no fields is populated, the client.env.container value MUST be entirely
// omitted.
func getContainerEnvInfo() *containerInfo {
	var runtime, orchestrator string
	if _, err := os.Stat(dockerEnvPath); !os.IsNotExist(err) {
		runtime = runtimeNameDocker
	}
	if v := os.Getenv(driverutil.EnvVarK8s); v != "" {
		orchestrator = orchestratorNameK8s
	}
	if runtime != "" || orchestrator != "" {
		return &containerInfo{
			runtime:      runtime,
			orchestrator: orchestrator,
		}
	}
	return nil
}

// appendClientAppName appends the application metadata to the dst. It is the
// responsibility of the caller to check that this appending does not cause dst
// to exceed any size limitations.
func appendClientAppName(dst []byte, name string) ([]byte, error) {
	if name == "" {
		return dst, nil
	}

	var idx int32
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "application")

	dst = bsoncore.AppendStringElement(dst, "name", name)

	return bsoncore.AppendDocumentEnd(dst, idx)
}

// appendClientDriver appends the driver metadata to dst. It is the
// responsibility of the caller to check that this appending does not cause dst
// to exceed any size limitations.
func appendClientDriver(dst []byte) ([]byte, error) {
	var idx int32
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "driver")

	dst = bsoncore.AppendStringElement(dst, "name", driverName)
	dst = bsoncore.AppendStringElement(dst, "version", version.Driver)

	return bsoncore.AppendDocumentEnd(dst, idx)
}

// appendClientEnv appends the environment metadata to dst. It is the
// responsibility of the caller to check that this appending does not cause dst
// to exceed any size limitations.
func appendClientEnv(dst []byte, omitNonName, omitDoc bool) ([]byte, error) {
	if omitDoc {
		return dst, nil
	}

	name := getFaasEnvName()
	container := getContainerEnvInfo()
	// Omit the entire 'env' if both name and container are empty because other
	// fields depend on either of them.
	if name == "" && container == nil {
		return dst, nil
	}

	var idx int32

	idx, dst = bsoncore.AppendDocumentElementStart(dst, "env")

	if name != "" {
		dst = bsoncore.AppendStringElement(dst, "name", name)
	}

	addMem := func(envVar string) []byte {
		mem := os.Getenv(envVar)
		if mem == "" {
			return dst
		}

		memInt64, err := strconv.ParseInt(mem, 10, 32)
		if err != nil {
			return dst
		}

		memInt32 := int32(memInt64)

		return bsoncore.AppendInt32Element(dst, "memory_mb", memInt32)
	}

	addRegion := func(envVar string) []byte {
		region := os.Getenv(envVar)
		if region == "" {
			return dst
		}

		return bsoncore.AppendStringElement(dst, "region", region)
	}

	addTimeout := func(envVar string) []byte {
		timeout := os.Getenv(envVar)
		if timeout == "" {
			return dst
		}

		timeoutInt64, err := strconv.ParseInt(timeout, 10, 32)
		if err != nil {
			return dst
		}

		timeoutInt32 := int32(timeoutInt64)
		return bsoncore.AppendInt32Element(dst, "timeout_sec", timeoutInt32)
	}

	if !omitNonName {
		// No other FaaS fields will be populated if the name is empty.
		switch name {
		case driverutil.EnvNameAWSLambda:
			dst = addMem(driverutil.EnvVarAWSLambdaFunctionMemorySize)
			dst = addRegion(driverutil.EnvVarAWSRegion)
		case driverutil.EnvNameGCPFunc:
			dst = addMem(driverutil.EnvVarFunctionMemoryMB)
			dst = addRegion(driverutil.EnvVarFunctionRegion)
			dst = addTimeout(driverutil.EnvVarFunctionTimeoutSec)
		case driverutil.EnvNameVercel:
			dst = addRegion(driverutil.EnvVarVercelRegion)
		}
	}

	if container != nil {
		var idxCntnr int32
		idxCntnr, dst = bsoncore.AppendDocumentElementStart(dst, "container")
		if container.runtime != "" {
			dst = bsoncore.AppendStringElement(dst, "runtime", container.runtime)
		}
		if container.orchestrator != "" {
			dst = bsoncore.AppendStringElement(dst, "orchestrator", container.orchestrator)
		}
		var err error
		dst, err = bsoncore.AppendDocumentEnd(dst, idxCntnr)
		if err != nil {
			return dst, err
		}
	}

	return bsoncore.AppendDocumentEnd(dst, idx)
}

// appendClientOS appends the OS metadata to dst. It is the responsibility of the
// caller to check that this appending does not cause dst to exceed any size
// limitations.
func appendClientOS(dst []byte, omitNonType bool) ([]byte, error) {
	var idx int32

	idx, dst = bsoncore.AppendDocumentElementStart(dst, "os")

	dst = bsoncore.AppendStringElement(dst, "type", runtime.GOOS)
	if !omitNonType {
		dst = bsoncore.AppendStringElement(dst, "architecture", runtime.GOARCH)
	}

	return bsoncore.AppendDocumentEnd(dst, idx)
}

// appendClientPlatform appends the platform metadata to dst. It is the
// responsibility of the caller to check that this appending does not cause dst
// to exceed any size limitations.
func appendClientPlatform(dst []byte) []byte {
	return bsoncore.AppendStringElement(dst, "platform", runtime.Version())
}

// encodeClientMetadata encodes the client metadata into a BSON document. maxLen
// is the maximum length the document can be. If the document exceeds maxLen,
// then an empty byte slice is returned. If there is not enough space to encode
// a document, the document is truncated and returned.
//
// This function attempts to build the following document. Fields are omitted to
// save space following the MongoDB Handshake.
//
//	{
//		application: {
//			name: "<string>"
//		},
//		driver: {
//			name: "<string>",
//			version: "<string>"
//		},
//		platform: "<string>",
//		os: {
//			type: "<string>",
//			name: "<string>",
//			architecture: "<string>",
//			version: "<string>"
//		},
//		env: {
//			name: "<string>",
//			timeout_sec: 42,
//			memory_mb: 1024,
//			region: "<string>",
//			container: {
//				runtime: "<string>",
//				orchestrator: "<string>"
//			}
//		}
//	}
func encodeClientMetadata(appname string, maxLen int) ([]byte, error) {
	dst := make([]byte, 0, maxLen)

	omitEnvDoc := false
	omitEnvNonName := false
	omitOSNonType := false
	omitEnvDocument := false
	truncatePlatform := false

retry:
	var idx int32
	idx, dst = bsoncore.AppendDocumentStart(dst)

	var err error
	dst, err = appendClientAppName(dst, appname)
	if err != nil {
		return nil, err
	}

	dst, err = appendClientDriver(dst)
	if err != nil {
		return nil, err
	}

	dst, err = appendClientOS(dst, omitOSNonType)
	if err != nil {
		return nil, err
	}

	if !truncatePlatform {
		dst = appendClientPlatform(dst)
	}

	if !omitEnvDocument {
		dst, err = appendClientEnv(dst, omitEnvNonName, omitEnvDoc)
		if err != nil {
			return nil, err
		}
	}

	dst, err = bsoncore.AppendDocumentEnd(dst, idx)
	if err != nil {
		return nil, err
	}

	if len(dst) > maxLen {
		// Implementers SHOULD cumulatively update fields in the
		// following order until the document is under the size limit
		//
		//    1. Omit fields from ``env`` except ``env.name``
		//    2. Omit fields from ``os`` except ``os.type``
		//    3. Omit the ``env`` document entirely
		//    4. Truncate ``platform``
		dst = dst[:0]

		if !omitEnvNonName {
			omitEnvNonName = true

			goto retry
		}

		if !omitOSNonType {
			omitOSNonType = true

			goto retry
		}

		if !omitEnvDoc {
			omitEnvDoc = true

			goto retry
		}

		if !truncatePlatform {
			truncatePlatform = true

			goto retry
		}

		// There is nothing left to update. Return an empty slice to
		// tell caller not to append a `client` document.
		return nil, nil
	}

	return dst, nil
}

// handshakeCommand appends all necessary command fields as well as client metadata, SASL supported mechs, and compression.
func (h *Hello) handshakeCommand(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst, err := h.command(dst, desc)
	if err != nil {
		return dst, err
	}

	if h.saslSupportedMechs != "" {
		dst = bsoncore.AppendStringElement(dst, "saslSupportedMechs", h.saslSupportedMechs)
	}
	if h.speculativeAuth != nil {
		dst = bsoncore.AppendDocumentElement(dst, "speculativeAuthenticate", h.speculativeAuth)
	}
	var idx int32
	idx, dst = bsoncore.AppendArrayElementStart(dst, "compression")
	for i, compressor := range h.compressors {
		dst = bsoncore.AppendStringElement(dst, strconv.Itoa(i), compressor)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, idx)

	clientMetadata, _ := encodeClientMetadata(h.appname, maxClientMetadataSize)

	// If the client metadata is empty, do not append it to the command.
	if len(clientMetadata) > 0 {
		dst = bsoncore.AppendDocumentElement(dst, "client", clientMetadata)
	}

	return dst, nil
}

// command appends all necessary command fields.
func (h *Hello) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	// Use "hello" if topology is LoadBalanced, API version is declared or server
	// has responded with "helloOk". Otherwise, use legacy hello.
	if h.loadBalanced || h.serverAPI != nil || desc.Server.HelloOK {
		dst = bsoncore.AppendInt32Element(dst, "hello", 1)
	} else {
		dst = bsoncore.AppendInt32Element(dst, handshake.LegacyHello, 1)
	}
	dst = bsoncore.AppendBooleanElement(dst, "helloOk", true)

	if tv := h.topologyVersion; tv != nil {
		var tvIdx int32

		tvIdx, dst = bsoncore.AppendDocumentElementStart(dst, "topologyVersion")
		dst = bsoncore.AppendObjectIDElement(dst, "processId", tv.ProcessID)
		dst = bsoncore.AppendInt64Element(dst, "counter", tv.Counter)
		dst, _ = bsoncore.AppendDocumentEnd(dst, tvIdx)
	}
	if h.maxAwaitTimeMS != nil {
		dst = bsoncore.AppendInt64Element(dst, "maxAwaitTimeMS", *h.maxAwaitTimeMS)
	}
	if h.loadBalanced {
		// The loadBalanced parameter should only be added if it's true. We should never explicitly send
		// loadBalanced=false per the load balancing spec.
		dst = bsoncore.AppendBooleanElement(dst, "loadBalanced", true)
	}

	return dst, nil
}

// Execute runs this operation.
func (h *Hello) Execute(ctx context.Context) error {
	if h.d == nil {
		return errors.New("a Hello must have a Deployment set before Execute can be called")
	}

	return h.createOperation().Execute(ctx)
}

// StreamResponse gets the next streaming Hello response from the server.
func (h *Hello) StreamResponse(ctx context.Context, conn driver.StreamerConnection) error {
	return h.createOperation().ExecuteExhaust(ctx, conn)
}

// isLegacyHandshake returns True if server API version is not requested and
// loadBalanced is False. If this is the case, then the drivers MUST use legacy
// hello for the first message of the initial handshake with the OP_QUERY
// protocol
func isLegacyHandshake(srvAPI *driver.ServerAPIOptions, loadbalanced bool) bool {
	return srvAPI == nil && !loadbalanced
}

func (h *Hello) createOperation() driver.Operation {
	op := driver.Operation{
		Clock:      h.clock,
		CommandFn:  h.command,
		Database:   "admin",
		Deployment: h.d,
		ProcessResponseFn: func(info driver.ResponseInfo) error {
			h.res = info.ServerResponse
			return nil
		},
		ServerAPI: h.serverAPI,
	}

	if isLegacyHandshake(h.serverAPI, h.loadBalanced) {
		op.Legacy = driver.LegacyHandshake
	}

	return op
}

// GetHandshakeInformation performs the MongoDB handshake for the provided connection and returns the relevant
// information about the server. This function implements the driver.Handshaker interface.
func (h *Hello) GetHandshakeInformation(ctx context.Context, _ address.Address, c driver.Connection) (driver.HandshakeInformation, error) {
	deployment := driver.SingleConnectionDeployment{C: c}

	op := driver.Operation{
		Clock:      h.clock,
		CommandFn:  h.handshakeCommand,
		Deployment: deployment,
		Database:   "admin",
		ProcessResponseFn: func(info driver.ResponseInfo) error {
			h.res = info.ServerResponse
			return nil
		},
		ServerAPI: h.serverAPI,
	}

	if isLegacyHandshake(h.serverAPI, h.loadBalanced) {
		op.Legacy = driver.LegacyHandshake
	}

	if err := op.Execute(ctx); err != nil {
		return driver.HandshakeInformation{}, err
	}

	info := driver.HandshakeInformation{
		Description: h.Result(c.Address()),
	}
	if speculativeAuthenticate, ok := h.res.Lookup("speculativeAuthenticate").DocumentOK(); ok {
		info.SpeculativeAuthenticate = speculativeAuthenticate
	}
	if serverConnectionID, ok := h.res.Lookup("connectionId").AsInt64OK(); ok {
		info.ServerConnectionID = &serverConnectionID
	}

	var err error

	// Cast to bson.Raw to lookup saslSupportedMechs to avoid converting from bsoncore.Value to bson.RawValue for the
	// StringSliceFromRawValue call.
	if saslSupportedMechs, lookupErr := bson.Raw(h.res).LookupErr("saslSupportedMechs"); lookupErr == nil {
		info.SaslSupportedMechs, err = bsonutil.StringSliceFromRawValue("saslSupportedMechs", saslSupportedMechs)
	}
	return info, err
}

// FinishHandshake implements the Handshaker interface. This is a no-op function because a non-authenticated connection
// does not do anything besides the initial Hello for a handshake.
func (h *Hello) FinishHandshake(context.Context, driver.Connection) error {
	return nil
}

// Authenticator sets the authenticator to use for this operation.
func (h *Hello) Authenticator(authenticator driver.Authenticator) *Hello {
	if h == nil {
		h = new(Hello)
	}

	h.authenticator = authenticator
	return h
}
