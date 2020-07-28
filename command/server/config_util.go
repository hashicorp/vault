// +build !enterprise

package server

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

type entConfig struct {
}

func (ec *entConfig) parseConfig(list *ast.ObjectList) error {
	return nil
}
