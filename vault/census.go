//go:build !enterprise

package vault

// CensusAgent is a stub for OSS
type CensusReporter struct{}

// setupCensusAgent is a stub for OSS.
func (c *Core) setupCensusAgent() error { return nil }
