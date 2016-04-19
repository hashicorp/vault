package aws

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var walRollbackMap = map[string]framework.WALRollbackFunc{
	"user": pathUserRollback,
}

func walRollback(req *logical.Request, kind string, data interface{}) error {
	f, ok := walRollbackMap[kind]
	if !ok {
		return fmt.Errorf("unknown type to rollback")
	}

	return f(req, kind, data)
}
