package aws

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var walRollbackMap = map[string]framework.WALRollbackFunc{
	"user": pathUserRollback,
}

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	if !b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformancePrimary) {
		return nil
	}

	f, ok := walRollbackMap[kind]
	if !ok {
		return fmt.Errorf("unknown type to rollback")
	}

	return f(ctx, req, kind, data)
}
