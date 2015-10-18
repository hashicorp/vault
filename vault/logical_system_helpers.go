package vault

import (
	"errors"
	"fmt"
	"time"
)

// tuneMount is used to set config on a mount point
func (b *SystemBackend) tuneMountTTLs(path string, meConfig *MountConfig, newDefault, newMax *time.Duration) error {
	if newDefault == nil && newMax == nil {
		return nil
	}
	if newDefault == nil && newMax != nil &&
		*newMax == meConfig.MaxLeaseTTL {
		return nil
	}
	if newMax == nil && newDefault != nil &&
		*newDefault == meConfig.DefaultLeaseTTL {
		return nil
	}
	if newMax != nil && newDefault != nil &&
		*newDefault == meConfig.DefaultLeaseTTL &&
		*newMax == meConfig.MaxLeaseTTL {
		return nil
	}

	if newMax != nil && newDefault != nil && *newMax < *newDefault {
		return fmt.Errorf("new backend max lease TTL of %d less than new backend default lease TTL of %d",
			int(newMax.Seconds()), int(newDefault.Seconds()))
	}

	if newMax != nil && newDefault == nil {
		if meConfig.DefaultLeaseTTL != 0 && *newMax < meConfig.DefaultLeaseTTL {
			return fmt.Errorf("new backend max lease TTL of %d less than backend default lease TTL of %d",
				int(newMax.Seconds()), int(meConfig.DefaultLeaseTTL.Seconds()))
		}
	}

	if newDefault != nil {
		if meConfig.MaxLeaseTTL == 0 {
			if newMax == nil && *newDefault > b.Core.maxLeaseTTL {
				return fmt.Errorf("new backend default lease TTL of %d greater than system max lease TTL of %d",
					int(newDefault.Seconds()), int(b.Core.maxLeaseTTL.Seconds()))
			}
		} else {
			if meConfig.MaxLeaseTTL < *newDefault {
				return fmt.Errorf("new backend default lease TTL of %d greater than backend max lease TTL of %d",
					int(newDefault.Seconds()), int(meConfig.MaxLeaseTTL.Seconds()))
			}
		}
	}

	if newMax != nil {
		meConfig.MaxLeaseTTL = *newMax
	}
	if newDefault != nil {
		meConfig.DefaultLeaseTTL = *newDefault
	}

	// Update the mount table
	if err := b.Core.persistMounts(b.Core.mounts); err != nil {
		return errors.New("failed to update mount table")
	}

	b.Core.logger.Printf("[INFO] core: tuned '%s'", path)

	return nil
}
