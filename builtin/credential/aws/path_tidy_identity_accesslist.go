// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathTidyIdentityAccessList() *framework.Path {
	return &framework.Path{
		Pattern: "tidy/identity-accesslist$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "identity-access-list",
			OperationVerb:   "tidy",
		},

		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": {
				Type:    framework.TypeDurationSecond,
				Default: 259200,
				Description: `The amount of extra time that must have passed beyond the identity's
expiration, before it is removed from the backend storage.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyIdentityAccessListUpdate,
			},
		},

		HelpSynopsis:    pathTidyIdentityAccessListSyn,
		HelpDescription: pathTidyIdentityAccessListDesc,
	}
}

// tidyAccessListIdentity is used to delete entries in the access list that are expired.
func (b *backend) tidyAccessListIdentity(ctx context.Context, req *logical.Request, safetyBuffer int) (*logical.Response, error) {
	// If we are a performance standby forward the request to the active node
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	if !atomic.CompareAndSwapUint32(b.tidyAccessListCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	s := req.Storage

	go func() {
		defer atomic.StoreUint32(b.tidyAccessListCASGuard, 0)

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		logger := b.Logger().Named("wltidy")

		bufferDuration := time.Duration(safetyBuffer) * time.Second

		doTidy := func() error {
			identities, err := s.List(ctx, identityAccessListStorage)
			if err != nil {
				return err
			}

			for _, instanceID := range identities {
				identityEntry, err := s.Get(ctx, identityAccessListStorage+instanceID)
				if err != nil {
					return fmt.Errorf("error fetching identity of instanceID %q: %w", instanceID, err)
				}

				if identityEntry == nil {
					return fmt.Errorf("identity entry for instanceID %q is nil", instanceID)
				}

				if identityEntry.Value == nil || len(identityEntry.Value) == 0 {
					return fmt.Errorf("found identity entry for instanceID %q but actual identity is empty", instanceID)
				}

				var result accessListIdentity
				if err := identityEntry.DecodeJSON(&result); err != nil {
					return err
				}

				if time.Now().After(result.ExpirationTime.Add(bufferDuration)) {
					if err := s.Delete(ctx, identityAccessListStorage+instanceID); err != nil {
						return fmt.Errorf("error deleting identity of instanceID %q from storage: %w", instanceID, err)
					}
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running access list tidy", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

// pathTidyIdentityAccessListUpdate is used to delete entries in the access list that are expired.
func (b *backend) pathTidyIdentityAccessListUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidyAccessListIdentity(ctx, req, data.Get("safety_buffer").(int))
}

const pathTidyIdentityAccessListSyn = `
Clean-up the access list instance identity entries.
`

const pathTidyIdentityAccessListDesc = `
When an instance identity is in the access list, the expiration time of the access list
entry is set based on the maximum 'max_ttl' value set on: the role, the role tag
and the backend's mount.

When this endpoint is invoked, all the entries that are expired will be deleted.
A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`
