// +build !enterprise

package server

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

type entConfig struct{}

func (ec *entConfig) parseConfig(list *ast.ObjectList) error {
	return nil
}

func (ec entConfig) Merge(ec2 entConfig) entConfig {
	result := entConfig{}
	return result
}

func (ec entConfig) Sanitized() map[string]interface{} {
	return nil
}
