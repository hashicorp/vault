// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	walRollbackMap := map[string]framework.WALRollbackFunc{
		"user": b.pathUserRollback,
	}

	if !b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
		return nil
	}

	f, ok := walRollbackMap[kind]
	if !ok {
		return fmt.Errorf("unknown type to rollback")
	}

	return f(ctx, req, kind, data)
}
