package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var (
	invalidateMFAConfig = func(context.Context, *SystemBackend, string) {}

	sysInvalidate = func(b *SystemBackend) func(context.Context, string) {
		return nil
	}

	getSystemSchemas = func() []func() *memdb.TableSchema { return nil }

	getEGPListResponseKeyInfo = func(*SystemBackend, *namespace.Namespace) map[string]interface{} { return nil }
	addSentinelPolicyData     = func(map[string]interface{}, *Policy) {}
	inputSentinelPolicyData   = func(*framework.FieldData, *Policy) *logical.Response { return nil }

	controlGroupUnwrap = func(context.Context, *SystemBackend, string, bool) (string, error) {
		return "", errors.New("control groups unavailable")
	}

	pathInternalUINamespacesRead = func(b *SystemBackend) framework.OperationFunc {
		return func(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
			// Short-circuit here if there's no client token provided
			if req.ClientToken == "" {
				return nil, fmt.Errorf("client token empty")
			}

			// Load the ACL policies so we can check for access and filter namespaces
			_, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(ctx, req)
			if err != nil {
				return nil, err
			}
			if entity != nil && entity.Disabled {
				b.logger.Warn("permission denied as the entity on the token is disabled")
				return nil, logical.ErrPermissionDenied
			}
			if te != nil && te.EntityID != "" && entity == nil {
				b.logger.Warn("permission denied as the entity on the token is invalid")
				return nil, logical.ErrPermissionDenied
			}

			return logical.ListResponse([]string{""}), nil
		}
	}

	pathLicenseRead = func(b *SystemBackend) framework.OperationFunc {
		return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
			return nil, nil
		}
	}

	pathLicenseUpdate = func(b *SystemBackend) framework.OperationFunc {
		return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
			return nil, nil
		}
	}

	entPaths = func(b *SystemBackend) []*framework.Path {
		return []*framework.Path{
			{
				Pattern: "replication/status",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
						resp := &logical.Response{
							Data: map[string]interface{}{
								"mode": "disabled",
							},
						}
						return resp, nil
					},
				},
			},
		}
	}

	checkRaw = func(b *SystemBackend, path string) error { return nil }
)

// tuneMount is used to set config on a mount point
func (b *SystemBackend) tuneMountTTLs(ctx context.Context, path string, me *MountEntry, newDefault, newMax time.Duration) error {
	zero := time.Duration(0)

	switch {
	case newDefault == zero && newMax == zero:
		// No checks needed

	case newDefault == zero && newMax != zero:
		// No default/max conflict, no checks needed

	case newDefault != zero && newMax == zero:
		// No default/max conflict, no checks needed

	case newDefault != zero && newMax != zero:
		if newMax < newDefault {
			return fmt.Errorf("backend max lease TTL of %d would be less than backend default lease TTL of %d", int(newMax.Seconds()), int(newDefault.Seconds()))
		}
	}

	origMax := me.Config.MaxLeaseTTL
	origDefault := me.Config.DefaultLeaseTTL

	me.Config.MaxLeaseTTL = newMax
	me.Config.DefaultLeaseTTL = newDefault

	// Update the mount table
	var err error
	switch {
	case strings.HasPrefix(path, credentialRoutePrefix):
		err = b.Core.persistAuth(ctx, b.Core.auth, &me.Local)
	default:
		err = b.Core.persistMounts(ctx, b.Core.mounts, &me.Local)
	}
	if err != nil {
		me.Config.MaxLeaseTTL = origMax
		me.Config.DefaultLeaseTTL = origDefault
		return fmt.Errorf("failed to update mount table, rolling back TTL changes")
	}
	if b.Core.logger.IsInfo() {
		b.Core.logger.Info("mount tuning of leases successful", "path", path)
	}

	return nil
}
