package pki

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigCluster(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/cluster",
		Fields: map[string]*framework.FieldSchema{
			"path": {
				Type: framework.TypeString,
				Description: `Canonical URI to this mount on this performance
replication cluster's external address. This is for resolving AIA URLs and
providing the {{cluster_path}} template parameter but might be used for other
purposes in the future.

This should only point back to this particular PR replica and should not ever
point to another PR cluster. It may point to any node in the PR replica,
including standby nodes, and need not always point to the active node.

For example: https://pr1.vault.example.com:8200/v1/pki`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathWriteCluster,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadCluster,
			},
		},

		HelpSynopsis:    pathConfigClusterHelpSyn,
		HelpDescription: pathConfigClusterHelpDesc,
	}
}

func (b *backend) pathReadCluster(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"path": cfg.Path,
		},
	}

	return resp, nil
}

func (b *backend) pathWriteCluster(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, err
	}

	cfg.Path = data.Get("path").(string)
	if !govalidator.IsURL(cfg.Path) {
		return nil, fmt.Errorf("invalid, non-URL path given to cluster: %v", cfg.Path)
	}

	if err := sc.writeClusterConfig(cfg); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"path": cfg.Path,
		},
	}

	return resp, nil
}

const pathConfigClusterHelpSyn = `
Set cluster-local configuration, including address to this PR cluster.
`

const pathConfigClusterHelpDesc = `
This path allows you to set cluster-local configuration, including the
URI to this performance replication cluster. This allows you to use
templated AIA URLs with /config/urls and /issuer/:issuer_ref, setting the
reference to the cluster's URI.

Only one address can be specified for any given cluster.
`
