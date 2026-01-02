package api

import (
    "crypto/tls"
    "errors"
	"fmt"

    "github.com/go-piv/piv-go/v2/piv"
)

type PIVConfig struct {
	Reader   string
	Slot     string
	PIN      string
	Intermediates [][]byte
}

func LoadPIVCertAndSigner(cfg PIVConfig) (tls.Certificate, *func() error, error) {
	slot := piv.SlotAuthentication
	switch cfg.Slot {
	case "", "9a": slot = piv.SlotAuthentication
	case "9c":     slot = piv.SlotSignature
	case "9d":     slot = piv.SlotKeyManagement
	case "9e":     slot = piv.SlotCardAuthentication
	default:
		return tls.Certificate{}, nil, fmt.Errorf("unsupported slot %q", cfg.Slot)
	}

	cards, err := piv.Cards()
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	var yk *piv.YubiKey
	for _, c := range cards {
		if cfg.Reader == "" || containsFold(c, cfg.Reader) {
			if yk, err = piv.Open(c); err == nil { break }
		}
	}
	if yk == nil {
		return tls.Certificate{}, nil, errors.New("no PIV reader found (or unable to open)")
	}
	close := yk.Close
	// defer yk.Close() // Would be good, but not sure we can ...?

	// Read the leaf certificate from the slot
	leaf, err := yk.Certificate(slot)
	if err != nil {
		return tls.Certificate{}, &close, fmt.Errorf("reading slot %s cert: %w", cfg.Slot, err)
	}

	// Construct a private key based on the key
	signer, err := yk.PrivateKey(slot, leaf.PublicKey, piv.KeyAuth{PIN: cfg.PIN})
	if err != nil {
		return tls.Certificate{}, &close, fmt.Errorf("getting signer for %s: %w", cfg.Slot, err)
	}

	// Construct the cert chain
	chain := make([][]byte, 0, 1+len(cfg.Intermediates))
	chain = append(chain, leaf.Raw)
	chain = append(chain, cfg.Intermediates...)

	return tls.Certificate{
		Certificate: chain,
		PrivateKey:  signer,
		Leaf:        leaf,
	}, &close, nil
}

func containsFold(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (stringEqualFold(s, sub) || containsFold(s[1:], sub)))
}
func stringEqualFold(a, b string) bool {
	return len(a) == len(b) &&
		(a == b || (len(a) > 0 && (a[0]|32) == (b[0]|32) && stringEqualFold(a[1:], b[1:])))
}
