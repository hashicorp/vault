package pluginidentityutil

import "errors"

var ErrPluginWorkloadIdentityUnsupported = errors.New("plugin workload identity not supported in Vault community edition")
