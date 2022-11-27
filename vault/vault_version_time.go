package vault

import "time"

type VaultVersion struct {
	TimestampInstalled time.Time
	Version            string
	BuildDate          string
}
