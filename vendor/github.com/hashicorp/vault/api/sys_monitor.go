package api

import (
	"bufio"
	"context"
)

// This function returns a channel that outputs strings containing the log messages
// coming from the server.
func (c *Sys) Monitor(ctx context.Context, logLevel string) (chan string, error) {
	r := c.c.NewRequest("GET", "/v1/sys/monitor")

	if logLevel == "" {
		r.Params.Add("log_level", "INFO")
	} else {
		r.Params.Add("log_level", logLevel)
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}

	logCh := make(chan string, 64)
	go func() {
		defer close(logCh)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)

		for {
			if ctx.Err() != nil {
				return
			}

			if !scanner.Scan() {
				return
			}

			// An empty string signals to the caller that
			// the scan is done, so make sure we only emit
			// that when the scanner says it's done, not if
			// we happen to ingest an empty line.
			if text := scanner.Text(); text != "" {
				logCh <- text
			}
		}
	}()

	return logCh, nil
}
