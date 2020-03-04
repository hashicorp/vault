package monitor

import (
	"fmt"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestMonitor_Start(t *testing.T) {
	t.Parallel()

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Level: log.Error,
	})

	m := New(512, logger, &log.LoggerOptions{
		Level: log.Debug,
	})

	logCh := m.Start()
	defer m.Stop()

	go func() {
		logger.Debug("test log")
		time.Sleep(10 * time.Millisecond)
	}()

	for {
		select {
		case log := <-logCh:
			require.Contains(t, string(log), "[DEBUG] test log")
			return
		case <-time.After(3 * time.Second):
			t.Fatal("Expected to receive from log channel")
		}
	}
}

// Ensure number of dropped messages are logged
func TestMonitor_DroppedMessages(t *testing.T) {
	t.Parallel()

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Level: log.Warn,
	})

	m := new(5, logger, &log.LoggerOptions{
		Level: log.Debug,
	})
	m.droppedDuration = 5 * time.Millisecond

	doneCh := make(chan struct{})
	defer close(doneCh)

	logCh := m.Start()

	for i := 0; i <= 100; i++ {
		logger.Debug(fmt.Sprintf("test message %d", i))
	}

	received := ""
	passed := make(chan struct{})
	go func() {
		for {
			select {
			case recv := <-logCh:
				received += string(recv)
				if strings.Contains(received, "[WARN] Monitor dropped") {
					close(passed)
				}
			}
		}
	}()

TEST:
	for {
		select {
		case <-passed:
			break TEST
		case <-time.After(2 * time.Second):
			require.Fail(t, "expected to see warn dropped messages")
		}
	}
}
