// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listener

import "net"

// Provider is used to provide a given capability when requested remotely. They
// must return a connection that is bridged or an error.
type Provider func(capability string, meta map[string]string, conn net.Conn) error
