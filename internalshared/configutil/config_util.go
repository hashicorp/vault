// +build !enterprise

package configutil

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

type EntSharedConfig struct{}

func (ec *EntSharedConfig) ParseConfig(list *ast.ObjectList) error {
	return nil
}

func ParseEntropy(result *SharedConfig, list *ast.ObjectList, blockName string) error {
	return nil
}
