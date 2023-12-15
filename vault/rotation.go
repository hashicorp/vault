package vault

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/logical"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/queue"
	"github.com/robfig/cron/v3"
)

const parseOptions = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow

var parser = cron.NewParser(parseOptions)

type RotationManager struct {
	logger log.Logger
	mu     sync.Mutex

	queue queue.PriorityQueue
	done  chan struct{}

	router *Router
}

func (rm *RotationManager) Start() error {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		rm.logger.Info("started ticker")
		for {
			// hack to solve weird double loading
			rm.mu.Lock()
			select {
			case <-rm.done:
				return
			case t := <-ticker.C:
				rm.logger.Info("time", "time", t.Format(time.RFC3339))
				rm.CheckQueue()
			}
			rm.mu.Unlock()
		}
	}()
	return nil
}

func (rm *RotationManager) CheckQueue() {
	for {
		now := time.Now()
		i, err := rm.queue.Pop()
		if err != nil {
			return
		}
		if i.Priority > now.Unix() {
			rm.queue.Push(i)
			break
		}

		// do rotation
		req := &logical.Request{
			Operation: logical.RotationOperation,
			Path:      "path",
		}
		_, err = rm.router.Route(context.Background(), req)
		if errors.Is(err, logical.ErrUnsupportedOperation) {
			rm.logger.Info("unsupported")
			continue
		}
	}
}

// A RotationSchedule is a way to store the requested rotation schedule of a credential
type RotationSchedule struct {
	s cron.Schedule
}

func ParseSchedule(s string) (*RotationSchedule, error) {
	c, err := parser.Parse(s)
	if err != nil {
		return nil, err
	}

	return &RotationSchedule{
		s: c,
	}, nil
}

func (c *Core) startRotation() error {
	logger := c.baseLogger.Named("rotation")
	c.AddLogger(logger)
	c.rotationManager = &RotationManager{
		logger: logger,
		queue:  queue.PriorityQueue{},
		done:   make(chan struct{}),
	}
	err := c.rotationManager.Start()
	if err != nil {
		return err
	}
	return nil

}
