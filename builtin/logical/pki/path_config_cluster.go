// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigCluster(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/cluster",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

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
			"aia_path": {
				Type: framework.TypeString,
				Description: `Optional URI to this mount's AIA distribution
point; may refer to an external non-Vault responder. This is for resolving AIA
URLs and providing the {{cluster_aia_path}} template parameter and will not
be used for other purposes. As such, unlike path above, this could safely
be an insecure transit mechanism (like HTTP without TLS).

For example: http://cdn.example.com/pr1/pki`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "cluster",
				},
				Callback: b.pathWriteCluster,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
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
							"aia_path": {
								Type: framework.TypeString,
								Description: `Optional URI to this mount's AIA distribution
point; may refer to an external non-Vault responder. This is for resolving AIA
URLs and providing the {{cluster_aia_path}} template parameter and will not
be used for other purposes. As such, unlike path above, this could safely
be an insecure transit mechanism (like HTTP without TLS).

For example: http://cdn.example.com/pr1/pki`,
							},
						},
					}},
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadCluster,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "cluster-configuration",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
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
								Required: true,
							},
							"aia_path": {
								Type: framework.TypeString,
								Description: `Optional URI to this mount's AIA distribution
point; may refer to an external non-Vault responder. This is for resolving AIA
URLs and providing the {{cluster_aia_path}} template parameter and will not
be used for other purposes. As such, unlike path above, this could safely
be an insecure transit mechanism (like HTTP without TLS).

For example: http://cdn.example.com/pr1/pki`,
							},
						},
					}},
				},
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
			"path":     cfg.Path,
			"aia_path": cfg.AIAPath,
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

	if value, ok := data.GetOk("path"); ok {
		cfg.Path = value.(string)

		// This field is required by ACME, if ever we allow un-setting in the
		// future, this code will need to verify that ACME is not enabled.
		if !govalidator.IsURL(cfg.Path) {
			return nil, fmt.Errorf("invalid, non-URL path given to cluster: %v", cfg.Path)
		}
	}

	if value, ok := data.GetOk("aia_path"); ok {
		cfg.AIAPath = value.(string)
		if !govalidator.IsURL(cfg.AIAPath) {
			return nil, fmt.Errorf("invalid, non-URL aia_path given to cluster: %v", cfg.AIAPath)
		}
	}

	if err := sc.writeClusterConfig(cfg); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"path":     cfg.Path,
			"aia_path": cfg.AIAPath,
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
