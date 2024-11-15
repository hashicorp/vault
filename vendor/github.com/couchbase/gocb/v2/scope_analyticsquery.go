package gocb

// AnalyticsQuery executes the analytics query statement on the server, constraining the query to the bucket and scope.
func (s *Scope) AnalyticsQuery(statement string, opts *AnalyticsOptions) (*AnalyticsResult, error) {
	return autoOpControl(s.analyticsController(), func(provider analyticsProvider) (*AnalyticsResult, error) {
		if opts == nil {
			opts = &AnalyticsOptions{}
		}

		return provider.AnalyticsQuery(statement, s, opts)
	})
}
