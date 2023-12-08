// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki_backend

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

type SystemViewGetter interface {
	System() logical.SystemView
}

type MountInfo interface {
	BackendUUID() string
}

type Logger interface {
	Logger() log.Logger
}

type CertificateCounter interface {
	IsInitialized() bool
	IncrementTotalCertificatesCount(certsCounted bool, newSerial string)
}
