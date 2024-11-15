package gocb

// Query executes the query statement on the server, constraining the query to the bucket and scope.
func (s *Scope) Query(statement string, opts *QueryOptions) (*QueryResult, error) {
	return autoOpControl(s.queryController(), func(provider queryProvider) (*QueryResult, error) {
		if opts == nil {
			opts = &QueryOptions{}
		}

		if opts.AsTransaction != nil {
			return s.getTransactions().singleQuery(statement, s, *opts)
		}

		return provider.Query(statement, s, opts)
	})
}
