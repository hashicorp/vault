package gore

// Config keeps some default configurations. Time is measured in second
var Config = &struct {
	ConnectTimeout  int
	RequestTimeout  int
	ReconnectTime   int
	PoolInitialSize int
	PoolMaximumSize int
}{
	ConnectTimeout:  5,
	RequestTimeout:  10,
	ReconnectTime:   2,
	PoolInitialSize: 5,
	PoolMaximumSize: 10,
}
