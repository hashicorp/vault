package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var enterprisePaths = []*framework.Path{
	{
		Pattern: "replication/status",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: enterpriseOnlyErrorOperationFunc,
			},
		},
	},
}
