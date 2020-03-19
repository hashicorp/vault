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

	m, _ := NewMonitor(512, logger, &log.LoggerOptions{
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
	case <-time.After(5 * time.Second):
		t.Fatal("Expected to receive from log channel")
	}
}

func TestMonitor_Start_Unbuffered(t *testing.T) {
	t.Parallel()

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Level: log.Error,
	})

	_, err := NewMonitor(0, logger, &log.LoggerOptions{
		Level: log.Debug,
	})

	if err == nil {
		t.Fatal("expected to get an error, but didn't")
	} else {
		if !strings.Contains(err.Error(), "greater than zero") {
			t.Fatal("expected an error about buf being greater than zero")
		}
	}
}

// Ensure number of dropped messages are logged
func TestMonitor_DroppedMessages(t *testing.T) {
	t.Parallel()

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Level: log.Warn,
	})

	m, _ := newMonitor(5, logger, &log.LoggerOptions{
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
