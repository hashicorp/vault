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

	router   *Router
	backends func() *[]MountEntry // list of logical and auth backends, remember to call RUnlock
}

func (rm *RotationManager) Start() error {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		rm.logger.Info("started ticker")
		for {
			rm.mu.Lock()
			select {
			case <-rm.done:
				return
			case t := <-ticker.C:
				rm.logger.Info("time", "time", t.Format(time.RFC3339))
				err := rm.CheckQueue()
				if err != nil {
					rm.logger.Error("check queue error", "err", err)
				}
			}
		}
	}()
	return nil
}

func (rm *RotationManager) CheckQueue() error {
	// loop runs forever, so break whenever you get to the first credential that doesn't need updating
	for {
		now := time.Now()
		i, err := rm.queue.Pop()
		if err != nil {
			rm.logger.Info("queue empty")
			return nil
		}
		if i.Priority > now.Unix() {
			err := rm.queue.Push(i)
			if err != nil {
				// this is pretty bad because we have no real way to fix it and save the item, but the Push operation only
				// errors on malformed items, which shouldn't be possible here
				return err
			}
			break // this item is not ripe yet, which means all later items are also unripe, so exit the check loop
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
		} else if err != nil {
			// requeue with backoff
			rm.logger.Info("other rotate error", "err", err)
			// TODO: We can either check the window here, or let the priority check above handle it
			i.Priority = i.Priority + 10
		}

		// success
		i.Priority = time.Now().Add(5 * time.Minute).Unix() // TODO: here we want to access the schedule and update the priority based on that
		err = rm.queue.Push(i)
		if err != nil {
			// again, this is bad because we can't really fix the item, but it also shouldn't happen because the item was good before
			return err
		}
	}

	return nil
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
