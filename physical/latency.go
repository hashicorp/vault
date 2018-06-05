package physical

import (
	"context"
	"math/rand"
	"time"

	log "github.com/hashicorp/go-hclog"
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

// Verify LatencyInjector satisfies the correct interfaces
var _ Backend = (*LatencyInjector)(nil)
var _ Transactional = (*TransactionalLatencyInjector)(nil)

// NewLatencyInjector returns a wrapped physical backend to simulate latency
func NewLatencyInjector(b Backend, latency time.Duration, jitter int, logger log.Logger) *LatencyInjector {
	if jitter < 0 || jitter > 100 {
		jitter = DefaultJitterPercent
	}
	logger.Info("creating latency injector")

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
func (l *LatencyInjector) Put(ctx context.Context, entry *Entry) error {
	l.addLatency()
	return l.backend.Put(ctx, entry)
}

// Get is a latent get request
func (l *LatencyInjector) Get(ctx context.Context, key string) (*Entry, error) {
	l.addLatency()
	return l.backend.Get(ctx, key)
}

// Delete is a latent delete request
func (l *LatencyInjector) Delete(ctx context.Context, key string) error {
	l.addLatency()
	return l.backend.Delete(ctx, key)
}

// List is a latent list request
func (l *LatencyInjector) List(ctx context.Context, prefix string) ([]string, error) {
	l.addLatency()
	return l.backend.List(ctx, prefix)
}

// Transaction is a latent transaction request
func (l *TransactionalLatencyInjector) Transaction(ctx context.Context, txns []*TxnEntry) error {
	l.addLatency()
	return l.Transactional.Transaction(ctx, txns)
}
