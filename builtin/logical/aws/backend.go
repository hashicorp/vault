// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

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
			// the following code is a modified version of the rotate-root method
			client, err := b.clientIAM(ctx, req.Storage)
			if err != nil {
				return err
			}
			if client == nil {
				return fmt.Errorf("nil IAM client")
			}

			b.clientMutex.Lock()
			defer b.clientMutex.Unlock()

			rawRootConfig, err := req.Storage.Get(ctx, "config/root")
			if err != nil {
				return err
			}
			if rawRootConfig == nil {
				return fmt.Errorf("no configuration found for config/root")
			}
			var config rootConfig
			if err := rawRootConfig.DecodeJSON(&config); err != nil {
				return fmt.Errorf("error reading root configuration: %w", err)
			}

			if config.AccessKey == "" || config.SecretKey == "" {
				return fmt.Errorf("cannot call config/rotate-root when either access_key or secret_key is empty")
			}

			var getUserInput iam.GetUserInput // empty input means get current user
			getUserRes, err := client.GetUserWithContext(ctx, &getUserInput)
			if err != nil {
				return fmt.Errorf("error calling GetUser: %w", err)
			}
			if getUserRes == nil {
				return fmt.Errorf("nil response from GetUser")
			}
			if getUserRes.User == nil {
				return fmt.Errorf("nil user returned from GetUser")
			}
			if getUserRes.User.UserName == nil {
				return fmt.Errorf("nil UserName returned from GetUser")
			}

			createAccessKeyInput := iam.CreateAccessKeyInput{
				UserName: getUserRes.User.UserName,
			}
			createAccessKeyRes, err := client.CreateAccessKeyWithContext(ctx, &createAccessKeyInput)
			if err != nil {
				return fmt.Errorf("error calling CreateAccessKey: %w", err)
			}
			if createAccessKeyRes.AccessKey == nil {
				return fmt.Errorf("nil response from CreateAccessKey")
			}
			if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
				return fmt.Errorf("nil AccessKeyId or SecretAccessKey returned from CreateAccessKey")
			}

			oldAccessKey := config.AccessKey

			config.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
			config.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

			newEntry, err := logical.StorageEntryJSON("config/root", config)
			if err != nil {
				return fmt.Errorf("error generating new config/root JSON: %w", err)
			}
			if err := req.Storage.Put(ctx, newEntry); err != nil {
				return fmt.Errorf("error saving new config/root: %w", err)
			}

			b.iamClient = nil
			b.stsClient = nil

			deleteAccessKeyInput := iam.DeleteAccessKeyInput{
				AccessKeyId: aws.String(oldAccessKey),
				UserName:    getUserRes.User.UserName,
			}
			_, err = client.DeleteAccessKeyWithContext(ctx, &deleteAccessKeyInput)
			if err != nil {
				return fmt.Errorf("error deleting old access key: %w", err)
			}

			return nil
		},
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	// Mutex to protect access to reading and writing policies
	roleMutex sync.RWMutex

	// Mutex to protect access to iam/sts clients and client configs
	clientMutex sync.RWMutex

	// iamClient and stsClient hold configured iam and sts clients for reuse, and
	// to enable mocking with AWS iface for tests
	iamClient iamiface.IAMAPI
	stsClient stsiface.STSAPI

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
// and returns it, setting it the internal variable
func (b *backend) clientIAM(ctx context.Context, s logical.Storage) (iamiface.IAMAPI, error) {
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

	iamClient, err := b.nonCachedClientIAM(ctx, s, b.Logger())
	if err != nil {
		return nil, err
	}
	b.iamClient = iamClient

	return b.iamClient, nil
}

func (b *backend) clientSTS(ctx context.Context, s logical.Storage) (stsiface.STSAPI, error) {
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
