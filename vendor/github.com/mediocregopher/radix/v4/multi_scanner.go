package radix

import (
	"context"
	"strings"
)

type multiScanner struct {
	multiClient MultiClient
	cfg         ScannerConfig

	clients     []Client
	currScanner Scanner
	lastErr     error
}

// NewMulti returns a Scanner which will scan over every primary instance in the
// MultiClient. This will panic if the ScanOpt's Command isn't "SCAN".
//
// NOTE this is primarily useful for scanning over all keys in a Cluster. It is
// not necessary to use this otherwise, unless you have implemented your own
// MultiClient which holds multiple primary Clients.
func (cfg ScannerConfig) NewMulti(mc MultiClient) Scanner {
	cfg = cfg.withDefaults()
	if strings.ToUpper(cfg.Command) != "SCAN" {
		panic("NewMulti can only perform SCAN operations")
	}

	clientsM, err := mc.Clients()
	if err != nil {
		return &multiScanner{lastErr: err}
	}

	clients := make([]Client, 0, len(clientsM))
	for _, replicaSet := range clientsM {
		clients = append(clients, replicaSet.Primary)
	}

	cs := &multiScanner{
		multiClient: mc,
		cfg:         cfg,
		clients:     clients,
	}
	cs.nextScanner()

	return cs
}

func (cs *multiScanner) closeCurr() {
	if cs.currScanner != nil {
		if err := cs.currScanner.Close(); err != nil && cs.lastErr == nil {
			cs.lastErr = err
		}
		cs.currScanner = nil
	}
}

func (cs *multiScanner) nextScanner() {
	cs.closeCurr()
	if len(cs.clients) == 0 {
		return
	}
	client := cs.clients[0]
	cs.clients = cs.clients[1:]
	cs.currScanner = cs.cfg.New(client)
}

func (cs *multiScanner) Next(ctx context.Context, res *string) bool {
	for {
		if cs.currScanner == nil {
			return false
		} else if out := cs.currScanner.Next(ctx, res); out {
			return true
		}
		cs.nextScanner()
	}
}

func (cs *multiScanner) Close() error {
	cs.closeCurr()
	return cs.lastErr
}
