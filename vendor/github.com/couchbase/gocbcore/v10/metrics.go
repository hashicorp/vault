package gocbcore

// Meter handles metrics information for SDK operations.
type Meter interface {
	Counter(name string, tags map[string]string) (Counter, error)
	ValueRecorder(name string, tags map[string]string) (ValueRecorder, error)
}

// Counter is used for incrementing a synchronous count metric.
type Counter interface {
	IncrementBy(num uint64)
}

// ValueRecorder is used for grouping synchronous count metrics.
type ValueRecorder interface {
	RecordValue(val uint64)
}
