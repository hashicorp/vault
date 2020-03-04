package api

import (
	"bufio"
	"context"
)

func (c *Sys) Monitor(logLevel string, stopCh chan struct{}) (chan string, error) {
	r := c.c.NewRequest("GET", "/v1/sys/monitor")

	if logLevel == "" {
		r.Params.Add("log_level", "INFO")
	} else {
		r.Params.Add("log_level", logLevel)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}

	logCh := make(chan string, 64)
	go func() {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)

		for {
			select {
			case <-stopCh:
				close(logCh)
				cancelFunc()
				return
			case <-ctx.Done():
				stopCh <- struct{}{}
				close(logCh)
				cancelFunc()
				return
			default:
			}

			if scanner.Scan() {
				// An empty string signals to the caller that
				// the scan is done, so make sure we only emit
				// that when the scanner says it's done, not if
				// we happen to ingest an empty line.
				if text := scanner.Text(); text != "" {
					logCh <- text
				} else {
					logCh <- " "
				}
			} else {
				// If Scan() returns false, that means the context deadline was exceeded, so
				// terminate this routine and start a new request.
				stopCh <- struct{}{}
				close(logCh)
				cancelFunc()
				return
			}
		}
	}()

	return logCh, nil
}
