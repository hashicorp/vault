package server

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
)

var(
	onEnterprise = false
	parseEntropy = parseEntropyOSS
)

func parseEntropyOSS(result *Config, list *ast.ObjectList, blockName string) error {
	return fmt.Errorf("%q is an enterprise feature", blockName)
}
