// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1package vault

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

		rm.logger.Debug("check", "window", re.RootCredential.Schedule.RotationWindow, "time", re.RootCredential.Schedule.NextVaultRotation)
		if !logical.DefaultScheduler.IsInsideRotationWindow(re.RootCredential.Schedule, now) {
			rm.logger.Debug("Not inside rotation window, pushing back to queue")
			err := rm.queue.Push(i)
			if err != nil {
				// this is pretty bad because we have no real way to fix it and save the item, but the Push operation only
				// errors on malformed items, which shouldn't be possible here
				return err
			}
			// don't break here, since the heap is keyed on priority, which is just timestamp.
			// it's possible for a later item to be ready for rotation, so we need to keep going
		}
		rm.logger.Debug("Item ready for rotation; making rotation request to sdk/backend")
		// do rotation
		req := &logical.Request{
			Operation: logical.RotationOperation,
			Path:      re.Path,
		}
		// TODO figure out how to get namespace with context here
		// ctx := namespace.ContextWithNamespace(rm.quitContext, n)
		rm.jobManager.AddJob(&rotationJob{
			rm:    rm,
			req:   req,
			entry: re,
		}, "best-queue-ever")
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

	rotationID := resp.RootCredential.RotationID
	var re *rotationEntry

	// either create a new rotationEntry or get the existing one
	if rotationID != "" {
		rm.logger.Debug("rotationID detected, this is an update", "rotationID", rotationID)
		entry, err := rm.queue.PopByKey(resp.RootCredential.RotationID)
		if err != nil {
			rm.logger.Warn("error popping item", "rotation_id", resp.RootCredential.RotationID, "err", err)
			return "", err
		}
		if entry != nil {
			// this is an update
			re = entry.Value.(*rotationEntry)
		}
	} else {
		// Create a rotation entry. We use TokenLength because that is what is used
		// by ExpirationManager

		// TODO: Check if we need to validate the root credential

		rotationRand, err := base62.Random(TokenLength)
		if err != nil {
			return "", err
		}
		ns, err := namespace.FromContext(ctx)
		rotationID = path.Join(req.Path, rotationRand)
		if ns.ID != namespace.RootNamespaceID {
			rotationID = fmt.Sprintf("%s.%s", rotationID, ns.ID)
		}
		re = &rotationEntry{
			RotationID:     rotationID,
			Path:           req.Path,
			Data:           resp.Data,
			RootCredential: resp.RootCredential,
			namespace:      ns,
		}
	}

	issueTime := time.Now()
	expireTime := time.Now()
	if resp.RootCredential.Schedule.Schedule != nil {
		expireTime = logical.DefaultScheduler.NextRotationTimeFromInput(resp.RootCredential.Schedule, time.Now())
		resp.RootCredential.Schedule.NextVaultRotation = expireTime
	}
	rm.logger.Debug("SCHEDULE", "VALUE", resp.RootCredential.Schedule.RotationSchedule)
	rm.logger.Debug("WINDOW", "VALUE", resp.RootCredential.Schedule.RotationWindow)
	rm.logger.Debug("TTL", "VALUE", resp.RootCredential.Schedule.TTL)

	re.ExpireTime = expireTime
	re.IssueTime = issueTime

	// lock and populate the queue
	// @TODO figure out why locking is leading to infinite loop
	// r.core.stateLock.Lock()

	rm.logger.Debug("Creating queue item")
	rm.logger.Debug("next rotation time", "exp", expireTime.Format(time.RFC3339))

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
	rm.logger.Debug("", "rotationID", re.RotationID)
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

// rotationJob implements fairshare.Job
//
// if you do queue management here you _must_ lock
type rotationJob struct {
	rm    *RotationManager
	req   *logical.Request
	entry *rotationEntry
}

// Execute is an implementation of fairshare.Job.Execute and in this case handles requesting rotation from
// the backend. It will return an error both in the case of a direct error, and in the case of certain kinds
// of error-shaped logical.Response returns.
func (j *rotationJob) Execute() error {
	j.rm.logger.Info("in execute")
	_, err := j.rm.router.Route(j.rm.quitContext, j.req)

	// TODO: clean up this branch
	if errors.Is(err, logical.ErrUnsupportedOperation) {
		j.rm.logger.Info("unsupported")
		return err
	} else if err != nil {
		// requeue with backoff
		j.rm.logger.Info("other rotate error", "err", err)
		return err
	}

	// TODO: inspect logical.Response for other error-y things (there may not be any)

	// success
	j.rm.logger.Debug("Successfully called rotate root code for backend")
	issueTime := time.Now()
	j.entry.RootCredential.Schedule.LastVaultRotation = issueTime
	expireTime := logical.DefaultScheduler.NextRotationTime(j.entry.RootCredential.Schedule)
	newEntry := &rotationEntry{
		RotationID:     j.entry.RotationID,
		Path:           j.entry.Path,
		Data:           j.entry.Data,
		RootCredential: j.entry.RootCredential,
		IssueTime:      issueTime,
		ExpireTime:     expireTime,
		namespace:      j.entry.namespace,
	}
	j.entry.RootCredential.Schedule.NextVaultRotation = newEntry.ExpireTime

	// lock and populate the queue
	j.rm.mu.Lock()

	item := &queue.Item{
		// will preserve same rotation ID, only updating Value, Priority with new rotation time
		Key:      newEntry.RotationID,
		Value:    newEntry,
		Priority: newEntry.ExpireTime.Unix(),
	}

	j.rm.logger.Debug("Pushing item into credential queue")
	j.rm.logger.Debug("will rotate at", "time", newEntry.ExpireTime.Format(time.RFC3339))

	if err := j.rm.queue.Push(item); err != nil {
		// TODO handle error
		j.rm.logger.Debug("Error pushing item into credential queue")
		return err
	}
	j.rm.mu.Unlock()

	return nil
}

// OnFailure implements the OnFailure interface method and requeues with a backoff when it happens
func (j *rotationJob) OnFailure(err error) {
	j.rm.logger.Info("rotation failed, requeuing", "error", err)

	err = j.rm.queue.Push(&queue.Item{
		Key:      j.entry.RotationID,
		Value:    j.entry,
		Priority: time.Now().Add(10 * time.Second).Unix(), // TODO: Configure this
	})
	// an error here is really bad because we can't really fix it and will lose the rotation entry as a result.
	if err != nil {
		j.rm.logger.Error("can't requeue an item", "id", j.entry.RotationID)
	}
}
