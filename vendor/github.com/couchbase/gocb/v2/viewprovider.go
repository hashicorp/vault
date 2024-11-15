package gocb

type viewProvider interface {
	ViewQuery(designDoc string, viewName string, opts *ViewOptions) (*ViewResult, error)
}

type viewRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}
