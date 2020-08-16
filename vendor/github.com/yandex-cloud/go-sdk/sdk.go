// Copyright (c) 2017 Yandex LLC. All rights reserved.
// Author: Alexey Baranov <baranovich@yandex-team.ru>

package ycsdk

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/endpoint"
	iampb "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
	"github.com/yandex-cloud/go-sdk/dial"
	apiendpoint "github.com/yandex-cloud/go-sdk/gen/apiendpoint"
	"github.com/yandex-cloud/go-sdk/gen/compute"
	"github.com/yandex-cloud/go-sdk/gen/iam"
	k8s "github.com/yandex-cloud/go-sdk/gen/kubernetes"
	gen_operation "github.com/yandex-cloud/go-sdk/gen/operation"
	"github.com/yandex-cloud/go-sdk/gen/resourcemanager"
	"github.com/yandex-cloud/go-sdk/gen/vpc"
	sdk_operation "github.com/yandex-cloud/go-sdk/operation"
	"github.com/yandex-cloud/go-sdk/pkg/grpcclient"
	"github.com/yandex-cloud/go-sdk/pkg/sdkerrors"
	"github.com/yandex-cloud/go-sdk/pkg/singleflight"
)

type Endpoint string

const (
	DefaultPageSize int64 = 1000

	ComputeServiceID            Endpoint = "compute"
	IAMServiceID                Endpoint = "iam"
	OperationServiceID          Endpoint = "operation"
	ResourceManagementServiceID Endpoint = "resource-manager"
	StorageServiceID            Endpoint = "storage"
	SerialSSHServiceID          Endpoint = "serialssh"
	// revive:disable:var-naming
	ApiEndpointServiceID Endpoint = "endpoint"
	// revive:enable:var-naming
	VpcServiceID        Endpoint = "vpc"
	KubernetesServiceID Endpoint = "managed-kubernetes"
)

// Config is a config that is used to create SDK instance.
type Config struct {
	// Credentials are used to authenticate the client. See Credentials for more info.
	Credentials Credentials
	// DialContextTimeout specifies timeout of dial on API endpoint that
	// is used when building an SDK instance.
	DialContextTimeout time.Duration
	// TLSConfig is optional tls.Config that one can use in order to tune TLS options.
	TLSConfig *tls.Config

	// Endpoint is an API endpoint of Yandex.Cloud against which the SDK is used.
	// Most users won't need to explicitly set it.
	Endpoint  string
	Plaintext bool
}

// SDK is a Yandex.Cloud SDK
type SDK struct {
	conf      Config
	cc        grpcclient.ConnContext
	endpoints struct {
		initDone bool
		mu       sync.Mutex
		ep       map[Endpoint]*endpoint.ApiEndpoint
	}

	initErr  error
	initCall singleflight.Call
	muErr    sync.Mutex
}

// Build creates an SDK instance
func Build(ctx context.Context, conf Config, customOpts ...grpc.DialOption) (*SDK, error) {
	if conf.Credentials == nil {
		return nil, errors.New("credentials required")
	}

	const defaultEndpoint = "api.cloud.yandex.net:443"
	if conf.Endpoint == "" {
		conf.Endpoint = defaultEndpoint
	}
	const DefaultTimeout = 20 * time.Second
	if conf.DialContextTimeout == 0 {
		conf.DialContextTimeout = DefaultTimeout
	}

	switch creds := conf.Credentials.(type) {
	case ExchangeableCredentials, NonExchangeableCredentials:
	default:
		return nil, fmt.Errorf("unsupported credentials type %T", creds)
	}
	var dialOpts []grpc.DialOption

	dialOpts = append(dialOpts, grpc.WithContextDialer(dial.NewProxyDialer(dial.NewDialer())))

	rpcCreds := newRPCCredentials(conf.Plaintext)
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(rpcCreds))
	if conf.DialContextTimeout > 0 {
		dialOpts = append(dialOpts, grpc.WithBlock(), grpc.WithTimeout(conf.DialContextTimeout)) // nolint
	}
	if conf.Plaintext {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		tlsConfig := conf.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		creds := credentials.NewTLS(tlsConfig)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}
	// Append custom options after default, to allow to customize dialer and etc.
	dialOpts = append(dialOpts, customOpts...)

	cc := grpcclient.NewLazyConnContext(grpcclient.DialOptions(dialOpts...))
	sdk := &SDK{
		cc:   cc,
		conf: conf,
	}
	rpcCreds.Init(sdk.CreateIAMToken)
	return sdk, nil
}

// Shutdown shutdowns SDK and closes all open connections.
func (sdk *SDK) Shutdown(ctx context.Context) error {
	return sdk.cc.Shutdown(ctx)
}

// WrapOperation wraps operation proto message to
func (sdk *SDK) WrapOperation(op *operation.Operation, err error) (*sdk_operation.Operation, error) {
	if err != nil {
		return nil, err
	}
	return sdk_operation.New(sdk.Operation(), op), nil
}

// IAM returns IAM object that is used to operate on Yandex Cloud Identity and Access Manager
func (sdk *SDK) IAM() *iam.IAM {
	return iam.NewIAM(sdk.getConn(IAMServiceID))
}

// Compute returns Compute object that is used to operate on Yandex Compute Cloud
func (sdk *SDK) Compute() *compute.Compute {
	return compute.NewCompute(sdk.getConn(ComputeServiceID))
}

