// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

// iamAPI is the subset of the AWS SDK v2 IAM client used by this backend. It is
// defined locally because the v2 SDK no longer ships generated interface
// packages (the v1 iamiface). *iam.Client satisfies this interface, and tests
// provide their own implementations.
type iamAPI interface {
	AddUserToGroup(ctx context.Context, params *iam.AddUserToGroupInput, optFns ...func(*iam.Options)) (*iam.AddUserToGroupOutput, error)
	AttachUserPolicy(ctx context.Context, params *iam.AttachUserPolicyInput, optFns ...func(*iam.Options)) (*iam.AttachUserPolicyOutput, error)
	CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	CreateUser(ctx context.Context, params *iam.CreateUserInput, optFns ...func(*iam.Options)) (*iam.CreateUserOutput, error)
	DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
	DeleteUser(ctx context.Context, params *iam.DeleteUserInput, optFns ...func(*iam.Options)) (*iam.DeleteUserOutput, error)
	DeleteUserPolicy(ctx context.Context, params *iam.DeleteUserPolicyInput, optFns ...func(*iam.Options)) (*iam.DeleteUserPolicyOutput, error)
	DetachUserPolicy(ctx context.Context, params *iam.DetachUserPolicyInput, optFns ...func(*iam.Options)) (*iam.DetachUserPolicyOutput, error)
	GetGroupPolicy(ctx context.Context, params *iam.GetGroupPolicyInput, optFns ...func(*iam.Options)) (*iam.GetGroupPolicyOutput, error)
	GetUser(ctx context.Context, params *iam.GetUserInput, optFns ...func(*iam.Options)) (*iam.GetUserOutput, error)
	ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
	ListAttachedGroupPolicies(ctx context.Context, params *iam.ListAttachedGroupPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedGroupPoliciesOutput, error)
	ListAttachedUserPolicies(ctx context.Context, params *iam.ListAttachedUserPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedUserPoliciesOutput, error)
	ListGroupPolicies(ctx context.Context, params *iam.ListGroupPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListGroupPoliciesOutput, error)
	ListGroupsForUser(ctx context.Context, params *iam.ListGroupsForUserInput, optFns ...func(*iam.Options)) (*iam.ListGroupsForUserOutput, error)
	ListUserPolicies(ctx context.Context, params *iam.ListUserPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error)
	PutUserPolicy(ctx context.Context, params *iam.PutUserPolicyInput, optFns ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error)
	RemoveUserFromGroup(ctx context.Context, params *iam.RemoveUserFromGroupInput, optFns ...func(*iam.Options)) (*iam.RemoveUserFromGroupOutput, error)
	TagUser(ctx context.Context, params *iam.TagUserInput, optFns ...func(*iam.Options)) (*iam.TagUserOutput, error)
}

// stsAPI is the subset of the AWS SDK v2 STS client used by this backend.
// *sts.Client satisfies this interface.
type stsAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
	GetFederationToken(ctx context.Context, params *sts.GetFederationTokenInput, optFns ...func(*sts.Options)) (*sts.GetFederationTokenOutput, error)
	GetSessionToken(ctx context.Context, params *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error)
}

const (
	rootConfigPath        = "config/root"
	minAwsUserRollbackAge = 5 * time.Minute
	operationPrefixAWS    = "aws"
	operationPrefixAWSASD = "aws-config"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(_ *logical.BackendConfig) *backend {
	var b backend
	b.minAllowableRotationPeriod = minAllowableRotationPeriod
	b.credRotationQueue = queue.New()
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				rootConfigPath,
				pathStaticCreds + "/",
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(&b),
			pathConfigRotateRoot(&b),
			pathConfigLease(&b),
			pathRoles(&b),
			pathListRoles(&b),
			pathStaticRoles(&b),
			pathListStaticRoles(&b),
			pathStaticCredentials(&b),
			pathUser(&b),
		},

		Secrets: []*framework.Secret{
			secretAccessKeys(&b),
		},

		InitializeFunc:    b.initialize,
		Invalidate:        b.invalidate,
		WALRollback:       b.walRollback,
		WALRollbackMinAge: minAwsUserRollbackAge,
		PeriodicFunc: func(ctx context.Context, req *logical.Request) error {
			if b.WriteSafeReplicationState() {
				return b.rotateExpiredStaticCreds(ctx, req)
			}
			return nil
		},
		RotateCredential: func(ctx context.Context, req *logical.Request) error {
			_, err := b.rotateRoot(ctx, req)
			return err
		},
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	// Function pointer used to override the IAM client creation for mocked testing
	// If set, this function will be called instead of creating real IAM clients
	nonCachedClientIAMFunc func(context.Context, logical.Storage, hclog.Logger, *staticRoleEntry) (iamAPI, error)

	// Mutex to protect access to reading and writing policies
	roleMutex sync.RWMutex

	// Mutex to protect access to iam/sts clients and client configs
	clientMutex sync.RWMutex

	// iamClient and stsClient hold configured iam and sts clients for reuse, and
	// to enable mocking with AWS iface for tests
	iamClient iamAPI
	stsClient stsAPI

	// the age of a static role's credential is tracked by a priority queue and handled
	// by the PeriodicFunc
	credRotationQueue *queue.PriorityQueue

	minAllowableRotationPeriod time.Duration
}

