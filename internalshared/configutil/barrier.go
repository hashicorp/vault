package configutil

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"time"
)

type BarrierConfig struct {
	KeyRotationMaxOperations int64 `hcl:"key_rotation_max_operations"`
	KeyRotationInterval      time.Duration
	KeyRotationIntervalRaw   interface{} `hcl:"key_rotation_interval"`
}

func parseBarrier(result *SharedConfig, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'barrier' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

	if result.Barrier == nil {
		result.Barrier = &BarrierConfig{}
	}

	if err := hcl.DecodeObject(&result.Barrier, item.Val); err != nil {
		return multierror.Prefix(err, "telemetry:")
	}

	if result.Barrier.KeyRotationIntervalRaw != nil {
		var err error
		if result.Barrier.KeyRotationInterval, err = parseutil.ParseDurationSecond(result.Barrier.KeyRotationIntervalRaw); err != nil {
			return err
		}
		result.Barrier.KeyRotationIntervalRaw = nil
	} else {
		result.Barrier.KeyRotationInterval = 0
	}
	return nil
}
