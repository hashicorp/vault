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

	bufferSizes := []int{512, 0}

	for _, b := range bufferSizes {
		b := b

		t.Run(fmt.Sprintf("Start_with_buffer_size_%d", b), func(t *testing.T) {
			t.Parallel()

			logger := log.NewInterceptLogger(&log.LoggerOptions{
				Level: log.Error,
			})

			m := NewMonitor(b, logger, &log.LoggerOptions{
				Level: log.Debug,
			})

			logCh := m.Start()
			defer m.Stop()

			go func() {
				logger.Debug("test log")
				time.Sleep(10 * time.Millisecond)
			}()

			select {
			case l := <-logCh:
				require.Contains(t, string(l), "[DEBUG] test log")
				return
			case <-time.After(3 * time.Second):
				t.Fatal("Expected to receive from log channel")
			}
		})
	}
}

// Ensure number of dropped messages are logged
func TestMonitor_DroppedMessages(t *testing.T) {
	t.Parallel()

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Level: log.Warn,
	})

	m := newMonitor(5, logger, &log.LoggerOptions{
		Level: log.Debug,
	})
	m.dropCheckInterval = 5 * time.Millisecond

	logCh := m.Start()
	defer m.Stop()

	for i := 0; i <= 100; i++ {
		logger.Debug(fmt.Sprintf("test message %d", i))
	}

	passed := make(chan struct{})
	go func() {
		for recv := range logCh {
			if strings.Contains(string(recv), "[WARN] Monitor dropped") {
				close(passed)
				return
			}
		}
	}()

	select {
	case <-passed:
	case <-time.After(2 * time.Second):
		require.Fail(t, "expected to see warn dropped messages")
	}
}
