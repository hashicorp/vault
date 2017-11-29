package physical

import (
	"math/rand"
	"time"

	log "github.com/mgutz/logxi/v1"
)

const (
	// DefaultJitterPercent is used if no cache size is specified for NewCache
	DefaultJitterPercent = 20
)

// LatencyInjector is used to add latency into underlying physical requests
type LatencyInjector struct {
	backend       Backend
	latency       time.Duration
	jitterPercent int
	random        *rand.Rand
}

// TransactionalLatencyInjector is the transactional version of the latency
// injector
type TransactionalLatencyInjector struct {
	*LatencyInjector
	Transactional
}

// NewLatencyInjector returns a wrapped physical backend to simulate latency
func NewLatencyInjector(b Backend, latency time.Duration, jitter int, logger log.Logger) *LatencyInjector {
	if jitter < 0 || jitter > 100 {
		jitter = DefaultJitterPercent
	}
	logger.Info("physical/latency: creating latency injector")

	return &LatencyInjector{
		backend:       b,
		latency:       latency,
		jitterPercent: jitter,
		random:        rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
	}
}

// NewTransactionalLatencyInjector creates a new transactional LatencyInjector
func NewTransactionalLatencyInjector(b Backend, latency time.Duration, jitter int, logger log.Logger) *TransactionalLatencyInjector {
	return &TransactionalLatencyInjector{
		LatencyInjector: NewLatencyInjector(b, latency, jitter, logger),
		Transactional:   b.(Transactional),
	}
}

func (l *LatencyInjector) addLatency() {
	// Calculate a value between 1 +- jitter%
	min := 100 - l.jitterPercent
	max := 100 + l.jitterPercent
	percent := l.random.Intn(max-min) + min
	latencyDuration := time.Duration(int(l.latency) * percent / 100)
	time.Sleep(latencyDuration)
}

// Put is a latent put request
func (l *LatencyInjector) Put(entry *Entry) error {
	l.addLatency()
	return l.backend.Put(entry)
}

// Get is a latent get request
func (l *LatencyInjector) Get(key string) (*Entry, error) {
	l.addLatency()
	return l.backend.Get(key)
}

// Delete is a latent delete request
func (l *LatencyInjector) Delete(key string) error {
	l.addLatency()
	return l.backend.Delete(key)
}

// List is a latent list request
func (l *LatencyInjector) List(prefix string) ([]string, error) {
	l.addLatency()
	return l.backend.List(prefix)
}

// Transaction is a latent transaction request
func (l *TransactionalLatencyInjector) Transaction(txns []*TxnEntry) error {
	l.addLatency()
	return l.Transactional.Transaction(txns)
}
