package gocb

type queryProvider interface {
	Query(statement string, s *Scope, opts *QueryOptions) (*QueryResult, error)
}

type queryRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
	PreparedName() (string, error)
	Endpoint() string
}
