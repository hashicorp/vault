package gocb

type analyticsProvider interface {
	AnalyticsQuery(statement string, scope *Scope, opts *AnalyticsOptions) (*AnalyticsResult, error)
}

type analyticsRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}
