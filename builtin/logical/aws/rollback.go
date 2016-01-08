package aws

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var rollbackMap = map[string]framework.RollbackFunc{
	"user": pathUserRollback,
}

func rollback(req *logical.Request, kind string, data interface{}) error {
	f, ok := rollbackMap[kind]
	if !ok {
		return fmt.Errorf("unknown type to rollback")
	}

	return f(req, kind, data)
}
