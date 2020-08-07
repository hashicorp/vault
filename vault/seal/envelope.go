package seal

import (
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

type Envelope struct {
	envelope *wrapping.Envelope
	once     sync.Once
}

func NewEnvelope() *Envelope {
	return &Envelope{}
}

func (e *Envelope) init() {
	e.envelope = new(wrapping.Envelope)
}

func (e *Envelope) Encrypt(plaintext, aad []byte) (*wrapping.EnvelopeInfo, error) {
	defer metrics.MeasureSince([]string{"seal", "envelope", "encrypt"}, time.Now())
	e.once.Do(e.init)

	return e.envelope.Encrypt(plaintext, aad)
}

func (e *Envelope) Decrypt(data *wrapping.EnvelopeInfo, aad []byte) ([]byte, error) {
	defer metrics.MeasureSince([]string{"seal", "envelope", "decrypt"}, time.Now())
	e.once.Do(e.init)

	return e.envelope.Decrypt(data, aad)
}
