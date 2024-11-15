package radix

import (
	"context"
	"strconv"
	"strings"

	"errors"

	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

// Scanner is used to iterate through the results of a SCAN call (or HSCAN,
// SSCAN, etc...)
//
// Once created, repeatedly call Next() on it to fill the passed in string
// pointer with the next result. Next will return false if there's no more
// results to retrieve or if an error occurred, at which point Close should be
// called to retrieve any error.
type Scanner interface {
	Next(context.Context, *string) bool
	Close() error
}

// ScannerConfig is used to create Scanner instances with particular settings. All
// fields are optional, all methods are thread-safe.
type ScannerConfig struct {
	// The scan command to do, e.g. "SCAN", "HSCAN", etc...
	//
	// Defaults to "SCAN".
	Command string

	// The key to perform the scan on. Only necessary when Command isn't "SCAN"
	Key string

	// An optional pattern to filter returned keys by
	Pattern string

	// An optional count hint to send to redis to indicate number of keys to
	// return per call. This does not affect the actual results of the scan
	// command, but it may be useful for optimizing certain datasets
	Count int

	// An optional type name to filter for values of the given type.
	// The type names are the same as returned by the "TYPE" command.
	// This if only available in Redis 6 or newer and only works with "SCAN".
	// If used with an older version of Redis or with a Command other than
	// "SCAN", scanning will fail.
	Type string
}

func (cfg ScannerConfig) withDefaults() ScannerConfig {
	if cfg.Command == "" {
		cfg.Command = "SCAN"
	}
	return cfg
}

func (cfg ScannerConfig) cmd(rcv interface{}, cursor string) Action {
	cmdStr := strings.ToUpper(cfg.Command)
	args := make([]string, 0, 8)
	if cmdStr != "SCAN" {
		args = append(args, cfg.Key)
	}

	args = append(args, cursor)
	if cfg.Pattern != "" {
		args = append(args, "MATCH", cfg.Pattern)
	}
	if cfg.Count > 0 {
		args = append(args, "COUNT", strconv.Itoa(cfg.Count))
	}
	if cfg.Type != "" {
		args = append(args, "TYPE", cfg.Type)
	}

	return Cmd(rcv, cmdStr, args...)
}

type scanner struct {
	Client
	cfg    ScannerConfig
	res    scanResult
	resIdx int
	err    error
}

// New creates a new Scanner instance which will iterate over the redis
// instance's Client using the ScannerConfig.
func (cfg ScannerConfig) New(c Client) Scanner {
	return &scanner{
		Client: c,
		cfg:    cfg.withDefaults(),
		res: scanResult{
			cur: "0",
		},
	}
}

func (s *scanner) Next(ctx context.Context, res *string) bool {
	for {
		if s.err != nil {
			return false
		}

		for s.resIdx < len(s.res.keys) {
			*res = s.res.keys[s.resIdx]
			s.resIdx++
			if *res != "" {
				return true
			}
		}

		if s.res.cur == "0" && s.res.keys != nil {
			return false
		}

		s.err = s.Client.Do(ctx, s.cfg.cmd(&s.res, s.res.cur))
		s.resIdx = 0
	}
}

func (s *scanner) Close() error {
	return s.err
}

type scanResult struct {
	cur  string
	keys []string
}

func (s *scanResult) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var ah resp3.ArrayHeader
	if err := ah.UnmarshalRESP(br, o); err != nil {
		return err
	} else if ah.NumElems != 2 {
		return errors.New("not enough parts returned")
	}

	var c resp3.BlobString
	if err := c.UnmarshalRESP(br, o); err != nil {
		return err
	}

	s.cur = c.S
	s.keys = s.keys[:0]

	return resp3.Unmarshal(br, &s.keys, o)
}
