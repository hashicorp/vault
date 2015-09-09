package vault

import (
	"errors"
	"fmt"
)

// tuneMount is used to set config on a mount point
func (b *SystemBackend) tuneMountTTLs(path string, meConfig, newConfig *MountConfig) error {
	if meConfig.MaxLeaseTTL == newConfig.MaxLeaseTTL &&
		meConfig.DefaultLeaseTTL == newConfig.DefaultLeaseTTL {
		return nil
	}

	if meConfig.DefaultLeaseTTL != 0 {
		if newConfig.MaxLeaseTTL < meConfig.DefaultLeaseTTL {
			if newConfig.DefaultLeaseTTL == 0 {
				return fmt.Errorf("New backend max lease TTL of %d less than backend default lease TTL of %d",
					newConfig.MaxLeaseTTL, meConfig.DefaultLeaseTTL)
			}
			if newConfig.MaxLeaseTTL < newConfig.DefaultLeaseTTL {
				return fmt.Errorf("New backend max lease TTL of %d less than new backend default lease TTL of %d",
					newConfig.MaxLeaseTTL, newConfig.DefaultLeaseTTL)
			}
		}
	}

	if meConfig.MaxLeaseTTL == 0 {
		if newConfig.DefaultLeaseTTL > b.Core.maxLeaseTTL {
			return fmt.Errorf("New backend default lease TTL of %d greater than system max lease TTL of %d",
				newConfig.DefaultLeaseTTL, b.Core.maxLeaseTTL)
		}
	} else {
		if meConfig.MaxLeaseTTL < newConfig.DefaultLeaseTTL {
			return fmt.Errorf("New backend default lease TTL of %d greater than backend max lease TTL of %d",
				newConfig.DefaultLeaseTTL, meConfig.MaxLeaseTTL)
		}
	}

	meConfig.MaxLeaseTTL = newConfig.MaxLeaseTTL
	meConfig.DefaultLeaseTTL = newConfig.DefaultLeaseTTL

	// Update the mount table
	if err := b.Core.persistMounts(b.Core.mounts); err != nil {
		return errors.New("failed to update mount table")
	}

	b.Core.logger.Printf("[INFO] core: tuned '%s'", path)

	return nil
}
