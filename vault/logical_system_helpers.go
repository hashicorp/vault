package vault

import (
	"fmt"
	"strings"
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
			if newMax == nil && *newDefault > meConfig.MaxLeaseTTL {
				return fmt.Errorf("new backend default lease TTL of %d greater than backend max lease TTL of %d",
					int(newDefault.Seconds()), int(meConfig.MaxLeaseTTL.Seconds()))
			}
		}
	}

	origMax := meConfig.MaxLeaseTTL
	origDefault := meConfig.DefaultLeaseTTL

	if newMax != nil {
		meConfig.MaxLeaseTTL = *newMax
	}
	if newDefault != nil {
		meConfig.DefaultLeaseTTL = *newDefault
	}

	// Update the mount table
	var err error
	switch {
	case strings.HasPrefix(path, "auth/"):
		err = b.Core.persistAuth(b.Core.auth)
	default:
		err = b.Core.persistMounts(b.Core.mounts)
	}
	if err != nil {
		meConfig.MaxLeaseTTL = origMax
		meConfig.DefaultLeaseTTL = origDefault
		return fmt.Errorf("failed to update mount table, rolling back TTL changes")
	}

	if b.Core.logger.IsInfo() {
		b.Core.logger.Info("core: mount tuning successful", "path", path)
	}

	return nil
}
