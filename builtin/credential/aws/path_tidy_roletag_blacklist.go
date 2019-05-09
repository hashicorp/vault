package awsauth

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathTidyRoletagBlacklist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy/roletag-blacklist$",
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 259200, // 72h
				Description: `The amount of extra time that must have passed beyond the roletag
expiration, before it is removed from the backend storage.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidyRoletagBlacklistUpdate,
		},

		HelpSynopsis:    pathTidyRoletagBlacklistSyn,
		HelpDescription: pathTidyRoletagBlacklistDesc,
	}
}

// tidyBlacklistRoleTag is used to clean-up the entries in the role tag blacklist.
func (b *backend) tidyBlacklistRoleTag(ctx context.Context, req *logical.Request, safetyBuffer int) (*logical.Response, error) {
	// If we are a performance standby forward the request to the active node
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	if !atomic.CompareAndSwapUint32(b.tidyBlacklistCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	s := req.Storage

	go func() {
		defer atomic.StoreUint32(b.tidyBlacklistCASGuard, 0)

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		logger := b.Logger().Named("bltidy")

		bufferDuration := time.Duration(safetyBuffer) * time.Second

		doTidy := func() error {
			tags, err := s.List(ctx, "blacklist/roletag/")
			if err != nil {
				return err
			}

			for _, tag := range tags {
				tagEntry, err := s.Get(ctx, "blacklist/roletag/"+tag)
				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("error fetching tag %q: {{err}}", tag), err)
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
					if err := s.Delete(ctx, "blacklist/roletag/"+tag); err != nil {
						return errwrap.Wrapf(fmt.Sprintf("error deleting tag %q from storage: {{err}}", tag), err)
					}
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running blacklist tidy", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

// pathTidyRoletagBlacklistUpdate is used to clean-up the entries in the role tag blacklist.
func (b *backend) pathTidyRoletagBlacklistUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidyBlacklistRoleTag(ctx, req, data.Get("safety_buffer").(int))
}

const pathTidyRoletagBlacklistSyn = `
Clean-up the blacklist role tag entries.
`

const pathTidyRoletagBlacklistDesc = `
When a role tag is blacklisted, the expiration time of the blacklist entry is
set based on the maximum 'max_ttl' value set on: the role, the role tag and the
backend's mount.

When this endpoint is invoked, all the entries that are expired will be deleted.
A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`
