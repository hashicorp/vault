package vault

import (
	"context"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/fairshare"
	"github.com/hashicorp/vault/helper/namespace"
)

type RotationManager struct {
	core   *Core
	router *Router
	idView *BarrierView
	logger log.Logger

	rotateFunc RotationStrategy
	jobManager *fairshare.JobManager
}

type RotationStrategy func(context.Context, *RotationManager, string, *namespace.Namespace)

// rotationJob should only be created through newRotationJob()
type rotationJob struct {
	rotationID string
	ns         *namespace.Namespace
	r          *RotationManager
	nsCtx      context.Context
	startTime  time.Time
}

func newRotationJob(nsCtx context.Context, rotationID string, ns *namespace.Namespace, r *RotationManager) (*rotationJob, error) {
	if rotationID == "" {
		return nil, fmt.Errorf("cannot have empty lease id")
	}
	if r == nil {
		return nil, fmt.Errorf("cannot have nil rotation manager")
	}
	if nsCtx == nil {
		return nil, fmt.Errorf("cannot have nil namespace context.Context")
	}

	return &rotationJob{
		rotationID: rotationID,
		ns:         ns,
		r:          r,
		nsCtx:      nsCtx,
		startTime:  time.Now(),
	}, nil
}

// NewRotationManger creates a new RotationManager that is backed
// using a given view, and uses the provided router for rotation.
func NewRotationManger(c *Core, view *BarrierView, r RotationStrategy, logger log.Logger) *RotationManager {
	managerLogger := logger.Named("job-manager")
	jobManager := fairshare.NewJobManager("rotate", getNumRotationWorkers(c, logger), managerLogger, c.metricSink)
	jobManager.Start()

	c.AddLogger(managerLogger)

	rot := &RotationManager{
		core:       c,
		router:     c.router,
		idView:     view.SubView(leaseViewPrefix),
		logger:     logger,
		rotateFunc: r,
		jobManager: jobManager,
	}

	if rot.logger == nil {
		opts := log.LoggerOptions{Name: "rotation_manager"}
		rot.logger = log.New(&opts)
	}

	return rot
}

func getNumRotationWorkers(_ *Core, _ log.Logger) int {
	// TODO: make this configurable? See getNumExpirationWorkers
	numWorkers := 200
	return numWorkers
}
