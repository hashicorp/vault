package gocb

// Ping will ping a list of services and verify they are active and
// responding in an acceptable period of time.
func (c *Cluster) Ping(opts *PingOptions) (*PingResult, error) {
	return autoOpControl(c.diagnosticsController(), func(provider diagnosticsProvider) (*PingResult, error) {
		if opts == nil {
			opts = &PingOptions{}
		}

		return provider.Ping(opts)
	})
}
