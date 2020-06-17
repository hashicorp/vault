package vault

import (
	"context"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// TODO: move emitMetrics into this file.

type kvMount struct {
	Namespace  *namespace.Namespace
	MountPoint string
	Version    string
	NumSecrets int
}

func (c *Core) findKvMounts() []*kvMount {
	mounts := make([]*kvMount, 0)

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	for _, entry := range c.mounts.Entries {
		if entry.Type == "kv" {
			version, ok := entry.Options["version"]
			if !ok {
				version = "1"
			}
			mounts = append(mounts, &kvMount{
				Namespace:  entry.namespace,
				MountPoint: entry.Path,
				Version:    version,
				NumSecrets: 0,
			})
		}
	}
	return mounts
}

func (c *Core) kvSecretGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	// Find all KV mounts
	mounts := c.findKvMounts()
	results := make([]metricsutil.GaugeLabelValues, len(mounts))

	// Context must have root namespace
	ctx = namespace.RootContext(ctx)

	// Route list requests to all the identified mounts.
	// (All of these will show up as activity in the vault.route metric.)
	// Then we have to explore each subdirectory.
	for i, m := range mounts {
		results[i].Labels = []metrics.Label{
			metricsutil.NamespaceLabel(m.Namespace),
			{"mount_point", m.MountPoint},
		}

		var subdirectories []string
		if m.Version == "1" {
			subdirectories = []string{m.MountPoint}
		} else {
			subdirectories = []string{m.MountPoint + "metadata/"}
		}

		for len(subdirectories) > 0 {
			// If collection was cancelled, return an empty array.
			select {
			case <-ctx.Done():
				return []metricsutil.GaugeLabelValues{}, nil
			default:
				break
			}

			currentDirectory := subdirectories[0]
			subdirectories = subdirectories[1:]

			listRequest := &logical.Request{
				Operation: logical.ListOperation,
				Path:      currentDirectory,
			}
			resp, err := c.router.Route(ctx, listRequest)
			if err != nil {
				c.logger.Error("failed to perform internal KV list", "mount_point", m.MountPoint, "error", err)
				// TODO: mark just one gauge as failed?
				return []metricsutil.GaugeLabelValues{}, err
			}
			rawKeys, ok := resp.Data["keys"]
			if !ok {
				continue
			}
			keys, ok := rawKeys.([]string)
			if !ok {
				c.logger.Error("keys are not a []string", "mount_point", m.MountPoint, "rawKeys", rawKeys)
				continue
			}
			for _, path := range keys {
				if path[len(path)-1] == '/' {
					subdirectories = append(subdirectories, currentDirectory+path)
				} else {
					m.NumSecrets += 1
				}
			}
		}

		results[i].Value = float32(m.NumSecrets)
	}

	return results, nil
}
