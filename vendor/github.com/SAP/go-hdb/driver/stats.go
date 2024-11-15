package driver

// StatsHistogram represents statistic data in a histogram structure.
type StatsHistogram struct {
	// Count holds the number of measurements
	Count uint64
	// Sum holds the sum of the measurements.
	Sum float64
	// Buckets contains the count of measurements belonging to a bucket where the
	// value of the measurement is less or equal the bucket map key.
	Buckets map[float64]uint64
}

// Stats contains driver statistics.
type Stats struct {
	// Gauges
	OpenConnections  int // The number of current established driver connections.
	OpenTransactions int // The number of current open driver transactions.
	OpenStatements   int // The number of current open driver database statements.
	// Counters
	ReadBytes       uint64 // Total bytes read by client connection.
	WrittenBytes    uint64 // Total bytes written by client connection.
	SessionConnects uint64 // Total number of session connects (switch users).
	// Time histograms (Sum and upper bounds in Unit)
	TimeUnit  string                     // Time unit
	ReadTime  *StatsHistogram            // Time spent on reading from connection.
	WriteTime *StatsHistogram            // Time spent on writing to connection.
	AuthTime  *StatsHistogram            // Time spent on authentication.
	SQLTimes  map[string]*StatsHistogram // Time spent on different SQL statements.
}