const backendHelp = `
The AWS backend dynamically generates AWS access keys for a set of
IAM policies. The AWS access keys have a configurable lease set and
are automatically revoked at the end of the lease.

After mounting this backend, credentials to generate IAM keys must
be configured with the "root" path and policies must be written using
the "roles/" endpoints before any access keys can be generated.
`

func (b *backend) invalidate(ctx context.Context, key string) {
	switch {
	case key == rootConfigPath:
		b.clearClients()
	}
}

// clearClients clears the backend's IAM and STS clients
func (b *backend) clearClients() {
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()
	b.iamClient = nil
	b.stsClient = nil
}

// clientIAM returns the configured IAM client. If nil, it constructs a new one
// and returns it, setting it the internal variable.
// entry is only needed when configuring the client to use for role assumption.
func (b *backend) clientIAM(ctx context.Context, s logical.Storage, entry *staticRoleEntry) (iamAPI, error) {
	b.clientMutex.RLock()
	if b.iamClient != nil {
		b.clientMutex.RUnlock()
		return b.iamClient, nil
	}

	// Upgrade the lock for writing
	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check client again, in the event that a client was being created while we
	// waited for Lock()
	if b.iamClient != nil {
		return b.iamClient, nil
	}

	iamClient, err := b.nonCachedClientIAM(ctx, s, b.Logger(), entry)
	if err != nil {
		return nil, err
	}

	b.iamClient = iamClient

	return b.iamClient, nil
}

func (b *backend) clientSTS(ctx context.Context, s logical.Storage) (stsAPI, error) {
	b.clientMutex.RLock()
	if b.stsClient != nil {
		b.clientMutex.RUnlock()
		return b.stsClient, nil
	}

	// Upgrade the lock for writing
	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check client again, in the event that a client was being created while we
	// waited for Lock()
	if b.stsClient != nil {
		return b.stsClient, nil
	}

	stsClient, err := b.nonCachedClientSTS(ctx, s, b.Logger())
	if err != nil {
		return nil, err
	}
	b.stsClient = stsClient

	return b.stsClient, nil
}

func (b *backend) initialize(ctx context.Context, request *logical.InitializationRequest) error {
	if !b.WriteSafeReplicationState() {
		b.Logger().Info("skipping populating rotation queue")
		return nil
	}
	b.Logger().Info("populating rotation queue")

	creds, err := request.Storage.List(ctx, pathStaticCreds+"/")
	if err != nil {
		return err
	}
	b.Logger().Debug(fmt.Sprintf("Adding %d items to the rotation queue", len(creds)))
	for _, roleName := range creds {
		if roleName == "" {
			continue
		}
		credPath := formatCredsStoragePath(roleName)
		credsEntry, err := request.Storage.Get(ctx, credPath)
		if err != nil {
			return fmt.Errorf("could not read credentials: %w", err)
		}
		if credsEntry == nil {
			continue
		}
		credentials := awsCredentials{}
		if err := credsEntry.DecodeJSON(&credentials); err != nil {
			return fmt.Errorf("failed to decode credentials: %w", err)
		}

		configEntry, err := request.Storage.Get(ctx, formatRoleStoragePath(roleName))
		if err != nil {
			return fmt.Errorf("could not read role: %w", err)
		}
		if configEntry == nil {
			continue
		}
		config := staticRoleEntry{}
		if err := configEntry.DecodeJSON(&config); err != nil {
			return fmt.Errorf("failed to decode role config: %w", err)
		}

		if credentials.Expiration == nil {
			expiration := time.Now().UTC().Add(config.RotationPeriod)
			credentials.Expiration = &expiration
			_, err := logical.StorageEntryJSON(credPath, creds)
			if err != nil {
				return fmt.Errorf("failed to marshal object to JSON: %w", err)
			}
			b.Logger().Debug("no known expiration time for credentials so resetting the expiration", "role", roleName, "new expiration", expiration)
		}

		err = b.credRotationQueue.Push(&queue.Item{
			Key:      config.Name,
			Value:    config,
			Priority: credentials.priority(config),
		})
		if err != nil {
			return fmt.Errorf("failed to add creds for role %s to queue: %w", roleName, err)
		}
	}
	return nil
}

// getNonCachedIAMClient returns an IAM client. In a test env, if a mocked client creation
// function is set (nonCachedClientIAMFunc), it will be used instead of the default client creation function.
// This allows us to mock AWS clients in tests.
func (b *backend) getNonCachedIAMClient(ctx context.Context, storage logical.Storage, cfg staticRoleEntry) (iamAPI, error) {
	if b.nonCachedClientIAMFunc != nil {
		return b.nonCachedClientIAMFunc(ctx, storage, b.Logger(), &cfg)
	}
	return b.nonCachedClientIAM(ctx, storage, b.Logger(), &cfg)
}
