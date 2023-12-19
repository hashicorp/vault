package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/vault/helper/fairshare"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/logical"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/queue"
)

const (
	fairshareRotationWorkersOverrideVar = "VAULT_CREDENTIAL_ROTATION_WORKERS"
)

type RotationManager struct {
	core   *Core
	logger log.Logger
	mu     sync.Mutex

	jobManager  *fairshare.JobManager
	queue       *queue.PriorityQueue
	done        chan struct{}
	quitContext context.Context

	router   *Router
	backends func() *[]MountEntry // list of logical and auth backends, remember to call RUnlock
}

// rotationEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type rotationEntry struct {
	RotationID     string                  `json:"rotation_id"`
	Path           string                  `json:"path"`
	Data           map[string]interface{}  `json:"data"`
	RootCredential *logical.RootCredential `json:"static_secret"`
	IssueTime      time.Time               `json:"issue_time"`
	ExpireTime     time.Time               `json:"expire_time"`

	namespace *namespace.Namespace
}

func (rm *RotationManager) Start() error {
	done := rm.done
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		rm.logger.Info("started ticker")
		for {
			// rm.mu.Lock()
			select {
			case <-done:
				rm.logger.Debug("done with loop; received from channel")
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

// Stop is used to prevent further automatic rotations.
func (rm *RotationManager) Stop() error {
	// Stop all the pending rotation timers
	rm.logger.Debug("stop triggered")
	defer rm.logger.Debug("finished stopping")

	rm.jobManager.Stop()

	// close done channel
	close(rm.done)

	return nil
}

func (rm *RotationManager) CheckQueue() error {
	// loop runs forever, so break whenever you get to the first credential that doesn't need updating
	for {
		now := time.Now()
		i, err := rm.queue.Pop()
		if err != nil {
			rm.logger.Info("automated rotation queue empty")
			return nil
		}

		if i.Priority > now.Unix() {
			rm.logger.Debug("Item not ready for rotation; adding back to queue")
			err := rm.queue.Push(i)
			if err != nil {
				// this is pretty bad because we have no real way to fix it and save the item, but the Push operation only
				// errors on malformed items, which shouldn't be possible here
				return err
			}
			break // this item is not ripe yet, which means all later items are also unripe, so exit the check loop
		}

		var re *rotationEntry
		entry, ok := i.Value.(*rotationEntry)
		if !ok {
			return fmt.Errorf("error parsing rotation entry from queue")
		}

		re = entry

		// TODO should we push the credential back into the queue if it is not in the rotation window?
		// if not in window, do we check the next credential?
		if !logical.DefaultScheduler.IsInsideRotationWindow(re.RootCredential.Schedule, now) {
			rm.logger.Debug("Not inside rotation window, pushing back to queue")
			err := rm.queue.Push(i)
			if err != nil {
				// this is pretty bad because we have no real way to fix it and save the item, but the Push operation only
				// errors on malformed items, which shouldn't be possible here
				return err
			}
			break
		}
		rm.logger.Debug("Item ready for rotation; making rotation request to sdk/backend")
		// do rotation
		req := &logical.Request{
			Operation: logical.RotationOperation,
			Path:      re.Path,
		}
		// TODO figure out how to get namespace with context here
		// ctx := namespace.ContextWithNamespace(rm.quitContext, n)
		_, err = rm.router.Route(rm.quitContext, req)
		if errors.Is(err, logical.ErrUnsupportedOperation) {
			rm.logger.Info("unsupported")
			continue
		} else if err != nil {
			// requeue with backoff
			rm.logger.Info("other rotate error", "err", err)
			// TODO figure out rollback procedure here, for now, pushing back into queue
			// one idea is to create a new rotation entry with updated priority and try later
		}

		// success
		rm.logger.Debug("Successfully called rotate root code for backend")
		issueTime := time.Now()
		newEntry := &rotationEntry{
			RotationID:     re.RotationID,
			Path:           re.Path,
			Data:           re.Data,
			RootCredential: re.RootCredential,
			IssueTime:      issueTime,
			// expires the next time the schedule is activated from the issue time
			ExpireTime: re.RootCredential.Schedule.Schedule.Next(issueTime),
			namespace:  re.namespace,
		}

		// lock and populate the queue
		rm.core.stateLock.Lock()

		item := &queue.Item{
			// will preserve same rotation ID, only updating Value, Priority with new rotation time
			Key:      newEntry.RotationID,
			Value:    newEntry,
			Priority: newEntry.ExpireTime.Unix(),
		}

		rm.logger.Debug("Pushing item into credential queue")

		if err := rm.queue.Push(item); err != nil {
			// TODO handle error
			rm.logger.Debug("Error pushing item into credential queue")
			return err
		}
		rm.core.stateLock.Unlock()
		if err != nil {
			// again, this is bad because we can't really fix the item, but it also shouldn't happen because the item was good before
			return err
		}
	}

	return nil
}

// Register takes a request and response with an associated StaticSecret. The
// secret gets assigned a RotationID and the management of the rotation is
// assumed by the rotation manager.
func (rm *RotationManager) Register(ctx context.Context, req *logical.Request, resp *logical.Response) (id string, retErr error) {
	rm.logger.Debug("Starting registration")

	// Ignore if there is no root cred
	if resp == nil || resp.RootCredential == nil {
		return "", nil
	}

	// TODO: Check if we need to validate the root credential

	// Create a rotation entry. We use TokenLength because that is what is used
	// by ExpirationManager
	rm.logger.Debug("Generating random rotation ID")
	rotationRand, err := base62.Random(TokenLength)
	if err != nil {
		return "", err
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", err
	}

	rotationID := path.Join(req.Path, rotationRand)

	if ns.ID != namespace.RootNamespaceID {
		rotationID = fmt.Sprintf("%s.%s", rotationID, ns.ID)
	}

	issueTime := time.Now()
	re := &rotationEntry{
		RotationID:     rotationID,
		Path:           req.Path,
		Data:           resp.Data,
		RootCredential: resp.RootCredential,
		IssueTime:      issueTime,
		// expires the next time the schedule is activated from the issue time
		ExpireTime: resp.RootCredential.Schedule.Schedule.Next(issueTime),
		namespace:  ns,
	}

	// lock and populate the queue
	// @TODO figure out why locking is leading to infinite loop
	// r.core.stateLock.Lock()

	rm.logger.Debug("Creating queue item")

	// @TODO for different cases, update rotation entry if it is already in queue
	// for now, assuming it is a fresh root credential and the schedule is not being updated
	item := &queue.Item{
		Key:      re.RotationID,
		Value:    re,
		Priority: re.ExpireTime.Unix(),
	}

	rm.logger.Debug("Pushing item into credential queue")

	if err := rm.queue.Push(item); err != nil {
		// TODO handle error
		rm.logger.Debug("Error pushing item into credential queue")
		return "", err
	}

	// r.core.stateLock.Unlock()
	return re.RotationID, nil
}

func getNumRotationWorkers(c *Core, l log.Logger) int {
	numWorkers := c.numExpirationWorkers

	workerOverride := os.Getenv(fairshareRotationWorkersOverrideVar)
	if workerOverride != "" {
		i, err := strconv.Atoi(workerOverride)
		if err != nil {
			l.Warn("vault rotation workers override must be an integer", "value", workerOverride)
		} else if i < 1 || i > 10000 {
			l.Warn("vault rotation workers override out of range", "value", i)
		} else {
			numWorkers = i
		}
	}

	return numWorkers
}

func (c *Core) startRotation() error {
	logger := c.baseLogger.Named("rotation-job-manager")

	jobManager := fairshare.NewJobManager("rotate", getNumRotationWorkers(c, logger), logger, c.metricSink)
	jobManager.Start()

	c.AddLogger(logger)
	c.rotationManager = &RotationManager{
		core:   c,
		logger: logger,
		// TODO figure out how to populate this if credentials already exist after unseal
		queue:       queue.New(),
		done:        make(chan struct{}),
		jobManager:  jobManager,
		quitContext: c.activeContext,
		router:      c.router,
	}
	err := c.rotationManager.Start()
	if err != nil {
		return err
	}
	return nil
}

// stopRotation is used to stop the rotation manager before
// sealing Vault.
func (c *Core) stopRotation() error {
	if c.rotationManager != nil {
		if err := c.rotationManager.Stop(); err != nil {
			return err
		}
		c.metricsMutex.Lock()
		defer c.metricsMutex.Unlock()
		c.rotationManager = nil
	}
	return nil
}
