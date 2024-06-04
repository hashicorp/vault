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

const (
	denyListRoletagStorage = "blacklist/roletag/"
)

func (b *backend) pathTidyRoletagDenyList() *framework.Path {
	return &framework.Path{
		Pattern: "tidy/roletag-denylist$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "role-tag-deny-list",
			OperationVerb:   "tidy",
		},

		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": {
				Type:    framework.TypeDurationSecond,
				Default: 259200, // 72h
				Description: `The amount of extra time that must have passed beyond the roletag
expiration, before it is removed from the backend storage.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyRoletagDenylistUpdate,
			},
		},

		HelpSynopsis:    pathTidyRoletagDenylistSyn,
		HelpDescription: pathTidyRoletagDenylistDesc,
	}
}

// tidyDenyListRoleTag is used to clean-up the entries in the role tag deny list.
func (b *backend) tidyDenyListRoleTag(ctx context.Context, req *logical.Request, safetyBuffer int) (*logical.Response, error) {
	// If we are a performance standby forward the request to the active node
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	if !atomic.CompareAndSwapUint32(b.tidyDenyListCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	s := req.Storage

	go func() {
		defer atomic.StoreUint32(b.tidyDenyListCASGuard, 0)

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		logger := b.Logger().Named("bltidy")

		bufferDuration := time.Duration(safetyBuffer) * time.Second

		doTidy := func() error {
			tags, err := s.List(ctx, denyListRoletagStorage)
			if err != nil {
				return err
			}

			for _, tag := range tags {
				tagEntry, err := s.Get(ctx, denyListRoletagStorage+tag)
				if err != nil {
					return fmt.Errorf("error fetching tag %q: %w", tag, err)
				}

				if tagEntry == nil {
					return fmt.Errorf("tag entry for tag %q is nil", tag)
				}

				if tagEntry.Value == nil || len(tagEntry.Value) == 0 {
					return fmt.Errorf("found entry for tag %q but actual tag is empty", tag)
				}

				var result roleTagBlacklistEntry
				if err := tagEntry.DecodeJSON(&result); err != nil {
					return err
				}

				if time.Now().After(result.ExpirationTime.Add(bufferDuration)) {
					if err := s.Delete(ctx, denyListRoletagStorage+tag); err != nil {
						return fmt.Errorf("error deleting tag %q from storage: %w", tag, err)
					}
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running deny list tidy", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

// pathTidyRoletagDenylistUpdate is used to clean-up the entries in the role tag deny list.
func (b *backend) pathTidyRoletagDenylistUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidyDenyListRoleTag(ctx, req, data.Get("safety_buffer").(int))
}

const pathTidyRoletagDenylistSyn = `
Clean-up the deny list role tag entries.
`

const pathTidyRoletagDenylistDesc = `
When a role tag is deny listed, the expiration time of the deny list entry is
set based on the maximum 'max_ttl' value set on: the role, the role tag and the
backend's mount.

When this endpoint is invoked, all the entries that are expired will be deleted.
A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`
