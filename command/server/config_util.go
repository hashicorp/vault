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

func parseEntropy(result *Config, list *ast.ObjectList, blockName string) error {
	return nil
}