// VPC returns VPC object that is used to operate on Yandex Virtual Private Cloud
func (sdk *SDK) VPC() *vpc.VPC {
	return vpc.NewVPC(sdk.getConn(VpcServiceID))
}

// MDB returns MDB object that is used to operate on Yandex Managed Databases
func (sdk *SDK) MDB() *MDB {
	return &MDB{sdk: sdk}
}

func (sdk *SDK) Serverless() *Serverless {
	return &Serverless{sdk: sdk}
}

func (sdk *SDK) Marketplace() *Marketplace {
	return &Marketplace{sdk: sdk}
}

// Operation gets OperationService client
func (sdk *SDK) Operation() *gen_operation.OperationServiceClient {
	group := gen_operation.NewOperation(sdk.getConn(OperationServiceID))
	return group.Operation()
}

// ResourceManager returns ResourceManager object that is used to operate on Folders and Clouds
func (sdk *SDK) ResourceManager() *resourcemanager.ResourceManager {
	return resourcemanager.NewResourceManager(sdk.getConn(ResourceManagementServiceID))
}

// revive:disable:var-naming

// ApiEndpoint gets ApiEndpointService client
func (sdk *SDK) ApiEndpoint() *apiendpoint.APIEndpoint {
	return apiendpoint.NewAPIEndpoint(sdk.getConn(ApiEndpointServiceID))
}

// revive:enable:var-naming

// Kubernetes returns Kubernetes object that is used to operate on Yandex Managed Kubernetes
func (sdk *SDK) Kubernetes() *k8s.Kubernetes {
	return k8s.NewKubernetes(sdk.getConn(KubernetesServiceID))
}

// AI returns AI object that is used to do AI stuff.
func (sdk *SDK) AI() *AI {
	return &AI{sdk: sdk}
}

func (sdk *SDK) Resolve(ctx context.Context, r ...Resolver) error {
	args := make([]func() error, len(r))
	for k, v := range r {
		resolver := v
		args[k] = func() error {
			return resolver.Run(ctx, sdk)
		}
	}
	return sdkerrors.CombineGoroutines(args...)
}

func (sdk *SDK) getConn(serviceID Endpoint) func(ctx context.Context) (*grpc.ClientConn, error) {
	return func(ctx context.Context) (*grpc.ClientConn, error) {
		if !sdk.initDone() {
			sdk.initCall.Do(func() interface{} {
				sdk.muErr.Lock()
				sdk.initErr = sdk.initConns(ctx)
				sdk.muErr.Unlock()
				return nil
			})
			if err := sdk.InitErr(); err != nil {
				return nil, err
			}
		}
		endpoint, endpointExist := sdk.Endpoint(serviceID)
		if !endpointExist {
			return nil, fmt.Errorf("server doesn't know service \"%v\". Known services: %v",
				serviceID,
				sdk.KnownServices())
		}
		return sdk.cc.GetConn(ctx, endpoint.Address)
	}
}

func (sdk *SDK) initDone() (b bool) {
	sdk.endpoints.mu.Lock()
	b = sdk.endpoints.initDone
	sdk.endpoints.mu.Unlock()
	return
}

func (sdk *SDK) KnownServices() []string {
	sdk.endpoints.mu.Lock()
	result := make([]string, 0, len(sdk.endpoints.ep))
	for k := range sdk.endpoints.ep {
		result = append(result, string(k))
	}
	sdk.endpoints.mu.Unlock()
	sort.Strings(result)
	return result
}

func (sdk *SDK) Endpoint(endpointName Endpoint) (ep *endpoint.ApiEndpoint, exist bool) {
	sdk.endpoints.mu.Lock()
	ep, exist = sdk.endpoints.ep[endpointName]
	sdk.endpoints.mu.Unlock()
	return
}

func (sdk *SDK) InitErr() error {
	sdk.muErr.Lock()
	defer sdk.muErr.Unlock()
	return sdk.initErr
}

func (sdk *SDK) initConns(ctx context.Context) error {
	discoveryConn, err := sdk.cc.GetConn(ctx, sdk.conf.Endpoint)
	if err != nil {
		return err
	}
	ec := endpoint.NewApiEndpointServiceClient(discoveryConn)
	const defaultEndpointPageSize = 100
	listResponse, err := ec.List(ctx, &endpoint.ListApiEndpointsRequest{
		PageSize: defaultEndpointPageSize,
	})
	if err != nil {
		return err
	}
	sdk.endpoints.mu.Lock()
	sdk.endpoints.ep = make(map[Endpoint]*endpoint.ApiEndpoint, len(listResponse.Endpoints))
	for _, e := range listResponse.Endpoints {
		sdk.endpoints.ep[Endpoint(e.Id)] = e
	}
	sdk.endpoints.initDone = true
	sdk.endpoints.mu.Unlock()
	return nil
}

func (sdk *SDK) CreateIAMToken(ctx context.Context) (*iampb.CreateIamTokenResponse, error) {
	creds := sdk.conf.Credentials
	switch creds := creds.(type) {
	case ExchangeableCredentials:
		req, err := creds.IAMTokenRequest()
		if err != nil {
			return nil, sdkerrors.WithMessage(err, "IAM token request build failed")
		}
		return sdk.IAM().IamToken().Create(ctx, req)
	case NonExchangeableCredentials:
		return creds.IAMToken(ctx)
	default:
		return nil, fmt.Errorf("credentials type %T is not supported yet", creds)
	}
}
