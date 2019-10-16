package server

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

var(
	onEnterprise = false
	parseEntropy = parseEntropyOSS
)

func parseEntropyOSS(result *Config, list *ast.ObjectList, blockName string) error {
	return nil
}
