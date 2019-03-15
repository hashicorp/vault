// +build !windows

package command

import (
	"os"
	"os/signal"
	"syscall"
)

// MakeSigUSR2Ch returns a channel that can be used for SIGUSR2
// goroutine logging. This channel will send a message for every
// SIGUSR2 received.
func MakeSigUSR2Ch() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGUSR2)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}
