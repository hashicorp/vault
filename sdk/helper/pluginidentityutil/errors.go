// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginidentityutil

import "errors"

var ErrPluginWorkloadIdentityUnsupported = errors.New("plugin workload identity not supported in Vault community edition")
