package seal

import (
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"sync"
	"time"
)

// SealWrapper contains a Wrapper and related information needed by the seal that uses it.
type SealWrapper struct {
	Wrapper  wrapping.Wrapper
	Priority int
	Name     string

	// sealConfigType is the KMS.Type of this wrapper. It is a string rather than a SealConfigType
	// to avoid a circular go package depency
	SealConfigType string

	// Disabled indicates, when true indicates that this wrapper should only be used for decryption.
	Disabled bool

	HcLock          sync.RWMutex
	LastHealthCheck time.Time
	LastSeenHealthy time.Time
	Healthy         bool
}
